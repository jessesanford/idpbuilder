#!/usr/bin/env python3
import argparse
import asyncio
import json
import re
import os
import shlex
import signal
import sys
import shutil
import subprocess
import threading
import time
import atexit
from dataclasses import dataclass
from datetime import datetime, timezone, timedelta
from pathlib import Path
from typing import AsyncIterator, Dict, Optional, Any, Tuple

from textual import events
from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.command import DiscoveryHit, Hit, Provider
from textual.app import SystemCommand
from textual.containers import Horizontal, Vertical
from textual.message import Message
from textual.reactive import reactive
from textual.widget import Widget
from textual.widgets import Footer, Header, Input, Static, Log, ListView, ListItem, Tree
from textual.containers import Vertical as VGroup
from textual.screen import ModalScreen
from rich.text import Text
import importlib.util


class ParamsScreen(ModalScreen[None]):
    def __init__(self, app_ref: "ClaudeOrchestratorApp") -> None:
        super().__init__()
        self.app_ref = app_ref
        self.input_command = Input(value=app_ref.config.command)
        self.input_params = Input(value=app_ref.config.env_params or "")

    def compose(self) -> ComposeResult:
        yield Vertical(
            Static("Edit command and params (-p)", id="title"),
            Static("Command:"),
            self.input_command,
            Static("Env params (CLAUDE_PARAMS):"),
            self.input_params,
            Horizontal(
                Static("Press Enter to save, Esc to cancel")
            ),
            id="params-modal",
        )

    async def on_input_submitted(self, _: Input.Submitted) -> None:
        self.app_ref.config.command = self.input_command.value
        self.app_ref.config.env_params = self.input_params.value or None
        self.app_ref._update_status()
        self.dismiss(None)

    async def on_key(self, event: events.Key) -> None:
        if event.key == "escape":
            self.dismiss(None)


class ScheduleScreen(ModalScreen[None]):
    def __init__(self, app_ref: "ClaudeOrchestratorApp") -> None:
        super().__init__()
        self.app_ref = app_ref
        self.input_dt = Input(placeholder="e.g., 12:15am EST 09/25/2025; empty to cancel")

    def compose(self) -> ComposeResult:
        yield Vertical(
            Static("Schedule one-time Claude start (Enter to save, Esc to cancel)", id="title"),
            self.input_dt,
            id="schedule-modal",
        )

    async def on_input_submitted(self, _: Input.Submitted) -> None:
        value = (self.input_dt.value or "").strip()
        self.app_ref._set_schedule_from_string(value)
        self.app_ref._update_status()
        self.dismiss(None)

    async def on_key(self, event: events.Key) -> None:
        if event.key == "escape":
            self.dismiss(None)


CONTINUE_MARKER = "CONTINUE-SOFTWARE-FACTORY=TRUE"
CONTINUE_TRUE_PATTERN = re.compile(r"CONTINUE-SOFTWARE-FACTORY\s*=\s*TRUE", re.IGNORECASE)
CONTINUE_FALSE_PATTERN = re.compile(r"CONTINUE-SOFTWARE-FACTORY\s*=\s*FALSE", re.IGNORECASE)


@dataclass
class RunnerConfig:
    command: str
    env_params: Optional[str] = None
    working_dir: str = os.getcwd()
    log_path: Optional[str] = None
    simulate: bool = False
    debug_raw: bool = False
    debug_verbose: bool = False
    sf2_project_dir: str = os.getcwd()


class ClaudeRunner:
    """Runs the Claude command, streams output lines, and detects continuation markers."""

    def __init__(self, config: RunnerConfig) -> None:
        self.config = config
        self._process: Optional[asyncio.subprocess.Process] = None
        self._stopping = False
        self.last_returncode: Optional[int] = None
        # PTY mode bookkeeping
        self._pty_pid: Optional[int] = None
        self._pty_fd: Optional[int] = None
        # Popen-based process tracking
        self._proc: Optional[subprocess.Popen] = None
        self._proc_pid: Optional[int] = None
        self._pgid: Optional[int] = None

    async def run_once(self) -> AsyncIterator[str]:
        env = os.environ.copy()
        if self.config.env_params:
            env["CLAUDE_PARAMS"] = self.config.env_params

        if self.config.simulate:
            for i in range(5):
                await asyncio.sleep(0.2)
                yield json.dumps({"type": "info", "id": i, "message": f"Simulated line {i}"})
            await asyncio.sleep(0.2)
            yield CONTINUE_MARKER
            return

        # Build argv from command string (no shell) unless tee is requested
        argv = shlex.split(self.config.command)
        if self.config.debug_verbose or self.config.debug_raw:
            yield f"[debug] using wrapper: none"
            yield f"[debug] exec: {' '.join(argv)}"
            claude_bin = argv[0] if argv else "(empty)"
            yield f"[debug] cwd: {self.config.working_dir}"
            yield f"[debug] PATH: {env.get('PATH','')}"
            yield f"[debug] CLAUDE_PARAMS: {env.get('CLAUDE_PARAMS','(unset)')}"
            yield f"[debug] which argv[0] exists: {shutil.which(claude_bin) or '(not on PATH)'}"

        # Launch process. If tee is requested and log_path is set, use a shell pipeline so stdout is also written to disk.
        if self.config.log_path:
            # We must use a shell here to construct the pipeline safely; quote argv for safety
            quoted = " ".join(shlex.quote(tok) for tok in argv)
            tee_cmd = f"set -o pipefail; {quoted} | tee -a {shlex.quote(self.config.log_path)}"
            proc = subprocess.Popen(
                ["bash", "-lc", tee_cmd],
                cwd=self.config.working_dir,
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                bufsize=1,
                universal_newlines=True,
                preexec_fn=os.setsid,
            )
        else:
            proc = subprocess.Popen(
                argv,
                cwd=self.config.working_dir,
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                bufsize=1,
                universal_newlines=True,
                preexec_fn=os.setsid,
            )
        self._proc = proc
        self._proc_pid = proc.pid
        try:
            self._pgid = os.getpgid(proc.pid)
        except Exception:
            self._pgid = None

        if self.config.debug_verbose:
            yield f"[debug] pid: {proc.pid}"

        loop = asyncio.get_event_loop()
        queue: asyncio.Queue[str] = asyncio.Queue()
        lines_out = 0

        def _stdout_reader() -> None:
            try:
                assert proc.stdout is not None
                for line in proc.stdout:
                    if self.config.debug_raw:
                        loop.call_soon_threadsafe(queue.put_nowait, f"[debug:chunk:stdout:{len(line)}] {line!r}")
                    loop.call_soon_threadsafe(queue.put_nowait, line.rstrip("\n\r"))
            except Exception:
                pass
            finally:
                loop.call_soon_threadsafe(queue.put_nowait, "__STDOUT_DONE__")

        def _stderr_reader() -> None:
            try:
                assert proc.stderr is not None
                for line in proc.stderr:
                    if self.config.debug_raw:
                        loop.call_soon_threadsafe(queue.put_nowait, f"[debug:chunk:stderr:{len(line)}] {line!r}")
                    loop.call_soon_threadsafe(queue.put_nowait, f"[stderr] {line.rstrip('\n\r')}")
            except Exception:
                pass
            finally:
                loop.call_soon_threadsafe(queue.put_nowait, "__STDERR_DONE__")

        t = threading.Thread(target=_stdout_reader, daemon=True)
        t.start()
        ts = threading.Thread(target=_stderr_reader, daemon=True)
        ts.start()

        # Drain queue asynchronously
        stdout_done = False
        stderr_done = False
        last_any = loop.time()
        while True:
            try:
                item = await asyncio.wait_for(queue.get(), timeout=1.0)
                if item == "__STDOUT_DONE__":
                    stdout_done = True
                    if self.config.debug_verbose:
                        yield "[debug] stdout reader done"
                    if stderr_done:
                        break
                    continue
                if item == "__STDERR_DONE__":
                    stderr_done = True
                    if self.config.debug_verbose:
                        yield "[debug] stderr reader done"
                    if stdout_done:
                        break
                    continue
                lines_out += 1
                last_any = loop.time()
                yield item
            except asyncio.TimeoutError:
                # Heartbeat
                if self.config.debug_verbose:
                    since = loop.time() - last_any
                    yield f"[debug] heartbeat: no new lines for {since:.1f}s"
                # Check if process has exited
                if proc.poll() is not None and stdout_done and stderr_done:
                    break

        # Wait for process to exit and capture return code
        try:
            proc.wait()
            self.last_returncode = proc.returncode
        except Exception:
            self.last_returncode = None

    async def stop(self) -> None:
        self._stopping = True
        # Attempt graceful then forceful termination of the whole process group
        try:
            proc = getattr(self, "_proc", None)
            if isinstance(proc, subprocess.Popen) and proc.poll() is None:
                pgid = getattr(self, "_pgid", None)
                try:
                    if isinstance(pgid, int):
                        os.killpg(pgid, signal.SIGINT)
                    elif isinstance(self._proc_pid, int):
                        os.kill(self._proc_pid, signal.SIGINT)
                except Exception:
                    pass
                await asyncio.sleep(0.5)
                if proc.poll() is None:
                    try:
                        if isinstance(pgid, int):
                            os.killpg(pgid, signal.SIGTERM)
                        else:
                            proc.terminate()
                    except Exception:
                        pass
                    await asyncio.sleep(0.5)
                if proc.poll() is None:
                    try:
                        if isinstance(pgid, int):
                            os.killpg(pgid, signal.SIGKILL)
                        else:
                            proc.kill()
                    except Exception:
                        pass
                try:
                    if proc.stdout:
                        proc.stdout.close()
                except Exception:
                    pass
                try:
                    if proc.stderr:
                        proc.stderr.close()
                except Exception:
                    pass
                try:
                    proc.wait(timeout=2)
                except Exception:
                    pass
        except Exception:
            pass

    def kill_now(self) -> None:
        """Best-effort immediate cleanup for use in atexit/signal paths (synchronous)."""
        try:
            proc = getattr(self, "_proc", None)
            if isinstance(proc, subprocess.Popen) and proc.poll() is None:
                pgid = getattr(self, "_pgid", None)
                try:
                    if isinstance(pgid, int):
                        os.killpg(pgid, signal.SIGTERM)
                    elif isinstance(self._proc_pid, int):
                        os.kill(self._proc_pid, signal.SIGTERM)
                except Exception:
                    pass
                time.sleep(0.2)
                if proc.poll() is None:
                    try:
                        if isinstance(pgid, int):
                            os.killpg(pgid, signal.SIGKILL)
                        else:
                            proc.kill()
                    except Exception:
                        pass
                try:
                    if proc.stdout:
                        proc.stdout.close()
                except Exception:
                    pass
                try:
                    if proc.stderr:
                        proc.stderr.close()
                except Exception:
                    pass
                try:
                    proc.wait(timeout=1)
                except Exception:
                    pass
        except Exception:
            pass


WHITESPACE_RE = re.compile(r"\s+")


def _pretty_json(data: Any) -> str:
    try:
        return json.dumps(data, indent=2, ensure_ascii=False, sort_keys=True)
    except Exception:
        return str(data)


def _safe_json_load(line: str) -> Any:
    try:
        return json.loads(line)
    except Exception:
        return {"raw": line}


