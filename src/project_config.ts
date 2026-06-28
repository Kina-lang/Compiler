import { parse } from "smol-toml";
import {
  KinaProjectConfigSchema,
  type IKinaProjectConfig,
} from "./schema/project_config";
import z from "zod";

export class KinaProjectConfig {
  static parse(fileContents: string): IKinaProjectConfig {
    try {
      const raw = parse(fileContents);

      return KinaProjectConfigSchema.parse(raw);
    } catch (error) {
      if (error instanceof z.ZodError)
        throw new Error(
          `Invalid project configuration: ${error.issues[0]?.message ?? error.message ?? "Unknown error"}`,
        );

      throw error;
    }
  }
}
