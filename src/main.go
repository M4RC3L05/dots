package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/core"
)

func main() {
	color.NoColor = false

	ok, err := commands.DiffCmd(commands.DiffCmdArgs{
		FromDir: "/home/main/.dotfiles/home",
		ToDir:   "/home/main",
		Extra: commands.DiffCmdArgsExtra{
			Logger: core.MakeLogger(),
		},
	})

	fmt.Printf("err: %v\n", err)
	fmt.Printf("ok: %v\n", ok)
}
