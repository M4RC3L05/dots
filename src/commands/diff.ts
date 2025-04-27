import { diff } from "@libs/diff";
import { exists, walk } from "@std/fs";
import type { Logger } from "../core/logger.ts";
import { blue, green, red } from "@std/fmt/colors";

type Args = {
  fromDir: string;
  toDir: string;
  extra?: {
    logger?: Logger;
    diff?: {
      color?: boolean;
    };
    color?: boolean;
  };
};

export default async ({ fromDir, toDir, extra = {} }: Args) => {
  const { logger, diff: diffOptions, color } = extra;

  let hasFilesWithChanges = false;

  for await (
    const file of walk(fromDir, { includeDirs: false, includeSymlinks: false })
  ) {
    if (!file.isFile) continue;

    const from = file.path;
    const to = from.replace(fromDir, toDir);

    if (!await exists(from, { isFile: true, isReadable: true })) {
      await logger?.log.warn(
        `File ${color ? blue(from) : from} is not readable, skipping...`,
      );

      continue;
    }

    if (!await exists(to, { isFile: true, isReadable: true })) {
      await logger?.log.warn(
        `File ${
          color ? blue(to) : to
        } does not exists or is not a file or is not readable, skipping...`,
      );

      continue;
    }

    await logger?.logNoNl.log(
      `Diffing ${color ? blue(from) : from} against ${
        color ? blue(to) : to
      } ...`,
    );

    const [fc1, fc2] = await Promise.all([
      Deno.readTextFile(from),
      Deno.readTextFile(to),
    ]);

    const diffStr = diff(fc1, fc2, { colors: diffOptions?.color ?? false });

    if (diffStr.length > 0) {
      await logger?.log.log(color ? red(" ✕") : " ✕");
      await logger?.log.log(diffStr);

      hasFilesWithChanges = true;
    } else {
      await logger?.log.log(color ? green(" ✓") : " ✓");
    }
  }

  return !hasFilesWithChanges;
};
