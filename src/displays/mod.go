package displays

import "github.com/m4rc3l05/dots/src/core"

type IDisplays interface {
	Environment(homedir string, dotfilesFilesDir string)
	Help()
	Version(version string)
}

type Displays struct {
	IDisplays

	Logger core.ILogger
}
