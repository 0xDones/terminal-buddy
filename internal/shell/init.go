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
`
