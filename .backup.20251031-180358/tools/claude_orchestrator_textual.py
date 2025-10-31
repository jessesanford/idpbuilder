#!/usr/bin/env python3

import argparse
import asyncio
import json
import os
import re
import shlex
import signal
import subprocess
from typing import Optional, List, Tuple, Union
CONTINUE_TRUE_PATTERN = re.compile(r"CONTINUE-SOFTWARE-FACTORY\s*=\s*TRUE", re.IGNORECASE)
CONTINUE_FALSE_PATTERN = re.compile(r"CONTINUE-SOFTWARE-FACTORY\s*=\s*FALSE", re.IGNORECASE)


from textual.app import App, ComposeResult
from textual.command import Provider, Hit, Hits, DiscoveryHit
from textual.widgets import Header, Footer, Static, Input
from textual.containers import Vertical, Horizontal, Container
from textual.reactive import reactive


def _clip(s: str, max_len: int = 240) -> str:
    return s if len(s) <= max_len else s[: max_len - 1] + "…"


def _extract_content(obj: dict) -> Optional[str]:
    message = obj.get("message") if isinstance(obj, dict) else None
    if isinstance(message, dict):
        content = message.get("content")
        if isinstance(content, str):
            return content
        if isinstance(content, list):
            parts: List[str] = []
            for it in content:
                if isinstance(it, str):
                    parts.append(it)
                elif isinstance(it, dict):
                    val = it.get("text") or it.get("content")
                    if isinstance(val, str):
                        parts.append(val)
            if parts:
                return " ".join(parts)
    for key in ("content", "text", "result"):
        val = obj.get(key)
        if isinstance(val, str):
            return val
    tool_result = obj.get("tool_result") or obj.get("toolOutput")
    if isinstance(tool_result, dict):
        return (
            tool_result.get("stdout")
            or tool_result.get("output")
            or tool_result.get("result")
            or tool_result.get("content")
        )
    return None


def _summarize(obj: dict) -> str:
    t = obj.get("type")
    subtype = obj.get("subtype") or obj.get("role")
    msg_id = obj.get("id") or obj.get("uuid")
    is_error = obj.get("is_error")
    duration = obj.get("duration_ms") or obj.get("latency_ms")
    base = []
    if msg_id:
        base.append(f"id={msg_id}")
    if t:
        base.append(f"{t}/{subtype}" if subtype else f"{t}")
    if is_error is True:
        base.append("ERROR")
    if duration is not None:
        base.append(f"{duration}ms")
    extra = []
    if t == "result":
        total_cost = obj.get("total_cost_usd")
        if isinstance(total_cost, (int, float)):
            extra.append(f"${total_cost:.4f}")
        text = _extract_content(obj)
        if isinstance(text, str) and "CONTINUE-SOFTWARE-FACTORY" in text:
            extra.append("decision")
    elif t in ("message", "user", "assistant"):
        text = _extract_content(obj)
        if text:
            extra.append(_clip(text))
    elif t and "tool" in str(t):
        tool_name = obj.get("tool") or obj.get("tool_name")
        if tool_name:
            extra.append(str(tool_name))
        text = _extract_content(obj)
        if text:
            extra.append(_clip(text))
    else:
        text = _extract_content(obj)
        if text:
            extra.append(_clip(text))
    summary = " | ".join(base)
    if extra:
        summary += " | " + " | ".join(extra)
    return summary or json.dumps(obj)


def _scan_strings_for_marker(value: Union[dict, list, str, int, float, None]) -> Tuple[bool, bool]:
    """Return (saw_true, saw_false) if any string within value contains the decision marker.
    Searches shallowly with limited depth to avoid heavy cost."""
    saw_true = False
    saw_false = False
    def _walk(v: Union[dict, list, str, int, float, None], depth: int = 0) -> None:
        nonlocal saw_true, saw_false
        if depth > 3:
            return
        if isinstance(v, str):
            if CONTINUE_TRUE_PATTERN.search(v):
                saw_true = True
            if CONTINUE_FALSE_PATTERN.search(v):
                saw_false = True
            return
        if isinstance(v, dict):
            for _, vv in v.items():
                _walk(vv, depth + 1)
            return
        if isinstance(v, list):
            for vv in v:
                _walk(vv, depth + 1)
            return
    _walk(value, 0)
    return saw_true, saw_false


