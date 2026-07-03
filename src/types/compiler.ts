export interface IKinaCompilerOptions {
  rootDir: string;
  buildDir: string;
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
