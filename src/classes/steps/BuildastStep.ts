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
    fileName: string,
    tokens: BaseToken[],
  ): Promise<FileNode> {
    return this.withMetrics("buildast", async () => {
      const ast = new KinaAST();
      const tree = ast.build(tokens);

      if (opts.debug?.emitAST)
        this._compiler.debugArtifactEmitter.add(
          fileName,
          "ast",
          JSON.stringify(tree.export(), null, 2),
          ".json",
        );

      return tree;
    });
  }
}
