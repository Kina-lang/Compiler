import { KinaLexer } from "@kina-lang/lexer";
import { KinaLogger } from "@kina-lang/utils";
import type { IKinaCompilerOptions } from "./types/compiler";
import path from "path";
import { existsSync } from "fs";
import { mkdir, readFile, rm, writeFile } from "fs/promises";
import { KinaASTParser } from "@kina-lang/ast";
import { KinaSemanticAnalyzer } from "@kina-lang/semantic-analyzer";

export class KinaCompiler {
  private readonly logger: KinaLogger = new KinaLogger(KinaCompiler.name);
  private readonly lexer: KinaLexer = new KinaLexer();

  constructor() {}

  public async compile(options: IKinaCompilerOptions) {
    this.logger.info(`Compiling ${options.name}@${options.version}`);

    const buildRoot = await this.prepareBuildDirectoryTree(options);

    const files = [options.entry];

    while (files.length > 0) {
      const file = files.shift()!;
      const fullPath = path.join(options.rootDir, file);

      const tokens = await this.lexer.process(
        file,
        await readFile(fullPath, "utf-8"),
      );

      await writeFile(
        path.join(
          buildRoot,
          "lexer",
          file.replaceAll("/", "$") + ".__lex.json",
        ),
        JSON.stringify(tokens, null, 2),
      );

      const ast = await new KinaASTParser(tokens, file).parse();

      await writeFile(
        path.join(buildRoot, "ast", file.replaceAll("/", "$") + ".__ast.json"),
        JSON.stringify(ast, null, 2),
      );

      await new KinaSemanticAnalyzer(ast).analyze();
    }
  }

  private async prepareBuildDirectoryTree(options: IKinaCompilerOptions) {
    this.logger.debug("Preparing build directory tree...");

    const buildRoot = path.join(
      options.buildDir,
      `${options.name}@${options.version}`,
    );
    this.logger.debug(`Build root: ${buildRoot}`);

    // Remove existing build dir, if it exists
    if (existsSync(buildRoot)) {
      this.logger.debug("Directory exists, removing...");
      await rm(buildRoot, { recursive: true });
    }

    // Create dir
    this.logger.debug("Creating directory...");
    await mkdir(buildRoot, { recursive: true });

    this.logger.debug("Creating lexer directory...");
    await mkdir(path.join(buildRoot, "lexer"), { recursive: true });

    this.logger.debug("Creating ast directory...");
    await mkdir(path.join(buildRoot, "ast"), { recursive: true });

    return buildRoot;
  }
}
