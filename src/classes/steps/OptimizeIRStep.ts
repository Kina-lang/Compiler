import type { IKinaCompilerOptions } from "../../types/compiler";
import { CommandRunner } from "../CommandRunner";
import type { KinaCompiler } from "../KinaCompiler";
import { CompilationStep } from "./_base";

export class OptimizeIRStep extends CompilationStep<string> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override execute(
    opts: IKinaCompilerOptions,
    ir: string,
    fileName: string,
  ): Promise<string> {
    return this.withMetrics("optimize-ir", async () => {
      const res = await CommandRunner.runWithPipe("opt", ["-O3", "-S"], ir);

      if (opts.debug?.emitOptimizedLLVM)
        this._compiler.debugArtifactEmitter.add(fileName, "opt-ir", res, ".ll");

      return res;
    });
  }
}
