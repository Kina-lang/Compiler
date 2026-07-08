import { BaseToken, KinaASI, KinaLexer } from "@kina-lang/lexer";
import { KinaAssertionError, KinaLogger } from "@kina-lang/utils";
import type { IKinaCompilerOptions } from "./types/compiler";
import path from "path";
import { existsSync, readFileSync } from "fs";
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
import Parser from "tree-sitter";
import C from "tree-sitter-c";
import type { KinaTypeTokenKind } from "@kina-lang/ast/src/types/types";
import { TypeTranslator } from "./TypeTranslator";
import { KinaProjectConfig } from "./project_config";

export class KinaCompiler {
  private readonly logger: KinaLogger = new KinaLogger(KinaCompiler.name);
  private readonly _options: IKinaCompilerOptions;
  private _buildRoot: string | null = null;
  private _includedFiles: string[] = [];

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

  public async compileIncludedFile(filePath: string, currentFilePath: string) {
    this.logger.info(`Compiling included file: ${filePath}`);

    const relativePath = path.relative(this._options.rootDir, filePath);
    const outObjPath = path.join(
      this._buildRoot!,
      `${relativePath.replaceAll("/", "$")}.o`,
    );

    const fileContents = await readFile(filePath, "utf-8");
    const tokens = await this.tokenize(relativePath, fileContents);
    const ast = await this.buildAST(relativePath, tokens);
    const scope = await this.semanticallyAnalyze(relativePath, ast, true);
    const ir = await this.buildIR(relativePath, ast, scope, true);
    const optIr = await this.optimizeIR(relativePath, ir);
    const { includedCFiles } = await this.processDirectives(filePath, ast);
    await this.compileIrObject(relativePath, optIr, outObjPath);

    for (const includedFile of [...includedCFiles]) {
      this.includeFile(includedFile);
    }

    this.includeFile(outObjPath);

    return { scope };
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

  private async semanticallyAnalyze(
    filePath: string,
    ast: FileNode,
    isIncluded: boolean = false,
  ) {
    const sa = new KinaSemanticAnalyzer();
    const scope = await sa.analyze(ast, this, filePath, isIncluded);

    if (this._options.debug?.emitSymbols)
      await this.emitDebugArtifact(
        `${filePath.replaceAll("/", "$")}.__symbols.json`,
        scope.export(),
      );

    return scope;
  }

  private async buildIR(
    filePath: string,
    ast: FileNode,
    scope: Scope,
    isIncluded: boolean = false,
  ) {
    const irBuilder = new KinaIRBuilder();
    const ir = irBuilder.build(ast, scope, isIncluded);

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
        ...this._includedFiles,
        "-O3",
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

  private async compileIrObject(filePath: string, ir: string, outPath: string) {
    await new Promise<void>((res) => {
      const proc = spawn("clang", [
        "-O3",
        "-c",
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

  public resolveIncludePath(includePath: string, currentFilePath: string) {
    const resolvedPath = path.resolve(
      path.dirname(currentFilePath),
      includePath,
    );

    if (!existsSync(resolvedPath))
      throw new KinaAssertionError(
        `Include file not found: ${includePath} (resolved to ${resolvedPath})`,
      );

    return resolvedPath;
  }

  public resolveNamespacedPath(includePath: string, currentFilePath: string) {
    const kinaModuleDir = path.join(
      this._options.rootDir,
      ".kina_modules",
      includePath.split(".")[0]!,
      includePath.split(".")[1]!,
    );
    if (!existsSync(kinaModuleDir))
      throw new KinaAssertionError(
        `Namespaced module not found: ${includePath}`,
      );

    const moduleConfigPath = path.join(kinaModuleDir, "kina.toml");
    if (!existsSync(moduleConfigPath))
      throw new KinaAssertionError(
        `Namespaced module config not found: ${includePath} (expected at ${moduleConfigPath})`,
      );

    const modConfig = KinaProjectConfig.parse(
      readFileSync(moduleConfigPath, "utf-8"),
    );
    if (!modConfig)
      throw new KinaAssertionError(
        `Failed to parse module config: ${moduleConfigPath}`,
      );

    const entryDirname = path.dirname(modConfig.package.entry);
    const fullRoot = path.join(kinaModuleDir, entryDirname);
    let filepath = fullRoot;

    for (const part of includePath.split(".").slice(2)) {
      filepath = path.join(filepath, part);
    }

    filepath = path.join(filepath, "lib.kin");

    if (!existsSync(filepath))
      throw new KinaAssertionError(
        `Include file not found: ${includePath} (resolved to ${filepath})`,
      );

    return filepath;
  }

  public getCSymbols(
    filePath: string,
  ): Record<
    string,
    { returnType: KinaTypeTokenKind; parameterTypes: KinaTypeTokenKind[] }
  > {
    const parser = new Parser();
    // @ts-ignore Module 'tree-sitter-c' has bad typing
    parser.setLanguage(C as any);

    const fileContents = readFileSync(filePath, "utf-8");
    const tree = parser.parse(fileContents);

    const symbols: Record<
      string,
      { returnType: KinaTypeTokenKind; parameterTypes: KinaTypeTokenKind[] }
    > = {};

    for (const node of tree.rootNode.children) {
      if (node.type !== "function_definition") continue;

      const declaratorNode = node.childForFieldName("declarator");
      if (!declaratorNode) continue;

      let currentDecl = declaratorNode;
      while (currentDecl.childForFieldName("declarator")) {
        currentDecl = currentDecl.childForFieldName("declarator")!;
      }

      const functionNameNode = currentDecl;
      if (!functionNameNode || functionNameNode.type !== "identifier") continue;

      const returnTypeNode = node.childForFieldName("type");
      if (!returnTypeNode) continue;

      let functionDeclaratorNode = declaratorNode;
      while (
        functionDeclaratorNode &&
        functionDeclaratorNode.type !== "function_declarator"
      ) {
        functionDeclaratorNode =
          functionDeclaratorNode.childForFieldName("declarator")!;
      }

      const parameterTypeNodes = functionDeclaratorNode
        ?.childForFieldName("parameters")
        ?.namedChildren.filter((n) => n.type === "parameter_declaration");
      if (!parameterTypeNodes) continue;

      const functionName = functionNameNode.text;

      let returnType = returnTypeNode.text;
      if (declaratorNode.type === "pointer_declarator") {
        returnType = `${returnType}*`;
      }

      const parameterTypes = parameterTypeNodes.map((n) => {
        const typeNode = n.childForFieldName("type");
        if (!typeNode)
          throw new KinaAssertionError(
            `Parameter declaration missing type: ${n.text}`,
          );

        const paramDecl = n.childForFieldName("declarator");
        const isPointer = paramDecl && paramDecl.text.startsWith("*");

        return isPointer ? `${typeNode.text}*` : typeNode.text;
      });

      symbols[functionName] = {
        returnType: TypeTranslator.cToKina(returnType),
        parameterTypes: parameterTypes.map((t) => TypeTranslator.cToKina(t)),
      };
    }

    return symbols;
  }

  public includeFile(filePath: string) {
    if (this._includedFiles.includes(filePath)) return;

    this._includedFiles.push(filePath);
  }
}
