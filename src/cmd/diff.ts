import { diff } from "@libs/diff";
import { exists, walk } from "@std/fs";
import { normalize, resolve } from "@std/path";

type Args = {
  fromDir: string;
  toDir: string;
};

export default async ({ fromDir, toDir }: Args) => {
  const te = new TextEncoder();

  fromDir = resolve(normalize(fromDir));
  toDir = resolve(normalize(toDir));

  const [fromDirExists, toDirExists] = await Promise.all([
    exists(fromDir, { isDirectory: true, isReadable: true }),
    exists(toDir, { isDirectory: true, isReadable: true }),
  ]);

  if (!fromDirExists) {
    console.log(
      `Path "${
        fromDir ?? ""
      }" does not exists or is not a directory or is not readable`,
    );

    return false;
  }

  if (!toDirExists) {
    console.log(
      `Path "${toDir}" does not exists or is not a directory or is not readable`,
    );

    return false;
  }

  let hasFilesWithChanges = false;

  for await (
    const file of walk(fromDir, { includeDirs: false, includeSymlinks: false })
  ) {
    if (!file.isFile) continue;

    const from = file.path;
    const to = from.replace(fromDir, toDir);

    if (!await exists(to, { isFile: true, isReadable: true })) {
      console.log(
        `File "${from}" does not exists or is not a file or is not readable on target "${to}", skipping...`,
      );

      continue;
    }

    await Deno.stdout.write(te.encode(`Diffing "${from}" against "${to}" ...`));

    const [fc1, fc2] = await Promise.all([
      Deno.readTextFile(from),
      Deno.readTextFile(to),
    ]);

    const diffStr = diff(fc1, fc2, { colors: true });

    if (diffStr.length > 0) {
      console.log(" ✕");
      console.log(diffStr);

      hasFilesWithChanges = true;
    } else {
      console.log(" ✓");
    }
  }

  return !hasFilesWithChanges;
};
