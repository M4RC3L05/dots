import meta from "../../deno.json" with { type: "json" };
import type { Logger } from "./logger.ts";

export const printHelp = ({ logger }: { logger: Logger }) => {
  logger.log.log(`
dots

Utility to manage your dotfiles, by keeping a given folder with a copy of the relevant dotfiles from user's home directory.
It allows you to adopt the dotfiles changes or override with local changes.

Usage: dots [OPTIONS] [COMMAND] [ARGS]

Options:
  --help, -h                              Display this help menu

  --version, -V                           Display version

  --dotfilesFilesDir <path>               Dotfiles files directory path to be used as the place where the ~/ will be mapped to.
                                          This directory should be version controlled in order to keep an history of the changes.
                                          It can also be controled with "DOTS_DOTFILES_FILES_DIR" env var, and takes precedence over the cmd flag.
                                          It defaults to "~/.dotfiles/home".

  --printEnvironment <true/false>         Prints homedir and dotfiles files dir values.

  --color, -c <true/false>                Colors output. Enabled by default.

Commands:
  diff                                    Diffs the user's dotfiles files with the ~/ files.

  adopt                                   Adopts changes from ~/ files to user's dotfiles files.
                                          A subpath of users home directory can be provided as an argument, in order to only apply part of the directories/files.
                                          It can be a subdirectory or a file.
    Args:
      path (optional)                     A path under the user's home directory to adopt from.

  apply                                   Apply changes from user's dotfiles files to ~/ files.
                                          A subpath of users dotfiles files directory can be provided as an argument, in order to only apply part of the directories/files.
                                          It can be a subdirectory or a file.
    Args:
      path (optional)                     A path under the user's dotfiles files directory.
  `.trim());
};

export const printEnvironment = (
  { homedir, dotfilesFilesDir, logger }: {
    homedir: string;
    logger: Logger;
    dotfilesFilesDir: string;
  },
) => {
  logger.log.log(`
-----------------------
Environment:

HOME:               ${homedir}
DOTFILES FILES DIR: ${dotfilesFilesDir}
-----------------------
  `.trim());
};

export const printVersion = ({ logger }: { logger: Logger }) => {
  logger.log.log(`v${meta.version}`);
};
