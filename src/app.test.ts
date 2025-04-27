import { afterEach, beforeEach, describe, it } from "@std/testing/bdd";
import { assertSpyCall } from "@std/testing/mock";
import {
  assertSpyCommandsCalls,
  assertSpyDisplaysCalls,
  assertSpyLoggerCalls,
  makeSpyCommands,
  makeSpyDisplays,
  makeSpyLogger,
  type SpyCommands,
  type SpyDisplays,
  type SpyLogger,
} from "./core/test-utils.ts";
import { main } from "./app.ts";
import { deepMerge } from "@std/collections";

type DeepPartial<T> = {
  [K in keyof T]?: T[K] extends object ? DeepPartial<T[K]> : T[K];
};

let workingDir: string;
let logger: SpyLogger;
let displays: SpyDisplays;
let commands: SpyCommands;

beforeEach(async () => {
  workingDir = await Deno.makeTempDir();
  logger = makeSpyLogger();
  displays = makeSpyDisplays();
  commands = makeSpyCommands();
});

afterEach(async () => {
  await Deno.remove(workingDir, { recursive: true });
});

const makeMainArgs = (
  args?: DeepPartial<Parameters<typeof main>[0]>,
) =>
  (deepMerge({
    dotfilesFilesDir: workingDir,
    homedir: workingDir,
    logger,
    commands,
    displays,
    cmdArgs: {
      rest: [],
      flags: {
        printEnvironment: false,
        color: false,
        help: false,
        version: false,
      },
    },
  }, args ?? {})) as Parameters<typeof main>[0];

describe("main()", () => {
  it("should print help if `help` flag is provided", async () => {
    await main(
      makeMainArgs({ cmdArgs: { flags: { help: true } } }),
    );

    assertSpyDisplaysCalls(displays, { help: 1 });
    assertSpyCommandsCalls(commands);
  });

  it("should print environment if `printEnvironment` flag is provided", async () => {
    await main(
      makeMainArgs({ cmdArgs: { flags: { printEnvironment: true } } }),
    );

    assertSpyDisplaysCalls(displays, { environment: 1 });
    assertSpyCommandsCalls(commands);
  });

  it("should print version if `version` flag is provided", async () => {
    await main(
      makeMainArgs({ cmdArgs: { flags: { version: true } } }),
    );

    assertSpyDisplaysCalls(displays, { version: 1 });
    assertSpyCommandsCalls(commands);
  });

  it("should run diff if `diff` command provided", async () => {
    await main(
      makeMainArgs({
        homedir: "foo",
        dotfilesFilesDir: "bar",
        cmdArgs: { rest: ["diff"] },
      }),
    );

    assertSpyDisplaysCalls(displays);
    assertSpyCommandsCalls(commands, { diff: 1 });
    assertSpyCall(commands.diff, 0, {
      args: [{
        fromDir: "bar",
        toDir: "foo",
        extra: {
          diff: { color: false },
          logger,
          color: false,
        },
      }],
    });
  });

  it("should run apply if `apply` command provided", async () => {
    await main(
      makeMainArgs({ cmdArgs: { rest: ["apply"] } }),
    );

    assertSpyDisplaysCalls(displays);
    assertSpyCommandsCalls(commands, { apply: 1 });
    assertSpyCall(commands.apply, 0, {
      args: [{
        from: undefined,
        extra: {
          homedir: workingDir,
          dotfilesFilesDir: workingDir,
          logger,
          color: false,
        },
      }],
    });
  });

  it("should run apply with first arg if `apply` command provided with an arg", async () => {
    await main(
      makeMainArgs({ cmdArgs: { rest: ["apply", "foo"] } }),
    );

    assertSpyDisplaysCalls(displays);
    assertSpyCommandsCalls(commands, { apply: 1 });
    assertSpyCall(commands.apply, 0, {
      args: [{
        from: "foo",
        extra: {
          homedir: workingDir,
          dotfilesFilesDir: workingDir,
          logger,
          color: false,
        },
      }],
    });
  });

  it("should run adopt if `adopt` command provided", async () => {
    await main(
      makeMainArgs({ cmdArgs: { rest: ["adopt"] } }),
    );

    assertSpyDisplaysCalls(displays);
    assertSpyCommandsCalls(commands, { adopt: 1 });
    assertSpyCall(commands.adopt, 0, {
      args: [{
        from: undefined,
        extra: {
          homedir: workingDir,
          dotfilesFilesDir: workingDir,
          logger,
          color: false,
        },
      }],
    });
  });

  it("should run adopt with first arg if `adopt` command provided with an arg", async () => {
    await main(
      makeMainArgs({ cmdArgs: { rest: ["adopt", "foo"] } }),
    );

    assertSpyDisplaysCalls(displays);
    assertSpyCommandsCalls(commands, { adopt: 1 });
    assertSpyCall(commands.adopt, 0, {
      args: [{
        from: "foo",
        extra: {
          homedir: workingDir,
          dotfilesFilesDir: workingDir,
          logger,
          color: false,
        },
      }],
    });
  });

  it("should print help if command is not supported", async () => {
    await main(
      makeMainArgs({ cmdArgs: { rest: ["foo"] } }),
    );

    assertSpyCommandsCalls(commands);
    assertSpyDisplaysCalls(displays, { help: 1 });
    assertSpyLoggerCalls(logger, { log: { warn: 1 } });
    assertSpyCall(logger.log.warn, 0, { args: ["Command foo not found"] });
  });
});
