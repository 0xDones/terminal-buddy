package clipboard

import cb "github.com/atotto/clipboard"

// Write copies text to the system clipboard.
// Returns an error if the clipboard is unavailable (e.g., no xclip/xsel on Linux).
func Write(text string) error {
	return cb.WriteAll(text)
}
