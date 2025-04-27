import { magenta } from "@std/fmt/colors";
import type { Commands } from "./commands/mod.ts";
import type { Logger } from "./core/logger.ts";
import type { Displays } from "./displays/mod.ts";

type Args = {
  cmdArgs: {
    flags: {
      help: boolean;
      version: boolean;
      color: boolean;
      printEnvironment: boolean;
    };
    rest: string[];
  };
  logger: Logger;
  homedir: string;
  dotfilesFilesDir: string;
  displays: Displays;
  commands: Commands;
};

export const main = async (
  { cmdArgs, logger, homedir, dotfilesFilesDir, displays, commands }: Args,
) => {
  const { flags, rest: restArgs } = cmdArgs;
  if (flags.help) {
    await displays.help({ logger, color: flags.color });

    return;
  }

  if (flags.printEnvironment) {
    await displays.environment({
      homedir,
      dotfilesFilesDir,
      logger,
      color: flags.color,
    });

    return;
  }

  if (flags.version) {
    await displays.version({ logger });

    return;
  }

  switch (restArgs[0]) {
    case "diff": {
      return await commands.diff({
        fromDir: dotfilesFilesDir,
        toDir: homedir,
        extra: { diff: { color: flags.color }, logger, color: flags.color },
      });
    }
    case "apply": {
      return await commands.apply({
        from: restArgs[1],
        extra: { homedir, dotfilesFilesDir, logger, color: flags.color },
      });
    }
    case "adopt": {
      return await commands.adopt({
        from: restArgs[1],
        extra: { homedir, dotfilesFilesDir, logger, color: flags.color },
      });
    }
    default: {
      await displays.help({ logger, color: flags.color });
      await logger.log.warn(
        `Command ${
          flags.color ? magenta(restArgs[0] ?? '""') : restArgs[0] ?? '""'
        } not found`,
      );

      return;
    }
  }
};
