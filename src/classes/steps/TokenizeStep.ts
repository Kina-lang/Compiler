import { KinaASI, KinaLexer, type BaseToken } from "@kina-lang/lexer";
import { CompilationStep } from "./_base";
import type { IKinaCompilerOptions } from "../../types/compiler";
import type { KinaCompiler } from "../KinaCompiler";

export class TokenizeStep extends CompilationStep<BaseToken[]> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override async execute(
    opts: IKinaCompilerOptions,
    fileName: string,
    fileContents: string,
  ): Promise<BaseToken[]> {
    return this.withMetrics("tokenize", async () => {
      const lexer = new KinaLexer({
        fileName: opts.entry,
        rootDir: opts.rootDir,
        skipUnknownTokens: false,
      });
      const asi = new KinaASI();

      const tokens = lexer.tokenize(fileContents);
      const mandatoryAndNewlines = lexer.filterMandatory(tokens, true);
      const asiProcessed = asi.process(mandatoryAndNewlines);
      const mandatoryAsiProcessed = lexer.filterMandatory(asiProcessed, false);

      if (opts.debug?.emitTokenized)
        this._compiler.debugArtifactEmitter.add(
          fileName,
          "lex",
          JSON.stringify(mandatoryAsiProcessed, null, 2),
          ".json",
        );

      return mandatoryAsiProcessed;
    });
  }
}
