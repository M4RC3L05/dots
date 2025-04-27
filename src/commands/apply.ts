import { exists, walk } from "@std/fs";
import { normalize, resolve } from "@std/path";
import { recreateFile } from "../core/utils.ts";
import type { Logger } from "../core/logger.ts";
import { blue, green, red } from "@std/fmt/colors";

type Args = {
  from?: string;
  extra: {
    homedir: string;
    dotfilesFilesDir: string;
    logger?: Logger;
    color?: boolean;
  };
};

export default async (
  { from, extra: { homedir, dotfilesFilesDir, logger, color } }: Args,
) => {
  if (from) {
    from = resolve(normalize(from));
  } else {
    from = dotfilesFilesDir;
  }

  if (!await exists(from, { isReadable: true })) {
    return new Error(
      `Path ${color ? blue(from) : from} does not exists or is not readable`,
    );
  }

  if (!from.includes(dotfilesFilesDir)) {
    return new Error(
      `Path ${color ? blue(from) : from} is not a subpath of ${
        color ? blue(dotfilesFilesDir) : dotfilesFilesDir
      }`,
    );
  }

  if (await exists(from, { isFile: true })) {
    const to = resolve(normalize(from.replace(dotfilesFilesDir, homedir)));

    await logger?.logNoNl?.log(
      `Applying ${color ? blue(from) : from} to ${color ? blue(to) : to} ...`,
    );

    try {
      await recreateFile(from, to);

      await logger?.log?.log(color ? green(" ✓") : " ✓");
    } catch (error) {
      await logger?.log?.log(color ? red(" ✕") : " ✕");

      return new Error(
        `Error applying ${color ? blue(from) : from} to ${
          color ? blue(to) : to
        }`,
        { cause: error },
      );
    }

    return true;
  }

  const errors: Error[] = [];

  for await (
    const file of walk(from, { includeDirs: false, includeSymlinks: false })
  ) {
    const to = resolve(normalize(file.path.replace(dotfilesFilesDir, homedir)));

    await logger?.logNoNl.log(
      `Applying ${color ? blue(file.path) : file.path} to ${
        color ? blue(to) : to
      } ...`,
    );

    try {
      await recreateFile(file.path, to);

      await logger?.log?.log(color ? green(" ✓") : " ✓");
    } catch (error) {
      await logger?.log?.log(color ? red(" ✕") : " ✕");

      errors.push(
        new Error(
          `Error applying ${color ? blue(file.path) : file.path} to ${
            color ? blue(to) : to
          }`,
          { cause: error },
        ),
      );
    }
  }

  if (errors.length > 0) {
    return new AggregateError(errors, "Error applying directory");
  }

  return true;
};
