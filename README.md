# terminal-buddy

A terminal UI for organizing and quickly accessing your shell commands. Browse, search, and select commands from a categorized list — then paste them directly into your shell prompt.

## Install

Requires Go 1.25.6+.

```bash
go install tb@latest
```

Or build from source:

```bash
git clone <repo-url> && cd terminal-buddy
go build -o tb .
```

## Usage

Launch the TUI:

```bash
tb
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑`/`k` `↓`/`j` | Navigate commands |
| `Tab` / `Shift+Tab` | Switch category tabs |
| `/` | Search |
| `Enter` | Select command (copies to clipboard, prints to stdout, exits) |
| `c` | Copy highlighted command to clipboard (stays in TUI) |
| `n` | Create new command |
| `e` | Edit selected command |
| `d` | Delete selected command |
| `q` / `Ctrl+C` | Quit |

## Shell Integration

Add one line to your shell config to bind `Ctrl+O` — pressing it opens `tb`, and the selected command is placed directly into your prompt for review before executing.

### Bash

Add to `~/.bashrc`:

```bash
eval "$(tb init bash)"
```

### Zsh

Add to `~/.zshrc`:

```bash
eval "$(tb init zsh)"
```

### Fish

Add to `~/.config/fish/config.fish`:

```fish
tb init fish | source
```

### Custom Keybinding

The default binding is `Ctrl+O`. Override it with the `TB_KEYBINDING` environment variable:

```bash
# Bash/Zsh — bind to Ctrl+K instead
export TB_KEYBINDING='\C-k'   # bash
export TB_KEYBINDING='^K'     # zsh

# Fish
set -gx TB_KEYBINDING \ck
```

## Configuration

Commands are stored in `~/.tb.yaml`. On first run, a default file is created with example commands. You can edit it directly or manage commands through the TUI.

```yaml
commands:
  - name: flush-dns
    description: Flush macOS DNS cache
    command: sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder
    category: network
  - name: docker-prune
    description: Remove all stopped containers, dangling images, and unused networks
    command: docker system prune -af
    category: docker
  - name: git-undo
    description: Undo last commit keeping changes staged
    command: git reset --soft HEAD~1
    category: git
```

Each command has four fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Short identifier shown in the list |
| `command` | Yes | The shell command to run |
| `description` | No | Longer explanation shown in the detail pane |
| `category` | No | Used for tab-based filtering |

## Requirements

- **Clipboard** (optional): On Linux, `xclip` or `xsel` must be installed for clipboard support. macOS and Windows work out of the box. If unavailable, `tb` still prints to stdout — clipboard is best-effort.
