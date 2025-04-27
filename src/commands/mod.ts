import type { default as diff } from "./diff.ts";
import type { default as apply } from "./apply.ts";
import type { default as adopt } from "./adopt.ts";

export type Commands = {
  diff: typeof diff;
  apply: typeof apply;
  adopt: typeof adopt;
};

export { default as diff } from "./diff.ts";
export { default as apply } from "./apply.ts";
export { default as adopt } from "./adopt.ts";
