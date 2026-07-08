import type { BuildTargetType } from "../classes/BuildTarget";

export interface IKinaCompilerOptions {
  rootDir: string;
  buildDir: string;
  target: BuildTargetType;
  entry: string;

  debug?: {
    emitTokenized?: boolean;
    emitAST?: boolean;
    emitSymbols?: boolean;
    emitLLVM?: boolean;
    emitOptimizedLLVM?: boolean;
  };

  name: string;
  version: string;
}
