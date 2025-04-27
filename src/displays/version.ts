import meta from "../../deno.json" with { type: "json" };
import type { Logger } from "./../core/logger.ts";

export default async ({ logger }: { logger: Logger }) => {
  await logger.log.log(`v${meta.version}`);
};
