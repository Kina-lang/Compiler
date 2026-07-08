import z from "zod";
import {
  KinaProjectConfigSchema,
  type IKinaProjectConfig,
} from "../schema/KinaProjectConfigSchema";
import { parse, stringify } from "smol-toml";

export class KinaProjectConfig {
  constructor() {}

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

  static stringify(config: IKinaProjectConfig): string {
    try {
      const parsed = KinaProjectConfigSchema.parse(config);

      return stringify(parsed);
    } catch (error) {
      if (error instanceof z.ZodError)
        throw new Error(
          `Invalid project configuration: ${error.issues[0]?.message ?? error.message ?? "Unknown error"}`,
        );

      throw error;
    }
  }
}
