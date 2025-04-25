import { copy, ensureDir } from "@std/fs";
import { dirname } from "@std/path";

export const recreateFile = async (from: string, to: string) => {
  await ensureDir(dirname(to));
  await copy(from, to, { overwrite: true });
};
