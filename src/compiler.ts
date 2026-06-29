import { KinaLexer } from "@kina-lang/lexer";
import { KinaLogger } from "@kina-lang/utils";
import type { IKinaCompilerOptions } from "./types/compiler";
import path from "path";
import { existsSync } from "fs";
import { mkdir, readFile, rm, writeFile } from "fs/promises";
import {
  EKinaASTNodeKind,
  KinaASTLiteralExpressionNode,
  KinaASTParser,
} from "@kina-lang/ast";
import { KinaSemanticAnalyzer } from "@kina-lang/semantic-analyzer";
import { KinaIRBuilder } from "@kina-lang/ir-builder";
import { spawn } from "child_process";
import type { KinaASTNode } from "@kina-lang/ast/src/nodes/_node";
import type { KinaASTIncludeDirectiveNode } from "@kina-lang/ast/src/nodes/includeDirective";

export class KinaCompiler {
  private readonly logger: KinaLogger = new KinaLogger(KinaCompiler.name);

  // TODO: Move this into language SDK dir on production build
  private static readonly RUNTIME_PATH = path.join(
    import.meta.filename,
    "../../../runtime/build/kina-runtime.a",
  );

  constructor() {}

  public async compile(options: IKinaCompilerOptions) {
    this.logger.info(`Compiling ${options.name}@${options.version}`);

    const buildRoot = await this.prepareBuildDirectoryTree(options);

    const files = [options.entry];

    const includes: string[] = [];

    while (files.length > 0) {
      const file = files.shift()!;
      const fullPath = path.join(options.rootDir, file);

      const tokens = await new KinaLexer().process(
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

      const symbolTable = await new KinaSemanticAnalyzer(ast).analyze();
      await writeFile(
        path.join(buildRoot, "sa", file.replaceAll("/", "$") + ".__sa.json"),
        JSON.stringify(symbolTable.toJson(), null, 2),
      );

      const { includedCFiles } = await this.processDirectives(fullPath, ast);
      includes.push(...includedCFiles);

      const ir = await new KinaIRBuilder(ast, symbolTable).build();

      await writeFile(
        path.join(buildRoot, "ir", file.replaceAll("/", "$") + ".__ir.ll"),
        ir,
      );
    }

    const outPath = path.join(
      buildRoot,
      path.basename(options.entry).split(".").slice(0, -1).join("."),
    );

    await new Promise<void>((res) => {
      const proc = spawn(
        "clang",
        [
          path.join(
            buildRoot,
            "ir",
            options.entry.replaceAll("/", "$") + ".__ir.ll",
          ),
          KinaCompiler.RUNTIME_PATH,
          ...includes,
          "-o",
          outPath,
        ],
        {
          stdio: "inherit",
        },
      );

      proc.on("exit", res);
    });

    return outPath;
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

    this.logger.debug("Creating ir directory...");
    await mkdir(path.join(buildRoot, "ir"), { recursive: true });

    this.logger.debug("Creating sa directory...");
    await mkdir(path.join(buildRoot, "sa"), { recursive: true });

    return buildRoot;
  }

  private async processDirectives(filePath: string, ast: KinaASTNode[]) {
    const includedCFiles: string[] = [];

    for (const node of ast) {
      switch (node.kind) {
        case EKinaASTNodeKind.IncludeDirective:
          includedCFiles.push(
            path.resolve(
              path.dirname(filePath),
              (node as KinaASTIncludeDirectiveNode).argument.value,
            ),
          );
          break;
        default:
          continue;
      }
    }

    return { includedCFiles };
  }
}
