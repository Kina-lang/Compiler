import { FileNode, KinaAST } from "@kina-lang/ast";
import { CompilationStep } from "./_base";
import type { IKinaCompilerOptions } from "../../types/compiler";
import type { BaseToken } from "@kina-lang/lexer";
import type { KinaCompiler } from "../KinaCompiler";

export class BuildastStep extends CompilationStep<FileNode> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override async execute(
    opts: IKinaCompilerOptions,
    tokens: BaseToken[],
  ): Promise<FileNode> {
    return this.withMetrics("buildast", async () => {
      const ast = new KinaAST();
      const tree = ast.build(tokens);

      return tree;
    });
  }
}
