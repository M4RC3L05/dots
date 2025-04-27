import { afterEach, beforeEach, describe, it } from "@std/testing/bdd";
import { assertSpyCall } from "@std/testing/mock";
import { assertEquals, assertInstanceOf } from "@std/assert";
import { copy, ensureDir } from "@std/fs";
import adopt from "./adopt.ts";
import {
  assertArrayIncludesRemoveMatch,
  assertSpyLoggerCalls,
  makeSpyLogger,
  type SpyLogger,
} from "../core/test-utils.ts";
import { dirname } from "@std/path/dirname";

let workingDir: string;
let logger: SpyLogger;

beforeEach(async () => {
  workingDir = await Deno.makeTempDir();
  logger = makeSpyLogger();
});

afterEach(async () => {
  await Deno.remove(workingDir, { recursive: true });
});

describe("adopt()", () => {
  it("should normalize `from` path", async () => {
    const result = await adopt({
      from: `${workingDir}/foo/bar/..`,
      extra: { homedir: workingDir, dotfilesFilesDir: workingDir },
    });

    assertEquals(
      (result as Error).message,
      `Path ${`${workingDir}/foo`} does not exists or is not readable`,
    );
  });

  it("should return an error if `from` does not exists", async () => {
    const result = await adopt({
      from: `${workingDir}/foo/`,
      extra: { homedir: workingDir, dotfilesFilesDir: workingDir },
    });

    assertInstanceOf(result, Error);
    assertEquals(
      result.message,
      `Path ${`${workingDir}/foo`} does not exists or is not readable`,
    );
  });

  it("should return an error if `from` is not readable", async () => {
    const from = await Deno.makeTempDir({ dir: workingDir });

    await Deno.chmod(from, 0o333);

    await using _ = {
      [Symbol.asyncDispose]: async () => {
        await Deno.chmod(from, 0o755);
      },
    };

    const result = await adopt({
      from,
      extra: { homedir: workingDir, dotfilesFilesDir: workingDir },
    });

    assertInstanceOf(result, Error);
    assertEquals(
      result.message,
      `Path ${from} does not exists or is not readable`,
    );
  });

  it("should return an error if `from` is not a subdirectory of `homedir`", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempDir({ dir: workingDir });

    const result = await adopt({
      from,
      extra: { homedir, dotfilesFilesDir },
    });

    assertInstanceOf(result, Error);
    assertEquals(
      result.message,
      `Path ${from} is not a subpath of ${homedir}`,
    );
  });

  it("should return an error if `from` is a subdirectory of `dotfilesFilesDir`", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: homedir });
    const from = await Deno.makeTempDir({ dir: dotfilesFilesDir });

    const result = await adopt({
      from,
      extra: { homedir, dotfilesFilesDir },
    });

    assertInstanceOf(result, Error);
    assertEquals(
      result.message,
      `Path ${from} can not be a subpath of ${dotfilesFilesDir}`,
    );
  });

  it("should adopt a single file", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: homedir });

    const result = await adopt({
      from,
      extra: { homedir, dotfilesFilesDir, logger },
    });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { logNoNl: { log: 1 }, log: { log: 1 } });
    assertSpyCall(logger.logNoNl.log, 0, {
      args: [
        `Adopting ${from} to ${from.replace(homedir, dotfilesFilesDir)} ...`,
      ],
    });
    assertSpyCall(logger.log.log, 0, { args: [" ✓"] });
  });

  it("should return an error if something appens while adopting file", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: homedir });
    const to = from.replace(homedir, dotfilesFilesDir);

    await using _ = await Deno.create(to);
    await Deno.chmod(to, 0o444);

    const result = await adopt({
      from,
      extra: { homedir, dotfilesFilesDir, logger },
    });

    assertInstanceOf(result, Error);
    assertEquals(result.message, `Error adopting ${from}`);
    assertInstanceOf(result.cause, Error);
    assertSpyLoggerCalls(logger, { logNoNl: { log: 1 }, log: { log: 1 } });
    assertSpyCall(logger.logNoNl.log, 0, {
      args: [
        `Adopting ${from} to ${to} ...`,
      ],
    });
    assertSpyCall(logger.log.log, 0, { args: [" ✕"] });
  });

  it("should fallback to `dotfilesFilesDir` if `from` is not provided", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });

    await Deno.chmod(dotfilesFilesDir, 0o333);

    await using _ = {
      [Symbol.asyncDispose]: async () => {
        await Deno.chmod(dotfilesFilesDir, 0o755);
      },
    };

    const result = await adopt({
      extra: { homedir, dotfilesFilesDir },
    });

    assertInstanceOf(result, Error);
    assertEquals(
      result.message,
      `Path ${dotfilesFilesDir} does not exists or is not readable`,
    );
  });

  it("should adopt a directory without providing `from`", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });
    const f1 = await Deno.makeTempFile({ dir: homedir, prefix: "1-" });
    const f2 = await Deno.makeTempFile({ dir: homedir, prefix: "2-" });

    await copy(f1, f1.replace(homedir, dotfilesFilesDir));
    await copy(f2, f2.replace(homedir, dotfilesFilesDir));

    const result = await adopt({
      extra: { homedir, dotfilesFilesDir, logger },
    });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { logNoNl: { log: 2 }, log: { log: 2 } });
    assertArrayIncludesRemoveMatch([
      ...logger.logNoNl.log.calls[0].args,
      ...logger.logNoNl.log.calls[1].args,
    ], [
      `Adopting ${f1} to ${f1.replace(homedir, dotfilesFilesDir)} ...`,
      `Adopting ${f2} to ${f2.replace(homedir, dotfilesFilesDir)} ...`,
    ]);
    assertArrayIncludesRemoveMatch([
      ...logger.log.log.calls[0].args,
      ...logger.log.log.calls[1].args,
    ], [
      " ✓",
      " ✓",
    ]);
  });

  it("should not adopt a directory without providing `from` if no files in `dotfilesFilesDir`", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });

    const result = await adopt({
      extra: { homedir, dotfilesFilesDir, logger },
    });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger);
  });

  it("should adopt a directory", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempDir({ dir: homedir });
    const f1 = await Deno.makeTempFile({ dir: from, prefix: "1-" });
    const f2 = await Deno.makeTempFile({ dir: from, prefix: "2-" });

    const result = await adopt({
      from,
      extra: { homedir, dotfilesFilesDir, logger },
    });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { logNoNl: { log: 2 }, log: { log: 2 } });
    assertArrayIncludesRemoveMatch([
      ...logger.logNoNl.log.calls[0].args,
      ...logger.logNoNl.log.calls[1].args,
    ], [
      `Adopting ${f1} to ${f1.replace(homedir, dotfilesFilesDir)} ...`,
      `Adopting ${f2} to ${f2.replace(homedir, dotfilesFilesDir)} ...`,
    ]);
    assertArrayIncludesRemoveMatch([
      ...logger.log.log.calls[0].args,
      ...logger.log.log.calls[1].args,
    ], [
      " ✓",
      " ✓",
    ]);
  });

  it("should return an error if something appens while adopting directory", async () => {
    const homedir = await Deno.makeTempDir({ dir: workingDir });
    const dotfilesFilesDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempDir({ dir: homedir });
    const f1 = await Deno.makeTempFile({ dir: from, prefix: "1-" });
    const f2 = await Deno.makeTempFile({ dir: from, prefix: "2-" });

    await ensureDir(dirname(f2.replace(homedir, dotfilesFilesDir)));
    await copy(f1, f2.replace(homedir, dotfilesFilesDir));
    await Deno.chmod(f2.replace(homedir, dotfilesFilesDir), 0o444);

    const result = await adopt({
      from,
      extra: { homedir, dotfilesFilesDir, logger },
    });

    assertInstanceOf(result, AggregateError);
    assertEquals(result.message, "Error adopting directory");
    assertEquals(result.errors.length, 1);
    assertInstanceOf(result.errors[0], Error);
    assertEquals(
      result.errors[0].message,
      `Error adopting ${f2}`,
    );
    assertInstanceOf(result.errors[0].cause, Error);
    assertSpyLoggerCalls(logger, { logNoNl: { log: 2 }, log: { log: 2 } });
    assertArrayIncludesRemoveMatch([
      ...logger.logNoNl.log.calls[0].args,
      ...logger.logNoNl.log.calls[1].args,
    ], [
      `Adopting ${f1} to ${f1.replace(homedir, dotfilesFilesDir)} ...`,
      `Adopting ${f2} to ${f2.replace(homedir, dotfilesFilesDir)} ...`,
    ]);
    assertArrayIncludesRemoveMatch([
      ...logger.log.log.calls[0].args,
      ...logger.log.log.calls[1].args,
    ], [
      " ✓",
      " ✕",
    ]);
  });
});
