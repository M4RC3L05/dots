import { afterEach, beforeEach, describe, it } from "@std/testing/bdd";
import { assertSpyCall } from "@std/testing/mock";
import { assertEquals } from "@std/assert";
import { copy } from "@std/fs";
import diff from "./diff.ts";
import {
  assertArrayIncludesRemoveMatch,
  assertSpyLoggerCalls,
  makeSpyLogger,
  type SpyLogger,
} from "../core/test-utils.ts";

let workingDir: string;
let logger: SpyLogger;

beforeEach(async () => {
  workingDir = await Deno.makeTempDir();
  logger = makeSpyLogger();
});

afterEach(async () => {
  await Deno.remove(workingDir, { recursive: true });
});

describe("diff()", () => {
  it("should not diff anything if no files exists on `fromDir`", async () => {
    const fromDir = await Deno.makeTempDir({ dir: workingDir });
    const toDir = await Deno.makeTempDir({ dir: workingDir });

    const result = await diff({ fromDir, toDir, extra: { logger } });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger);
  });

  it("should not diff anything if no readable files exists on `fromDir`", async () => {
    const fromDir = await Deno.makeTempDir({ dir: workingDir });
    const toDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: fromDir });

    await Deno.chmod(from, 0o300);

    const result = await diff({ fromDir, toDir, extra: { logger } });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { log: { warn: 1 } });
    assertSpyCall(logger.log.warn, 0, {
      args: [
        `File ${from} is not readable, skipping...`,
      ],
    });
  });

  it("should not diff anything if no matching file exists on `toDir`", async () => {
    const fromDir = await Deno.makeTempDir({ dir: workingDir });
    const toDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: fromDir });
    const to = from.replace(fromDir, toDir);

    const result = await diff({ fromDir, toDir, extra: { logger } });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { log: { warn: 1 } });
    assertSpyCall(logger.log.warn, 0, {
      args: [
        `File ${to} does not exists or is not a file or is not readable, skipping...`,
      ],
    });
  });

  it("should not diff anything if matching file is not readable `toDir`", async () => {
    const fromDir = await Deno.makeTempDir({ dir: workingDir });
    const toDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: fromDir });
    const to = from.replace(fromDir, toDir);

    await copy(from, to);
    await Deno.chmod(to, 0o300);

    const result = await diff({ fromDir, toDir, extra: { logger } });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { log: { warn: 1 } });
    assertSpyCall(logger.log.warn, 0, {
      args: [
        `File ${to} does not exists or is not a file or is not readable, skipping...`,
      ],
    });
  });

  it("should diff and not show diferences if files are the same", async () => {
    const fromDir = await Deno.makeTempDir({ dir: workingDir });
    const toDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: fromDir });
    const to = from.replace(fromDir, toDir);

    await Deno.writeTextFile(from, "foo");
    await copy(from, to);

    const result = await diff({ fromDir, toDir, extra: { logger } });

    assertEquals(result, true);
    assertSpyLoggerCalls(logger, { log: { log: 1 }, logNoNl: { log: 1 } });
    assertSpyCall(logger.logNoNl.log, 0, {
      args: [`Diffing ${from} against ${to} ...`],
    });
    assertSpyCall(logger.log.log, 0, { args: [" ✓"] });
  });

  it("should diff and not show diferences if files are the same", async () => {
    const fromDir = await Deno.makeTempDir({ dir: workingDir });
    const toDir = await Deno.makeTempDir({ dir: workingDir });
    const from = await Deno.makeTempFile({ dir: fromDir, prefix: "1-" });
    const to = from.replace(fromDir, toDir);
    const from2 = await Deno.makeTempFile({ dir: fromDir, prefix: "2-" });
    const to2 = from2.replace(fromDir, toDir);

    await Deno.writeTextFile(from, "foo");
    await Deno.writeTextFile(from2, "foo");
    await copy(from, to);
    await copy(from2, to2);
    await Deno.writeTextFile(to, "foa");

    const result = await diff({
      fromDir,
      toDir,
      extra: { logger },
    });

    assertEquals(result, false);
    assertSpyLoggerCalls(logger, { log: { log: 3 }, logNoNl: { log: 2 } });
    assertArrayIncludesRemoveMatch([
      ...logger.logNoNl.log.calls[0].args,
      ...logger.logNoNl.log.calls[1].args,
    ], [
      `Diffing ${from} against ${to} ...`,
      `Diffing ${from2} against ${to2} ...`,
    ]);
    assertArrayIncludesRemoveMatch([
      ...logger.log.log.calls[0].args,
      ...logger.log.log.calls[1].args,
      ...logger.log.log.calls[2].args,
    ], [
      " ✓",
      " ✕",
      "--- a\n" +
      "+++ b\n" +
      "@@ -1 +1 @@\n" +
      "-foo\n" +
      "\\ No newline at end of file\n" +
      "+foa\n" +
      "\\ No newline at end of file",
    ]);
  });
});
