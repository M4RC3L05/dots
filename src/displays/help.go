package displays

import (
	"strings"

	"github.com/fatih/color"
)

func (d Displays) Help() {
	d.Logger.Lognl(strings.TrimSpace(`
%s

Utility to manage your dotfiles, by keeping a given folder with a copy of the relevant dotfiles from user's home directory.
It allows you to adopt the dotfiles changes or override with local changes.

Usage: %s %s %s %s

%s:
  --help, -h                              Display this help menu

  --version                               Display version

  --dotfilesFilesDir <path>               Dotfiles files directory path to be used as the place where the ~/ will be mapped to.
                                          This directory should be version controlled in order to keep an history of the changes.
                                          It can also be controled with "DOTS_DOTFILES_FILES_DIR" env var, and takes precedence over the cmd flag.
                                          It defaults to "~/.dotfiles/home".

  --printEnv <true/false>                 Prints homedir and dotfiles files dir values.

  --color <true/false>                    Colors output. Enabled by default.

%s:
  diff                                    Diffs the user's dotfiles files with the ~/ files.

  adopt                                   Adopts changes from ~/ files to user's dotfiles files.
                                          A subpath of users home directory can be provided as an argument, in order to only apply part of the directories/files.
                                          It can be a subdirectory or a file.
    %s:
      path (optional)                     A path under the user's home directory to adopt from.

  apply                                   Apply changes from user's dotfiles files to ~/ files.
                                          A subpath of users dotfiles files directory can be provided as an argument, in order to only apply part of the directories/files.
                                          It can be a subdirectory or a file.
    %s:
      path (optional)                     A path under the user's dotfiles files directory.
`), color.MagentaString("dots"), color.MagentaString("dots"), color.GreenString("[OPTIONS]"), color.MagentaString("[COMMAND]"), color.YellowString("[ARGS]"), color.GreenString("Options"), color.MagentaString("Command"), color.YellowString("Args"), color.YellowString("Args"))
}
