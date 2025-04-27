import { copy, ensureDir, exists } from "@std/fs";
import { dirname, normalize, resolve } from "@std/path";
import type { Logger } from "./logger.ts";

export const recreateFile = async (from: string, to: string) => {
  await ensureDir(dirname(to));
  await copy(from, to, { overwrite: true });
};

export const resolveHomeDir = async () => {
  const homeDirFromEnv = Deno.build.os === "windows"
    ? Deno.env.get("USERPROFILE")
    : Deno.env.get("HOME");

  if (!homeDirFromEnv) {
    return new Error("Could not determine users home directory");
  }

  const homeDir = resolve(normalize(homeDirFromEnv));

  if (
    !homeDir ||
    !await exists(homeDir, { isDirectory: true, isReadable: true })
  ) {
    return new Error(
      `Homedir "${homeDir}" does not exists or is not a directory or is not readable`,
    );
  }

  return homeDir;
};

export const resolveDotfilesFilesDir = async (
  { dotfilesFilesDirPath }: {
    homedir: string;
    dotfilesFilesDirPath?: string;
  },
) => {
  const dir = Deno.env.get("DOTS_DOTFILES_FILES_DIR") ?? dotfilesFilesDirPath;

  if (!dir) {
    return new Error("Dotfiles files directory path not provided.");
  }

  const dotfilesFilesDir = resolve(normalize(dir));

  if (
    !dotfilesFilesDir ||
    !await exists(dotfilesFilesDir, { isDirectory: true, isReadable: true })
  ) {
    return Error(
      `Dotfiles "${dotfilesFilesDir}" does not exists or is not a directory or is not readable`,
    );
  }

  return dotfilesFilesDir;
};

export const makeHandleExpectedError =
  (logger: Logger) => async (error: Error | AggregateError) => {
    if (error instanceof AggregateError) {
      await logger.log.error(error.message);

      for (const err of error.errors) {
        if (!(err instanceof Error)) continue;

        await logger.log.log(`     > ${err.message}`);

        if (err.cause) {
          if (typeof err.cause === "string") {
            await logger.log.log(`       » ${err.cause}`);
          } else if (
            typeof err.cause === "object" &&
            "message" in err.cause &&
            typeof err.cause.message === "string"
          ) {
            await logger.log.log(`       » ${err.cause.message}`);
          }
        }
      }
    } else {
      await logger.log.error(error.message);

      if (error.cause) {
        if (typeof error.cause === "string") {
          await logger.log.log(`     » ${error.cause}`);
        } else if (
          typeof error.cause === "object" &&
          "message" in error.cause &&
          typeof error.cause.message === "string"
        ) {
          await logger.log.log(`     » ${error.cause.message}`);
        }
      }
    }

    Deno.exitCode = 1;
  };

export const makeHandleUnknownError =
  (logger: Logger) => async (error: unknown) => {
    await logger.log.error("Unexpected error occurred...");
    await logger.log.log(
      Deno.inspect(error, {
        colors: true,
        depth: 1000,
        compact: false,
        breakLength: 200,
      }),
    );

    Deno.exit(1);
  };
