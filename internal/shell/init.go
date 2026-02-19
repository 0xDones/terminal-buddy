package shell

import "fmt"

// InitScript returns the shell integration script for the given shell.
// Supported shells: "bash", "zsh", "fish".
func InitScript(shellName string) (string, error) {
	switch shellName {
	case "bash":
		return bashScript, nil
	case "zsh":
		return zshScript, nil
	case "fish":
		return fishScript, nil
	default:
		return "", fmt.Errorf("unsupported shell: %q (supported: bash, zsh, fish)", shellName)
	}
}

const bashScript = `# tb shell integration (bash)
__tb_widget() {
  local selected
  selected="$(command tb)"
  if [[ -n "$selected" ]]; then
    READLINE_LINE="$selected"
    READLINE_POINT=${#selected}
  fi
}
bind -x "\"${TB_KEYBINDING:-\\C-o}\": __tb_widget"

# Wrapper so that running "tb" directly also prefills the command line.
# Bash cannot set READLINE_LINE outside bind -x, so we add to history instead.
tb() {
  if [[ $# -gt 0 ]]; then
    command tb "$@"
    return
  fi
  local selected
  selected="$(command tb)"
  if [[ -n "$selected" ]]; then
    history -s "$selected"
    echo "$selected"
    echo "(copied to clipboard · press ↑ to edit/run)"
  fi
}
`

const zshScript = `# tb shell integration (zsh)
__tb_widget() {
  local selected
  selected="$(command tb)"
  if [[ -n "$selected" ]]; then
    BUFFER="$selected"
    CURSOR=${#selected}
  fi
  zle redisplay
}
zle -N __tb_widget
bindkey "${TB_KEYBINDING:-^O}" __tb_widget

# Wrapper so that running "tb" directly also prefills the command line.
# print -z pushes text onto the zsh line editor buffer.
tb() {
  if [[ $# -gt 0 ]]; then
    command tb "$@"
    return
  fi
  local selected
  selected="$(command tb)"
  if [[ -n "$selected" ]]; then
    print -z -- "$selected"
  fi
}
`

const fishScript = `# tb shell integration (fish)
function __tb_widget
  set -l selected (command tb)
  if test -n "$selected"
    commandline -r -- $selected
  end
  commandline -f repaint
end
bind (set -q TB_KEYBINDING; and echo $TB_KEYBINDING; or echo \co) __tb_widget

# Wrapper so that running "tb" directly also prefills the command line.
function tb --wraps=tb
  if count $argv > /dev/null
    command tb $argv
    return
  end
  set -l selected (command tb)
  if test -n "$selected"
    commandline -r -- $selected
    commandline -f repaint
  end
end
`
