import z from "zod";

export const KinaProjectConfigSchema = z.object({
  package: z.object({
    name: z.string().min(1, "Package name cannot be empty"),
    version: z.string().regex(/^\d+\.\d+\.\d+$/, "Must be semantic versioning"),
    entry: z.string(),
  }),
});

export type IKinaProjectConfig = z.infer<typeof KinaProjectConfigSchema>;
