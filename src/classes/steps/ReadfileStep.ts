import path from "path";
import type { IKinaCompilerOptions } from "../../types/compiler";
import { CompilationStep } from "./_base";
import { readFile } from "fs/promises";
import type { KinaCompiler } from "../KinaCompiler";

export class ReadfileStep extends CompilationStep<{
  filePath: string;
  fileContent: string;
}> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override async execute(
    opts: IKinaCompilerOptions,
    filePath?: string,
  ): Promise<{ filePath: string; fileContent: string }> {
    return this.withMetrics("readfile", async () => {
      const fullPath = path.join(opts.rootDir, filePath ?? opts.entry);

      return {
        filePath: fullPath,
        fileContent: await readFile(fullPath, { encoding: "utf-8" }),
      };
    });
  }
}