def _shorten_id(ev_id: Optional[str]) -> str:
    if not ev_id:
        return "—"
    try:
        s = str(ev_id)
    except Exception:
        return "—"
    return s[-8:]


def _collapse(s: str, limit: int = 180) -> str:
    if not s:
        return ""
    s = WHITESPACE_RE.sub(" ", s.replace("\n", " ").replace("\r", " ")).strip()
    if len(s) > limit:
        return s[:limit - 1] + "…"
    return s


def _extract_text_from_content(content: Any) -> str:
    if isinstance(content, str):
        return content
    if isinstance(content, list):
        parts = []
        for item in content:
            if isinstance(item, dict):
                text = item.get("text") or item.get("content") or item.get("delta") or item.get("value")
                if isinstance(text, str) and text:
                    parts.append(text)
        return " ".join(parts)
    if isinstance(content, dict):
        text = content.get("text") or content.get("content")
        if isinstance(text, str):
            return text
    return ""


def _count_images(content: Any) -> int:
    count = 0
    if isinstance(content, list):
        for item in content:
            if isinstance(item, dict):
                t = item.get("type")
                if t in {"image", "input_image"}:
                    count += 1
                src = item.get("source")
                if isinstance(src, dict) and src.get("type") in {"base64", "url"}:
                    count += 0
    elif isinstance(content, dict):
        if content.get("type") in {"image", "input_image"}:
            count += 1
    return count


def _detect_type_and_id(data: Dict[str, Any]) -> Tuple[str, Optional[str]]:
    ev_type = data.get("type") or data.get("event") or data.get("kind")
    message = data.get("message") or {}
    # Tolerate message being a bare string in simulated/headless data
    if isinstance(message, str):
        message = {"content": message}
    role = message.get("role") or data.get("role")
    ev_id = (
        data.get("id")
        or message.get("id")
        or data.get("message_id")
        or data.get("tool_use_id")
        or data.get("command_id")
    )
    if role in {"assistant", "user", "system"} and not ev_type:
        ev_type = role
    if ev_type in {"message_start", "message_delta", "message_stop"} and role in {"assistant", "user", "system"}:
        ev_type = role
    if not ev_type:
        if "tool" in data or data.get("type") == "tool_use":
            ev_type = "tool"
        elif "command" in data or data.get("command_name"):
            ev_type = "command"
        elif "error" in data or data.get("status") == "error":
            ev_type = "error"
        elif "usage" in data:
            ev_type = "usage"
        else:
            ev_type = "event"
    if ev_type == "tool_result":
        ev_type = "tool"
    return ev_type, ev_id


def _extract_content_and_meta(data: Dict[str, Any], ev_type: str) -> Tuple[str, Dict[str, Any]]:
    meta: Dict[str, Any] = {}
    usage = data.get("usage") or data.get("token_usage")
    if isinstance(usage, dict):
        meta["usage"] = {
            "in": usage.get("input_tokens") or usage.get("inputTokens") or usage.get("prompt_tokens"),
            "out": usage.get("output_tokens") or usage.get("outputTokens") or usage.get("completion_tokens"),
        }
    stop_reason = (
        data.get("stop_reason")
        or data.get("finish_reason")
        or data.get("stopReason")
    )
    if stop_reason:
        meta["stop_reason"] = stop_reason

    message = data.get("message") or {}
    if isinstance(message, str):
        message = {"content": message}
    candidates = [
        message.get("content"),
        data.get("content"),
        data.get("delta"),
        data.get("text"),
    ]
    content = ""
    image_hint = 0
    for c in candidates:
        text = _extract_text_from_content(c)
        if text:
            content = text
            break
        image_hint = image_hint or _count_images(c)
    if not content and image_hint:
        content = "[image]" if image_hint == 1 else f"[image x{image_hint}]"

    if ev_type in {"tool", "tool_use"}:
        tool_name = data.get("tool") or data.get("name") or message.get("name")
        meta["tool_name"] = tool_name
        args = data.get("input") or data.get("args") or data.get("parameters") or data.get("command_args")
        if args is None:
            result = data.get("result") or data.get("output")
            if not content:
                content = _extract_text_from_content(result)
            if not content and isinstance(result, (str, int, float)):
                content = str(result)
        else:
            try:
                content = content or json.dumps(args, ensure_ascii=False)
            except Exception:
                content = content or str(args)

    if ev_type.startswith("command") or ev_type == "command":
        cmd_name = data.get("command") or data.get("command_name") or data.get("name")
        meta["command_name"] = cmd_name
        cmd_args = data.get("args") or data.get("input") or data.get("commandline") or data.get("cmd")
        if isinstance(cmd_args, str):
            content = content or cmd_args
        elif cmd_args is not None:
            try:
                content = content or json.dumps(cmd_args, ensure_ascii=False)
            except Exception:
                content = content or str(cmd_args)

    if ev_type == "error":
        err = data.get("error") or data.get("message") or data.get("detail")
        status = data.get("status") or data.get("exit_code") or data.get("code")
        if isinstance(err, dict):
            err = err.get("message") or err.get("error") or str(err)
        meta["error_summary"] = _collapse(f"{err or 'error'}" + (f" (exit={status})" if status is not None else ""), 120)
        if not content:
            content = meta["error_summary"]

    return _collapse(content, 180), meta


def _format_usage(usage: Optional[Dict[str, Any]]) -> str:
    if not usage:
        return ""
    parts = []
    if usage.get("in") is not None:
        parts.append(f"in={usage['in']}")
    if usage.get("out") is not None:
        parts.append(f"out={usage['out']}")
    return f" ({' '.join(parts)})" if parts else ""


def _format_summary(ev_type: str, ev_id: Optional[str], content: str, meta: Dict[str, Any], verbose: bool) -> str:
    sid = _shorten_id(ev_id)
    usage = _format_usage(meta.get("usage")) if verbose else ""
    stop = f" stop={meta['stop_reason']}" if verbose and meta.get("stop_reason") else ""
    if ev_type in {"assistant", "user", "system"}:
        body = content or ""
        return f"{ev_type} id={sid} — {body}{usage}{stop}".strip()
    if ev_type in {"tool", "tool_use"}:
        name = meta.get("tool_name") or "tool"
        args = content or ""
        return f"tool id={sid} {name} — {args}".strip()
    if ev_type.startswith("command") or ev_type == "command":
        name = meta.get("command_name") or "command"
        cmd = content or ""
        return f"command id={sid} {name} — {cmd}".strip()
    if ev_type == "error":
        err = meta.get("error_summary") or content or "error"
        return f"error id={sid}: {err}".strip()
    if ev_type == "usage":
        return f"usage id={sid}{usage}".strip()
    return f"{ev_type} id={sid}{usage}{stop}".strip()


class MessageParser:
    """Parses streaming lines to extract summaries and details, Claude Code style."""

    def __init__(self, summary_verbose: bool = True, buffer_char_cap: int = 8192) -> None:
        self.summary_verbose = summary_verbose
        self._buffer_char_cap = max(512, int(buffer_char_cap))
        self._buffers: Dict[str, list[str]] = {}
        self._roles: Dict[str, str] = {}
        self._sidecar_usage: Dict[str, Dict[str, Any]] = {}

    def parse_line(self, line: str) -> Dict[str, Any]:
        if line.strip() == CONTINUE_MARKER:
            return {
                "type": "marker",
                "summary": CONTINUE_MARKER,
                "detail_raw": CONTINUE_MARKER,
                "detail_pretty": CONTINUE_MARKER,
            }
        data = _safe_json_load(line)
        raw: Dict[str, Any] = data if isinstance(data, dict) else {}
        raw_type = raw.get("type") or raw.get("event") or raw.get("kind")
        ev_type, ev_id = _detect_type_and_id(raw)
        content, meta = _extract_content_and_meta(raw, ev_type)

        if ev_id and not meta.get("usage") and ev_id in self._sidecar_usage:
            meta["usage"] = self._sidecar_usage.get(ev_id)

        _msg = raw.get("message") or {}
        if isinstance(_msg, str):
            _msg = {"content": _msg}
        role = _msg.get("role") or raw.get("role")
        if ev_id and raw_type == "message_start" and role:
            self._roles[ev_id] = role

        def _extract_raw_delta_text(raw_ev: Dict[str, Any]) -> str:
            d = raw_ev.get("delta") or raw_ev.get("text")
            if isinstance(d, str):
                return d
            if isinstance(d, dict):
                t = d.get("text") or d.get("delta") or d.get("content")
                return t if isinstance(t, str) else ""
            return ""

        is_delta = raw_type == "message_delta" or (raw.get("delta") is not None and role in {"assistant", "user", "system"})
        if ev_id and is_delta:
            piece = _extract_raw_delta_text(raw)
            buf = self._buffers.setdefault(ev_id, [])
            if piece:
                buf.append(piece)
            elif content:
                buf.append(content)
            total = sum(len(p) for p in buf)
            while total > self._buffer_char_cap and buf:
                dropped = buf.pop(0)
                total -= len(dropped)

        is_stop = raw_type in {"message_stop", "finish"}
        if ev_id and is_stop:
            buf = self._buffers.pop(ev_id, None)
            if buf:
                content = _collapse("".join(buf), 120)
            self._roles.pop(ev_id, None)
            self._sidecar_usage.pop(ev_id, None)

        if ev_type == "usage" and ev_id and meta.get("usage"):
            self._sidecar_usage[ev_id] = meta["usage"]

        summary = _format_summary(ev_type, ev_id, content, meta, self.summary_verbose)

        # Build a compact extras text for expansion (so users see helpful fields without full JSON)
        extras_lines: list[str] = []
        if meta.get("usage"):
            u = meta["usage"]
            parts = []
            if u.get("in") is not None:
                parts.append(f"in={u['in']}")
            if u.get("out") is not None:
                parts.append(f"out={u['out']}")
            extras_lines.append(f"usage: {' '.join(parts)}")
        if meta.get("stop_reason"):
            extras_lines.append(f"stop: {meta['stop_reason']}")
        # Tool / command details (truncate long JSON)
        try:
            if ev_type in {"tool", "tool_use"}:
                name = meta.get("tool_name") or "tool"
                extras_lines.append(f"tool: {name}")
                args = raw.get("input") or raw.get("args") or raw.get("parameters")
                if args is not None:
                    try:
                        preview = json.dumps(args, ensure_ascii=False)
                    except Exception:
                        preview = str(args)
                    if len(preview) > 240:
                        preview = preview[:239] + "…"
                    extras_lines.append(f"args: {preview}")
                result = raw.get("result") or raw.get("output")
                if result is not None and not content:
                    rprev = str(result)
                    if len(rprev) > 240:
                        rprev = rprev[:239] + "…"
                    extras_lines.append(f"result: {rprev}")
            if ev_type == "command" or (isinstance(raw.get("command"), str)):
                cname = meta.get("command_name") or raw.get("command")
                if cname:
                    extras_lines.append(f"command: {cname}")
                cargs = raw.get("args") or raw.get("cmd") or raw.get("commandline")
                if cargs is not None:
                    try:
                        cap = json.dumps(cargs, ensure_ascii=False) if not isinstance(cargs, str) else cargs
                    except Exception:
                        cap = str(cargs)
                    if len(cap) > 240:
                        cap = cap[:239] + "…"
                    extras_lines.append(f"args: {cap}")
            if ev_type == "error" and meta.get("error_summary"):
                extras_lines.append(f"error: {meta['error_summary']}")
        except Exception:
            pass

        extras_text = "\n".join(extras_lines) if extras_lines else ""
        return {
            "type": ev_type,
            "id": ev_id,
            "raw_event": raw_type,
            "summary": summary,
            "detail_pretty": _pretty_json(data),
            "detail_raw": line,
            "usage": meta.get("usage"),
            "stop_reason": meta.get("stop_reason"),
            "extras_text": extras_text,
        }

    def set_summary_verbose(self, verbose: bool) -> None:
        self.summary_verbose = verbose