def _load_state_machine(sm_path: str) -> Optional[dict]:
    try:
        with open(sm_path, "r", encoding="utf-8") as f:
            return json.load(f)
    except Exception:
        return None


def _load_current_state(state_path: str) -> Optional[str]:
    try:
        with open(state_path, "r", encoding="utf-8") as f:
            data = json.load(f)
            for key in ("current_state", "state", "current"):
                if key in data and isinstance(data[key], str):
                    return data[key]
    except Exception:
        return None


def _render_state_graph(sm: Optional[dict], current_state: Optional[str]) -> List[str]:
    if not isinstance(sm, dict):
        return ["State machine not found or invalid"]
    raw_states = sm.get("states") or sm.get("States") or sm.get("nodes")
    names: List[str] = []
    if isinstance(raw_states, dict):
        names = list(raw_states.keys())
    elif isinstance(raw_states, list):
        for it in raw_states:
            if isinstance(it, str):
                names.append(it)
            elif isinstance(it, dict):
                n = it.get("name") or it.get("id") or it.get("state")
                if isinstance(n, str):
                    names.append(n)
    names = sorted(set(names))
    out: List[str] = ["State Machine"]
    if not names:
        out.append("(no states)")
        return out
    cell_w = max(8, min(16, max((len(n) for n in names), default=8) + 2))
    # Escape brackets for Rich markup rendering in Textual by doubling them
    row = " ".join(f"[[{n[: cell_w - 2].ljust(cell_w - 2)}]]" for n in names)
    out.append(row)
    if current_state and current_state in names:
        idx = names.index(current_state)
        start_col = idx * (cell_w + 1)
        pointer = " " * max(0, start_col + (cell_w // 2) - 1) + "^ YOU ARE HERE"
        out.append(pointer)
    return out


class GraphWidget(Static):
    async def update_graph(self, sm_path: str, state_path: str) -> None:
        sm = _load_state_machine(sm_path)
        cur = _load_current_state(state_path)
        header = f"Graph: {sm_path} (state: {state_path})"
        body = "\n".join(_render_state_graph(sm, cur))
        self.update(header + "\n" + body)


class StreamWidget(Static):
    async def set_text(self, text: str) -> None:
        self.update(text)

class OrchestratorApp(App):
    COMMANDS = App.COMMANDS | set()
    CSS = """
    Screen { layout: vertical; }
    #status { height: 1; }
    #help { height: 1; }
    #banner { height: 1; }
    #graph { height: 8; overflow: hidden auto; border: round; padding: 0 1; }
    #stream { overflow: auto; border: round; padding: 0 1; }
    """

    wrap = reactive(False)
    details = reactive(False)
    search = reactive("")

    def __init__(self, cmd: List[str], sm_path: str, state_path: str, log_file: Optional[str], *, details_default: bool = False):
        super().__init__()
        self.cmd = cmd
        self.sm_path = sm_path
        self.state_path = state_path
        self.log_file = log_file
        self._details_default = details_default
        self.proc: Optional[asyncio.subprocess.Process] = None
        self.should_continue: bool = False
        self.msg_count: int = 0
        self.continue_count: int = 0
        self.help = Static(id="help")
        self.banner = Static(id="banner")
        self.graph = GraphWidget(id="graph")
        self.stream = StreamWidget(id="stream")
        self.status = Static(id="status")
        self.input: Optional[Input] = None
        self._stream_text: str = ""
        self._stream_max_chars: int = 200_000
        self._log_disabled_message_shown: bool = False
        # Preserve messages for re-render (summary/detail/search/wrap)
        self._messages: List[object] = []

    def compose(self) -> ComposeResult:
        yield Header()
        with Container():
            yield self.status
            yield self.banner
            yield self.help
        yield self.graph
        yield self.stream
        yield Footer()

    async def on_mount(self) -> None:
        await self._reload_graph()
        # Apply initial detail mode if requested
        try:
            if self._details_default:
                self.details = True
        except Exception:
            pass
        # Show the command being executed so the screen isn't blank
        try:
            cmd_str = " ".join(shlex.quote(p) for p in self.cmd)
            await self._append_stream(f"Executing: {cmd_str}\n")
        except Exception:
            pass
        # Show key help
        self.help.update("Keys: q quit | d details | / search | w wrap | r reload | Ctrl+P: Command Palette")
        # Run the orchestrator loop in the background so UI stays responsive
        asyncio.create_task(self._run_forever())

    async def _run_forever(self) -> None:
        # Loop based on CONTINUE-SOFTWARE-FACTORY decision in the last result
        while True:
            await self._run_once()
            if not self.should_continue:
                break
            # Prepare for next iteration
            try:
                await self._append_stream("\n[continue] Starting next iteration...\n")
            except Exception:
                pass
            await self._reload_graph()

    async def _reload_graph(self) -> None:
        try:
            await self.graph.update_graph(self.sm_path, self.state_path)
        except Exception:
            # fallback inline render
            sm = _load_state_machine(self.sm_path)
            cur = _load_current_state(self.state_path)
            header = f"Graph: {self.sm_path} (state: {self.state_path})"
            body = "\n".join(_render_state_graph(sm, cur))
            self.graph.update(header + "\n" + body)

    async def _append_stream(self, text: str) -> None:
        out = f"{self._stream_text}{text}"
        if len(out) > self._stream_max_chars:
            out = out[-self._stream_max_chars:]
        rendered = self._apply_search_and_wrap(out)
        await self.stream.set_text(rendered)
        self._stream_text = out

    def _refresh_stream_render(self) -> None:
        # Re-apply wrap and search highlighting
        rendered = self._apply_search_and_wrap(self._stream_text)
        self.stream.update(rendered)

    def _apply_search_and_wrap(self, text: str) -> str:
        out = text
        # Wrap by terminal width (best-effort)
        if self.wrap:
            try:
                width = max(20, (self.size.width or 100) - 4)
                lines: List[str] = []
                for line in out.splitlines():
                    while len(line) > width:
                        lines.append(line[:width])
                        line = line[width:]
                    lines.append(line)
                out = "\n".join(lines)
            except Exception:
                pass
        # Search highlight
        if self.search:
            try:
                pattern = re.compile(self.search, re.IGNORECASE)
                out = pattern.sub(lambda m: f"[bold]{m.group(0)}[/]", out)
            except Exception:
                pass
        return out

    async def _run_once(self) -> None:
        if self.proc is not None:
            return
        # Start process
        self.proc = await asyncio.create_subprocess_exec(
            *self.cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )

        last_result_raw: Optional[str] = None
        saw_true_any = False
        saw_false_any = False
        saw_true_in_last_result = False
        saw_false_in_last_result = False
        log_fh = open(self.log_file, "a", encoding="utf-8") if self.log_file else None

        def safe_log_write(data: str) -> None:
            nonlocal log_fh
            if not log_fh:
                return
            try:
                log_fh.write(data)
            except OSError as e:
                # ENOSPC or similar: disable logging, inform user once
                try:
                    if not self._log_disabled_message_shown:
                        self._log_disabled_message_shown = True
                        self.status.update(f"Logging disabled: {e}")
                except Exception:
                    pass
                try:
                    log_fh.close()
                except Exception:
                    pass
                log_fh = None

        try:
            assert self.proc.stdout is not None
            buffer = ""
            while True:
                try:
                    chunk = await asyncio.wait_for(self.proc.stdout.read(4096), timeout=0.2)
                except asyncio.TimeoutError:
                    if self.proc.returncode is not None:
                        break
                    continue
                if not chunk:
                    if self.proc.returncode is not None:
                        break
                    continue
                if isinstance(chunk, (bytes, bytearray)):
                    chunk = chunk.decode("utf-8", errors="ignore")
                safe_log_write(chunk)
                buffer += chunk
                # Process complete lines
                while True:
                    nl = buffer.find("\n")
                    if nl == -1:
                        break
                    line = buffer[: nl + 1]
                    buffer = buffer[nl + 1 :]
                    try:
                        obj = json.loads(line)
                        if isinstance(obj, dict):
                            view = json.dumps(obj, indent=2) + "\n" if self.details else _summarize(obj) + "\n"
                            await self._append_stream(view)
                            # Store message model for re-rendering
                            self._messages.append(obj)
                            self.msg_count += 1
                            self._update_status()
                            # scan markers
                            t_true, t_false = _scan_strings_for_marker(obj)
                            if t_true:
                                saw_true_any = True
                            if t_false:
                                saw_false_any = True
                            if obj.get("type") == "result":
                                last_result_raw = line
                                saw_true_in_last_result = t_true
                                saw_false_in_last_result = t_false
                        else:
                            await self._append_stream(line)
                            self._messages.append(line)
                    except Exception:
                        await self._append_stream(line)
                        self._messages.append(line)

            # Flush any remaining buffered json without newline
            tail = buffer.strip()
            if tail:
                try:
                    obj = json.loads(tail)
                    if isinstance(obj, dict):
                        view = json.dumps(obj, indent=2) + "\n" if self.details else _summarize(obj) + "\n"
                        await self._append_stream(view)
                        self._messages.append(obj)
                        self.msg_count += 1
                        self._update_status()
                        t_true, t_false = _scan_strings_for_marker(obj)
                        if t_true:
                            saw_true_any = True
                        if t_false:
                            saw_false_any = True
                        if obj.get("type") == "result":
                            last_result_raw = tail
                            saw_true_in_last_result = t_true
                            saw_false_in_last_result = t_false
                    else:
                        await self._append_stream(tail + "\n")
                        self._messages.append(tail + "\n")
                except Exception:
                    await self._append_stream(tail + "\n")
                    self._messages.append(tail + "\n")

            # flush stderr
            if self.proc.stderr is not None:
                err = await self.proc.stderr.read()
                if isinstance(err, (bytes, bytearray)):
                    err = err.decode("utf-8", errors="ignore")
                if err:
                    safe_log_write(err)
                    await self._append_stream(err)
        finally:
            # Allow subsequent runs
            self.proc = None
            if log_fh:
                try:
                    log_fh.close()
                except Exception:
                    pass

        # decide continuation
        # Decide with precedence: explicit FALSE in last result > TRUE in last result > any FALSE > any TRUE
        self.should_continue = False
        if saw_false_in_last_result:
            self.should_continue = False
        elif saw_true_in_last_result:
            self.should_continue = True
        elif saw_false_any:
            self.should_continue = False
        elif saw_true_any:
            self.should_continue = True
        # Update banner
        try:
            if self.should_continue:
                self.banner.update("[green bold]CONTINUE-SOFTWARE-FACTORY = TRUE[/]")
                self.continue_count += 1
            else:
                self.banner.update("")
        except Exception:
            pass
        self._update_status()

    def _update_status(self) -> None:
        cont = "YES" if self.should_continue else "NO"
        try:
            self.status.update(f"Msgs: {self.msg_count} | Continue: {cont}")
        except Exception:
            pass

    # Action methods used by the command palette
    def toggle_details(self) -> None:
        self.details = not self.details
        self._update_status()

    def details_on(self) -> None:
        self.details = True
        self._update_status()

    def details_off(self) -> None:
        self.details = False
        self._update_status()

    async def prompt_search(self) -> None:
        self.status.update("Search regex: ")
        self.input = Input(placeholder="regex", id="search")
        await self.mount(self.input)
        await self.input.focus()

    def toggle_wrap(self) -> None:
        self.wrap = not self.wrap
        self._refresh_stream_render()

    async def reload_graph_command(self) -> None:
        await self._reload_graph()
        await self._append_stream("[reload] graph/state reloaded\n")

    async def on_key(self, event) -> None:
        key = event.key
        if key == "q":
            await self.action_quit()
        elif key == "w":
            self.wrap = not self.wrap
            self._refresh_stream_render()
        elif key == "d":
            self.details = not self.details
            # Inform user; detail affects new messages going forward
            mode = "ON" if self.details else "OFF"
            self.status.update(f"Msgs: {self.msg_count} | Continue: {'YES' if self.should_continue else 'NO'} | Details: {mode}")
        elif key == "r":
            await self._reload_graph()
            # Acknowledge reload
            await self._append_stream("[reload] graph/state reloaded\n")
        elif key == "/":
            # Simple prompt in status
            self.status.update("Search regex: ")
            self.input = Input(placeholder="regex", id="search")
            await self.mount(self.input)
            await self.input.focus()
        elif key == ":":
            # Command prompt
            self.status.update("Command (:details | :summary | :toggle): ")
            self.input = Input(placeholder=":details", id="command")
            await self.mount(self.input)
            await self.input.focus()
        elif key == "p":
            # Simple palette as help overlay
            self.help.update("Palette: :details | :summary | :toggle | / search | w wrap | r reload | q quit")

    async def on_input_submitted(self, message: Input.Submitted) -> None:
        if message.input.id == "search":
            self.search = message.value
            await message.input.remove()
            self.input = None
            self.status.update("")
            # Immediately re-highlight existing buffer
            try:
                self._refresh_stream_render()
            except Exception:
                pass
        elif message.input.id == "command":
            cmd = (message.value or "").strip().lstrip(":").lower()
            await message.input.remove()
            self.input = None
            if cmd in ("details", "detail", "on", "true"):
                self.details = True
            elif cmd in ("summary", "off", "false"):
                self.details = False
            elif cmd in ("toggle",):
                self.details = not self.details
            mode = "ON" if self.details else "OFF"
            self._update_status()
            # brief ack in help bar
            try:
                self.help.update(f"Keys: q d / w r : | Details {mode}")
            except Exception:
                pass


class OrchestratorCommands(Provider):
    async def search(self, query: str) -> Hits:
        matcher = self.matcher(query)
        app = self.app
        assert isinstance(app, OrchestratorApp)
        commands = [
            ("Details: Toggle", app.toggle_details, "Toggle detailed JSON output"),
            ("Details: On", app.details_on, "Show full JSON for new messages"),
            ("Details: Off", app.details_off, "Show summarized messages"),
            ("Search: Prompt", app.prompt_search, "Enter a regex to highlight matches"),
            ("Wrap: Toggle", app.toggle_wrap, "Toggle line wrapping in stream"),
            ("Reload: Graph/State", app.reload_graph_command, "Reload state-machine and current state"),
            ("Quit", app.action_quit, "Exit the application"),
        ]
        for label, callback, help_text in commands:
            score = matcher.match(label)
            if score > 0:
                yield Hit(score, matcher.highlight(label), callback, help=help_text)

    async def discover(self) -> list[DiscoveryHit]:
        app = self.app
        assert isinstance(app, OrchestratorApp)
        return [
            DiscoveryHit("Details: Toggle", app.toggle_details, help="Toggle detailed JSON output"),
            DiscoveryHit("Search: Prompt", app.prompt_search, help="Enter a regex to highlight matches"),
            DiscoveryHit("Wrap: Toggle", app.toggle_wrap, help="Toggle line wrapping"),
            DiscoveryHit("Reload: Graph/State", app.reload_graph_command, help="Reload graph"),
        ]


# Register the command provider
OrchestratorApp.COMMANDS = OrchestratorApp.COMMANDS | {OrchestratorCommands}


def build_claude_command(claude_bin: str, prompt: str, skip_permissions: bool, verbose: bool, extra: Optional[List[str]]) -> List[str]:
    cmd = [claude_bin]
    if skip_permissions:
        cmd.append("--dangerously-skip-permissions")
    cmd.extend(["-p", prompt])
    if verbose:
        cmd.append("--verbose")
    cmd.extend(["--output-format", "stream-json"])
    if extra:
        cmd.extend(extra)
    return cmd


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--claude-bin", default=os.environ.get("CLAUDE_BIN", "claude"))
    parser.add_argument("--prompt", "-p", default=os.environ.get("CLAUDE_PROMPT", "/continue-orchestrating"))
    parser.add_argument("--no-skip-permissions", action="store_true")
    parser.add_argument("--no-verbose", action="store_true")
    parser.add_argument("--log-file")
    parser.add_argument("--state-machine-path", default="./state-machines/software-factory-3.0-state-machine.json")
    parser.add_argument("--state-path", default="./orchestrator-state-v3.json")
    parser.add_argument("--details", action="store_true", help="Start with full JSON details on")
    parser.add_argument("--", dest="extra", nargs=argparse.REMAINDER)
    args = parser.parse_args()

    extra = [] if args.extra is None else args.extra
    if extra and extra[0] == "--":
        extra = extra[1:]
    cmd = build_claude_command(
        claude_bin=args.claude_bin,
        prompt=args.prompt,
        skip_permissions=not args.no_skip_permissions,
        verbose=not args.no_verbose,
        extra=extra,
    )

    app = OrchestratorApp(cmd, args.state_machine_path, args.state_path, args.log_file, details_default=args.details)
    app.run()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())


