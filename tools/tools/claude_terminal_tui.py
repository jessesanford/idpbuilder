#!/usr/bin/env python3
import argparse
import os
import pty
import fcntl
import termios
import tty
import selectors
import shlex
import signal
from typing import Optional

from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.containers import Horizontal
from textual.reactive import reactive
from textual.widgets import Header, Footer, Static
from rich.text import Text


# Reuse the rich state panel if available
try:
    from tools.claude_orchestrator_tui import StatePanel  # type: ignore
except Exception:
    StatePanel = None  # type: ignore

"""
Bittty/Textual-TTY integration

We prefer to embed the bittty-powered terminal if available (via textual-tty),
and fall back to our built-in PTY terminal otherwise. textual-tty is the
Textual widget frontend for the bittty emulator. See:
  - https://github.com/bitplane/bittty
"""
_BITTTY_AVAILABLE = False
_BitttyTerminalWidget = None
try:
    # Try common import paths for textual-tty's Terminal widget
    try:
        from textual_tty.terminal import Terminal as _BitttyTerminalWidget  # type: ignore
    except Exception:
        try:
            from textual_tty.widgets.terminal import Terminal as _BitttyTerminalWidget  # type: ignore
        except Exception:
            try:
                from textual_tty.widgets import Terminal as _BitttyTerminalWidget  # type: ignore
            except Exception:
                _BitttyTerminalWidget = None
    _BITTTY_AVAILABLE = _BitttyTerminalWidget is not None
except Exception:
    _BitttyTerminalWidget = None
    _BITTTY_AVAILABLE = False


