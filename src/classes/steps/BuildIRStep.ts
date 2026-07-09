import type { FileNode } from "@kina-lang/ast";
import type { IKinaCompilerOptions } from "../../types/compiler";
import type { KinaCompiler } from "../KinaCompiler";
import { CompilationStep } from "./_base";
import type { Scope } from "@kina-lang/semantic-analyzer";
import { KinaIRBuilder } from "@kina-lang/ir-builder";

export class BuildIRStep extends CompilationStep<string> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override execute(
    opts: IKinaCompilerOptions,
    ast: FileNode,
    scope: Scope,
    fileName: string,
    isIncluded: boolean = false,
  ): Promise<string> {
    return this.withMetrics("build-ir", async () => {
      const builder = new KinaIRBuilder();
      const ir = builder.build(ast, scope, isIncluded);

      if (opts.debug?.emitLLVM)
        this._compiler.debugArtifactEmitter.add(fileName, "ir", ir, ".ll");

      return ir;
    });
  }
}
