import { normalize, resolve } from "@std/path";
import { parseArgs } from "@std/cli";
import * as cmd from "./cmd/mod.ts";
import meta from "./../deno.json" with { type: "json" };
import { exists } from "@std/fs";

const printHelp = () => {
  console.log(`
dots

Utility to manage your dotfiles, by keeping a given folder with a copy of the relevant dotfiles from user's home directory.
It allows you to adopt the dotfiles changes or override with local changes.

Usage: dots [OPTIONS] [COMMAND] [ARGS]

Options:
  --help, -h                  Display this help menu
  --version, -V               Display version
  --dotfilesFilesDir          Dotfiles files directory path to be used as the place where the ~/ will be mapped to.
                              This directory should be version controlled in order to keep an history of the changes.
                              It can also be controled with "DOTS_DOTFILES_FILES_DIR" env var, and takes precedence over the cmd flag.
                              It defaults to "~/.dotfiles/home".
  --printEnvironment          Prints homedir and dotfiles files dir values.

Commands:
  diff                        Diffs the user's dotfiles files with the ~/ files.
  adopt [PATH (optional)]     Adopts changes from ~/ files to user's dotfiles files.
                              A subpath of users home directory can be provided as an argument, in order to only apply part of the directories/files.
                              It can be a subdirectory or a file.
  apply [PATH (optional)]     Apply changes from user's dotfiles files to ~/ files.
                              A subpath of users dotfiles files directory can be provided as an argument, in order to only apply part of the directories/files.
                              It can be a subdirectory or a file.
  `.trim());
};

const printEnvironment = (
  { homedir, dotfilesFilesDir }: { homedir: string; dotfilesFilesDir: string },
) => {
  console.log(`
-----------------------
Environment:

HOME:               ${homedir}
DOTFILES FILES DIR: ${dotfilesFilesDir}
-----------------------
  `.trim());
};

const printVersion = () => {
  console.log(`v${meta.version}`);
};

const resolveHomeDir = async () => {
  const homeDirFromEnv = Deno.build.os === "windows"
    ? Deno.env.get("USERPROFILE")
    : Deno.env.get("HOME");

  if (!homeDirFromEnv) {
    console.log("Could not determine users home diretory");

    return false;
  }

  const homeDir = resolve(normalize(homeDirFromEnv));

  if (
    !homeDir ||
    !await exists(homeDir, { isDirectory: true, isReadable: true })
  ) {
    console.log(
      `Homedir "${homeDir}" does not exists or is not a directory or is not readable`,
    );

    return false;
  }

  return homeDir;
};

const resolveDotfilesFilesDir = async (
  { dotfilesFilesDirPath }: {
    homedir: string;
    dotfilesFilesDirPath?: string;
  },
) => {
  const dir = Deno.env.get("DOTS_DOTFILES_FILES_DIR") ?? dotfilesFilesDirPath;

  if (!dir) {
    console.log("Dotfiles files directory path not provided.");

    return false;
  }

  const dotfilesFilesDir = resolve(normalize(dir));

  if (
    !dotfilesFilesDir ||
    !await exists(dotfilesFilesDir, { isDirectory: true, isReadable: true })
  ) {
    console.log(
      `Dotfiles "${dotfilesFilesDir}" does not exists or is not a directory or is not readable`,
    );

    return false;
  }

  return dotfilesFilesDir;
};

if (import.meta.main) {
  globalThis.addEventListener("unhandledrejection", (e) => {
    console.log("Unexpected error occurred...");
    console.dir(Deno.inspect(e.reason, {
      colors: true,
      depth: 1000,
      compact: false,
      breakLength: 200,
    }));

    Deno.exit(1);
  });

  globalThis.addEventListener("error", (e) => {
    console.log("Unexpected error occurred...");
    console.log(
      Deno.inspect(e.error, {
        colors: true,
        depth: 1000,
        compact: false,
        breakLength: 200,
      }),
    );

    Deno.exit(1);
  });

  const homedir = await resolveHomeDir();

  if (!homedir) Deno.exit(1);

  const args = parseArgs(Deno.args, {
    alias: {
      help: "h",
      version: "V",
    },
    string: ["dotfilesFilesDir"],
    boolean: ["help", "version", "printEnvironment"],
    default: {
      dotfilesFilesDir: resolve(normalize(`${homedir}/.dotfiles/home`)),
    },
  });

  const dotfilesFilesDir = await resolveDotfilesFilesDir({
    homedir: homedir,
    dotfilesFilesDirPath: args.dotfilesFilesDir,
  });

  if (!dotfilesFilesDir) Deno.exit(1);

  if (args.printEnvironment) {
    printEnvironment({ homedir, dotfilesFilesDir });
  }

  switch (args._[0]) {
    case "diff": {
      const ok = await cmd.diff({ fromDir: dotfilesFilesDir, toDir: homedir });

      if (!ok) Deno.exit(1);

      break;
    }
    case "apply": {
      const ok = await cmd.apply({
        from: args._[1] as string,
        extra: { homedir, dotfilesFilesDir },
      });

      if (!ok) Deno.exit(1);

      break;
    }
    case "adopt": {
      const ok = await cmd.adopt({
        from: args._[1] as string,
        extra: { homedir, dotfilesFilesDir },
      });

      if (!ok) Deno.exit(1);

      break;
    }
    default: {
      if (args.h || args.help) {
        printHelp();
        Deno.exit(0);
      }

      if (args.version || args.V) {
        printVersion();
        Deno.exit(0);
      }

      printHelp();
      Deno.exit(0);
    }
  }
}
