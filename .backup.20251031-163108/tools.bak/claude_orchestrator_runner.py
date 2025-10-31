#!/usr/bin/env python3

import argparse
import json
import os
import re
import shlex
import signal
import subprocess
import sys
import curses
from typing import List, Optional, TextIO, Tuple
import threading
try:
    import pytermgui as ptg  # type: ignore
except Exception:
    ptg = None
try:
    from prompt_toolkit.application import Application  # type: ignore
    from prompt_toolkit.layout import Layout, HSplit, Window as PTKWindow  # type: ignore
    from prompt_toolkit.layout.controls import BufferControl  # type: ignore
    from prompt_toolkit.buffer import Buffer  # type: ignore
    from prompt_toolkit.key_binding import KeyBindings  # type: ignore
    from prompt_toolkit.layout.dimension import Dimension  # type: ignore
    ptk_available = True
except Exception:
    ptk_available = False


CONTINUE_TRUE_PATTERN = re.compile(r"CONTINUE-SOFTWARE-FACTORY\s*=\s*TRUE", re.IGNORECASE)
CONTINUE_FALSE_PATTERN = re.compile(r"CONTINUE-SOFTWARE-FACTORY\s*=\s*FALSE", re.IGNORECASE)


def build_claude_command(
    claude_bin: str,
    prompt: str,
    skip_permissions: bool,
    verbose: bool,
    extra_args: Optional[List[str]] = None,
) -> List[str]:
    cmd: List[str] = [claude_bin]
    if skip_permissions:
        cmd.append("--dangerously-skip-permissions")
    cmd.extend(["-p", prompt])
    if verbose:
        cmd.append("--verbose")
    cmd.extend(["--output-format", "stream-json"])
    if extra_args:
        cmd.extend(extra_args)
    return cmd


def _format_json_for_output(raw_line: str, pretty: bool) -> str:
    try:
        obj = json.loads(raw_line)
        if pretty:
            return json.dumps(obj, indent=2, ensure_ascii=False) + "\n"
        else:
            return raw_line
    except Exception:
        return raw_line


def _clip(s: str, max_len: int = 240) -> str:
    if len(s) <= max_len:
        return s
    return s[: max_len - 1] + "…"


def _first_non_empty(*values: Optional[str]) -> Optional[str]:
    for v in values:
        if isinstance(v, str) and v.strip():
            return v
    return None


def _extract_content(obj: dict) -> Optional[str]:
    # Try common fields for message content
    message = obj.get("message") if isinstance(obj, dict) else None
    if isinstance(message, dict):
        content = message.get("content")
        if isinstance(content, str):
            return content
        if isinstance(content, list):
            # join text parts if array
            parts = []
            for it in content:
                if isinstance(it, str):
                    parts.append(it)
                elif isinstance(it, dict):
                    val = it.get("text") or it.get("content")
                    if isinstance(val, str):
                        parts.append(val)
            if parts:
                return " ".join(parts)
    # Direct fields
    for key in ("content", "text", "result"):
        val = obj.get(key)
        if isinstance(val, str):
            return val
    # Tool output
    tool_result = obj.get("tool_result") or obj.get("toolOutput")
    if isinstance(tool_result, dict):
        return _first_non_empty(
            tool_result.get("stdout"),
            tool_result.get("output"),
            tool_result.get("result"),
            tool_result.get("content"),
        )
    return None


def _summarize_json_line(obj: dict) -> str:
    t = obj.get("type")
    subtype = obj.get("subtype") or obj.get("role")
    msg_id = obj.get("id") or obj.get("uuid")
    is_error = obj.get("is_error")
    duration = obj.get("duration_ms") or obj.get("latency_ms")

    base = []
    if msg_id:
        base.append(f"id={msg_id}")
    if t:
        if subtype:
            base.append(f"{t}/{subtype}")
        else:
            base.append(f"{t}")
    if is_error is True:
        base.append("ERROR")
    if duration is not None:
        base.append(f"{duration}ms")

    # Additional fields heuristics per type
    extra = []
    if t == "result":
        # Show cost if available and decision flag
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
        # tool_start/tool_result
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
        summary = summary + " | " + " | ".join(extra)
    return summary or json.dumps(obj)


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
            # Try common shapes
            for key in ("current_state", "state", "current", "CurrentState"):
                if isinstance(data, dict) and key in data and isinstance(data[key], str):
                    return data[key]
            # Nested shapes like {"machine": {"state": "..."}}
            machine = data.get("machine") if isinstance(data, dict) else None
            if isinstance(machine, dict):
                for key in ("current_state", "state", "current"):
                    if key in machine and isinstance(machine[key], str):
                        return machine[key]
    except Exception:
        return None


