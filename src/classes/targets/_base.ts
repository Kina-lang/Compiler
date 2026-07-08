export abstract class CompilationTarget {
  abstract buildObjectFileFromLLVM(
    input: string,
    output: string,
  ): Promise<void>;

  abstract buildObjectFileFromC(input: string, output: string): Promise<void>;

  abstract buildOutput(
    includedFiles: Set<string>,
    buildRoot: string,
  ): Promise<string[]>;
}
