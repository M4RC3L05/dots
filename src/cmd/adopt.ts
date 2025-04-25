import { exists, walk } from "@std/fs";
import { normalize, resolve } from "@std/path";
import { recreateFile } from "../utils.ts";

type Args = {
  from: string;
  extra: {
    homedir: string;
    dotfilesFilesDir: string;
  };
};

export default async (
  { from, extra: { homedir, dotfilesFilesDir } }: Args,
) => {
  const te = new TextEncoder();

  if (from) {
    from = resolve(normalize(from));
  }

  if (from && !await exists(from, { isReadable: true })) {
    console.log(`Path "${from}" does not exists or is not readable`);

    return false;
  }

  if (from && !from.startsWith(homedir)) {
    console.log(`Path "${from}" is not a subpath of "${homedir}"`);

    return false;
  }

  if (from && from.startsWith(dotfilesFilesDir)) {
    console.log(`Path "${from}" can not be a subpath of "${dotfilesFilesDir}"`);

    return false;
  }

  if (from && await exists(from, { isFile: true })) {
    const to = resolve(normalize(from.replace(homedir, dotfilesFilesDir)));

    await Deno.stdout.write(te.encode(`Adopting ${from} to ${to} ...`));

    try {
      await recreateFile(from, to);

      console.log(" ✓");
    } catch (error) {
      console.log(" ✕");

      throw error;
    }

    return true;
  }

  if (!from) {
    from = dotfilesFilesDir;
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

    if (!await exists(origin, { isFile: true, isReadable: true })) {
      console.log(
        `File "${origin}" does not exists or is not a file or is not readable, skipping...`,
      );

      continue;
    }

    await Deno.stdout.write(
      te.encode(`Adopting ${origin} to ${destionation} ...`),
    );

    try {
      await recreateFile(origin, destionation);

      console.log(" ✓");
    } catch (error) {
      console.log(" ✕");

      errors.push(error as Error);
    }
  }

  if (errors.length > 0) {
    throw new AggregateError(errors, "Error applying multiple files");
  }

  return true;
};
