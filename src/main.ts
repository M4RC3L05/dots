import { parseArgs } from "@std/cli";
import { normalize, resolve } from "@std/path";
import { main } from "./app.ts";
import { makeLogger } from "./core/logger.ts";
import {
  makeHandleExpectedError,
  makeHandleUnknownError,
  resolveDotfilesFilesDir,
  resolveHomeDir,
} from "./core/utils.ts";
import * as displays from "./displays/mod.ts";
import * as commands from "./commands/mod.ts";

const run = async () => {
  const args = parseArgs(Deno.args, {
    alias: { help: "h", version: "V", color: "c" },
    string: ["dotfilesFilesDir"],
    boolean: ["help", "version", "printEnvironment", "color"],
    default: { color: true },
  });

  const outputColor = Deno.noColor ? false : args.c ?? args.color;
  const logger = makeLogger({ color: outputColor });
  const handleExpectedError = makeHandleExpectedError(logger);
  const handleUnknownError = makeHandleUnknownError(logger);

  globalThis.addEventListener(
    "unhandledrejection",
    (e) => handleUnknownError(e.reason),
  );

  globalThis.addEventListener("error", (e) => handleUnknownError(e.error));

  const homedir = await resolveHomeDir();

  if (homedir instanceof Error) {
    await handleExpectedError(homedir);

    return;
  }

  const dotfilesFilesDir = await resolveDotfilesFilesDir({
    homedir: homedir,
    dotfilesFilesDirPath: args.dotfilesFilesDir ??
      resolve(normalize(`${homedir}/.dotfiles/home`)),
  });

  if (dotfilesFilesDir instanceof Error) {
    await handleExpectedError(dotfilesFilesDir);

    return;
  }

  const result = await main({
    commands,
    displays,
    cmdArgs: {
      rest: args._.filter((x) => typeof x === "string"),
      flags: {
        color: outputColor,
        help: args.help,
        printEnvironment: args.printEnvironment,
        version: args.V ?? args.version,
      },
    },
    logger,
    dotfilesFilesDir,
    homedir,
  });

  if (result instanceof Error) {
    await handleExpectedError(result);

    return;
  }

  if (result === false) Deno.exitCode = 1;
};

if (import.meta.main) {
  await run();
}
