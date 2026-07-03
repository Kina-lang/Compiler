import z from "zod";

export const KinaProjectConfigSchema = z.object({
  package: z.object({
    name: z.string().min(1, "Package name cannot be empty"),
    version: z.string().regex(/^\d+\.\d+\.\d+$/, "Must be semantic versioning"),
    entry: z.string(),
  }),
  debug: z
    .object({
      emitTokenized: z.boolean(),
      emitAST: z.boolean(),
      emitSymbols: z.boolean(),
      emitLLVM: z.boolean(),
      emitOptimizedLLVM: z.boolean(),
    })
    .partial()
    .optional(),
});

export type IKinaProjectConfig = z.infer<typeof KinaProjectConfigSchema>;
