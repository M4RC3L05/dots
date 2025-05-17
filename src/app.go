package src

import (
	"github.com/fatih/color"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/core"
	"github.com/m4rc3l05/dots/src/displays"
)

type CmdFlagsArgs struct {
	Help             bool
	Version          bool
	Color            bool
	PrintEnvironment bool
}

type CmdArgs struct {
	Flags CmdFlagsArgs
	Rest  []string
}

type Args struct {
	CmdArgs          CmdArgs
	Version          string
	Logger           core.ILogger
	Homedir          string
	DotfilesFilesDir string
	Displays         displays.IDisplays
	Commands         commands.ICommands
}

func resolveForm(from []string, fallback string) string {
	if len(from) <= 1 {
		return fallback
	}

	return from[1]
}

func resolveCmd(rest []string) string {
	if len(rest) > 0 {
		return rest[0]
	}

	return ""
}

func App(args Args) (bool, error) {
	if args.CmdArgs.Flags.Help {
		args.Displays.Help()

		return true, nil
	}

	if args.CmdArgs.Flags.PrintEnvironment {
		args.Displays.Environment(args.Homedir, args.DotfilesFilesDir)

		return true, nil
	}

	if args.CmdArgs.Flags.Version {
		args.Displays.Version(args.Version)

		return true, nil
	}

	cmd := resolveCmd(args.CmdArgs.Rest)

	switch cmd {
	case "diff":
		{
			return args.Commands.Diff(commands.DiffArgs{
				FromDir: args.DotfilesFilesDir,
				ToDir:   args.Homedir,
			})
		}

	case "apply":
		{
			return args.Commands.Apply(commands.ApplyArgs{
				From: resolveForm(args.CmdArgs.Rest, args.DotfilesFilesDir),
				Extra: commands.ApplyArgsExtra{
					Homedir:          args.Homedir,
					DotfilesFilesDir: args.DotfilesFilesDir,
				},
			})
		}

	case "adopt":
		{
			return args.Commands.Adopt(commands.AdoptArgs{
				From: resolveForm(args.CmdArgs.Rest, args.DotfilesFilesDir),
				Extra: commands.AdoptArgsExtra{
					Homedir:          args.Homedir,
					DotfilesFilesDir: args.DotfilesFilesDir,
				},
			})
		}
	default:
		{
			args.Displays.Help()
			args.Logger.Warnnl("Command %s not found", color.MagentaString(cmd))
			return false, nil
		}
	}
}
