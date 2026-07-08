import path from "path";
import type { IKinaCompilerOptions } from "../../types/compiler";
import { BuildTarget } from "../BuildTarget";
import type { KinaCompiler } from "../KinaCompiler";
import { CompilationStep } from "./_base";

export class CompileStep extends CompilationStep<string> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override async execute(
    opts: IKinaCompilerOptions,
    ir: string,
    filePath: string,
    objectFileDirectoryPath: string,
  ): Promise<string> {
    return this.withMetrics("compile", async () => {
      const target = BuildTarget.getTarget(opts.target);
      const outPath = path.join(
        objectFileDirectoryPath,
        filePath.replaceAll("_", "__").replaceAll("/", "_") + ".o",
      );

      await target.buildObjectFileFromLLVM(ir, outPath);

      return outPath;
    });
  }
}
