#!/usr/bin/env bash
set -e

echo "📦 Installing Byobu (tmux wrapper) ..."

# Update apt cache once per run
sudo apt-get update -y

# Install byobu and tmux
sudo apt-get install -y byobu tmux

# Ensure byobu uses tmux as backend and is enabled for current user (remoteUser)
byobu-select-backend tmux || true
byobu-ctrl-a emacs || true
mkdir -p "$HOME/.config/byobu"
# Enable on login shells
byobu-enable || true
# Add convenient aliases if not present
if ! grep -q "^alias bby=" ~/.bashrc 2>/dev/null; then
  echo "alias bby=\"byobu\"" >> ~/.bashrc
fi

# Idempotent Byobu/tmux config to work well inside VS Code + macOS
CFG="$HOME/.config/byobu"
mkdir -p "$CFG"

# 1) Disable dynamic palette to avoid OSC 10/11 garbage
cat > "$CFG/color" <<'EOF'
BACKGROUND=k
FOREGROUND=w
MONOCHROME=1
EOF

# 2) Mac/VS Code friendly keybindings: no Alt or F-keys; use tmux defaults
cat > "$CFG/keybindings.tmux" <<'EOF'
set -g prefix C-b
unbind-key -n C-a
unbind-key -n F2
unbind-key -n F3
unbind-key -n F4
unbind-key -n F6
unbind-key -n F7
unbind-key -n F8
unbind-key -n F9
unbind-key -n F11
unbind-key -n F12
# Ensure new window works with both c and Ctrl-c
bind-key -T prefix c   new-window -c "#{pane_current_path}"
bind-key -T prefix C-c new-window -c "#{pane_current_path}"
set -s escape-time 0
EOF

# 3) Tmux defaults: use zsh if available, correct quoting, good TERM
ZSH_PATH=$(command -v zsh || true)
if [ -n "$ZSH_PATH" ]; then
  DEFAULT_SHELL="$ZSH_PATH"
else
  DEFAULT_SHELL="/bin/bash"
fi
cat > "$CFG/.tmux.conf" <<EOF
set -g default-shell $DEFAULT_SHELL
set -g default-command "$DEFAULT_SHELL -l"
set -g default-terminal "tmux-256color"
set -s escape-time 0
EOF

# 4) Zsh guard: simplify prompt and disable title/VS Code OSC inside tmux/byobu
if [ -f "$HOME/.zshrc" ] && ! grep -q "# --- Byobu/tmux + VS Code fix ---" "$HOME/.zshrc"; then
  cat >> "$HOME/.zshrc" <<'EOF'
# --- Byobu/tmux + VS Code fix ---
if [ -n "$TMUX" ] || [ -n "$BYOBU_BACKEND" ]; then
  export DISABLE_AUTO_TITLE=true
  export VSCODE_SHELL_INTEGRATION=0
  preexec() { :; }
  precmd()  { :; }
  PROMPT='%n@%m %~ %# '
  RPROMPT=
fi
# --- end ---
EOF
fi

echo "✅ Byobu installation and basic setup completed. Use 'byobu' to start."


