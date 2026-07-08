import z from "zod";

export const KinaProjectConfigSchema = z.object({
  package: z.object({
    name: z.string().min(1, "Package name cannot be empty"),
    author: z
      .string()
      .min(3, "Author name must be at least 3 characters long")
      .optional(),
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
  dependencies: z
    .array(
      z.object({
        name: z.string().min(1, "Dependency name cannot be empty"),
        author: z
          .string()
          .min(3, "Dependency author name must be at least 3 characters long"),
        version: z
          .string()
          .regex(/^\d+\.\d+\.\d+$/, "Must be semantic versioning"),
        source: z.string(),
      }),
    )
    .optional(),
});

export type IKinaProjectConfig = z.infer<typeof KinaProjectConfigSchema>;
