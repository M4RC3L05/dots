import { assertSpyCalls, type Spy, spy } from "@std/testing/mock";
import type { LogFn, Logger } from "./logger.ts";
import type { Displays } from "../displays/mod.ts";
import type { Commands } from "../commands/mod.ts";
import { assertArrayIncludes, assertEquals } from "@std/assert";

export type SpyNested<
  T,
  // deno-lint-ignore no-explicit-any
  E extends (...args: any) => any = (...args: any) => any,
> = {
  [K in keyof T]: T[K] extends E ? Spy<T[K], Parameters<T[K]>, ReturnType<T[K]>>
    : T[K] extends object ? SpyNested<T[K], E>
    : never;
};

export type SpyNestedCallNumber<
  T,
  // deno-lint-ignore no-explicit-any
  E extends (...args: any) => any = (...args: any) => any,
> = Partial<
  {
    [K in keyof T]: T[K] extends E ? number
      : T[K] extends object ? SpyNestedCallNumber<T[K], E>
      : never;
  }
>;

export type SpyLogger = SpyNested<Logger, LogFn>;

export const makeSpyLogger = (): SpyLogger => ({
  log: {
    debug: spy(),
    info: spy(),
    warn: spy(),
    error: spy(),
    log: spy(),
  },
  logNoNl: {
    debug: spy(),
    info: spy(),
    warn: spy(),
    error: spy(),
    log: spy(),
  },
});

type SpyLoggerCallNumber = SpyNestedCallNumber<Logger, LogFn>;

export const assertSpyLoggerCalls = (
  logger: SpyLogger,
  callNumber?: SpyLoggerCallNumber,
) => {
  assertSpyCalls(logger.log.debug, callNumber?.log?.debug ?? 0);
  assertSpyCalls(logger.log.error, callNumber?.log?.error ?? 0);
  assertSpyCalls(logger.log.info, callNumber?.log?.info ?? 0);
  assertSpyCalls(logger.log.log, callNumber?.log?.log ?? 0);
  assertSpyCalls(logger.log.warn, callNumber?.log?.warn ?? 0);
  assertSpyCalls(logger.logNoNl.debug, callNumber?.logNoNl?.debug ?? 0);
  assertSpyCalls(logger.logNoNl.error, callNumber?.logNoNl?.error ?? 0);
  assertSpyCalls(logger.logNoNl.info, callNumber?.logNoNl?.info ?? 0);
  assertSpyCalls(logger.logNoNl.log, callNumber?.logNoNl?.log ?? 0);
  assertSpyCalls(logger.logNoNl.warn, callNumber?.logNoNl?.warn ?? 0);
};

export type SpyDisplays = SpyNested<Displays>;

export const makeSpyDisplays = (): SpyDisplays => ({
  environment: spy(),
  help: spy(),
  version: spy(),
});

type SpyDisplaysCallNumber = SpyNestedCallNumber<Displays>;

export const assertSpyDisplaysCalls = (
  displays: SpyDisplays,
  callNumber?: SpyDisplaysCallNumber,
) => {
  assertSpyCalls(displays.environment, callNumber?.environment ?? 0);
  assertSpyCalls(displays.help, callNumber?.help ?? 0);
  assertSpyCalls(displays.version, callNumber?.version ?? 0);
};

export type SpyCommands = SpyNested<Commands>;

export const makeSpyCommands = (): SpyCommands => ({
  adopt: spy(),
  apply: spy(),
  diff: spy(),
});

type SpyCommandsCallNumber = SpyNestedCallNumber<Commands>;

export const assertSpyCommandsCalls = (
  commands: SpyCommands,
  callNumber?: SpyCommandsCallNumber,
) => {
  assertSpyCalls(commands.adopt, callNumber?.adopt ?? 0);
  assertSpyCalls(commands.apply, callNumber?.apply ?? 0);
  assertSpyCalls(commands.diff, callNumber?.diff ?? 0);
};

export const assertArrayIncludesRemoveMatch = <T>(
  actual: T[],
  expected: T[],
) => {
  assertEquals(
    actual.length,
    expected.length,
    `Arrays are not the same size; actual: ${actual.length}, expected: ${expected.length}`,
  );

  for (let i = expected.length - 1; i >= 0; i -= 1) {
    assertArrayIncludes(actual, [expected[i]]);
    actual.splice(actual.indexOf(expected[i]), 1);
    expected.splice(i, 1);
  }
};