def _render_state_graph(sm: Optional[dict], current_state: Optional[str]) -> List[Tuple[str, int]]:
    lines: List[Tuple[str, int]] = []
    if not isinstance(sm, dict):
        return [("State machine not found or invalid", curses.A_BOLD)]

    # Extract states in a robust way
    raw_states = (
        sm.get("states")
        or sm.get("States")
        or sm.get("nodes")
        or sm.get("Nodes")
        or sm.get("StatesList")
    )
    state_names: List[str] = []
    if isinstance(raw_states, dict):
        state_names = list(raw_states.keys())
    elif isinstance(raw_states, list):
        for it in raw_states:
            if isinstance(it, str):
                state_names.append(it)
            elif isinstance(it, dict):
                name = it.get("name") or it.get("id") or it.get("state")
                if isinstance(name, str):
                    state_names.append(name)

    # Extract transitions
    raw_transitions = (
        sm.get("transitions")
        or sm.get("Transitions")
        or sm.get("edges")
        or sm.get("Edges")
    )
    transitions: List[Tuple[str, Optional[str], str]] = []  # (from, symbol, to)
    if isinstance(raw_transitions, dict):
        for frm, outs in raw_transitions.items():
            if isinstance(outs, dict):
                for symbol, to_state in outs.items():
                    if isinstance(to_state, str):
                        transitions.append((frm, str(symbol), to_state))
            elif isinstance(outs, list):
                for to_state in outs:
                    if isinstance(to_state, str):
                        transitions.append((frm, None, to_state))
    elif isinstance(raw_transitions, list):
        for t in raw_transitions:
            if isinstance(t, dict):
                frm = t.get("from") or t.get("source")
                to = t.get("to") or t.get("target")
                symbol = t.get("on") or t.get("event")
                if isinstance(frm, str) and isinstance(to, str):
                    transitions.append((frm, symbol if isinstance(symbol, str) else None, to))

    # If transitions not provided, infer from per-state 'on' maps like XState
    if not transitions and isinstance(raw_states, dict):
        for frm, node in raw_states.items():
            if isinstance(node, dict):
                on_map = node.get("on") or node.get("On")
                if isinstance(on_map, dict):
                    for symbol, to in on_map.items():
                        if isinstance(to, str):
                            transitions.append((frm, str(symbol), to))

    # Build lines
    lines.append(("State Machine", curses.A_BOLD))
    unique_states = list(sorted(set(state_names)))
    if unique_states:
        # Render a single-row map of states
        cell_w = max(8, min(16, max((len(n) for n in unique_states), default=8) + 2))
        row_parts: List[str] = []
        current_idx = None
        for idx, name in enumerate(unique_states):
            if current_state and name == current_state:
                current_idx = idx
            padded = name[: cell_w - 2].ljust(cell_w - 2)
            row_parts.append(f"[{padded}]")
        map_row = " ".join(row_parts)
        lines.append((map_row, curses.A_NORMAL))
        # Pointer row to "highlight" current state position
        if current_idx is not None:
            # Compute start column of current block in the composed string
            block_len = cell_w
            sep_len = 1  # space between blocks
            start_col = current_idx * (block_len + sep_len)
            pointer = " " * max(0, start_col + (block_len // 2) - 1) + "^ YOU ARE HERE"
            attr = (curses.color_pair(3) | curses.A_BOLD) if curses.has_colors() else curses.A_BOLD
            lines.append((pointer, attr))
    else:
        lines.append(("(no states)", curses.A_DIM))

    lines.append(("", curses.A_NORMAL))
    lines.append(("Transitions", curses.A_BOLD))
    if transitions:
        # Show current state's outgoing transitions first for quick context
        if current_state:
            for frm, symbol, to in transitions:
                if frm == current_state:
                    label = f"{frm} -[{symbol}]-> {to}" if symbol else f"{frm} --> {to}"
                    attr = (curses.color_pair(1) | curses.A_BOLD) if curses.has_colors() else curses.A_BOLD
                    lines.append((label, attr))
        # Then list the rest dimmed
        for frm, symbol, to in transitions:
            if current_state and frm == current_state:
                continue
            label = f"{frm} -[{symbol}]-> {to}" if symbol else f"{frm} --> {to}"
            lines.append((label, curses.A_DIM))
    else:
        lines.append(("(no transitions)", curses.A_DIM))
    return lines


class _CursesUI:
    def __init__(self, title: str):
        self.title = title
        # Store (text, attr) pairs so we can highlight certain lines
        self.lines: List[tuple[str, int]] = []
        self.stdscr = None
        # Graph window state
        self.graph_visible: bool = True
        self.graph_lines: List[Tuple[str, int]] = []
        # Help overlay
        self.help_visible: bool = False
        # Wrapping
        self.wrap_enabled: bool = False
        # Search
        self.search_query: Optional[str] = None
        self.search_case_sensitive: bool = False
        # Detail toggle (show full JSON)
        self.details_enabled: bool = False
        # Reload flag (for graph/state)
        self._request_reload: bool = False

    def __enter__(self):
        self.stdscr = curses.initscr()
        curses.noecho()
        curses.cbreak()
        self.stdscr.keypad(True)
        self.stdscr.nodelay(True)
        try:
            curses.start_color()
            curses.use_default_colors()
            # Color pairs: 1=green, 2=red, 3=yellow
            curses.init_pair(1, curses.COLOR_GREEN, -1)
            curses.init_pair(2, curses.COLOR_RED, -1)
            curses.init_pair(3, curses.COLOR_YELLOW, -1)
        except Exception:
            pass
        return self

    def __exit__(self, exc_type, exc, tb):
        try:
            if self.stdscr is not None:
                self.stdscr.keypad(False)
            curses.nocbreak()
            curses.echo()
            curses.endwin()
        except Exception:
            # Ensure terminal is restored even on errors
            pass

    def append_text(self, text: str):
        # Split multiline text to discrete display lines
        for part in text.splitlines():
            self.lines.append((part, curses.A_NORMAL))
        self._render()

    def append_error_text(self, text: str):
        attr = curses.color_pair(2) | curses.A_BOLD if curses.has_colors() else curses.A_BOLD
        for part in text.splitlines():
            self.lines.append((f"[stderr] {part}", attr))
        self._render()

    def append_decision(self, continue_count: int, is_continue: bool):
        if is_continue:
            msg = f"Decision: CONTINUE  (continues so far = {continue_count})"
            attr = curses.color_pair(1) | curses.A_BOLD if curses.has_colors() else curses.A_BOLD
        else:
            msg = "Decision: STOP"
            attr = curses.color_pair(2) | curses.A_BOLD if curses.has_colors() else curses.A_BOLD
        self.lines.append((msg, attr))
        self._render()

    def set_graph_lines(self, graph_lines: List[Tuple[str, int]]):
        self.graph_lines = graph_lines
        self._render()

    def _render(self):
        if self.stdscr is None:
            return
        try:
            height, width = self.stdscr.getmaxyx()
            # Header and footer consume 2 lines
            content_height = max(0, height - 2)

            self.stdscr.erase()
            header = f"{self.title}  |  q: quit  g: toggle graph  w: wrap  d: details  /: search  r: reload  ?: help"
            self.stdscr.addnstr(0, 0, header, width, curses.A_BOLD)

            # Help overlay
            if self.help_visible:
                help_lines = [
                    ("Keyboard Shortcuts", curses.A_BOLD),
                    ("q: quit", curses.A_NORMAL),
                    ("g: toggle graph window", curses.A_NORMAL),
                    ("w: toggle wrap streamed lines", curses.A_NORMAL),
                    ("r: reload state machine + current state", curses.A_NORMAL),
                    ("?: toggle this help", curses.A_NORMAL),
                ]
                for idx, (line, attr) in enumerate(help_lines[:content_height]):
                    self.stdscr.addnstr(1 + idx, 0, line, width, attr)
                footer = f"Lines: {len(self.lines)}"
                self.stdscr.addnstr(height - 1, 0, footer, width)
                self.stdscr.refresh()
                return

            # Graph window at top
            top_rows_used = 0
            if self.graph_visible and self.graph_lines:
                graph_rows = min(max(3, content_height // 3), min(len(self.graph_lines), content_height))
                for idx in range(graph_rows):
                    line, attr = self.graph_lines[idx]
                    self.stdscr.addnstr(1 + idx, 0, line, width, attr)
                top_rows_used = graph_rows

            # Draw streaming content below graph
            remaining_height = max(0, content_height - top_rows_used)
            start = max(0, len(self.lines) - remaining_height)
            visible = self.lines[start:]

            # Wrap lines if enabled
            row_idx = 0
            for (line, attr) in visible:
                if row_idx >= remaining_height:
                    break
                if self.wrap_enabled and width > 1:
                    text = line
                    while text and row_idx < remaining_height:
                        segment = text[: width - 1]
                        self._addnstr_highlighted(1 + top_rows_used + row_idx, 0, segment, width, attr)
                        text = text[width - 1 :]
                        row_idx += 1
                else:
                    self._addnstr_highlighted(1 + top_rows_used + row_idx, 0, line, width, attr)
                    row_idx += 1

            footer = f"Lines: {len(self.lines)}"
            self.stdscr.addnstr(height - 1, 0, footer, width)
            self.stdscr.refresh()
        except Exception:
            # Ignore render errors (e.g., during resize race)
            pass

    def should_quit(self) -> bool:
        if self.stdscr is None:
            return False
        try:
            ch = self.stdscr.getch()
            if ch in (ord('q'), ord('Q')):
                return True
            if ch in (ord('g'), ord('G')):
                self.graph_visible = not self.graph_visible
                self._render()
            if ch in (ord('w'), ord('W')):
                self.wrap_enabled = not self.wrap_enabled
                self._render()
            if ch in (ord('?'),):
                self.help_visible = not self.help_visible
                self._render()
            if ch in (ord('r'), ord('R')):
                self._request_reload = True
            if ch == ord('/'):
                q = self._prompt_line("Search regex: ")
                if q is not None:
                    self.search_query = q.strip() or None
                    self._render()
            if ch in (ord('d'), ord('D')):
                self.details_enabled = not self.details_enabled
                self._render()
        except Exception:
            return False
        return False
class _PTGUI:
    def __init__(self, title: str, search_query: Optional[str] = None, search_case_sensitive: bool = False):
        if ptg is None:
            raise RuntimeError("pytermgui is not installed. Install with: pip install pytermgui")
        self.title = title
        self.lines: List[tuple[str, int]] = []
        self.graph_lines: List[Tuple[str, int]] = []
        self.wrap_enabled: bool = False
        self.search_query = search_query
        self.search_case_sensitive = search_case_sensitive
        self.manager = ptg.WindowManager()
        # Build UI widgets
        self.header = ptg.Window("[bold]" + self.title + "[/]\n(q to quit; g graph; w wrap; r reload; ? help; / search)" )
        self.graph_label = ptg.Label("")
        self.graph_win = ptg.Window(self.graph_label)
        self.text_label = ptg.Label("")
        self.text_win = ptg.Window(self.text_label)
        # Layout
        self.manager.layout.add_slot("Header")
        self.manager.layout.add_slot("Graph")
        self.manager.layout.add_slot("Body")
        self.manager.add(self.header, assign="Header")
        self.manager.add(self.graph_win, assign="Graph")
        self.manager.add(self.text_win, assign="Body")

    def __enter__(self):
        self.manager.__enter__()
        return self

    def __exit__(self, exc_type, exc, tb):
        try:
            self.manager.__exit__(exc_type, exc, tb)
        except Exception:
            pass

    def _render_lines_to_text(self) -> str:
        def highlight(s: str) -> str:
            if not self.search_query:
                return s
            try:
                import re
                flags = 0 if self.search_case_sensitive else re.IGNORECASE
                pattern = re.compile(self.search_query, flags)
                def repl(m: re.Match) -> str:
                    return "[210 bold]" + m.group(0) + "[/]"
                return pattern.sub(repl, s)
            except Exception:
                return s
        return "\n".join(highlight(line) for (line, _attr) in self.lines)

    def _render_graph_to_text(self) -> str:
        return "\n".join(line for (line, _attr) in self.graph_lines)

    def append_text(self, text: str):
        for part in text.splitlines():
            self.lines.append((part, curses.A_NORMAL))
        self.text_label.value = self._render_lines_to_text()
        self.manager.redraw()

    def append_error_text(self, text: str):
        for part in text.splitlines():
            self.lines.append((f"[stderr] {part}", curses.A_BOLD))
        self.text_label.value = self._render_lines_to_text()
        self.manager.redraw()

    def append_decision(self, continue_count: int, is_continue: bool):
        if is_continue:
            msg = f"Decision: CONTINUE  (continues so far = {continue_count})"
        else:
            msg = "Decision: STOP"
        self.lines.append((msg, curses.A_BOLD))
        self.text_label.value = self._render_lines_to_text()
        self.manager.redraw()

    def set_graph_lines(self, graph_lines: List[Tuple[str, int]]):
        # Drop attributes; PTG uses markup, keep plain text
        self.graph_lines = graph_lines
        self.graph_label.value = self._render_graph_to_text()
        self.manager.redraw()

    def should_quit(self) -> bool:
        return False

    def consume_reload_requested(self) -> bool:
        return False


class _PTKUI:
    def __init__(self, title: str, search_query: Optional[str] = None, search_case_sensitive: bool = False):
        if not ptk_available:
            raise RuntimeError("prompt_toolkit is not installed. Install with: pip install prompt_toolkit")
        self.title = title
        self.search_query = search_query
        self.search_case_sensitive = search_case_sensitive
        self.lines: List[str] = []
        self.graph_lines: List[str] = []
        self.wrap_enabled: bool = False
        self._request_reload: bool = False
        self._should_exit: bool = False

        self.header_buf = Buffer(read_only=True)
        self.graph_buf = Buffer(read_only=True)
        self.text_buf = Buffer(read_only=True)

        self.header_buf.text = f"{self.title}\n(q to quit; g graph; w wrap; r reload; / search)\n"

        root = HSplit([
            PTKWindow(BufferControl(buffer=self.header_buf), height=Dimension(min=2, max=2)),
            PTKWindow(BufferControl(buffer=self.graph_buf), height=Dimension(preferred=10)),
            PTKWindow(BufferControl(buffer=self.text_buf))
        ])
        kb = KeyBindings()

        @kb.add('q')
        def _(event):
            self._should_exit = True
            event.app.exit()

        @kb.add('r')
        def _(event):
            self._request_reload = True

        @kb.add('w')
        def _(event):
            self.wrap_enabled = not self.wrap_enabled
            self._render_text()

        @kb.add('/')
        def _(event):
            # Enter a simple search prompt in the header
            from prompt_toolkit.shortcuts import prompt
            try:
                q = prompt('Search regex: ')
                if q:
                    self.search_query = q
                    self._render_text()
            except Exception:
                pass

        self.app = Application(layout=Layout(root), key_bindings=kb, full_screen=True)
        self._thread: Optional[threading.Thread] = None

    def __enter__(self):
        def run():
            try:
                self.app.run()
            except Exception:
                pass
        self._thread = threading.Thread(target=run, daemon=True)
        self._thread.start()
        return self

    def __exit__(self, exc_type, exc, tb):
        try:
            self.app.exit()
        except Exception:
            pass

    def _highlight(self, s: str) -> str:
        if not self.search_query:
            return s
        try:
            import re
            flags = 0 if self.search_case_sensitive else re.IGNORECASE
            pattern = re.compile(self.search_query, flags)
            return pattern.sub(lambda m: f"{m.group(0)}", s)
        except Exception:
            return s

    def _wrap_lines(self, text: str) -> str:
        if not self.wrap_enabled:
            return text
        try:
            import shutil
            width = shutil.get_terminal_size((100, 20)).columns
            out_lines: List[str] = []
            for line in text.splitlines():
                while len(line) > width:
                    out_lines.append(line[:width])
                    line = line[width:]
                out_lines.append(line)
            return "\n".join(out_lines)
        except Exception:
            return text

    def _render_text(self):
        text = "\n".join(self.lines)
        text = self._highlight(text)
        text = self._wrap_lines(text)
        self.text_buf.text = text

    def append_text(self, text: str):
        for part in text.splitlines():
            self.lines.append(part)
        self._render_text()

    def append_error_text(self, text: str):
        for part in text.splitlines():
            self.lines.append(f"[stderr] {part}")
        self._render_text()

    def append_decision(self, continue_count: int, is_continue: bool):
        msg = f"Decision: CONTINUE (continues so far = {continue_count})" if is_continue else "Decision: STOP"
        self.lines.append(msg)
        self._render_text()

    def set_graph_lines(self, graph_lines: List[Tuple[str, int]]):
        # Convert to plain text
        self.graph_lines = [line for (line, _attr) in graph_lines]
        self.graph_buf.text = "\n".join(self.graph_lines)

    def should_quit(self) -> bool:
        return self._should_exit

    def consume_reload_requested(self) -> bool:
        if self._request_reload:
            self._request_reload = False
            return True
        return False

    def _prompt_line(self, prompt: str) -> Optional[str]:
        if self.stdscr is None:
            return None
        try:
            curses.echo()
            self.stdscr.nodelay(False)
            height, width = self.stdscr.getmaxyx()
            self.stdscr.addnstr(height - 1, 0, (prompt + " " * width)[:width], width)
            self.stdscr.move(height - 1, min(len(prompt), width - 1))
            s = self.stdscr.getstr(height - 1, len(prompt), max(1, width - len(prompt) - 1))
            return s.decode('utf-8', errors='ignore') if isinstance(s, bytes) else None
        except Exception:
            return None
        finally:
            try:
                curses.noecho()
                if self.stdscr is not None:
                    self.stdscr.nodelay(True)
            except Exception:
                pass

    def _addnstr_highlighted(self, y: int, x: int, text: str, width: int, base_attr: int):
        if self.stdscr is None:
            return
        if not self.search_query:
            self.stdscr.addnstr(y, x, text, width, base_attr)
            return
        try:
            import re
            flags = 0 if self.search_case_sensitive else re.IGNORECASE
            pattern = re.compile(self.search_query, flags)
            pos = 0
            for m in pattern.finditer(text):
                pre = text[pos:m.start()]
                self.stdscr.addnstr(y, x, pre, width - (x), base_attr)
                x += len(pre)
                hl = m.group(0)
                attr = curses.A_REVERSE | curses.A_BOLD
                self.stdscr.addnstr(y, x, hl, width - (x), attr)
                x += len(hl)
                pos = m.end()
                if x >= width:
                    return
            rest = text[pos:]
            if rest and x < width:
                self.stdscr.addnstr(y, x, rest, width - (x), base_attr)
        except Exception:
            self.stdscr.addnstr(y, x, text, width, base_attr)
    def consume_reload_requested(self) -> bool:
        if self._request_reload:
            self._request_reload = False
            return True
        return False


def stream_process_and_detect_continue(
    proc: subprocess.Popen,
    *,
    ui: Optional[_CursesUI],
    pretty: bool,
    decision_continue_count: Optional[int] = None,
    log_file: Optional[TextIO] = None,
    on_reload: Optional[callable] = None,
) -> bool:
    """
    Stream stdout to our stdout and mirror stderr to stderr. Detect continuation
    ONLY from the final JSON object whose "type" is "result". We continue only
    if that final result contains CONTINUE-SOFTWARE-FACTORY=TRUE. If it contains
    an explicit FALSE, or lacks the marker, we stop.
    Returns True if we should continue, False otherwise.
    """
    last_result_raw: Optional[str] = None

    assert proc.stdout is not None
    try:
        for line in proc.stdout:
            # Tee raw stdout to log file if requested
            if log_file is not None:
                try:
                    log_file.write(line)
                except Exception:
                    pass
            # Output handling
            if ui is not None:
                # Summarized or detailed line based on toggle
                if getattr(ui, 'details_enabled', False):
                    ui.append_text(_format_json_for_output(line, pretty))
                else:
                    try:
                        obj = json.loads(line)
                        if isinstance(obj, dict):
                            ui.append_text(_summarize_json_line(obj) + "\n")
                        else:
                            ui.append_text(_format_json_for_output(line, pretty))
                    except Exception:
                        ui.append_text(_format_json_for_output(line, pretty))
                if ui.should_quit():
                    raise KeyboardInterrupt
                if ui.consume_reload_requested() and on_reload is not None:
                    try:
                        on_reload()
                    except Exception:
                        pass
            else:
                # stdout mode
                formatted = _format_json_for_output(line, pretty)
                sys.stdout.write(formatted)
                # Ensure newline if pretty-printed JSON lacked trailing newline
                if not formatted.endswith("\n"):
                    sys.stdout.write("\n")
                sys.stdout.flush()

            # Track only the last object where type == "result"
            try:
                obj = json.loads(line)
                if isinstance(obj, dict) and obj.get("type") == "result":
                    last_result_raw = line
            except Exception:
                pass
    finally:
        # Drain and forward any remaining stderr content after stdout completes
        if proc.stderr is not None:
            err = proc.stderr.read()
            if err:
                # Tee stderr to log as well
                if log_file is not None:
                    try:
                        log_file.write(err)
                    except Exception:
                        pass
                if ui is not None:
                    ui.append_error_text(err)
                else:
                    sys.stderr.write(err)
                    sys.stderr.flush()

    # If we never saw a final result object, do not continue
    if not last_result_raw:
        if ui is not None:
            ui.append_decision(continue_count=0, is_continue=False)
        return False

    # Evaluate only the last result payload
    if CONTINUE_FALSE_PATTERN.search(last_result_raw):
        if ui is not None:
            ui.append_decision(continue_count=0, is_continue=False)
        return False
    if CONTINUE_TRUE_PATTERN.search(last_result_raw):
        # Use provided decision count if available, otherwise 0
        if ui is not None:
            ui.append_decision(continue_count=(decision_continue_count or 0), is_continue=True)
        return True
    if ui is not None:
        ui.append_decision(continue_count=0, is_continue=False)
    return False


def run_once(cmd: List[str], *, ui: Optional[_CursesUI], pretty: bool, decision_continue_count: Optional[int], log_file: Optional[TextIO], on_reload: Optional[callable]) -> int:
    """Run a single claude invocation, stream output, and return exit code."""
    # Use line-buffered text mode for stdout; merge no streams to preserve JSON on stdout
    proc = subprocess.Popen(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        bufsize=1,
        universal_newlines=True,
    )

    # Ensure child terminates if we receive SIGINT/SIGTERM
    def _forward_signal(signum, frame):
        try:
            proc.send_signal(signum)
        except Exception:
            pass

    # Temporarily set signal handlers while this process runs
    old_int = signal.getsignal(signal.SIGINT)
    old_term = signal.getsignal(signal.SIGTERM)
    signal.signal(signal.SIGINT, _forward_signal)
    signal.signal(signal.SIGTERM, _forward_signal)

    try:
        should_continue = stream_process_and_detect_continue(
            proc,
            ui=ui,
            pretty=pretty,
            decision_continue_count=decision_continue_count,
            log_file=log_file,
            on_reload=on_reload,
        )
        returncode = proc.wait()
    finally:
        # Restore original handlers
        signal.signal(signal.SIGINT, old_int)
        signal.signal(signal.SIGTERM, old_term)

    # If the model ran out of tokens or errored, return non-zero
    if returncode != 0:
        return returncode if returncode is not None else 1

    # Encode continuation decision in exit code semantics for the outer loop
    return 0 if should_continue else 2


def main() -> int:
    parser = argparse.ArgumentParser(
        description=(
            "Continuously invoke 'claude' with a prompt (default: /continue-orchestrating), "
            "streaming JSON output to stdout and re-running while the orchestrator emits "
            "CONTINUE-SOFTWARE-FACTORY=TRUE."
        )
    )
    parser.add_argument(
        "--claude-bin",
        default=os.environ.get("CLAUDE_BIN", "claude"),
        help="Path to the claude binary (default: 'claude' or $CLAUDE_BIN)",
    )
    parser.add_argument(
        "--prompt",
        "-p",
        default=os.environ.get("CLAUDE_PROMPT", "/continue-orchestrating"),
        help=(
            "Prompt/command to pass to '-p'. Examples: '/continue-orchestrating' or "
            "'/continue-orchestration'."
        ),
    )
    parser.add_argument(
        "--no-skip-permissions",
        action="store_true",
        help="Do not pass --dangerously-skip-permissions to claude",
    )
    parser.add_argument(
        "--no-verbose",
        action="store_true",
        help="Do not pass --verbose to claude",
    )
    parser.add_argument(
        "--once",
        action="store_true",
        help="Run only a single iteration (useful for testing)",
    )
    parser.add_argument(
        "--",
        dest="extra",
        nargs=argparse.REMAINDER,
        help="Additional args to pass through to claude (place after --)",
    )
    parser.add_argument(
        "--pretty",
        action="store_true",
        help="Pretty-print JSON output as it streams",
    )
    parser.add_argument(
        "--tui",
        action="store_true",
        help="Render a simple curses TUI to view streaming messages (press q to quit)",
    )
    parser.add_argument(
        "--title",
        default="Claude Orchestrator Runner",
        help="Title text for the TUI header",
    )
    parser.add_argument(
        "--ui",
        choices=["auto", "curses", "ptg", "ptk"],
        default="auto",
        help="Choose UI: ptk (prompt_toolkit), ptg (pytermgui), curses, or auto (prefer ptk)",
    )
    parser.add_argument(
        "--search",
        help="Highlight messages matching this regex (pytermgui UI)",
    )
    parser.add_argument(
        "--search-case-sensitive",
        action="store_true",
        help="Make --search case sensitive",
    )
    parser.add_argument(
        "--state-machine-path",
        default="./state-machines/software-factory-3.0-state-machine.json",
        help="Path to the software-factory state-machine JSON (for TUI graph)",
    )
    parser.add_argument(
        "--state-path",
        default="./orchestrator-state-v3.json",
        help="Path to the current orchestrator state JSON (for TUI graph)",
    )
    parser.add_argument(
        "--log-file",
        help="Path to a file to tee raw JSON stdout/stderr into",
    )
    parser.add_argument(
        "--log-append",
        action="store_true",
        help="Append to the log file instead of overwriting",
    )

    args = parser.parse_args()

    extra_args = [] if args.extra is None else args.extra
    # If user provided a bare '--', argparse includes it; remove if present
    if extra_args and extra_args[0] == "--":
        extra_args = extra_args[1:]

    cmd = build_claude_command(
        claude_bin=args.claude_bin,
        prompt=args.prompt,
        skip_permissions=not args.no_skip_permissions,
        verbose=not args.no_verbose,
        extra_args=extra_args,
    )

    sys.stderr.write("Starting claude orchestrator runner. Press Ctrl+C to stop.\n")
    sys.stderr.flush()

    iteration = 0
    continue_count = 0
    # Setup TUI once for the entire session, if requested
    ui: Optional[object] = None
    if args.tui:
        try:
            if (args.ui == "ptk" or (args.ui == "auto" and ptk_available)):
                ui = _PTKUI(args.title, search_query=args.search, search_case_sensitive=args.search_case_sensitive)
            elif (args.ui == "ptg" or (args.ui == "auto" and ptg is not None)):
                ui = _PTGUI(args.title, search_query=args.search, search_case_sensitive=args.search_case_sensitive)
            else:
                ui = _CursesUI(args.title)
            ui.__enter__()
        except Exception as e:
            sys.stderr.write(f"Failed to initialize TUI ({e}); falling back to stdout mode.\n")
            sys.stderr.flush()
            ui = None
    log_fh: Optional[TextIO] = None
    if args.log_file:
        try:
            mode = "a" if args.log_append else "w"
            log_fh = open(args.log_file, mode, encoding="utf-8")
        except Exception as e:
            sys.stderr.write(f"Failed to open log file {args.log_file}: {e}\n")
            sys.stderr.flush()
            log_fh = None

    try:
        while True:
            iteration += 1
            sys.stderr.write(f"\n--- Iteration {iteration} ---\n")
            sys.stderr.write("Executing: " + " ".join(shlex.quote(p) for p in cmd) + "\n")
            sys.stderr.flush()

            # If TUI is active, pass the next continue count for decision highlighting
            next_decision_count = continue_count + 1
            # Define on_reload to refresh graph/state on demand
            def _on_reload():
                if ui is not None:
                    sm = _load_state_machine(args.state_machine_path)
                    cur = _load_current_state(args.state_path)
                    ui.set_graph_lines(_render_state_graph(sm, cur))

            # initial load before first run
            if iteration == 1 and ui is not None:
                _on_reload()

            rc = run_once(
                cmd,
                ui=ui,
                pretty=args.pretty,
                decision_continue_count=next_decision_count,
                log_file=log_fh,
                on_reload=_on_reload if ui is not None else None,
            )

            if rc == 0:
                # Continue
                continue_count += 1
                if args.once:
                    sys.stderr.write("--once set; stopping after one successful iteration.\n")
                    sys.stderr.flush()
                    return 0
                continue
            elif rc == 2:
                sys.stderr.write(
                    "Stop condition met: CONTINUE-SOFTWARE-FACTORY not TRUE (or explicitly FALSE).\n"
                )
                sys.stderr.flush()
                return 0
            else:
                sys.stderr.write(f"claude exited with non-zero status {rc}; stopping.\n")
                sys.stderr.flush()
                return rc
    except KeyboardInterrupt:
        sys.stderr.write("Interrupted by user (Ctrl+C).\n")
        sys.stderr.flush()
        return 130
    finally:
        if log_fh is not None:
            try:
                log_fh.flush()
                log_fh.close()
            except Exception:
                pass
        if ui is not None:
            try:
                ui.__exit__(None, None, None)
            except Exception:
                pass


if __name__ == "__main__":
    raise SystemExit(main())


