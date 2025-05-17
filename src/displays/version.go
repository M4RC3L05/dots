package displays

import (
	"github.com/m4rc3l05/dots/src/core"
)

func PrintVersion(version string, logger core.ILogger) {
	logger.Lognl("v%s", version)
}
