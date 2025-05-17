package commands

import "github.com/m4rc3l05/dots/src/core"

type ICommands interface {
	Adopt(args AdoptArgs) (bool, error)
	Diff(args DiffArgs) (bool, error)
	Apply(args ApplyArgs) (bool, error)
}

type Commands struct {
	ICommands

	Logger core.ILogger
}
