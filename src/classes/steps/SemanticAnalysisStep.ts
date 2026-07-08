import { KinaSemanticAnalyzer, type Scope } from "@kina-lang/semantic-analyzer";
import { CompilationStep } from "./_base";
import type { IKinaCompilerOptions } from "../../types/compiler";
import type { FileNode } from "@kina-lang/ast";
import type { KinaCompiler } from "../KinaCompiler";

export class SemanticAnalysisStep extends CompilationStep<Scope> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override execute(
    opts: IKinaCompilerOptions,
    ast: FileNode,
    filePath: string,
    isIncluded: boolean = false,
  ): Promise<Scope> {
    return this.withMetrics("semantic-analysis", async () => {
      const sa = new KinaSemanticAnalyzer();
      const scope = sa.analyze(ast, this._compiler, filePath, isIncluded);

      return scope;
    });
  }
}
