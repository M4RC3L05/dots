package displays

import (
	"strings"

	"github.com/fatih/color"
	"github.com/m4rc3l05/dots/src/core"
)

func PrintEnvironment(
	homedir string,
	dotfilesFilesDir string,
	logger core.ILogger,
) {
	logger.Lognl(strings.TrimSpace(`
-----------------------
Environment:

HOME:               %s
DOTFILES FILES DIR: %s
-----------------------
`), color.BlueString(homedir), color.BlueString(dotfilesFilesDir))
}