class IterationBanner(Static):
    """Full-width colored banner to indicate iteration."""

    def show_banner(self, iteration: int) -> None:
        self.update(f"[bold white on green] CONTINUATION #{iteration} [/]")
        self.styles.background = "green"
        self.styles.color = "white"


class StatePanel(Static):
    """Right-side panel that delegates rendering to state-machine-viewer.py."""

    def __init__(self, project_dir: str) -> None:
        super().__init__()
        self.project_dir = project_dir or os.getcwd()
        self.visible = False
        self.render_mode = "full"  # "full" or "zoom" (mapped to viewer)
        self._scroll_offset = 0
        self._viewer = self._load_viewer()
        self._last_full_lines: list[str] = []
        # Allow focus so mouse wheel and future key events can target this pane
        try:
            self.can_focus = True
        except Exception:
            pass

    def set_project_dir(self, path: str) -> None:
        self.project_dir = path or os.getcwd()

    def toggle(self) -> None:
        self.visible = not self.visible
        self.display = "block" if self.visible else "none"

    def _load_viewer(self):
        try:
            viewer_path = Path(__file__).with_name("state-machine-viewer.py")
            spec = importlib.util.spec_from_file_location("state_machine_viewer_mod", str(viewer_path))
            if spec and spec.loader:
                mod = importlib.util.module_from_spec(spec)
                spec.loader.exec_module(mod)  # type: ignore
                ViewerClass = getattr(mod, "StateMachineViewer", None)
                return ViewerClass() if ViewerClass else None
        except Exception:
            return None
        return None

    def _load_state(self):
        # Kept for existing callers; returns orchestrator data and empty machine
        try:
            prev = os.getcwd()
            os.chdir(self.project_dir)
            import json as _json
            data = {}
            p = Path("orchestrator-state-v3.json")
            if p.exists():
                data = _json.loads(p.read_text(encoding="utf-8"))
            return data, {}
        except Exception:
            return {}, {}
        finally:
            try:
                os.chdir(prev)
            except Exception:
                pass

    def _render_via_viewer(self) -> Text:
        if self._viewer is None:
            return Text("state-machine-viewer.py missing; cannot render.")
        prev = os.getcwd()
        try:
            os.chdir(self.project_dir)
            ok = bool(self._viewer.load_data())
        finally:
            os.chdir(prev)
        if not ok:
            return Text("Unable to load state files.")
        # Map render modes to viewer views
        if self.render_mode == "zoom":
            try:
                self._viewer.set_view('focused')
            except Exception:
                pass
            # Focused view returns a single string
            focused = getattr(self._viewer, "render_focused_view", None)
            if callable(focused):
                return Text.from_ansi(focused())
            return Text("Focused view unavailable")
        # For non-zoom: prefer journey view; fall back to overview
        try:
            cur = getattr(self._viewer, 'current_view', 'journey')
            if cur not in ('journey', 'overview'):
                self._viewer.set_view('journey')
        except Exception:
            pass
        get_lines = getattr(self._viewer, "get_lines", None)
        if callable(get_lines):
            all_lines = get_lines()
        else:
            all_lines = []
            journey = getattr(self._viewer, "render_journey_view", None)
            overview = getattr(self._viewer, "render_overview_view", None)
            if callable(journey):
                all_lines = journey()
            elif callable(overview):
                all_lines = overview()
            else:
                return Text("No suitable viewer render method found.")
        try:
            viewport = max(5, int(self.size.height))
        except Exception:
            viewport = 60
        self._scroll_offset = max(0, min(self._scroll_offset, max(0, len(all_lines) - viewport)))
        window = all_lines[self._scroll_offset:self._scroll_offset + viewport]
        # Track lines for scroll calculations in actions
        self._last_full_lines = all_lines
        return Text.from_ansi("\n".join(window))

    def render_ascii(self) -> str:
        # Maintain signature for any external use; delegate to viewer
        return self._render_via_viewer().plain

    def refresh_panel(self) -> None:
        try:
            self.update(self._render_via_viewer())
        except Exception as exc:
            self.update(f"[error rendering state panel] {exc}")

    # --- Native mouse wheel scrolling for convenience ---
    def on_mouse_scroll_up(self, event: events.MouseScrollUp) -> None:  # type: ignore[override]
        # Scroll a few lines per wheel notch
        self.scroll_up(3)

    def on_mouse_scroll_down(self, event: events.MouseScrollDown) -> None:  # type: ignore[override]
        self.scroll_down(3)

    # ---- Graph helpers ----
    def _extract_states_transitions(self, machine: dict):
        states = machine.get("states") or machine.get("States") or {}
        transitions = machine.get("transitions") or machine.get("Transitions") or {}
        node_names: list[str] = []
        edges: list[tuple[str, Optional[str], str]] = []  # (from, symbol, to)

        # Helper to record a node name
        def add_node(name: Optional[str]):
            if isinstance(name, str) and name not in node_names:
                node_names.append(name)

        # Helper to add edge(s)
        def add_edge(frm: Optional[str], symbol: Optional[str], to: Optional[str]):
            if isinstance(frm, str) and isinstance(to, str):
                edges.append((frm, symbol if isinstance(symbol, str) else None, to))
                add_node(frm)
                add_node(to)

        # 1) Top-level transitions map or list
        if isinstance(transitions, dict):
            for frm, outs in transitions.items():
                if isinstance(outs, dict):
                    for symbol, to in outs.items():
                        if isinstance(to, str):
                            add_edge(frm, str(symbol), to)
                        elif isinstance(to, list):
                            for item in to:
                                if isinstance(item, str):
                                    add_edge(frm, str(symbol), item)
                                elif isinstance(item, dict):
                                    add_edge(frm, str(symbol), item.get("target") or item.get("to"))
                elif isinstance(outs, list):
                    for to in outs:
                        if isinstance(to, str):
                            add_edge(frm, None, to)
        elif isinstance(transitions, list):
            for t in transitions:
                if isinstance(t, dict):
                    add_edge(t.get("from") or t.get("source"), t.get("on") or t.get("event"), t.get("to") or t.get("target"))

        # 2) States list/map for names (top-level generic)
        if isinstance(states, dict):
            for name, node in states.items():
                add_node(name)
        elif isinstance(states, list):
            for it in states:
                if isinstance(it, str):
                    add_node(it)
                elif isinstance(it, dict):
                    add_node(it.get("name") or it.get("id") or it.get("state"))

        # 2b) Support alternative graph shapes: nodes/edges at top-level
        nodes_alt = machine.get("nodes") or machine.get("Nodes")
        if isinstance(nodes_alt, list):
            for it in nodes_alt:
                if isinstance(it, str):
                    add_node(it)
                elif isinstance(it, dict):
                    add_node(it.get("name") or it.get("label") or it.get("id") or it.get("state"))

        edges_alt = machine.get("edges") or machine.get("Edges")
        if isinstance(edges_alt, list):
            for e in edges_alt:
                if isinstance(e, dict):
                    add_edge(e.get("from") or e.get("source"), e.get("label") or e.get("on") or e.get("event"), e.get("to") or e.get("target"))

        # 2c) Software Factory shape: agents.orchestrator.states with valid_transitions
        try:
            agents = machine.get("agents") or {}
            orch = agents.get("orchestrator") or {}
            orch_states = orch.get("states") or {}
            if isinstance(orch_states, dict):
                for s_name, s_def in orch_states.items():
                    add_node(s_name)
                    if isinstance(s_def, dict):
                        vts = s_def.get("valid_transitions")
                        if isinstance(vts, list):
                            for tgt in vts:
                                if isinstance(tgt, str):
                                    add_edge(s_name, None, tgt)
        except Exception:
            pass

        # 2d) Software Factory shape: transition_matrix.orchestrator
        try:
            tmat = machine.get("transition_matrix") or {}
            orch_tm = tmat.get("orchestrator") or {}
            if isinstance(orch_tm, dict):
                for frm, outs in orch_tm.items():
                    add_node(frm)
                    if isinstance(outs, list):
                        for to in outs:
                            if isinstance(to, str):
                                add_edge(frm, None, to)
        except Exception:
            pass

        # 3) Fallback: derive edges from per-state "on" maps (XState-style), recursively
        def walk_state(name: str, node: dict):
            if not isinstance(node, dict):
                return
            on_map = node.get("on") or node.get("On") or node.get("transitions")
            if isinstance(on_map, dict):
                for ev, target in on_map.items():
                    if isinstance(target, str):
                        add_edge(name, str(ev), target)
                    elif isinstance(target, list):
                        for item in target:
                            if isinstance(item, str):
                                add_edge(name, str(ev), item)
                            elif isinstance(item, dict):
                                add_edge(name, str(ev), item.get("target") or item.get("to"))
            # Nested states
            nested = node.get("states") or node.get("States")
            if isinstance(nested, dict):
                for child_name, child_node in nested.items():
                    add_node(child_name)
                    walk_state(child_name, child_node)
            elif isinstance(nested, list):
                for it in nested:
                    if isinstance(it, dict):
                        child_name = it.get("name") or it.get("id") or it.get("state")
                        add_node(child_name)
                        walk_state(child_name or "", it)

        if isinstance(states, dict):
            for name, node in states.items():
                walk_state(name, node)

        return node_names, edges

    # Build ASCII graph text for the state machine
    def _build_graph_text(self, machine: dict, current_state: str, mode: str) -> str:
        if not isinstance(machine, dict):
            return ""
        nodes, edges = self._extract_states_transitions(machine)
        if not nodes and not edges:
            return ""
        # 1) Build layered layout (BFS) starting from a likely root
        from collections import deque, defaultdict
        graph = defaultdict(list)
        for s, sym, t in edges:
            graph[s].append((sym, t))
        start = None
        for candidate in ("INIT", "START", "WAVE_START", current_state):
            if candidate in nodes:
                start = candidate
                break
        if start is None and nodes:
            start = nodes[0]
        layers = []
        seen = set()
        q = deque([start] if start else [])
        seen |= set(q)
        while q:
            level_size = len(q)
            layer = []
            for _ in range(level_size):
                n = q.popleft()
                layer.append(n)
                for _sym, t in graph.get(n, [])[:8]:
                    if t not in seen:
                        seen.add(t)
                        q.append(t)
            if layer:
                layers.append(layer)
        # Add any remaining nodes not connected to start
        remaining = [n for n in nodes if n not in seen]
        if remaining:
            layers.append(remaining)

        # If zoom mode, reduce to n-1, n, n+1 layers around the current state
        if mode == "zoom" and layers:
            cur_idx = 0
            for i, layer in enumerate(layers):
                if current_state in layer:
                    cur_idx = i
                    break
            start_i = max(0, cur_idx - 1)
            end_i = min(len(layers), cur_idx + 2)
            layers = layers[start_i:end_i]

        # 2) Prepare canvas and draw boxes row by row
        try:
            term_w = max(80, int(self.size.width))
        except Exception:
            import shutil as _sh
            term_w = _sh.get_terminal_size((160, 40)).columns
        PAD_X = 1
        # Adaptive box width based on available width to make full graph fit in side pane
        # Compute an adaptive NAME_W so the entire widest row fits the panel
        # Estimate max nodes per row by taking the widest layer
        widest = max((len(layer) for layer in layers), default=1)
        # Reserve some padding and column spacing
        COL_SP = 2 if term_w < 60 else (3 if term_w < 90 else 4)
        # NAME_W derived from available width divided by nodes, clamped to [8, 24]
        avail = max(20, term_w - 2 * PAD_X - (widest - 1) * COL_SP)
        NAME_W = max(8, min(24, (avail // max(1, widest)) - 2))
        BOX_W = NAME_W + 2
        ROW_GAP = max(3, int(self._row_gap))  # vertical spacing between rows (larger = more spread)

        def new_row() -> list[str]:
            return [" "] * term_w

        canvas: list[list[str]] = []

        def draw_text(r: int, c: int, text: str):
            if r < 0 or c < 0:
                return
            while r >= len(canvas):
                canvas.append(new_row())
            row = canvas[r]
            for i, ch in enumerate(text):
                x = c + i
                if 0 <= x < term_w:
                    row[x] = ch

        def draw_hline(r: int, x1: int, x2: int):
            if x2 < x1:
                x1, x2 = x2, x1
            while r >= len(canvas):
                canvas.append(new_row())
            for x in range(max(0, x1), min(term_w, x2 + 1)):
                canvas[r][x] = "─"

        def draw_vline(c: int, y1: int, y2: int):
            if y2 < y1:
                y1, y2 = y2, y1
            while y2 >= len(canvas):
                canvas.append(new_row())
            for y in range(max(0, y1), y2 + 1):
                if 0 <= c < term_w:
                    canvas[y][c] = "│"

        # Break layers into rows that fit horizontally
        rows: list[list[str]] = []
        for layer in layers:
            x_needed = 0
            cur: list[str] = []
            for n in layer:
                need = (BOX_W if not cur else BOX_W + COL_SP)
                if PAD_X + x_needed + need > term_w - PAD_X and cur:
                    rows.append(cur)
                    cur = [n]
                    x_needed = BOX_W
                else:
                    cur.append(n)
                    x_needed += need
            if cur:
                rows.append(cur)

        # Compute centers and draw boxes
        centers_per_row: list[list[int]] = []
        y = 0
        for idx_row, row in enumerate(rows):
            y = idx_row * ROW_GAP
            total_w = len(row) * BOX_W + (len(row) - 1) * COL_SP
            start_x = max(PAD_X, (term_w - total_w) // 2)
            centers = [start_x + i * (BOX_W + COL_SP) + BOX_W // 2 for i in range(len(row))]
            centers_per_row.append(centers)
            # draw each box
            for name, cx in zip(row, centers):
                label = name[:NAME_W]
                # Use non-markup brackets so Textual/Rich doesn't strip content
                left, right = ("⟪", "⟫") if name == current_state else ("⟦", "⟧")
                box = f"{left}{label.center(NAME_W)}{right}"
                draw_text(y, cx - BOX_W // 2, box)
            # spacer rows are accounted for by ROW_GAP

        # Draw connectors between adjacent rows
        for i in range(len(rows) - 1):
            row_src, row_tgt = rows[i], rows[i + 1]
            cs_src, cs_tgt = centers_per_row[i], centers_per_row[i + 1]
            idx_tgt = {n: j for j, n in enumerate(row_tgt)}
            # y position for connectors between row i and i+1
            y_box = i * ROW_GAP
            y_next_box = (i + 1) * ROW_GAP
            y_conn = y_box + max(2, ROW_GAP // 2)
            for s_idx, s_name in enumerate(row_src):
                cx_s = cs_src[s_idx]
                # vertical from source box down to connector line
                draw_vline(cx_s, y_box + 1, y_conn)
                outs = [(sym, t) for (ss, sym, t) in edges if ss == s_name and t in idx_tgt]
                for (sym, t) in outs:
                    cx_t = cs_tgt[idx_tgt[t]]
                    # horizontal connector
                    draw_hline(y_conn, cx_s, cx_t)
                    draw_text(y_conn, cx_t, ">")
                    if sym:
                        lab = f"[{sym}]"
                        mid = (cx_s + cx_t) // 2
                        draw_text(y_conn, max(PAD_X, mid - len(lab) // 2), lab)
                    # vertical into target box
                    if y_next_box - (y_conn + 1) > 0:
                        draw_vline(cx_t, y_conn + 1, y_next_box - 1)

        # YOU ARE HERE pointer under current row
        for i, row in enumerate(rows):
            if current_state in row:
                idx = row.index(current_state)
                cx = centers_per_row[i][idx]
                draw_text(i * ROW_GAP + max(1, ROW_GAP - 1), max(0, cx - 8), "▲ YOU ARE HERE")
                break

        lines = ["".join(r).rstrip() for r in canvas]

        # 3) Non-adjacent edges and loops
        others = []
        layer_pos = {}
        for i, layer in enumerate(layers):
            for n in layer:
                layer_pos[n] = i
        for s, sym, t in edges:
            if s in layer_pos and t in layer_pos and abs(layer_pos[t]-layer_pos[s]) <= 1:
                continue
            arrow = f" -[{sym}]-> " if sym else " --> "
            others.append(f"{s}{arrow}{t}")
        if others:
            lines.append("Other Edges:")
            for e in others[:50]:
                lines.append("  " + e)

        return "\n".join(lines[:600]) + "\n"

    # Alias without underscore for future refs
    def build_graph_text(self, machine: dict, current_state: str, mode: str) -> str:  # pragma: no cover
        return self._build_graph_text(machine, current_state, mode)

    # --- Scrolling controls ---
    def scroll_up(self, n: int = 1) -> None:
        self._scroll_offset = max(0, self._scroll_offset - max(1, n))
        self.refresh_panel()

    def scroll_down(self, n: int = 1) -> None:
        total = len(self._last_full_lines) if self._last_full_lines else 0
        try:
            viewport = max(5, int(self.size.height))
        except Exception:
            viewport = 60
        max_off = max(0, total - viewport)
        self._scroll_offset = min(max_off, self._scroll_offset + max(1, n))
        self.refresh_panel()

    def scroll_top(self) -> None:
        self._scroll_offset = 0
        self.refresh_panel()

    def scroll_bottom(self) -> None:
        total = len(self._last_full_lines) if self._last_full_lines else 0
        try:
            viewport = max(5, int(self.size.height))
        except Exception:
            viewport = 60
        self._scroll_offset = max(0, total - viewport)
        self.refresh_panel()


class StatePlot(Widget):
    pass


class StatusTop(Static):
    def update_status(self, cwd: str, claude_cmd: str, script_cmd: str, scheduled: str | None = None) -> None:
        sched = f"  |  Scheduled: [bold red]{scheduled}[/]" if scheduled else ""
        self.update(
            f"Dir: [bold]{cwd}[/]  |  Claude Cmd: [italic]{claude_cmd}[/]  |  Cmd: [italic]{script_cmd}[/]{sched}"
        )


class StatusBottom(Static):
    def update_counts(self, messages: int, continuations: int, commit: str | None = None) -> None:
        suffix = f"  |  Commit: [bold]{commit}[/]" if commit else ""
        self.update(f"Messages: [bold]{messages}[/]  |  Continuations: [bold]{continuations}[/]{suffix}")


class StreamList(ListView):
    """List-based stream that supports grouping and wrapping-friendly rendering.

    Provides compatibility helpers: write_line(), clear(), max_lines attribute.
    """

    def __init__(self, *args, **kwargs) -> None:
        super().__init__(*args, **kwargs)
        self._max_items: int = 10000
        # compatibility toggles
        self.wrap: bool = True

    # compatibility with previous Log API
    @property
    def max_lines(self) -> int:
        return self._max_items

    @max_lines.setter
    def max_lines(self, value: int) -> None:  # type: ignore
        self._max_items = max(100, int(value))
        self._trim()

    def set_max_items(self, value: int) -> None:
        self._max_items = max(100, int(value))
        self._trim()

    def _trim(self) -> None:
        # Remove from the top if we exceed the cap
        try:
            while len(self.children) > self._max_items:
                self.children[0].remove()
        except Exception:
            pass

    def write_line(self, text: str) -> None:
        # Use ExpandableItem to prevent rows from flexing to full height
        item = ExpandableItem(text, "", expanded=False)
        self.append(item)
        self._trim()

    def append_info(self, text: str) -> None:
        self.write_line(f"[cyan]{text}[/]")

    def append_error(self, text: str) -> None:
        self.write_line(f"[red]{text}[/]")

    def append_debug(self, text: str) -> None:
        self.write_line(f"[dim][debug][/]: {text}")

    def append_banner(self, text: str) -> None:
        self.write_line(f"[bold white on green]{text}[/]")

    def append_raw(self, text: str) -> None:
        self.write_line(text)

    def add_parsed(self, parsed: Dict[str, Any], show_detail: bool, pretty: bool) -> None:
        prefix = ""
        raw_event = parsed.get("raw_event")
        if raw_event == "message_start":
            prefix = "▶ "
        elif raw_event == "message_delta":
            prefix = "… "
        elif raw_event in {"message_stop", "finish"}:
            prefix = "■ "
        summary = f"{prefix}{parsed.get('summary','')}"
        # Build normalized compact preview
        extras = parsed.get("extras_text") or "no additional info"
        detail_primary = parsed.get("detail_pretty") if pretty else parsed.get("detail_raw")
        preview_lines: list[str] = []
        if extras.strip() and extras.strip().lower() != "no additional info":
            preview_lines.append(extras)
        if isinstance(detail_primary, str) and detail_primary.strip():
            lines = detail_primary.splitlines()
            max_detail_lines = 80
            preview_lines.extend(lines[:max_detail_lines])
            if len(lines) > max_detail_lines:
                preview_lines.append("… (truncated)")
        detail = "\n".join(preview_lines) if preview_lines else "no additional info"
        item = ExpandableItem(summary, detail, expanded=bool(show_detail))
        self.append(item)
        self._trim()

    def clear(self) -> None:  # type: ignore[override]
        try:
            super().clear()
        except Exception:
            # Fallback: remove children manually
            for ch in list(self.children):
                try:
                    ch.remove()
                except Exception:
                    pass

    # Toggling helpers
    def toggle_index(self, index: int) -> None:
        try:
            items = [ch for ch in self.children if isinstance(ch, ListItem)]
            if 0 <= index < len(items):
                item = items[index]
                if isinstance(item, ExpandableItem):
                    item.toggle()
        except Exception:
            pass

    def toggle_current(self) -> None:
        idx = None
        try:
            idx = getattr(self, "index", None)
        except Exception:
            idx = None
        if isinstance(idx, int):
            self.toggle_index(idx)

    def expand_all(self) -> None:
        try:
            for ch in self.children:
                if isinstance(ch, ExpandableItem):
                    ch.set_expanded(True)
        except Exception:
            pass

    def collapse_all(self) -> None:
        try:
            for ch in self.children:
                if isinstance(ch, ExpandableItem):
                    ch.set_expanded(False)
        except Exception:
            pass

    # Clicking within the list toggles the currently selected item
    def on_click(self, event) -> None:  # type: ignore[override]
        try:
            from textual import events as _events
            if isinstance(event, _events.Click):
                # Ensure focus to update selection highlight, then toggle
                try:
                    self.focus()
                except Exception:
                    pass
                self.toggle_current()
                event.stop()
        except Exception:
            # Fallback: always attempt to toggle current
            self.toggle_current()


class StreamTree(Tree[str]):
    """Tree-based stream with auto-sized expandable detail nodes."""

    def __init__(self, *args, **kwargs) -> None:
        # Give the root a visible label so the pane isn't blank when empty
        super().__init__("Logs:", *args, **kwargs)
        self._max_items: int = 10000
        self._autoscroll: bool = True
        self._user_scrolled_up: bool = False
        # Expand root by default so users see the container immediately
        try:
            self.root.expand()
        except Exception:
            pass

    def set_max_items(self, value: int) -> None:
        self._max_items = max(100, int(value))
        self._trim()

    def _trim(self) -> None:
        try:
            while len(self.root.children) > self._max_items:
                self.root.children[0].remove()
        except Exception:
            pass

    # --- Auto-scroll helpers ---
    @property
    def autoscroll(self) -> bool:
        return self._autoscroll

    def set_autoscroll(self, value: bool) -> None:
        self._autoscroll = bool(value)
        # If enabling autoscroll, reset manual override and jump to bottom
        if self._autoscroll:
            self._user_scrolled_up = False
            self.force_scroll_bottom()

    def force_scroll_bottom(self) -> None:
        try:
            # Prefer a native helper if available
            self.scroll_end(animate=False)  # type: ignore[arg-type]
        except Exception:
            try:
                # Fallback: large y to ensure bottom
                self.scroll_to(y=10**9)
            except Exception:
                pass
        # After explicit bottom, allow autoscroll again
        self._user_scrolled_up = False

    def _maybe_autoscroll(self) -> None:
        try:
            if self._autoscroll and not self._user_scrolled_up:
                self.force_scroll_bottom()
        except Exception:
            pass

    def clear(self) -> None:  # type: ignore[override]
        try:
            self.root.clear()
        except Exception:
            pass
        # Keep autoscroll behavior after clear
        self._maybe_autoscroll()

    # Optimize bulk expansions: never expand more than one heavy node at a time
    def on_tree_node_expanded(self, event) -> None:  # type: ignore[override]
        try:
            # If many siblings are expanded, collapse older ones to keep DOM light
            expanded = [ch for ch in self.root.children if getattr(ch, 'is_expanded', False)]
            MAX_EXPANDED = 3
            if len(expanded) > MAX_EXPANDED:
                for n in expanded[:-MAX_EXPANDED]:
                    try:
                        n.collapse()
                    except Exception:
                        pass
        except Exception:
            pass

    # Simple helpers
    def append_info(self, text: str) -> None:
        self.root.add(f"[cyan]{text}[/]")
        self._trim()
        self._maybe_autoscroll()

    def append_error(self, text: str) -> None:
        self.root.add(f"[red]{text}[/]")
        self._trim()
        self._maybe_autoscroll()

    def append_debug(self, text: str) -> None:
        self.root.add(f"[dim][debug][/]: {text}")
        self._trim()
        self._maybe_autoscroll()

    def append_banner(self, text: str) -> None:
        self.root.add(f"[bold white on green]{text}[/]")
        self._trim()
        self._maybe_autoscroll()

    def append_raw(self, text: str) -> None:
        self.root.add(text)
        self._trim()
        self._maybe_autoscroll()

    def add_parsed(self, parsed: Dict[str, Any], show_detail: bool, pretty: bool) -> None:
        prefix = ""
        raw_event = parsed.get("raw_event")
        if raw_event == "message_start":
            prefix = "▶ "
        elif raw_event == "message_delta":
            prefix = "… "
        elif raw_event in {"message_stop", "finish"}:
            prefix = "■ "
        summary = f"{prefix}{parsed.get('summary','')}"
        node = self.root.add(summary)
        # Build a normalized compact preview block: extras + pretty/raw JSON
        extras = parsed.get("extras_text") or "no additional info"
        detail_primary = parsed.get("detail_pretty") if pretty else parsed.get("detail_raw")
        # Build child lines so Tree renders all content, not just first line
        child_lines: list[str] = []
        # Show extras header first if meaningful
        if extras.strip() and extras.strip().lower() != "no additional info":
            child_lines.append(extras)
        # Then append a condensed JSON preview (first N lines)
        if isinstance(detail_primary, str) and detail_primary.strip():
            lines = detail_primary.splitlines()
            max_detail_lines = 80
            child_lines.extend(lines[:max_detail_lines])
            if len(lines) > max_detail_lines:
                child_lines.append("… (truncated)")
        if not child_lines:
            child_lines.append("no additional info")
        # Add lines as separate children
        for ln in child_lines[:300]:
            node.add(ln)
        # Avoid expanding many heavy nodes at once; expand only when explicitly requested
        if show_detail:
            try:
                node.expand()
            except Exception:
                pass
        self._trim()
        self._maybe_autoscroll()

    def expand_all(self) -> None:
        try:
            for n in list(self.root.children):
                n.expand()
        except Exception:
            pass

    def collapse_all(self) -> None:
        try:
            for n in list(self.root.children):
                n.collapse()
        except Exception:
            pass

    # Mark that the user intentionally scrolled up; temporarily disable auto-scroll
    def on_mouse_scroll_up(self, event) -> None:  # type: ignore[override]
        try:
            self._user_scrolled_up = True
        except Exception:
            pass

    # If the user scrolls down to the bottom, re-enable following
    def on_mouse_scroll_down(self, event) -> None:  # type: ignore[override]
        try:
            # Heuristic: if currently near bottom, consider it bottom and re-enable
            near_bottom = False
            try:
                max_y = getattr(self, "virtual_size", None)
                cur_y = getattr(self, "scroll_y", None)
                if isinstance(max_y, tuple) and isinstance(cur_y, int):
                    _, total_h = max_y
                    if total_h - cur_y < 20:
                        near_bottom = True
            except Exception:
                pass
            if near_bottom:
                self._user_scrolled_up = False
        except Exception:
            pass


class ExpandableItem(ListItem):
    """A ListItem that shows a summary and can expand to reveal details.

    - Click, Enter, or Space toggles expansion.
    - Uses Static widgets so long text wraps naturally within the pane width.
    """

    def __init__(self, summary_text: str, detail_text: str | None, expanded: bool = False) -> None:
        self._base_summary = summary_text
        self.summary_widget = Static(summary_text, classes="stream-summary")
        self.detail_widget = Static(detail_text or "", classes="stream-detail")
        try:
            self.summary_widget.wrap = True  # type: ignore[attr-defined]
            self.detail_widget.wrap = True  # type: ignore[attr-defined]
        except Exception:
            pass
        # Wrap-friendly layout: vertical group with summary then detail
        container = VGroup(self.summary_widget, self.detail_widget)
        super().__init__(container)
        self._expanded = False
        try:
            self.can_focus = True
        except Exception:
            pass
        # Prevent this row from expanding to fill the viewport
        try:
            self.styles.flex_grow = 0
            self.styles.height = None
            container.styles.flex_grow = 0
            container.styles.height = None
            self.summary_widget.styles.height = None
            self.detail_widget.styles.height = None
        except Exception:
            pass
        # Start collapsed unless expanded requested
        self.set_expanded(expanded)
        self._apply_summary()

    def set_expanded(self, expanded: bool) -> None:
        self._expanded = bool(expanded)
        try:
            self.detail_widget.display = ("block" if self._expanded else "none")
        except Exception:
            pass
        self._apply_summary()
        try:
            # Force a layout refresh so wrapping recalculates
            self.refresh(layout=True)
        except Exception:
            pass

    def _apply_summary(self) -> None:
        try:
            caret = "▾" if self._expanded else "▸"
            self.summary_widget.update(f"{caret} {self._base_summary}")
        except Exception:
            pass

    def toggle(self) -> None:
        self.set_expanded(not self._expanded)

    # Input handlers to toggle expansion
    def on_click(self, event) -> None:  # type: ignore[override]
        try:
            from textual import events as _events
            if isinstance(event, _events.Click):
                self.toggle()
                event.stop()
                return
        except Exception:
            pass
        self.toggle()

    # Rely on parent list click to toggle selection; keep item-level click for redundancy

    def on_key(self, event) -> None:  # type: ignore[override]
        try:
            from textual import events as _events
            if isinstance(event, _events.Key) and (event.key in ("enter", "space")):
                self.toggle()
                event.stop()
                return
        except Exception:
            pass
        try:
            key = getattr(event, "key", None)
            if key in ("enter", "space"):
                self.toggle()
        except Exception:
            pass

    # Ensure the ListView measures this row to the number of lines needed
    def get_content_height(self, container, viewport_width: int, viewport_height: int) -> int:  # type: ignore[override]
        # Estimate height conservatively: summary + detail text lines (min 1)
        try:
            s = str(self.summary_widget.renderable)
        except Exception:
            s = ""
        try:
            d = str(self.detail_widget.renderable) if self._expanded and self.detail_widget.display != "none" else ""
        except Exception:
            d = ""
        s_lines = max(1, len(s.splitlines()) or 1)
        d_lines = max(0, len(d.splitlines()))
        return max(1, s_lines + d_lines)

class ClaudeOrchestratorApp(App):
    CSS = """
    Screen { layout: vertical; }
    #top { height: 3; }
    #banner { height: 1; }
    #log { height: 1fr; overflow: auto; }
    #bottom { height: 1; }
    #state { width: 80; border: round green; padding: 0 1; height: 1fr; overflow: auto; }
    .stream-summary { text-wrap: wrap; width: 1fr; }
    .stream-detail { text-wrap: wrap; width: 1fr; color: #a0a0a0; }
    """

    BINDINGS = [
        Binding("ctrl+p", "command_palette", "Command Palette"),
        Binding("ctrl+q", "quit", "Quit"),
    ]

    running = reactive(False)
    messages_count = reactive(0)
    continuations_count = reactive(0)
    detail = reactive(False)
    pretty = reactive(False)
    wrap_long = reactive(True)
    scheduled_start_at = reactive(None)

    def __init__(self, config: RunnerConfig) -> None:
        super().__init__()
        self.config = config
        self.runner = ClaudeRunner(config)
        self.parser = MessageParser()
        self._task: Optional[asyncio.Task] = None
        self._log_file = open(self.config.log_path, "a", buffering=1) if self.config.log_path else None
        self.wrap_enabled = False
        self.state_fullscreen = False
        self._state_refresh_timer = None
        self._state_refresh_once = None
        self._commit_short = self._get_commit_short()
        # Capture the command used to invoke this TUI process (not the Claude subprocess)
        self._tui_invocation_cmd = self._get_tui_invocation()
        # Periodic schedule checker
        self._schedule_check_timer = None
        # Ensure we always clean up subprocesses on interpreter exit
        try:
            atexit.register(self._atexit_cleanup)
        except Exception:
            pass
        # Track the most recent actual Claude command we spawned
        self._current_claude_cmd: str = self.config.command

    def compose(self) -> ComposeResult:
        yield Header(show_clock=True)
        with Vertical():
            self.top = StatusTop(id="top")
            yield self.top
            self.banner = IterationBanner(id="banner")
            yield self.banner
            with Horizontal():
                # Switch to Tree-based stream to avoid row flex height issues
                self.stream = StreamTree(id="log")
                # Ensure the list can receive focus first for keyboard toggles
                try:
                    self.stream.can_focus = True
                except Exception:
                    pass
                yield self.stream
                # Right side state panel (hidden by default)
                self.state_panel = StatePanel(self.config.sf2_project_dir)
                self.state_panel.id = "state"
                self.state_panel.toggle()  # default show
                yield self.state_panel
                # Plot view removed
            self.bottom = StatusBottom(id="bottom")
            yield self.bottom
        yield Footer()

    async def on_mount(self) -> None:
        # Textual v6 uses get_system_commands; no explicit registration needed.
        self._update_status()
        # Initial render of state panel so it isn't empty on first load
        if hasattr(self, "state_panel"):
            # Try to pre-load viewer data then refresh immediately
            try:
                if getattr(self.state_panel, "_viewer", None):
                    try:
                        self.state_panel._viewer.load_data()
                    except Exception:
                        pass
                self.state_panel.refresh_panel()
            except Exception:
                pass
            # Schedule a deferred refresh once layout has stabilized
            try:
                self.set_timer(0.5, self.state_panel.refresh_panel)
            except Exception:
                try:
                    def _one_shot():
                        try:
                            self.state_panel.refresh_panel()
                        finally:
                            try:
                                if self._state_refresh_once is not None:
                                    self._state_refresh_once.stop()
                            except Exception:
                                pass
                    self._state_refresh_once = self.set_interval(0.5, _one_shot)
                except Exception:
                    pass
        # Configure circular buffer for the stream list
        try:
            self.stream.set_max_items(10000)
        except Exception:
            pass
        # Periodically refresh the state graph pane every 30 seconds
        try:
            self._state_refresh_timer = self.set_interval(30.0, self._refresh_state_panel_periodic)
        except Exception:
            self._state_refresh_timer = None
        # Periodically check scheduled start time
        try:
            self._schedule_check_timer = self.set_interval(30.0, self._check_scheduled_start)
        except Exception:
            self._schedule_check_timer = None

    def _update_status(self) -> None:
        cwd = self.config.working_dir
        # Show the TUI's own invocation command, not the Claude subprocess command
        script_cmd = self._tui_invocation_cmd
        claude_cmd = getattr(self, "_current_claude_cmd", script_cmd)
        # Prepare scheduled time for display (local time)
        sched_local = None
        try:
            when = self.scheduled_start_at
            if isinstance(when, datetime):
                sched_local = when.astimezone().strftime("%I:%M:%S %p %Z %m/%d/%Y").lstrip("0")
        except Exception:
            sched_local = None
        self.top.update_status(cwd, claude_cmd, script_cmd, sched_local)
        self.bottom.update_counts(self.messages_count, self.continuations_count, self._commit_short)

    def _get_commit_short(self) -> str | None:
        try:
            import subprocess, pathlib
            repo = pathlib.Path(__file__).resolve().parents[1]
            out = subprocess.check_output(["git", "-C", str(repo), "rev-parse", "--short", "HEAD"], text=True)
            return out.strip()
        except Exception:
            return None

    def _get_tui_invocation(self) -> str:
        """Reconstruct the command used to launch this TUI process.

        Preference order:
        - /proc/self/cmdline on Linux for the exact argv
        - sys.executable + sys.argv as a portable fallback
        - script filename as last resort
        """
        try:
            # Linux: read the raw null-separated command line
            cmdline_path = "/proc/self/cmdline"
            if os.path.exists(cmdline_path):
                with open(cmdline_path, "rb") as f:
                    parts = f.read().split(b"\0")
                parts = [p.decode("utf-8", "replace") for p in parts if p]
                if parts:
                    return " ".join(shlex.quote(p) for p in parts)
        except Exception:
            pass
        try:
            argv = [sys.executable] + list(sys.argv)
            return " ".join(shlex.quote(p) for p in argv)
        except Exception:
            pass
        try:
            return str(Path(__file__).resolve())
        except Exception:
            return "tui"

    def action_command_palette(self) -> None:
        # Built-in action opens the palette
        super().action_command_palette()

    def get_system_commands(self, screen) -> list[SystemCommand]:
        return [
            SystemCommand("Start", "Start Claude with current parameters", self.action_start),
            SystemCommand("Stop", "Stop Claude", self.action_stop),
            SystemCommand("Restart", "Restart Claude", self.action_restart),
            SystemCommand("Set Params", "Edit parameters (-p env or string)", self.action_set_params),
            SystemCommand("Schedule Start", "Set or cancel a one-time future start", self.action_schedule_start),
            SystemCommand("Toggle Detail", "Toggle detailed message view", self.action_toggle_detail),
            SystemCommand("Toggle Pretty", "Toggle pretty JSON for detail view", self.action_toggle_pretty),
            SystemCommand("Toggle Wrap", "Toggle line wrapping in the log view", self.action_toggle_wrap),
            SystemCommand("Toggle Debug Raw", "Emit raw chunk debug lines", self.action_toggle_debug_raw),
            SystemCommand("Use Mock Command", "Point command to mock_claude_emitter.py", self.action_use_mock),
            SystemCommand("Use Legacy Runner", "Delegate to tools/claude_orchestrator_runner.py", self.action_use_legacy_runner),
            SystemCommand("Toggle State Panel", "Show/Hide the state machine panel", self.action_toggle_state_panel),
            SystemCommand("Refresh State Panel", "Re-read state and redraw ASCII", self.action_refresh_state_panel),
            SystemCommand("Set SF2 Project Dir", "Choose the software-factory project directory", self.action_set_project_dir),
            SystemCommand("State Fullscreen", "Expand state panel to full screen", self.action_state_fullscreen),
            SystemCommand("State Zoomed", "Show zoomed-in state around current state", self.action_state_zoomed),
            SystemCommand("State Full Graph", "Switch state panel to full graph mode", self.action_state_fullgraph),
            SystemCommand("State Scroll Up", "Scroll state panel up", self.action_state_scroll_up),
            SystemCommand("State Scroll Down", "Scroll state panel down", self.action_state_scroll_down),
            SystemCommand("State Page Up", "Page up state panel", self.action_state_page_up),
            SystemCommand("State Page Down", "Page down state panel", self.action_state_page_down),
            SystemCommand("State Scroll Top", "Jump to top of state panel", self.action_state_scroll_top),
            SystemCommand("State Scroll Bottom", "Jump to bottom of state panel", self.action_state_scroll_bottom),
            
            SystemCommand("State Increase Spacing", "Increase vertical spacing between rows", self.action_state_rowgap_inc),
            SystemCommand("State Decrease Spacing", "Decrease vertical spacing between rows", self.action_state_rowgap_dec),
            SystemCommand("Cycle State View", "Cycle journey ↔ focused ↔ overview", self.action_state_cycle_view),
            SystemCommand("Reload State Data", "Reload orchestrator/state-machine JSON", self.action_state_reload_state),
            SystemCommand("State Fit To Pane", "Auto-fit graph to current pane and scroll top", self.action_state_fit_pane),
            SystemCommand("Set Log Buffer Size", "Change max log lines (circular buffer)", self.action_set_log_buffer_size),
            SystemCommand("Toggle Auto-Scroll", "Toggle auto-scroll to bottom on new messages", self.action_toggle_autoscroll),
            SystemCommand("Scroll To Bottom", "Force scroll to bottom of log", self.action_scroll_to_bottom),
            SystemCommand("Clear", "Clear message log", self.action_clear),
            SystemCommand("Simulate", "Toggle simulation mode", self.action_toggle_simulation),
            SystemCommand("Expand All", "Expand all list items (show details)", self.action_expand_all),
            SystemCommand("Collapse All", "Collapse all list items (hide details)", self.action_collapse_all),
            SystemCommand("Quit", "Exit the app", self.action_quit),
        ]

    async def action_start(self) -> None:
        if self.running:
            return
        # Preflight check: ensure binary exists
        try:
            first_token = shlex.split(self.config.command)[0]
        except Exception:
            first_token = "claude"
        if not self.config.simulate and shutil.which(first_token) is None:
            self.stream.append_error(f"Error: command not found: {first_token}. Ensure it is on PATH.")
            return
        self.running = True
        self._task = asyncio.create_task(self._run_loop())

    async def action_stop(self) -> None:
        if not self.running:
            return
        await self.runner.stop()
        self.running = False
        if self._task:
            try:
                await asyncio.wait_for(self._task, timeout=2)
            except Exception:
                pass

    async def action_restart(self) -> None:
        await self.action_stop()
        await self.action_start()

    async def action_set_params(self) -> None:
        self.push_screen(ParamsScreen(self))

    async def action_schedule_start(self) -> None:
        self.push_screen(ScheduleScreen(self))

    async def action_toggle_detail(self) -> None:
        self.detail = not self.detail

    async def action_toggle_pretty(self) -> None:
        self.pretty = not self.pretty

    async def action_toggle_wrap(self) -> None:
        # Toggle wrap flag; fallback to emitting a note if attribute unsupported
        self.wrap_long = not self.wrap_long
        updated = False
        try:
            self.stream.wrap = self.wrap_long  # type: ignore[attr-defined]
            updated = True
        except Exception:
            updated = False
        self.stream.append_debug(f"wrap -> {self.wrap_long}{'' if updated else ' (Log compat)'}")

    async def action_expand_all(self) -> None:
        try:
            self.stream.expand_all()
        except Exception:
            pass

    async def action_collapse_all(self) -> None:
        try:
            self.stream.collapse_all()
        except Exception:
            pass

    async def action_clear(self) -> None:
        self.stream.clear()
        self.messages_count = 0
        self._update_status()

    async def action_toggle_state_panel(self) -> None:
        self.state_panel.toggle()
        self.state_panel.refresh_panel()

    async def action_refresh_state_panel(self) -> None:
        self.state_panel.refresh_panel()

    async def action_set_project_dir(self) -> None:
        prompt = Input(placeholder="Enter path to software-factory project dir (contains orchestrator-state-v3.json)")
        container = Vertical(Static("Set SF2 Project Directory:"), prompt)
        await self.mount(container, before=self.query_one(Footer))
        await prompt.focus()
        # Read current value immediately (simple flow)
        new_dir = (prompt.value or "").strip() or self.config.sf2_project_dir
        self.config.sf2_project_dir = new_dir
        self.state_panel.set_project_dir(new_dir)
        self.state_panel.refresh_panel()
        await container.remove()

    async def action_set_log_buffer_size(self) -> None:
        # Prompt for an integer size for the circular log buffer
        prompt = Input(placeholder="Enter max log lines (e.g., 10000)")
        container = Vertical(Static("Set Log Buffer Size:"), prompt)
        await self.mount(container, before=self.query_one(Footer))
        await prompt.focus()
        # Read immediately (simple flow)
        raw = (prompt.value or "").strip()
        await container.remove()
        try:
            size = int(raw) if raw else 10000
            size = max(100, min(500000, size))
            try:
                self.stream.set_max_items(size)
                self.stream.append_debug(f"log buffer -> {size}")
            except Exception:
                self.stream.append_error("Unable to set log buffer size")
        except Exception:
            self.stream.write_line(f"[red]Invalid number:[/] {raw}")

    async def action_toggle_autoscroll(self) -> None:
        try:
            cur = bool(getattr(self.stream, "autoscroll", True))
            self.stream.set_autoscroll(not cur)
            self.stream.append_debug(f"auto-scroll -> {not cur}")
        except Exception:
            self.stream.append_error("Unable to toggle auto-scroll")

    async def action_scroll_to_bottom(self) -> None:
        try:
            self.stream.force_scroll_bottom()
        except Exception:
            pass

    async def action_state_fullscreen(self) -> None:
        # Toggle fullscreen for state panel
        self.state_fullscreen = not self.state_fullscreen
        if self.state_fullscreen:
            self.stream.display = "none"
            self.state_panel.display = "block"
            self.state_panel.styles.width = "1fr"
            self.state_panel.styles.height = "1fr"
            # In fullscreen, default to full graph and reset scroll
            self.state_panel.render_mode = "full"
            self.state_panel.scroll_top()
        else:
            self.stream.display = "block"
            # Restore side pane layout
            self.state_panel.display = "block"
            self.state_panel.styles.width = 80
            self.state_panel.styles.height = "1fr"
        # Force re-render after geometry change
        self.state_panel.refresh_panel()

    async def action_state_zoomed(self) -> None:
        self.state_panel.render_mode = "zoom"
        try:
            self.state_panel._viewer.set_view('focused')
        except Exception:
            pass
        self.state_panel.refresh_panel()

    async def action_state_fullgraph(self) -> None:
        self.state_panel.render_mode = "full"
        try:
            self.state_panel._viewer.set_view('journey')
        except Exception:
            pass
        self.state_panel.refresh_panel()

    async def action_state_scroll_up(self) -> None:
        try:
            self.state_panel._viewer.handle_key('UP')
        except Exception:
            pass
        self.state_panel.scroll_up(1)

    async def action_state_scroll_down(self) -> None:
        try:
            self.state_panel._viewer.handle_key('DOWN')
        except Exception:
            pass
        self.state_panel.scroll_down(1)

    async def action_state_page_up(self) -> None:
        try:
            viewport = max(5, int(self.state_panel.size.height))
        except Exception:
            viewport = 20
        try:
            self.state_panel._viewer.handle_key('PAGEUP')
        except Exception:
            pass
        # Inform viewer too
        self.state_panel.scroll_up(viewport - 2)

    async def action_state_page_down(self) -> None:
        try:
            viewport = max(5, int(self.state_panel.size.height))
        except Exception:
            viewport = 20
        try:
            self.state_panel._viewer.handle_key('PAGEDOWN')
        except Exception:
            pass
        try:
            pass
        except Exception:
            pass
        self.state_panel.scroll_down(viewport - 2)

    async def action_state_scroll_top(self) -> None:
        try:
            self.state_panel._viewer.handle_key('HOME')
        except Exception:
            pass
        self.state_panel.scroll_top()

    async def action_state_scroll_bottom(self) -> None:
        try:
            self.state_panel._viewer.handle_key('END')
        except Exception:
            pass
        self.state_panel.scroll_bottom()

    async def action_state_cycle_view(self) -> None:
        try:
            self.state_panel._viewer.handle_key('v')
        except Exception:
            # Fallback: toggle between full and zoom
            self.state_panel.render_mode = 'zoom' if self.state_panel.render_mode != 'zoom' else 'full'
        # Sync render mode with viewer if available
        try:
            cur = getattr(self.state_panel._viewer, 'current_view', 'journey')
            self.state_panel.render_mode = 'zoom' if cur == 'focused' else 'full'
        except Exception:
            pass
        # Reset local scroll and refresh
        try:
            self.state_panel._scroll_offset = 0
        except Exception:
            pass
        self.state_panel.refresh_panel()

    async def action_state_reload_state(self) -> None:
        try:
            self.state_panel._viewer.load_data()
        except Exception:
            pass
        self.state_panel.refresh_panel()

    async def action_state_rowgap_inc(self) -> None:
        try:
            self.state_panel._row_gap = min(20, int(self.state_panel._row_gap) + 1)
        except Exception:
            self.state_panel._row_gap = 8
        self.state_panel.refresh_panel()

    async def action_state_rowgap_dec(self) -> None:
        try:
            self.state_panel._row_gap = max(3, int(self.state_panel._row_gap) - 1)
        except Exception:
            self.state_panel._row_gap = 6
        self.state_panel.refresh_panel()

    async def action_state_plot_toggle(self) -> None:
        # Toggle plotext widget visibility and keep layout consistent
        if PlotextPlot is None:
            self.stream.write_line("[yellow]textual-plotext not installed. Install with: pip install textual-plotext[/]")
            return
        showing_plot = getattr(self.state_plot, "display", "none") != "none"
        if showing_plot:
            # Switch to ASCII
            self.state_plot.display = "none"
            self.state_panel.display = "block"
            # Restore side width if not fullscreen
            if not self.state_fullscreen:
                self.state_panel.styles.width = 48
            self.state_panel.refresh_panel()
        else:
            # Switch to Plot view
            self.state_plot.display = "block"
            self.state_panel.display = "none"
            if self.state_fullscreen:
                self.state_plot.styles.width = "1fr"
                self.state_plot.styles.height = "1fr"
            else:
                self.state_plot.styles.width = 48
                self.state_plot.styles.height = "1fr"
            try:
                data, _m = self.state_panel._load_state()
                cur = str(data.get("current_state") or data.get("state") or data.get("CurrentState") or "?")
                ok = self.state_plot.refresh_plot(cur)
                if not ok:
                    self.state_plot.draw_demo()
            except Exception:
                try:
                    self.state_plot.draw_demo()
                except Exception:
                    pass

    async def action_state_fit_pane(self) -> None:
        # Reset spacing and scroll, then re-render to fit current pane width
        try:
            # Choose a moderate default row gap for readability
            self.state_panel._row_gap = 5
        except Exception:
            pass
        self.state_panel.scroll_top()
        self.state_panel.refresh_panel()

    async def action_state_plot_demo(self) -> None:
        if PlotextPlot is None:
            self.stream.write_line("[yellow]textual-plotext not installed. Install with: pip install textual-plotext[/]")
            return
        # Ensure plot view is visible and sized
        self.state_plot.display = "block"
        self.state_panel.display = "none"
        if self.state_fullscreen:
            self.state_plot.styles.width = "1fr"
            self.state_plot.styles.height = "1fr"
        else:
            self.state_plot.styles.width = 48
            self.state_plot.styles.height = "1fr"
        try:
            self.state_plot.draw_demo()
        except Exception:
            pass

    async def action_toggle_simulation(self) -> None:
        self.config.simulate = not self.config.simulate
        await self.action_restart()

    async def action_cycle_wrapper(self) -> None:
        order = ["none", "stdbuf", "unbuffer", "script"]
        try:
            idx = order.index((self.config.wrapper_mode or "none").lower())
        except ValueError:
            idx = 0
        self.config.wrapper_mode = order[(idx + 1) % len(order)]
        self.stream.append_debug(f"wrapper -> {self.config.wrapper_mode}")
        await self.action_restart()

    async def action_toggle_debug_raw(self) -> None:
        self.config.debug_raw = not self.config.debug_raw
        self.stream.append_debug(f"debug_raw -> {self.config.debug_raw}")
        await self.action_restart()

    async def action_use_mock(self) -> None:
        mock_path = str(Path(__file__).with_name("mock_claude_emitter.py"))
        self.config.command = f"python3 {shlex.quote(mock_path)} --mode sf2_exact --emit-continuation"
        self.stream.append_debug(f"command -> {self.config.command}")
        await self.action_restart()

    async def action_use_legacy_runner(self) -> None:
        legacy_path = str(Path(__file__).with_name("tools").joinpath("claude_orchestrator_runner.py"))
        # Delegate: run legacy runner; it will handle continuations and emit JSON on stdout
        self.config.command = f"python3 {shlex.quote(legacy_path)}"
        self.stream.write_line(f"[debug] command -> {self.config.command}")
        await self.action_restart()

    async def action_quit(self) -> None:
        await self.action_stop()
        # Belt-and-suspenders: ensure process group is gone
        try:
            self.runner.kill_now()
        except Exception:
            pass
        if self._log_file is not None:
            try:
                self._log_file.close()
            except Exception:
                pass
        self.exit()

    async def on_shutdown(self) -> None:
        await self.action_stop()
        try:
            self.runner.kill_now()
        except Exception:
            pass
        if self._log_file is not None:
            try:
                self._log_file.close()
            except Exception:
                pass
        # Cancel periodic state refresh timer if present
        try:
            if self._state_refresh_timer is not None:
                self._state_refresh_timer.stop()
        except Exception:
            pass
        # Cancel schedule checker
        try:
            if self._schedule_check_timer is not None:
                self._schedule_check_timer.stop()
        except Exception:
            pass

    # Minimal compatibility API for tests: return command-like objects with .label
    def get_system_commands(self, *_args, **_kwargs):  # pragma: no cover - test helper
        class _Cmd:
            def __init__(self, label, help, action):
                self.label = label
                self.help = help
                self.action = action
        return [
            _Cmd("Start", "Start Claude", self.action_start),
            _Cmd("Stop", "Stop Claude", self.action_stop),
            _Cmd("Schedule Start", "Schedule Claude start time", self._set_schedule_from_string),
        ]

    # --- Scheduling helpers ---
    def _set_schedule_from_string(self, value: str) -> None:
        value = (value or "").strip()
        if not value:
            self.scheduled_start_at = None
            try:
                self.stream.append_info("Claude start schedule cleared")
            except Exception:
                pass
            return
        dt = self._parse_flexible_datetime(value)
        if dt is None:
            try:
                self.stream.append_error(f"Unable to parse date/time: {value}")
            except Exception:
                pass
            return
        # ensure timezone; default to local -> convert to UTC
        if dt.tzinfo is None:
            try:
                local_tz = datetime.now().astimezone().tzinfo or timezone.utc
                dt = dt.replace(tzinfo=local_tz)
            except Exception:
                dt = dt.replace(tzinfo=timezone.utc)
        dt_utc = dt.astimezone(timezone.utc)
        self.scheduled_start_at = dt_utc
        try:
            disp = dt_utc.astimezone().strftime("%I:%M:%S %p %Z %m/%d/%Y").lstrip("0")
            self.stream.append_info(f"Claude start scheduled for {disp}")
        except Exception:
            pass

    def _parse_flexible_datetime(self, text: str) -> Optional[datetime]:
        s = (text or "").strip()
        # Relative: now+5m, now+2h, now+1d
        try:
            low = s.lower()
            if low.startswith("now+"):
                num = ""
                unit = ""
                for ch in s[4:]:
                    if ch.isdigit():
                        num += ch
                    else:
                        unit += ch
                n = int(num) if num else 0
                delta = timedelta()
                if unit.lower().startswith("m"):
                    delta = timedelta(minutes=n)
                elif unit.lower().startswith("h"):
                    delta = timedelta(hours=n)
                elif unit.lower().startswith("d"):
                    delta = timedelta(days=n)
                return datetime.now().astimezone() + delta
        except Exception:
            pass

        # Pattern A/B: 12:15am EST 2025/09/25  OR  12:15am EST 09/25/2025
        try:
            import re as _re
            pat_y_m_d = r"^\s*(\d{1,2}):(\d{2})(am|pm)\s*([A-Za-z]{2,4})\s*(\d{4})[-/](\d{1,2})[-/](\d{1,2})\s*$"
            pat_m_d_y = r"^\s*(\d{1,2}):(\d{2})(am|pm)\s*([A-Za-z]{2,4})\s*(\d{1,2})[-/](\d{1,2})[-/](\d{2,4})\s*$"
            m = _re.match(pat_y_m_d, s, _re.IGNORECASE)
            if m:
                hh, mm, ap, tzs, year, mo, d = m.groups()
            else:
                m = _re.match(pat_m_d_y, s, _re.IGNORECASE)
                if m:
                    hh, mm, ap, tzs, mo, d, year = m.groups()
                else:
                    m = None
            if m:
                hour = int(hh) % 12
                if ap.lower() == "pm":
                    hour += 12
                year_i = int(year) if len(str(year)) == 4 else 2000 + int(year)
                dt_naive = datetime(year_i, int(mo), int(d), hour, int(mm))
                tz_map = {"UTC": 0, "Z": 0, "EST": -5, "EDT": -4, "CST": -6, "CDT": -5, "MST": -7, "MDT": -6, "PST": -8, "PDT": -7}
                off = tz_map.get(tzs.upper())
                if off is not None:
                    return dt_naive.replace(tzinfo=timezone(timedelta(hours=off)))
                # Attach local timezone if unknown abbreviation
                try:
                    local_tz = datetime.now().astimezone().tzinfo or timezone.utc
                    return dt_naive.replace(tzinfo=local_tz)
                except Exception:
                    return dt_naive.replace(tzinfo=timezone.utc)
        except Exception:
            pass

        # Try several standard formats
        for fmt in (
            "%Y-%m-%d %H:%M:%S%z",
            "%Y-%m-%d %H:%M:%S",
            "%Y-%m-%d %H:%M%z",
            "%Y-%m-%d %H:%M",
            "%m/%d/%Y %I:%M%p",
            "%m/%d/%Y %I:%M %p",
            "%m/%d/%Y %H:%M",
            "%I:%M %p %m/%d/%Y",
            "%I:%M%p %m/%d/%Y",
            "%H:%M %m/%d/%Y",
        ):
            try:
                dt = datetime.strptime(s.replace("Z", "+0000"), fmt)
                return dt
            except Exception:
                pass
        return None

    def _check_scheduled_start(self) -> None:
        try:
            when = self.scheduled_start_at
            if not isinstance(when, datetime):
                return
            now = datetime.now(timezone.utc)
            if now >= when:
                # Fire once
                self.scheduled_start_at = None
                try:
                    self.stream.append_info("Scheduled time reached; starting Claude now")
                except Exception:
                    pass
                try:
                    asyncio.create_task(self.action_start())
                except Exception:
                    pass
                self._update_status()
        except Exception:
            pass

    async def _run_loop(self) -> None:
        try:
            while self.running:
                # Set the current displayed Claude command to what we are about to spawn
                self._current_claude_cmd = self.config.command
                saw_marker = False
                counted_iter = False
                last_decision: Optional[bool] = None  # None = unknown, True/False per latest flag
                self.stream.append_info(f"Starting: {self.config.command}  cwd={self.config.working_dir}")
                self._update_status()
                async for line in self.runner.run_once():
                    if self._log_file is not None and not self.config.log_path:
                        # Avoid double-writing when tee is active
                        self._log_file.write(f"{datetime.utcnow().isoformat()} {line}\n")
                    # Opportunistic flag detection directly from raw line content
                    try:
                        if CONTINUE_TRUE_PATTERN.search(line):
                            last_decision = True
                        elif CONTINUE_FALSE_PATTERN.search(line):
                            last_decision = False
                    except Exception:
                        pass
                    parsed = self.parser.parse_line(line)
                    if parsed["type"] == "marker":
                        # Exact marker guarantees continuation; count once and show banner
                        if not counted_iter:
                            self.continuations_count += 1
                            counted_iter = True
                            self.banner.show_banner(self.continuations_count)
                            self.stream.append_banner("===== CONTINUATION =====")
                            if hasattr(self, "state_panel"):
                                self.state_panel.refresh_panel()
                            self._update_status()
                        last_decision = True
                        saw_marker = True
                        continue
                    self.messages_count += 1
                    # Add as expandable item in the StreamList
                    self.stream.add_parsed(parsed, show_detail=self.detail, pretty=self.pretty)
                    self._update_status()

                # If we finished an iteration without marker, stop
                rc = self.runner.last_returncode
                self.stream.append_info(f"Process exited with return code {rc if rc is not None else 'n/a'}")
                if not self.running:
                    break
                # Decide continuation using precedence: explicit FALSE wins, else TRUE if seen (marker or regex)
                should_continue = False
                if last_decision is False:
                    should_continue = False
                elif last_decision is True or saw_marker:
                    should_continue = True
                else:
                    should_continue = False

                if should_continue:
                    # If we didn't already count (from explicit marker), count now and show banner
                    if not counted_iter:
                        self.continuations_count += 1
                        counted_iter = True
                        self.banner.show_banner(self.continuations_count)
                        self.stream.append_banner("===== CONTINUATION =====")
                        if hasattr(self, "state_panel"):
                            self.state_panel.refresh_panel()
                        self._update_status()
                    # Announce continuation once at decision point
                    try:
                        self.stream.append_info("CONTINUING - DECISION WAS CONTINUE-SOFTWARE-FACTORY=TRUE")
                    except Exception:
                        pass
                    # Set next command to direct claude continuation as requested
                    next_cmd = "claude --dangerously-skip-permissions -p '/continue-orchestrating' --verbose --output-format stream-json"
                    self.config.command = next_cmd
                    self._current_claude_cmd = next_cmd
                    self._update_status()
                    # Recreate runner to reset process state
                    self.runner = ClaudeRunner(self.config)
                    continue
                else:
                    # Announce stop if we explicitly saw FALSE
                    try:
                        if last_decision is False:
                            self.stream.append_info("STOPPING - CONTINUE-SOFTWARE-FACTORY=FALSE")
                    except Exception:
                        pass
                    break
        except Exception as exc:
            self.stream.write_line(f"[red]Error:[/] {exc}")
        finally:
            self.running = False

    def _refresh_state_panel_periodic(self) -> None:
        try:
            if hasattr(self, "state_panel"):
                self.state_panel.refresh_panel()
        except Exception:
            pass

    def _atexit_cleanup(self) -> None:
        # Synchronous cleanup for interpreter shutdown
        try:
            self.runner.kill_now()
        except Exception:
            pass
        try:
            if self._log_file is not None:
                self._log_file.close()
        except Exception:
            pass


def build_default_config(args: Optional[argparse.Namespace] = None) -> RunnerConfig:
    # Default to legacy runner
    # __file__ is now inside tools/, so legacy runner sits alongside us
    legacy_path = str(Path(__file__).with_name("claude_orchestrator_runner.py"))
    default_cmd = f"python3 {shlex.quote(legacy_path)}"
    cmd = args.command if args and getattr(args, "command", None) else default_cmd
    log_path = None
    if args and args.log_file:
        log_path = args.log_file
    # Wrapper removed; keep placeholder for compatibility
    wrapper_mode = "none"
    debug_raw = bool(getattr(args, "debug_raw", False))
    working_dir = args.cwd if args and getattr(args, "cwd", None) else os.getcwd()
    # Default SF2 project dir: prefer CWD of user, not tools/
    sf2_project_dir = getattr(args, "sf2_project_dir", working_dir)
    use_login_shell = bool(getattr(args, "shell_login", False))
    debug_verbose = bool(getattr(args, "debug_verbose", False))
    return RunnerConfig(
        command=cmd,
        env_params=os.environ.get("CLAUDE_PARAMS"),
        log_path=log_path,
        debug_raw=debug_raw,
        working_dir=working_dir,
        debug_verbose=debug_verbose,
        sf2_project_dir=sf2_project_dir,
    )


async def run_headless(config: RunnerConfig, view: str) -> int:
    """Run the orchestrator without TUI, streaming output to stdout.

    Args:
        config: Runner settings.
        view: One of 'raw', 'summary', 'detail', 'pretty'.

    Returns:
        Exit code of the final process.
    """
    runner = ClaudeRunner(config)
    parser = MessageParser()
    messages = 0
    continuations = 0
    logf = open(config.log_path, "a", buffering=1) if (config.log_path and False) else None
    try:
        while True:
            saw_marker = False
            print(f"[starting] {config.command}  cwd={config.working_dir}", flush=True)
            async for line in runner.run_once():
                # When log_path is set, stdout is already mirrored by tee; avoid double logging here
                if view == "raw":
                    print(line, flush=True)
                    if line.strip() == CONTINUE_MARKER:
                        continuations += 1
                        print(f"===== CONTINUATION #{continuations} =====", flush=True)
                        saw_marker = True
                    else:
                        messages += 1
                    continue
                parsed = parser.parse_line(line)
                if parsed["type"] == "marker":
                    continuations += 1
                    print(f"===== CONTINUATION #{continuations} =====", flush=True)
                    saw_marker = True
                    continue
                messages += 1
                if view == "summary":
                    print(parsed["summary"], flush=True)
                elif view == "detail":
                    print(parsed["summary"], flush=True)
                    print(parsed["detail_raw"], flush=True)
                elif view == "pretty":
                    print(parsed["summary"], flush=True)
                    print(parsed["detail_pretty"], flush=True)
                else:
                    print(line, flush=True)

            print(f"[exit] return code {runner.last_returncode}", flush=True)
            if not saw_marker:
                break
            # In simulation mode, avoid looping indefinitely on continuation markers
            if getattr(config, "simulate", False):
                break
            runner = ClaudeRunner(config)
        return runner.last_returncode or 0
    finally:
        if logf is not None:
            try:
                logf.close()
            except Exception:
                pass


def main() -> None:
    parser = argparse.ArgumentParser(description="Claude Orchestrator TUI (Textual)")
    parser.add_argument("--log-file", help="Path to write raw output log; omit to disable logging")
    parser.add_argument("--simulate", action="store_true", help="Run in simulation mode")
    parser.add_argument("--mode", choices=["tui", "headless"], default="tui", help="Run with Textual UI or headless streaming")
    parser.add_argument("--view", choices=["raw", "summary", "detail", "pretty"], default="raw", help="Headless output view")
    parser.add_argument("--debug-raw", action="store_true", help="Print raw chunk debug lines from stdout reader")
    parser.add_argument("--command", help="Override the exact command to run (bypass default claude invocation)")
    parser.add_argument("--cwd", help="Working directory to execute in (defaults to current)")
    parser.add_argument("--debug-verbose", action="store_true", help="Emit extra debug: env, which, pid, heartbeats")
    parser.add_argument("--sf2-project-dir", help="Path to software-factory project (contains orchestrator-state-v3.json)")
    args = parser.parse_args([]) if False else parser.parse_args()
    config = build_default_config(args)
    if args.simulate:
        config.simulate = True
    if args.mode == "headless":
        try:
            code = asyncio.run(run_headless(config, args.view))
        except KeyboardInterrupt:
            code = 130
        sys.exit(code)
    else:
        app = ClaudeOrchestratorApp(config)
        app.run()


if __name__ == "__main__":
    main()


