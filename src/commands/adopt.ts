import { exists, walk } from "@std/fs";
import { normalize, resolve } from "@std/path";
import { blue, green, red } from "@std/fmt/colors";
import { recreateFile } from "../core/utils.ts";
import type { Logger } from "../core/logger.ts";

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

  if (from !== dotfilesFilesDir) {
    if (!from.startsWith(homedir)) {
      return new Error(
        `Path ${color ? blue(from) : from} is not a subpath of ${
          color ? blue(homedir) : homedir
        }`,
      );
    }

    if (from.startsWith(dotfilesFilesDir)) {
      return new Error(
        `Path ${color ? blue(from) : from} can not be a subpath of ${
          color ? blue(dotfilesFilesDir) : dotfilesFilesDir
        }`,
      );
    }
  }

  if (await exists(from, { isFile: true })) {
    const to = resolve(normalize(from.replace(homedir, dotfilesFilesDir)));

    await logger?.logNoNl.log(
      `Adopting ${color ? blue(from) : from} to ${color ? blue(to) : to} ...`,
    );

    try {
      await recreateFile(from, to);

      await logger?.log.log(color ? green(" ✓") : " ✓");
    } catch (error) {
      await logger?.log.log(color ? red(" ✕") : " ✕");

      return new Error(`Error adopting ${color ? blue(from) : from}`, {
        cause: error,
      });
    }

    return true;
  }

  const errors: Error[] = [];

  for await (
    const file of walk(from, { includeDirs: false, includeSymlinks: false })
  ) {
    let origin: string;
    let destionation: string;

    if (file.path.startsWith(dotfilesFilesDir)) {
      origin = file.path.replace(dotfilesFilesDir, homedir);
      destionation = file.path;
    } else {
      origin = file.path;
      destionation = file.path.replace(homedir, dotfilesFilesDir);
    }

    await logger?.logNoNl.log(
      `Adopting ${color ? blue(origin) : origin} to ${
        color ? blue(destionation) : destionation
      } ...`,
    );

    try {
      await recreateFile(origin, destionation);

      await logger?.log.log(color ? green(" ✓") : " ✓");
    } catch (error) {
      await logger?.log.log(color ? red(" ✕") : " ✕");

      errors.push(
        new Error(`Error adopting ${color ? blue(origin) : origin}`, {
          cause: error,
        }),
      );
    }
  }

  if (errors.length > 0) {
    return new AggregateError(errors, "Error adopting directory");
  }

  return true;
};
