import { blue } from "@std/fmt/colors";
import type { Logger } from "./../core/logger.ts";

export default async (
  { homedir, dotfilesFilesDir, logger, color }: {
    homedir: string;
    logger: Logger;
    dotfilesFilesDir: string;
    color?: boolean;
  },
) => {
  await logger.log.log(`
-----------------------
Environment:

HOME:               ${color ? blue(homedir) : homedir}
DOTFILES FILES DIR: ${color ? blue(dotfilesFilesDir) : dotfilesFilesDir}
-----------------------
  `.trim());
};
