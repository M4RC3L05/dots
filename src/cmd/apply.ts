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
  } else {
    from = dotfilesFilesDir;
  }

  if (!await exists(from, { isReadable: true })) {
    console.log(`Path "${from}" does not exists or is not readable`);

    return false;
  }

  if (!from.includes(dotfilesFilesDir)) {
    console.log(`Path "${from}" is not a subpath of "${dotfilesFilesDir}"`);

    return false;
  }

  if (await exists(from, { isFile: true })) {
    const to = resolve(normalize(from.replace(dotfilesFilesDir, homedir)));

    await Deno.stdout.write(te.encode(`Applying ${from} to ${to} ...`));

    try {
      await recreateFile(from, to);

      console.log(" ✓");
    } catch (error) {
      console.log(" ✕");

      throw error;
    }

    return true;
  }

  const errors: Error[] = [];

  for await (
    const file of walk(from, { includeDirs: false, includeSymlinks: false })
  ) {
    const to = resolve(normalize(file.path.replace(dotfilesFilesDir, homedir)));

    await Deno.stdout.write(te.encode(`Applying ${file.path} to ${to} ...`));

    try {
      await recreateFile(file.path, to);

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
