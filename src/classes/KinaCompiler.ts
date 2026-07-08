import type { IKinaCompilerOptions } from "../types/compiler";
import { KinaAssertionError, KinaLogger } from "@kina-lang/utils";
import { CompilerMetrics } from "./CompilerMetrics";
import { getSteps } from "./steps";
import { BuildTarget, type BuildTargetType } from "./BuildTarget";
import { PathResolver } from "./PathResolver";
import path from "path";
import { CSymbols } from "./CSymbols";
import { TypeTranslator } from "./TypeTranslator";
import { IncludeManager } from "./IncludeManager";
import { readFile } from "fs/promises";
import type { Scope } from "@kina-lang/semantic-analyzer";

export class KinaCompiler {
  private readonly _config: IKinaCompilerOptions;
  private readonly _logger: KinaLogger = new KinaLogger(KinaCompiler.name);
  private _metrics: CompilerMetrics = new CompilerMetrics();
  private _buildRoot: string | null = null;
  private _objDir: string | null = null;
  private readonly _steps = getSteps(this);
  private readonly _includesCache: Map<string, Scope> = new Map([]);
  private readonly _cIncludesCache: Set<string> = new Set([]);

  public readonly pathResolver = new PathResolver(this);
  public readonly cSymbols = new CSymbols(this);
  public readonly typeTranslator = new TypeTranslator();
  public readonly includeManager = new IncludeManager();

  private static readonly RUNTIME_DIR = path.join(
    import.meta.filename,
    "../../../../runtime/build",
  );

  constructor(config: IKinaCompilerOptions) {
    this._config = config;
  }

  public get metrics(): CompilerMetrics {
    return this._metrics;
  }

  public get logger(): KinaLogger {
    return this._logger;
  }

  public get config(): IKinaCompilerOptions {
    return this._config;
  }

  private getRuntimePath(target: BuildTargetType): string {
    const runtimePath = path.join(
      KinaCompiler.RUNTIME_DIR,
      target,
      "kina-runtime.a",
    );

    return runtimePath;
  }

  async compile(): Promise<string[]> {
    if (!BuildTarget.isSupported(this._config.target))
      throw new KinaAssertionError(
        `Unsupported build target: ${this._config.target}`,
      );

    this._logger.info(
      `Starting compilation of ${this._config.name}@v${this._config.version} (target: ${this._config.target})...`,
    );
    this._logger.info(`Compiling ${this._config.entry}...`);
    this._metrics.capture("total");

    const { buildRoot, objDir } = await this._steps.PrepareFS.execute(
      this._config,
    );

    this._buildRoot = buildRoot;
    this._objDir = objDir;

    const readfileStepResult = await this._steps.Readfile.execute(this._config);
    const tokenizeStepResult = await this._steps.Tokenize.execute(
      this._config,
      readfileStepResult.fileContent,
    );
    const buildastStepResult = await this._steps.Buildast.execute(
      this._config,
      tokenizeStepResult,
    );
    const semanticAnalysisStepResult =
      await this._steps.SemanticAnalysis.execute(
        this._config,
        buildastStepResult,
        readfileStepResult.filePath,
      );
    const buildirStepResult = await this._steps.BuildIR.execute(
      this._config,
      buildastStepResult,
      semanticAnalysisStepResult,
    );
    const optimizeirStepResult = await this._steps.OptimizeIR.execute(
      this._config,
      buildirStepResult,
    );
    const compileStepResult = await this._steps.Compile.execute(
      this._config,
      optimizeirStepResult,
      this._config.entry,
      objDir,
    );

    // Add runtime into linker inputs
    this.includeManager.add(this.getRuntimePath(this._config.target));
    this.includeManager.add(compileStepResult);

    const outPath = await BuildTarget.getTarget(
      this._config.target,
    ).buildOutput(this.includeManager.getAll(), buildRoot);

    this._metrics.capture("total");
    const compilationTime = this._metrics.calculateDelta("total");

    this._logger.info(
      `Compilation finished in ${compilationTime.toFixed(2)}ms`,
    );

    return outPath;
  }

  // TODO: Add caching
  async compileIncluded(file: string, currentFile: string) {
    const previousMetrics = this._metrics;
    this._metrics = new CompilerMetrics();

    const relativeFilePath = path.relative(this._config.rootDir, file);

    if (this._includesCache.has(relativeFilePath))
      return { scope: this._includesCache.get(relativeFilePath)! };

    this._logger.info(`Compiling ${relativeFilePath}...`);
    this._metrics.capture("total");

    const readfileStepResult = await this._steps.Readfile.execute(
      this._config,
      relativeFilePath,
    );
    const tokenizeStepResult = await this._steps.Tokenize.execute(
      this._config,
      readfileStepResult.fileContent,
    );
    const buildastStepResult = await this._steps.Buildast.execute(
      this._config,
      tokenizeStepResult,
    );
    const semanticAnalysisStepResult =
      await this._steps.SemanticAnalysis.execute(
        this._config,
        buildastStepResult,
        readfileStepResult.filePath,
        true,
      );
    const buildirStepResult = await this._steps.BuildIR.execute(
      this._config,
      buildastStepResult,
      semanticAnalysisStepResult,
      true,
    );
    const optimizeirStepResult = await this._steps.OptimizeIR.execute(
      this._config,
      buildirStepResult,
    );
    const compileStepResult = await this._steps.Compile.execute(
      this._config,
      optimizeirStepResult,
      relativeFilePath,
      this._objDir!,
    );

    this.includeManager.add(compileStepResult);

    this._metrics.capture("total");
    this._metrics = previousMetrics;

    this._includesCache.set(relativeFilePath, semanticAnalysisStepResult);

    return { scope: semanticAnalysisStepResult };
  }

  async compileIncludedC(file: string, currentFile: string) {
    const previousMetrics = this._metrics;
    this._metrics = new CompilerMetrics();

    const relativeFilePath = path.relative(this._config.rootDir, file);
    const outPath = path.join(
      this._objDir!,
      relativeFilePath.replaceAll("_", "__").replaceAll("/", "_") + ".o",
    );

    if (this._cIncludesCache.has(relativeFilePath)) return;

    this._logger.info(`Compiling ${relativeFilePath}...`);
    this._metrics.capture("total");

    const compileStepResult = await BuildTarget.getTarget(
      this._config.target,
    ).buildObjectFileFromC(
      await readFile(file, { encoding: "utf-8" }),
      outPath,
    );

    this.includeManager.add(outPath);

    this._metrics.capture("total");
    this._metrics = previousMetrics;

    this._cIncludesCache.add(relativeFilePath);
  }
}
