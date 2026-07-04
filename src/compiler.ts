import { BaseToken, KinaASI, KinaLexer } from "@kina-lang/lexer";
import { KinaAssertionError, KinaLogger } from "@kina-lang/utils";
import type { IKinaCompilerOptions } from "./types/compiler";
import path from "path";
import { existsSync } from "fs";
import { mkdir, readFile, rm, writeFile } from "fs/promises";
import {
  FileNode,
  IncludeDirectiveNode,
  KinaAST,
  NodeKind,
} from "@kina-lang/ast";
import { KinaSemanticAnalyzer } from "@kina-lang/semantic-analyzer";
import type { Scope } from "@kina-lang/semantic-analyzer";
import { KinaIRBuilder } from "@kina-lang/ir-builder";
import { spawn } from "child_process";

export class KinaCompiler {
  private readonly logger: KinaLogger = new KinaLogger(KinaCompiler.name);
  private readonly _options: IKinaCompilerOptions;
  private _buildRoot: string | null = null;

  // TODO: Move this into language SDK dir on production build
  private static readonly RUNTIME_PATH = path.join(
    import.meta.filename,
    "../../../runtime/build/kina-runtime.a",
  );

  constructor(options: IKinaCompilerOptions) {
    this._options = options;
  }

  public async compile() {
    this.logger.info(
      `Compiling ${this._options.name}@${this._options.version}`,
    );
    const s_time = performance.now();

    const buildRoot = await this.prepareBuildDirectoryTree();
    this._buildRoot = buildRoot;

    const files = [this._options.entry];

    const outPath = path.join(
      buildRoot,
      path.basename(this._options.entry).split(".").slice(0, -1).join("."),
    );

    while (files.length > 0) {
      const file = files.shift()!;
      const fullPath = path.join(this._options.rootDir, file);

      const fileContents = await readFile(fullPath, "utf-8");
      const tokens = await this.tokenize(file, fileContents);
      const ast = await this.buildAST(file, tokens);
      const scope = await this.semanticallyAnalyze(file, ast);
      const ir = await this.buildIR(file, ast, scope);
      const optIr = await this.optimizeIR(file, ir);
      const { includedCFiles } = await this.processDirectives(fullPath, ast);
      await this.compileIr(file, optIr, includedCFiles, outPath);
    }

    const e_time = performance.now();
    this.logger.info(
      `Compilation finished in ${(e_time - s_time).toFixed(2)}ms`,
    );

    return outPath;
  }

  private async tokenize(filePath: string, fileContents: string) {
    const lexer = new KinaLexer({
      fileName: filePath,
      rootDir: this._options.rootDir,
      skipUnknownTokens: false,
    });

    const rawTokens = lexer.tokenize(fileContents);
    const mandatoryAndNewlines = lexer.filterMandatory(rawTokens, true);
    const asiProcessed = new KinaASI().process(mandatoryAndNewlines);
    const mandatoryAsiProcessed = lexer.filterMandatory(asiProcessed, false);

    if (this._options.debug?.emitTokenized)
      await this.emitDebugArtifact(
        `${filePath.replaceAll("/", "$")}.__tokens.json`,
        mandatoryAsiProcessed.map((t) => t.export()),
      );

    return mandatoryAsiProcessed;
  }

  private async buildAST(filePath: string, tokens: BaseToken[]) {
    const ast = new KinaAST();

    const tree = ast.build(tokens);

    if (this._options.debug?.emitAST)
      await this.emitDebugArtifact(
        `${filePath.replaceAll("/", "$")}.__ast.json`,
        tree.export(),
      );

    return tree;
  }

  private async semanticallyAnalyze(filePath: string, ast: FileNode) {
    const sa = new KinaSemanticAnalyzer();
    const scope = sa.analyze(ast);

    if (this._options.debug?.emitSymbols)
      await this.emitDebugArtifact(
        `${filePath.replaceAll("/", "$")}.__symbols.json`,
        scope.export(),
      );

    return scope;
  }

  private async buildIR(filePath: string, ast: FileNode, scope: Scope) {
    const irBuilder = new KinaIRBuilder();
    const ir = irBuilder.build(ast, scope);

    if (this._options.debug?.emitLLVM)
      await this.emitDebugArtifact(
        `${filePath.replaceAll("/", "$")}.__ir.ll`,
        ir,
      );

    return ir;
  }

  private async optimizeIR(filePath: string, ir: string) {
    const optimizedIR = await new Promise<string>((res) => {
      const proc = spawn("opt", ["-O3", "-S"]);

      proc.on("spawn", () => {
        proc.stdin.write(ir);
        proc.stdin.end();
      });

      let optimizedIR = "";
      proc.stdout.on("data", (data) => {
        optimizedIR += data.toString();
      });

      proc.stderr.on("data", (data) => {
        this.logger.error(data.toString());
      });

      proc.on("close", (code) => {
        if (code !== null && code !== 0)
          throw new KinaAssertionError(`opt exited with code ${code}`);

        res(optimizedIR);
      });
    });

    if (this._options.debug?.emitOptimizedLLVM)
      await this.emitDebugArtifact(
        `${filePath.replaceAll("/", "$")}.__opt.ll`,
        optimizedIR,
      );

    return optimizedIR;
  }

  private async compileIr(
    filePath: string,
    ir: string,
    includedFiles: string[],
    outPath: string,
  ) {
    await new Promise<void>((res) => {
      const proc = spawn("clang", [
        KinaCompiler.RUNTIME_PATH,
        ...includedFiles,
        "-x",
        "ir",
        "-",
        "-o",
        outPath,
      ]);

      proc.on("spawn", () => {
        proc.stdin.write(ir);
        proc.stdin.end();
      });

      proc.stderr.on("data", (data) => {
        this.logger.error(data.toString());
      });

      proc.on("close", (code) => {
        if (code !== null && code !== 0)
          throw new KinaAssertionError(`clang exited with code ${code}`);

        res();
      });
    });
  }

  private async emitDebugArtifact(
    name: string,
    content: string | Record<string, any> | Array<any>,
  ) {
    const str =
      typeof content === "string" ? content : JSON.stringify(content, null, 2);

    this.logger.debug(`Emitting debug artifact: ${name}`);
    const p = path.join(this._buildRoot!, "debug", name);

    await writeFile(p, str, "utf-8");
  }

  private async prepareBuildDirectoryTree() {
    this.logger.debug("Preparing build directory tree...");

    const buildRoot = path.join(
      this._options.buildDir,
      `${this._options.name}@${this._options.version}`,
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

    this.logger.debug("Creating debug directory...");
    await mkdir(path.join(buildRoot, "debug"), { recursive: true });

    return buildRoot;
  }

  private async processDirectives(filePath: string, ast: FileNode) {
    const includedCFiles: string[] = [];

    for (const node of ast.nodes) {
      if (node.kind != NodeKind.IncludeDirective) continue;

      const includeNode = node;
      const includePath = path.resolve(
        path.dirname(filePath),
        (includeNode as IncludeDirectiveNode).path,
      );
      includedCFiles.push(includePath);
    }

    return { includedCFiles };
  }
}