class PtyTerminal(Static):
    """A PTY-backed terminal with a minimal ANSI/VT processor and fixed screen buffer.

    This avoids scroll by repainting a character grid (rows x cols), similar to a real terminal.
    """

    running = reactive(False)

    def __init__(self, command: str, *, cwd: Optional[str] = None) -> None:
        super().__init__(id="terminal")
        self.command = command
        self.cwd = cwd or os.getcwd()
        self._master_fd: Optional[int] = None
        self._pid: Optional[int] = None
        try:
            self.can_focus = True
        except Exception:
            pass
        # Screen buffer
        self.rows: int = 24
        self.cols: int = 80
        self.buf: list[list[str]] = [[" "] * self.cols for _ in range(self.rows)]
        self.cur_x: int = 0
        self.cur_y: int = 0
        self._esc: bool = False
        self._csi: str = ""

    def compose(self) -> ComposeResult:  # type: ignore[override]
        yield from ()

    def on_mount(self) -> None:  # type: ignore[override]
        self._resize_to_widget()
        self._spawn()
        self.set_interval(0.02, self._pump)  # periodic read
        # Re-evaluate size shortly after layout settles
        self.set_timer(0.2, self._resize_to_widget)
        self.set_timer(0.6, self._resize_to_widget)

    def on_resize(self) -> None:  # type: ignore[override]
        self._resize_to_widget()

    def _resize_to_widget(self) -> None:
        try:
            w = max(20, int(self.size.width))
            h = max(5, int(self.size.height))
        except Exception:
            w, h = 80, 24
        if w != self.cols or h != self.rows:
            self.cols = w
            self.rows = h
            self.buf = [[" "] * self.cols for _ in range(self.rows)]
            self.cur_x = 0
            self.cur_y = 0
            self._apply_winsize()
            self._render_screen()

    def _spawn(self) -> None:
        # Create a new PTY and fork+exec the command
        pid, master_fd = pty.fork()
        if pid == 0:
            # Child: exec the shell/command
            try:
                os.chdir(self.cwd)
            except Exception:
                pass
            argv = shlex.split(self.command) if self.command else ["bash", "-l"]
            try:
                os.execvp(argv[0], argv)
            except Exception as exc:
                # If exec fails, print error and exit
                os.write(1, f"exec failed: {exc}\r\n".encode())
                os._exit(1)
        else:
            # Parent: configure PTY as non-blocking
            self._pid = pid
            self._master_fd = master_fd
            fl = fcntl.fcntl(master_fd, fcntl.F_GETFL)
            fcntl.fcntl(master_fd, fcntl.F_SETFL, fl | os.O_NONBLOCK)
            self._apply_winsize()  # send initial size to child
            self.running = True
            self._write_status(f"Started: {self.command}  cwd={self.cwd}")

    def _pump(self) -> None:
        if not self.running or self._master_fd is None:
            return
        try:
            while True:
                try:
                    data = os.read(self._master_fd, 4096)
                except BlockingIOError:
                    break
                if not data:
                    self.running = False
                    break
                # Normalize newlines for TextLog, but keep most terminal output intact
                text = data.decode(errors="replace")
                self._process_input(text)
                self._render_screen()
        except Exception as exc:
            self._write_status(f"Terminal read error: {exc}")
            self.running = False

    def send(self, s: str) -> None:
        if self._master_fd is None:
            return
        try:
            os.write(self._master_fd, s.encode())
        except Exception as exc:
            self._write_status(f"send failed: {exc}")

    # Basic key forwarding (reserve global keys in App BINDINGS)
    def process_key(self, event) -> bool:
        """Process a Textual key event, forwarding to the PTY.

        Returns True if the key was handled (consumed).
        """
        key = getattr(event, "key", "")
        if key in {"ctrl+p", "ctrl+q"}:
            return False
        # Control keys (Ctrl+<letter>)
        if key.startswith("ctrl+") and len(key) == 6:
            letter = key[-1]
            code = ord(letter.upper()) - 64
            if 1 <= code <= 26:
                self.send(chr(code))
                return True
        # Printable
        ch = getattr(event, "character", None)
        if ch:
            self.send(ch)
            return True
        # Navigation and function keys
        translate = {
            "enter": "\r",
            "return": "\r",
            "tab": "\t",
            "backspace": "\x7f",
            "escape": "\x1b",
            "left": "\x1b[D",
            "right": "\x1b[C",
            "up": "\x1b[A",
            "down": "\x1b[B",
            "home": "\x1bOH",
            "end": "\x1bOF",
            "pageup": "\x1b[5~",
            "pagedown": "\x1b[6~",
            "delete": "\x1b[3~",
            "insert": "\x1b[2~",
            "f1": "\x1bOP",
            "f2": "\x1bOQ",
            "f3": "\x1bOR",
            "f4": "\x1bOS",
            "f5": "\x1b[15~",
            "f6": "\x1b[17~",
            "f7": "\x1b[18~",
            "f8": "\x1b[19~",
            "f9": "\x1b[20~",
            "f10": "\x1b[21~",
            "f11": "\x1b[23~",
            "f12": "\x1b[24~",
        }
        seq = translate.get(key)
        if seq:
            self.send(seq)
            return True
        return False

    def on_key(self, event) -> None:  # type: ignore[override]
        # Ensure focus and forward
        try:
            if not self.has_focus:
                self.focus()
        except Exception:
            pass
        if self.process_key(event):
            event.stop()
    # Mouse tracking (basic SGR mode) when app enables it via DECSET 1006
    def on_mouse_down(self, event) -> None:  # type: ignore[override]
        if getattr(self, "_mouse_sgr", False):
            self._send_mouse_sgr(event.x, event.y, 0, press=True)
            event.stop()

    def on_mouse_up(self, event) -> None:  # type: ignore[override]
        if getattr(self, "_mouse_sgr", False):
            self._send_mouse_sgr(event.x, event.y, 0, press=False)
            event.stop()

    def on_mouse_move(self, event) -> None:  # type: ignore[override]
        if getattr(self, "_mouse_sgr", False):
            self._send_mouse_sgr(event.x, event.y, 32, press=False)  # motion flag
            event.stop()

    # --- Minimal VT processing with basic SGR color handling ---
    def _process_input(self, s: str) -> None:
        i = 0
        while i < len(s):
            ch = s[i]
            if self._esc:
                if ch == "[":
                    self._csi = ""
                    i += 1
                    # collect until alpha
                    while i < len(s) and not s[i].isalpha():
                        self._csi += s[i]
                        i += 1
                    if i < len(s):
                        cmd = s[i]
                        self._handle_csi(self._csi, cmd)
                        i += 1
                    else:
                        # incomplete
                        break
                else:
                    # unsupported ESC sequence
                    i += 1
                self._esc = False
                continue
            if ch == "\x1b":
                self._esc = True
                i += 1
            elif ch == "\r":
                self.cur_x = 0
                i += 1
            elif ch == "\n":
                self.cur_y += 1
                if self.cur_y >= self.rows:
                    self._scroll_up()
                    self.cur_y = self.rows - 1
                i += 1
            elif ch == "\b":
                self.cur_x = max(0, self.cur_x - 1)
                i += 1
            else:
                # printable
                if 32 <= ord(ch) <= 126 or ch == "\t":
                    if ch == "\t":
                        spaces = 8 - (self.cur_x % 8)
                        for _ in range(spaces):
                            self._put_char(" ")
                    else:
                        self._put_char(ch)
                i += 1

    def _put_char(self, ch: str) -> None:
        if 0 <= self.cur_y < self.rows and 0 <= self.cur_x < self.cols:
            # Store (char, style) for colored rendering
            self.buf[self.cur_y][self.cur_x] = (ch, getattr(self, "_style_str", ""))
        self.cur_x += 1
        if self.cur_x >= self.cols:
            self.cur_x = 0
            self.cur_y += 1
            if self.cur_y >= self.rows:
                self._scroll_up()
                self.cur_y = self.rows - 1

    def _scroll_up(self) -> None:
        self.buf.pop(0)
        self.buf.append([(" ", getattr(self, "_style_str", "")) for _ in range(self.cols)])

    def _handle_csi(self, params: str, cmd: str) -> None:
        # Parse params as semicolon-separated ints (1-based defaults)
        def ints(default: int = 1) -> list[int]:
            if not params:
                return [default]
            try:
                return [int(p) if p else default for p in params.split(";")]
            except Exception:
                return [default]
        if cmd in ("H", "f"):
            vals = ints()
            y = (vals[0] - 1) if len(vals) >= 1 else 0
            x = (vals[1] - 1) if len(vals) >= 2 else 0
            self.cur_y = max(0, min(self.rows - 1, y))
            self.cur_x = max(0, min(self.cols - 1, x))
        elif cmd == "J":
            vals = ints(0)
            mode = vals[0]
            if mode == 2:  # clear screen
                self.buf = [[" "] * self.cols for _ in range(self.rows)]
                self.cur_x = self.cur_y = 0
            elif mode == 0:  # clear to end of screen
                for y in range(self.cur_y, self.rows):
                    start = self.cur_x if y == self.cur_y else 0
                    for x in range(start, self.cols):
                        self.buf[y][x] = " "
            elif mode == 1:  # clear from start of screen to cursor
                for y in range(0, self.cur_y + 1):
                    end = self.cur_x if y == self.cur_y else self.cols - 1
                    for x in range(0, end + 1):
                        self.buf[y][x] = " "
            elif mode == 3:  # erase saved lines (treat as full clear)
                self.buf = [[" "] * self.cols for _ in range(self.rows)]
                self.cur_x = self.cur_y = 0
        elif cmd == "K":
            vals = ints(0)
            mode = vals[0]
            if mode == 0:  # clear to end of line
                for x in range(self.cur_x, self.cols):
                    self.buf[self.cur_y][x] = " "
            elif mode == 2:  # clear entire line
                for x in range(0, self.cols):
                    self.buf[self.cur_y][x] = " "
        elif cmd == "m":
            # SGR - basic styling
            self._handle_sgr(params)
        elif cmd == "h":  # DEC private mode set (if starts with ?)
            if params.startswith("?"):
                self._handle_decset(params[1:], True)
        elif cmd == "l":  # DEC private mode reset
            if params.startswith("?"):
                self._handle_decset(params[1:], False)
        elif cmd == "s":  # save cursor
            self._saved = (self.cur_x, self.cur_y)
        elif cmd == "u":  # restore cursor
            try:
                self.cur_x, self.cur_y = self._saved
            except Exception:
                pass
        else:
            # Unsupported; ignore
            pass

    def _render_screen(self) -> None:
        # Build rich Text per line with styles
        out_lines: list[Text] = []
        for row in self.buf:
            t = Text()
            last_style = None
            seg_chars = []
            for cell in row:
                if isinstance(cell, tuple):
                    ch, style = cell
                else:
                    ch, style = cell, ""
                if style != last_style and seg_chars:
                    t.append("".join(seg_chars), style=last_style)
                    seg_chars = []
                last_style = style
                seg_chars.append(ch)
            if seg_chars:
                t.append("".join(seg_chars), style=last_style)
            # Some versions of Rich return None from rstrip(); normalize
            try:
                _ = t.rstrip()
            except Exception:
                pass
            out_lines.append(t)
        try:
            self.update(Text("\n").join(out_lines))
        except Exception:
            # Fallback: render plain text if join fails
            self.update("\n".join([lt.plain if hasattr(lt, "plain") else str(lt) for lt in out_lines]))

    def _write_status(self, msg: str) -> None:
        # Write a status message at the bottom line
        s = f"[ {msg} ]"
        self.cur_y = self.rows - 1
        self.cur_x = 0
        style = getattr(self, "_style_str", "")
        for x in range(self.cols):
            ch = s[x] if x < len(s) else " "
            self.buf[self.cur_y][x] = (ch, style)
        self._render_screen()

    # --- SGR helper ---
    def _handle_sgr(self, params: str) -> None:
        # Maintain current style string like "bold red on black"
        if not hasattr(self, "_sgr"):
            self._sgr = {"fg": None, "bg": None, "bold": False}
        if not params:
            params = "0"
        parts = params.split(";")
        i = 0
        while i < len(parts):
            p = parts[i] or "0"
            try:
                code = int(p)
            except ValueError:
                i += 1
                continue
            if code == 0:
                self._sgr = {"fg": None, "bg": None, "bold": False}
                i += 1
                continue
            if code == 1:
                self._sgr["bold"] = True
                i += 1
                continue
            if code == 22:
                self._sgr["bold"] = False
                i += 1
                continue
            if code == 39:
                self._sgr["fg"] = None
                i += 1
                continue
            if code == 49:
                self._sgr["bg"] = None
                i += 1
                continue
            # 8/16 colors
            if 30 <= code <= 37:
                self._sgr["fg"] = self._ansi_basic_color(code - 30, bright=False)
                i += 1
                continue
            if 40 <= code <= 47:
                self._sgr["bg"] = self._ansi_basic_color(code - 40, bright=False)
                i += 1
                continue
            if 90 <= code <= 97:
                self._sgr["fg"] = self._ansi_basic_color(code - 90, bright=True)
                i += 1
                continue
            if 100 <= code <= 107:
                self._sgr["bg"] = self._ansi_basic_color(code - 100, bright=True)
                i += 1
                continue
            # 256-color or truecolor: 38/48 ; 5;n or 2;r;g;b
            if code in (38, 48):
                is_fg = (code == 38)
                mode = parts[i+1] if i + 1 < len(parts) else ""
                if mode == "5" and i + 2 < len(parts):
                    try:
                        n = int(parts[i+2])
                        color = self._ansi_256_color(n)
                        if is_fg:
                            self._sgr["fg"] = color
                        else:
                            self._sgr["bg"] = color
                    except Exception:
                        pass
                    i += 3
                    continue
                elif mode == "2" and i + 4 < len(parts):
                    try:
                        r = int(parts[i+2]); g = int(parts[i+3]); b = int(parts[i+4])
                        color = f"#{r:02x}{g:02x}{b:02x}"
                        if is_fg:
                            self._sgr["fg"] = color
                        else:
                            self._sgr["bg"] = color
                    except Exception:
                        pass
                    i += 5
                    continue
                else:
                    i += 1
                    continue
            i += 1
        # Compose rich style string
        parts = []
        if self._sgr["bold"]:
            parts.append("bold")
        if self._sgr["fg"]:
            parts.append(self._sgr["fg"])  # e.g., "red"
        if self._sgr["bg"]:
            parts.append(f"on {self._sgr['bg']}")
        self._style_str = " ".join(parts)

    @staticmethod
    def _ansi_basic_color(idx: int, bright: bool) -> str:
        names = ["black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"]
        base = names[idx % 8]
        return f"bright_{base}" if bright else base

    @staticmethod
    def _ansi_256_color(n: int) -> str:
        # Map 256-color index to rich color string; use hex for 6x6x6 cube and grayscale
        n = max(0, min(255, n))
        if n < 16:
            # basic palette via names
            names = [
                "#000000","#800000","#008000","#808000","#000080","#800080","#008080","#c0c0c0",
                "#808080","#ff0000","#00ff00","#ffff00","#0000ff","#ff00ff","#00ffff","#ffffff",
            ]
            return names[n]
        if 16 <= n <= 231:
            n -= 16
            r = (n // 36) % 6
            g = (n // 6) % 6
            b = n % 6
            def comp(v: int) -> int:
                return 55 + v * 40 if v > 0 else 0
            return f"#{comp(r):02x}{comp(g):02x}{comp(b):02x}"
        # grayscale 232–255
        level = 8 + (n - 232) * 10
        level = max(0, min(255, level))
        return f"#{level:02x}{level:02x}{level:02x}"

    def _handle_decset(self, params: str, enable: bool) -> None:
        # Handle mouse tracking modes: 1002 (button), 1003 (any), 1006 (SGR)
        for p in params.split(";"):
            if p in ("1002", "1003", "1006"):
                if p == "1006":
                    self._mouse_sgr = enable
                # Note: we treat 1002/1003 as aliases; only SGR encoding is sent
                self._mouse_enabled = enable
            elif p in ("47", "1047", "1048", "1049"):
                # Alternate screen buffer on/off -> emulate by clearing screen and homing
                self.buf = [[" "] * self.cols for _ in range(self.rows)]
                self.cur_x = self.cur_y = 0

    def _send_mouse_sgr(self, x: int, y: int, flags: int, press: bool) -> None:
        # Translate widget coords to terminal cells (1-based)
        try:
            col = max(1, min(self.cols, x))
            row = max(1, min(self.rows, y))
        except Exception:
            return
        btn = flags
        # 0 press, 1 release in SGR; encode as CSI < btn ; col ; row M/m
        ch = "M" if press else "m"
        seq = f"\x1b[<" + str(btn) + ";" + str(col) + ";" + str(row) + ch
        self.send(seq)

    def _apply_winsize(self) -> None:
        # Push current rows/cols to the child PTY so TUIs format correctly
        if self._master_fd is None:
            return
        try:
            import struct
            winsz = struct.pack("HHHH", self.rows, self.cols, 0, 0)
            fcntl.ioctl(self._master_fd, termios.TIOCSWINSZ, winsz)
            if self._pid:
                try:
                    os.kill(self._pid, signal.SIGWINCH)
                except Exception:
                    pass
        except Exception:
            pass


class BitttyTerminal(Static):
    """Wrapper for a bittty-powered terminal via textual-tty, if available.

    Exposes a send(text) method similar to PtyTerminal so the app logic can
    remain the same across backends.
    """

    def __init__(self, *, cwd: Optional[str], command: Optional[str]) -> None:
        if not _BITTTY_AVAILABLE or _BitttyTerminalWidget is None:
            raise RuntimeError("textual-tty (bittty) is not installed")
        super().__init__(id="terminal")
        self._cwd = cwd or os.getcwd()
        self._command = command or ""
        self._term = _BitttyTerminalWidget()  # type: ignore
        try:
            self.can_focus = True
        except Exception:
            pass

    def compose(self) -> ComposeResult:  # type: ignore[override]
        yield self._term

    def on_mount(self) -> None:  # type: ignore[override]
        try:
            if self._cwd and self._cwd != os.getcwd():
                self.send(f"cd {shlex.quote(self._cwd)}\n")
        except Exception:
            pass
        if self._command:
            try:
                self.set_timer(0.3, lambda: self.send(self._command + "\n"))
            except Exception:
                try:
                    self.send(self._command + "\n")
                except Exception:
                    pass
        try:
            self.focus()
        except Exception:
            pass

    def send(self, s: str) -> None:
        term = getattr(self, "_term", None)
        if term is None:
            return
        for attr in ("send", "send_text", "write", "write_text", "input"):
            fn = getattr(term, attr, None)
            if callable(fn):
                try:
                    fn(s)
                    return
                except Exception:
                    continue
        try:
            for ch in s:
                fn = getattr(term, "send", None)
                if callable(fn):
                    fn(ch)
        except Exception:
            pass

    def process_key(self, event) -> bool:  # type: ignore[override]
        return False


class PexpectTerminal(Static):
    """Pexpect-backed terminal that renders line output and forwards keystrokes.

    This is a simplified backend intended to mitigate duplication issues some
    TUIs exhibit under embedded terminals by reading the child's stdout as text
    lines rather than a raw PTY framebuffer. It's not a full terminal emulator
    and is best for simple interactive programs.
    """

    running = reactive(False)

    def __init__(self, command: str, *, cwd: Optional[str] = None) -> None:
        super().__init__(id="terminal")
        self.command = command or "bash"
        self.cwd = cwd or os.getcwd()
        self._buf: list[str] = []
        self._child = None
        try:
            self.can_focus = True
        except Exception:
            pass

    def compose(self) -> ComposeResult:  # type: ignore[override]
        yield from ()

    def on_mount(self) -> None:  # type: ignore[override]
        try:
            import pexpect
        except Exception as exc:
            self.update(f"pexpect not installed: {exc}")
            return
        try:
            self._child = pexpect.spawn("/bin/bash", ["-lc", self.command], cwd=self.cwd, encoding="utf-8", timeout=None)
            self.running = True
            # Poll output periodically; pexpect's read_nonblocking prevents blocking UI
            self.set_interval(0.05, self._pump)
            try:
                self.focus()
            except Exception:
                pass
        except Exception as exc:
            self.update(f"Failed to start: {exc}")

    def _pump(self) -> None:
        if not self.running or self._child is None:
            return
        try:
            try:
                chunk = self._child.read_nonblocking(size=4096, timeout=0)
            except Exception:
                chunk = ""
            if chunk:
                self._append_text(chunk)
            if not self._child.isalive():
                self.running = False
        except Exception:
            self.running = False

    def _append_text(self, text: str) -> None:
        # Keep a capped buffer of lines and render
        self._buf.extend(text.splitlines(True))
        if len(self._buf) > 8000:
            self._buf = self._buf[-8000:]
        try:
            self.update("".join(self._buf))
        except Exception:
            pass

    def send(self, s: str) -> None:
        try:
            if self._child is not None:
                self._child.send(s)
        except Exception:
            pass

    def process_key(self, event) -> bool:  # type: ignore[override]
        key = getattr(event, "key", "")
        if key in {"ctrl+p", "ctrl+q"}:
            return False
        ch = getattr(event, "character", None)
        if ch:
            self.send(ch)
            return True
        translate = {
            "enter": "\r",
            "return": "\r",
            "tab": "\t",
            "backspace": "\x7f",
        }
        seq = translate.get(key)
        if seq:
            self.send(seq)
            return True
        return False

class PyteTerminal(Static):
    """Pyte-based terminal emulator rendered to Textual using Rich text.

    This provides a lightweight terminal emulation. It spawns a PTY child
    process and feeds output through pyte's Screen+Stream to maintain a text
    framebuffer which we render line-by-line.
    """

    running = reactive(False)

    def __init__(self, command: str, *, cwd: Optional[str] = None) -> None:
        super().__init__(id="terminal")
        self.command = command or "bash -l"
        self.cwd = cwd or os.getcwd()
        self._master_fd: Optional[int] = None
        self._pid: Optional[int] = None
        self._cols = 80
        self._rows = 24
        self._screen = None
        self._stream = None
        try:
            self.can_focus = True
        except Exception:
            pass

    def compose(self) -> ComposeResult:  # type: ignore[override]
        yield from ()

    def on_mount(self) -> None:  # type: ignore[override]
        try:
            import pyte
        except Exception as exc:
            self.update(f"pyte not installed: {exc}")
            return
        # Initialize screen/stream
        self._screen = pyte.Screen(self._cols, self._rows)
        self._screen.set_mode(pyte.modes.LNM)
        self._stream = pyte.Stream(self._screen)
        # Spawn child in PTY
        pid, master_fd = pty.fork()
        if pid == 0:
            try:
                os.chdir(self.cwd)
            except Exception:
                pass
            argv = shlex.split(self.command)
            try:
                os.execvp(argv[0], argv)
            except Exception as exc:
                os.write(1, f"exec failed: {exc}\r\n".encode())
                os._exit(1)
        else:
            self._pid = pid
            self._master_fd = master_fd
            fl = fcntl.fcntl(master_fd, fcntl.F_GETFL)
            fcntl.fcntl(master_fd, fcntl.F_SETFL, fl | os.O_NONBLOCK)
            self._apply_winsize()
            self.running = True
            self.set_interval(0.02, self._pump)
            try:
                self.focus()
            except Exception:
                pass

    def on_resize(self) -> None:  # type: ignore[override]
        try:
            w = max(20, int(self.size.width))
            h = max(5, int(self.size.height))
        except Exception:
            w, h = 80, 24
        if w != self._cols or h != self._rows:
            self._cols = w
            self._rows = h
            self._apply_winsize()

    def _apply_winsize(self) -> None:
        if self._master_fd is None:
            return
        try:
            import struct
            winsz = struct.pack("HHHH", self._rows, self._cols, 0, 0)
            fcntl.ioctl(self._master_fd, termios.TIOCSWINSZ, winsz)
            if self._pid:
                try:
                    os.kill(self._pid, signal.SIGWINCH)
                except Exception:
                    pass
        except Exception:
            pass

    def _pump(self) -> None:
        if not self.running or self._master_fd is None or self._stream is None or self._screen is None:
            return
        try:
            try:
                data = os.read(self._master_fd, 4096)
            except BlockingIOError:
                data = b""
            if not data:
                # Check if process still alive
                try:
                    if os.waitpid(self._pid or 0, os.WNOHANG)[0] != 0:
                        self.running = False
                except Exception:
                    pass
                return
            text = data.decode(errors="replace")
            try:
                self._stream.feed(text)
            except Exception:
                pass
            # Render pyte screen to plain text lines
            lines = []
            for y in range(self._screen.lines):  # type: ignore[attr-defined]
                row = self._screen.buffer.get(y)  # type: ignore[attr-defined]
                if row is None:
                    lines.append("")
                    continue
                # row is a dict of x->Char; reconstruct
                max_x = self._cols
                chars = []
                for x in range(max_x):
                    cell = row.get(x)
                    ch = getattr(cell, 'data', ' ') if cell else ' '
                    chars.append(ch)
                lines.append("".join(chars).rstrip())
            try:
                self.update("\n".join(lines))
            except Exception:
                pass
        except Exception:
            self.running = False

    def send(self, s: str) -> None:
        if self._master_fd is None:
            return
        try:
            os.write(self._master_fd, s.encode())
        except Exception:
            pass

    def process_key(self, event) -> bool:  # type: ignore[override]
        key = getattr(event, "key", "")
        if key in {"ctrl+p", "ctrl+q"}:
            return False
        # Control keys (Ctrl+<letter>)
        if key.startswith("ctrl+") and len(key) == 6:
            letter = key[-1]
            code = ord(letter.upper()) - 64
            if 1 <= code <= 26:
                self.send(chr(code))
                return True
        # Printable
        ch = getattr(event, "character", None)
        if ch:
            self.send(ch)
            return True
        # Navigation and function keys
        translate = {
            "enter": "\r",
            "return": "\r",
            "tab": "\t",
            "backspace": "\x7f",
            "escape": "\x1b",
            "left": "\x1b[D",
            "right": "\x1b[C",
            "up": "\x1b[A",
            "down": "\x1b[B",
            "home": "\x1bOH",
            "end": "\x1bOF",
            "pageup": "\x1b[5~",
            "pagedown": "\x1b[6~",
            "delete": "\x1b[3~",
            "insert": "\x1b[2~",
            "f1": "\x1bOP",
            "f2": "\x1bOQ",
            "f3": "\x1bOR",
            "f4": "\x1bOS",
            "f5": "\x1b[15~",
            "f6": "\x1b[17~",
            "f7": "\x1b[18~",
            "f8": "\x1b[19~",
            "f9": "\x1b[20~",
            "f10": "\x1b[21~",
            "f11": "\x1b[23~",
            "f12": "\x1b[24~",
        }
        seq = translate.get(key)
        if seq:
            self.send(seq)
            return True
        return False

    def on_key(self, event) -> None:  # type: ignore[override]
        # Ensure focus and forward
        try:
            if not self.has_focus:
                self.focus()
        except Exception:
            pass
        if self.process_key(event):
            event.stop()


class ClaudeTerminalApp(App):
    CSS = """
    Screen { layout: vertical; }
    #top { height: 1; }
    #panes { height: 1fr; }
    #terminal { height: 1fr; width: 1fr; overflow: hidden; }
    #state { width: 80; border: round green; padding: 0 1; height: 1fr; overflow: auto; }
    """

    BINDINGS = [
        Binding("ctrl+p", "command_palette", "Command Palette", show=True),
        Binding("ctrl+q", "quit", "Quit", show=True),
    ]

    def __init__(self, command: str, continue_cmd: str, sf2_project_dir: Optional[str], *, backend: str = "auto", auto_continue: bool = True, wrap_claude_script: bool = True, wrap_tmux: bool = True, tmux_safe_ui: bool = False) -> None:
        super().__init__()
        self._command = command
        self._continue_cmd = continue_cmd
        self._sf2_dir = sf2_project_dir or os.getcwd()
        self._backend = backend  # "auto" | "bittty" | "pty"
        self._auto_continue = bool(auto_continue)
        self._wrap_claude_script = bool(wrap_claude_script)
        self._wrap_tmux = bool(wrap_tmux)
        self._tmux_safe_ui = bool(tmux_safe_ui)
        # Debounce flags to prevent duplicate sends
        self._sent_start_cmd = False
        self._sent_continue_cmd = False

    def compose(self) -> ComposeResult:  # type: ignore[override]
        # Avoid per-second header updates when using bittty/textual-tty backend,
        # which can cause excessive re-layout and duplicate painting.
        show_clock = not (getattr(self, "_backend", "auto") in ("bittty", "textual-tty"))
        yield Header(show_clock=show_clock)
        # Auto-wrap 'claude' commands before choosing backend so it applies everywhere
        try:
            first_token = (self._command or "").strip().split()[0]
        except Exception:
            first_token = ""
        if first_token == "claude":
            try:
                orig = self._command
                if self._wrap_tmux:
                    # Launch tmux with UTF-8/256color and set inner exports before running Claude
                    inner = f"export LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8; export NCURSES_NO_UTF8_ACS=1; {orig}"
                    if self._tmux_safe_ui:
                        setup = "tmux -u -2 new-session -Ad -s sf; " \
                                "tmux -u -2 set -g -q assume-paste-time 0; " \
                                "tmux -u -2 set -g default-terminal screen-256color; " \
                                "tmux -u -2 set -g -a terminal-overrides ',*:smcup@:rmcup@'"
                        full = setup + "; tmux -u -2 new -As sf bash -lc " + shlex.quote(inner)
                        self._command = f"bash -lc {shlex.quote(full)}"
                    else:
                        tmux_cmd = f"tmux -u -2 new -As sf bash -lc {shlex.quote(inner)}"
                        self._command = tmux_cmd
                elif self._wrap_claude_script:
                    wrapped = f"script -qfec {shlex.quote(orig)} /dev/null"
                    self._command = f"bash -lc {shlex.quote(wrapped + '; exec bash')}"
            except Exception:
                pass

        with Horizontal(id="panes"):
            backend = self._backend
            if backend == "auto":
                backend = "bittty" if _BITTTY_AVAILABLE else "pty"
                # Heuristic for Claude: prefer bittty when tmux wrapping, else PTY
                try:
                    first = (self._command or "").strip().split()[0]
                except Exception:
                    first = ""
                if first == "claude":
                    if self._wrap_tmux and _BITTTY_AVAILABLE:
                        backend = "bittty"
                    else:
                        backend = "pty"
            if backend == "pexpect":
                self.term = PexpectTerminal(self._command)
                self._backend = "pexpect"
            elif backend == "pyte":
                self.term = PyteTerminal(self._command)
                self._backend = "pyte"
            elif backend in ("bittty", "textual-tty") and _BITTTY_AVAILABLE:
                try:
                    # Embed textual-tty Terminal directly to avoid wrapper repaint issues
                    self.term = _BitttyTerminalWidget()  # type: ignore
                    # Keep a note of selected backend
                    self._backend = "bittty"
                except Exception:
                    self.term = PtyTerminal(self._command)
                    self._backend = "pty"
            else:
                self.term = PtyTerminal(self._command)
                self._backend = "pty"
            yield self.term
            if StatePanel is not None:
                self.state_panel = StatePanel(self._sf2_dir)
                self.state_panel.id = "state"
                yield self.state_panel
        yield Footer()

    def on_mount(self) -> None:  # type: ignore[override]
        # Helper to send text into either backend
        def _send(text: str) -> None:
            t = getattr(self, "term", None)
            if t is None:
                return
            # PtyTerminal exposes send directly
            fn = getattr(t, "send", None)
            if callable(fn):
                try:
                    fn(text)
                    return
                except Exception:
                    pass
            # textual-tty Terminal may use different method names
            for name in ("send", "send_text", "write", "write_text", "input"):
                fn = getattr(t, name, None)
                if callable(fn):
                    try:
                        fn(text)
                        return
                    except Exception:
                        continue
        # If using bittty backend, run command in shell and then send continue
        def _maybe_send(text: str, flag_attr: str) -> None:
            try:
                if getattr(self, flag_attr):
                    return
            except Exception:
                pass
            _send(text)
            try:
                setattr(self, flag_attr, True)
            except Exception:
                pass

        if self._backend in ("bittty", "textual-tty", "pexpect", "pyte"):
            # Start the requested command once
            self.set_timer(0.3, lambda: _maybe_send(self._command + "\n", "_sent_start_cmd"))
            # Optionally send orchestrator continue once
            if self._auto_continue:
                self.set_timer(1.0, lambda: _maybe_send(self._continue_cmd + "\n", "_sent_continue_cmd"))
        else:
            # PTY backend already execs the command; optionally send continue once
            if self._auto_continue:
                self.set_timer(1.0, lambda: _maybe_send(self._continue_cmd + "\n", "_sent_continue_cmd"))
        # Periodically refresh the state panel if present
        if StatePanel is not None:
            self.set_interval(30.0, lambda: self.state_panel.refresh_panel())
        # Ensure terminal has focus so it receives keyboard input
        try:
            self.term.focus()
        except Exception:
            pass
        # Forward app-level keys to terminal unless reserved
        # Keep the periodic tick only for PTY backend to reduce repaint churn
        if self._backend != "bittty":
            self.set_interval(0.01, lambda: None)

    def on_key(self, event) -> None:  # type: ignore[override]
        # Global key handler: forward only when terminal is NOT focused (fallback)
        key = getattr(event, "key", "")
        if key in {"ctrl+p", "ctrl+q"}:
            return
        try:
            if getattr(self, "focused", None) is not self.term:
                if hasattr(self, "term") and hasattr(self.term, "process_key") and callable(getattr(self.term, "process_key")):
                    if self.term.process_key(event):
                        event.stop()
                        return
        except Exception:
            pass


def main() -> None:
    parser = argparse.ArgumentParser(description="Claude Terminal TUI (Textual)")
    parser.add_argument("--command", help="Command to run in the terminal", default=os.environ.get("CLAUDE_CMD", "claude"))
    parser.add_argument("--continue-cmd", help="Command to send after start", default="/continue-orchestrating")
    parser.add_argument("--sf2-project-dir", help="Software Factory project directory for the state panel")
    parser.add_argument("--backend", choices=["auto", "bittty", "textual-tty", "pty", "pexpect", "pyte"], default="auto", help="Terminal backend to use")
    parser.add_argument("--no-auto-continue", action="store_true", help="Disable auto-sending the continue command after startup")
    parser.add_argument("--no-script-wrap", action="store_true", help="Disable wrapping 'claude' in script -qfec to reduce duplicate paints")
    parser.add_argument("--no-wrap-tmux", action="store_true", help="Disable wrapping 'claude' in tmux (-u -2 new -As sf)")
    parser.add_argument("--tmux-safe-ui", action="store_true", help="Harden tmux UI: disable alt-screen, set screen-256color, UTF-8, and ACS fallback")
    args = parser.parse_args()
    app = ClaudeTerminalApp(
        args.command,
        args.continue_cmd,
        args.sf2_project_dir,
        backend=args.backend,
        auto_continue=(not args.no_auto_continue),
        wrap_claude_script=(not args.no_script_wrap),
        wrap_tmux=(not args.no_wrap_tmux),
        tmux_safe_ui=bool(args.tmux_safe_ui),
    )
    app.run()


if __name__ == "__main__":
    main()


