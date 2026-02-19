# terminal-buddy

A terminal UI for organizing and quickly accessing your shell commands. Browse, search, and select commands from a categorized list — then paste them directly into your shell prompt.

## Install

### From GitHub Releases (recommended)

Download the latest binary for your platform from [Releases](https://github.com/0xDones/terminal-buddy/releases):

```bash
# Linux (amd64)
curl -Lo /usr/local/bin/tb https://github.com/0xDones/terminal-buddy/releases/latest/download/tb_linux_amd64
chmod +x /usr/local/bin/tb

# macOS (Apple Silicon)
curl -Lo /usr/local/bin/tb https://github.com/0xDones/terminal-buddy/releases/latest/download/tb_darwin_arm64
chmod +x /usr/local/bin/tb
```

### Build from Source

Requires Go 1.25.6+ and Make.

```bash
git clone https://github.com/0xDones/terminal-buddy.git && cd terminal-buddy
make build
```

The binary version is derived automatically from git tags (`git describe --tags`).

## Usage

Launch the TUI:

```bash
tb
```

Check version:

```bash
tb version
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑`/`k` `↓`/`j` | Navigate commands |
| `Tab` / `Shift+Tab` | Switch category tabs |
| `/` | Search |
| `Enter` | Select command (exits and prefills your prompt) |
| `c` | Copy highlighted command to clipboard (stays in TUI) |
| `n` | Create new command |
| `e` | Edit selected command |
| `d` | Delete selected command |
| `q` / `Ctrl+C` | Quit |

All keybindings are customizable — see [Configuration](#custom-keybindings).

## Shell Integration

**Shell integration is required for `tb` to work properly.** Without it, selected commands are only printed to stdout instead of being placed on your command line. Adding the init line below gives you two things:

1. **`Ctrl+O` keybind** — opens `tb` from anywhere in your prompt, selected command is inserted at the cursor.
2. **`tb` wrapper function** — typing `tb` directly also prefills your prompt with the selected command (zsh uses `print -z`, bash adds to history, fish uses `commandline`).

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

### Custom Keybindings

Override default key mappings by adding a `keybindings` section to `~/.tb.yaml`. Only the keys you want to change need to be specified — omitted keys keep their defaults.

```yaml
keybindings:
  up: ["up", "k"]
  down: ["down", "j"]
  next_tab: ["tab"]
  prev_tab: ["shift+tab"]
  search: ["/"]
  clear_esc: ["esc"]
  select: ["enter"]
  copy: ["c"]
  quit: ["q", "ctrl+c"]
  create: ["n"]
  edit: ["e"]
  delete: ["d"]
```

Keys use [BubbleTea key identifiers](https://pkg.go.dev/github.com/charmbracelet/bubbletea#KeyMsg): `"up"`, `"down"`, `"tab"`, `"shift+tab"`, `"enter"`, `"esc"`, `"ctrl+c"`, or any single character like `"k"`, `"/"`, `"q"`.

## Requirements

- **Shell integration** (strongly recommended): See [Shell Integration](#shell-integration). Without it, `tb` can only print the selected command to stdout.
- **Clipboard** (optional): On Linux, `xclip` or `xsel` must be installed for the `c` (copy) key to work. macOS and Windows work out of the box.
