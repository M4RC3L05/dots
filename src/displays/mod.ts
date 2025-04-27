import type { default as environment } from "./environment.ts";
import type { default as help } from "./help.ts";
import type { default as version } from "./version.ts";

export type Displays = {
  environment: typeof environment;
  help: typeof help;
  version: typeof version;
};

export { default as environment } from "./environment.ts";
export { default as help } from "./help.ts";
export { default as version } from "./version.ts";
