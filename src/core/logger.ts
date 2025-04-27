import { cyan, magenta, red, yellow } from "@std/fmt/colors";

export type LogFn = (msg: string) => Promise<void>;

export type Log = {
  debug: LogFn;
  info: LogFn;
  warn: LogFn;
  error: LogFn;
  log: LogFn;
};

export type Logger = {
  log: Log;
  logNoNl: Log;
};

const logLevels = {
  DEBUG: "DBG",
  INFO: "INF",
  WARN: "WRN",
  ERROR: "ERR",
} as const;

const logLevelColorMap = {
  [logLevels.DEBUG]: magenta,
  [logLevels.INFO]: cyan,
  [logLevels.WARN]: yellow,
  [logLevels.ERROR]: red,
};

type LogLevels = typeof logLevels;

const loggerTxtEncoder = new TextEncoder();

const log = async (
  { level, msg, nl, color }: {
    level?: LogLevels[keyof LogLevels];
    nl: boolean;
    msg: string;
    color: boolean;
  },
) => {
  await Deno.stdout.write(
    loggerTxtEncoder.encode(
      `${
        level
          ? (color ? logLevelColorMap[level](`${level}: `) : `${level}: `)
          : ""
      }${msg}${nl ? "\n" : ""}`,
    ),
  );
};

export const makeLogger = ({ color }: { color?: boolean }) => {
  return {
    log: {
      debug: (msg) =>
        log({ level: logLevels.DEBUG, color: !!color, msg, nl: true }),
      info: (msg) =>
        log({ level: logLevels.INFO, color: !!color, msg, nl: true }),
      warn: (msg) =>
        log({ level: logLevels.WARN, color: !!color, msg, nl: true }),
      error: (msg) =>
        log({ level: logLevels.ERROR, color: !!color, msg, nl: true }),
      log: (msg) => log({ color: !!color, msg, nl: true }),
    },
    logNoNl: {
      debug: (msg) =>
        log({ level: logLevels.DEBUG, color: !!color, msg, nl: false }),
      info: (msg) =>
        log({ level: logLevels.INFO, color: !!color, msg, nl: false }),
      warn: (msg) =>
        log({ level: logLevels.WARN, color: !!color, msg, nl: false }),
      error: (msg) =>
        log({ level: logLevels.ERROR, color: !!color, msg, nl: false }),
      log: (msg) => log({ color: !!color, msg, nl: false }),
    },
  } as Logger;
};
