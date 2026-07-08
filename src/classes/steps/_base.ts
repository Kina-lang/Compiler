import type { IKinaCompilerOptions } from "../../types/compiler";
import type { KinaCompiler } from "../KinaCompiler";

export abstract class CompilationStep<T> {
  protected readonly _compiler: KinaCompiler;

  constructor(compiler: KinaCompiler) {
    this._compiler = compiler;
  }

  abstract execute(opts: IKinaCompilerOptions, ...args: any): Promise<T>;

  protected async withMetrics<T>(
    stepName: string,
    fn: () => Promise<T>,
  ): Promise<T> {
    this._compiler.metrics.capture(stepName);

    return fn().finally(() => {
      this._compiler.metrics.capture(stepName);
      this._compiler.logger.debug(
        `Step ${stepName} took ${this._compiler.metrics
          .calculateDelta(stepName)
          .toFixed(2)}ms`,
      );
    });
  }
}
