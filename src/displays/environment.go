package displays

import (
	"strings"

	"github.com/fatih/color"
)

func (d Displays) Environment(homedir string, dotfilesFilesDir string) {
	d.Logger.Lognl(strings.TrimSpace(`
-----------------------
Environment:

HOME:               %s
DOTFILES FILES DIR: %s
-----------------------
`), color.BlueString(homedir), color.BlueString(dotfilesFilesDir))
}
