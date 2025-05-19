package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/m4rc3l05/dots/src"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/core"
	"github.com/m4rc3l05/dots/src/displays"
)

var Version string = "v2.0.0"

func handlePanic(logger core.ILogger) {
	if r := recover(); r != nil {
		logger.Errornl("Recovered from panic")

		if err, ok := r.(error); ok {
			core.LogErrors(logger, err, 0)
		}

		os.Exit(1)
	}
}

func main() {
	logger := core.MakeLogger()
	displays := displays.Displays{Logger: logger}
	commands := commands.Commands{Logger: logger}

	defer handlePanic(logger)

	homedir, err := core.ResolveHomeDir()
	if err != nil {
		core.LogErrors(logger, err, 0)

		os.Exit(1)
	}

	dotfilesFilesDirFallback, err := filepath.Abs(filepath.Join(homedir, ".dotfiles", "home"))
	if err != nil {
		core.LogErrors(logger, err, 0)

		os.Exit(1)
	}

	flag.CommandLine.SetOutput(os.Stdout)

	dotfilesFilesDirFlag := flag.String(
		"dotfilesFilesDir",
		dotfilesFilesDirFallback,
		"Set dotfiles files dir",
	)
	helpFlag := flag.Bool("help", false, "Show help menu")
	versionFlag := flag.Bool("version", false, "Show version")
	printEnvironmentFlag := flag.Bool("printEnv", false, "Show environment")
	colorFlag := flag.Bool("color", true, "Print with color")

	flag.Usage = func() {
		displays.Help()
	}

	flag.Parse()

	dotfilesFilesDir, err := core.ResolveDotfilesFilesDir(dotfilesFilesDirFlag)
	if err != nil {
		core.LogErrors(logger, err, 0)

		os.Exit(1)
	}

	color.NoColor = !*colorFlag

	ok, err := src.App(src.Args{
		CmdArgs: src.CmdArgs{
			Flags: src.CmdFlagsArgs{
				Help:             *helpFlag,
				Version:          *versionFlag,
				Color:            *colorFlag,
				PrintEnvironment: *printEnvironmentFlag,
			},
			Rest: flag.Args(),
		},
		Version:          Version,
		Logger:           logger,
		Homedir:          homedir,
		DotfilesFilesDir: dotfilesFilesDir,
		Displays:         displays,
		Commands:         commands,
	})
	if err != nil {
		core.LogErrors(logger, err, 0)
	}

	if !ok {
		os.Exit(1)
	}
}
