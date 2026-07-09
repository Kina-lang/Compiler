import path from "path";
import { CommandRunner } from "../CommandRunner";
import { CompilationTarget } from "./_base";

export class X86_64UnknownLinuxGnuTarget extends CompilationTarget {
  constructor() {
    super();
  }

  override async buildObjectFileFromLLVM(
    inputCode: string,
    outputPath: string,
  ): Promise<void> {
    const stdout = await CommandRunner.runWithPipe(
      "clang",
      ["-O3", "-c", "-x", "ir", "-", "-o", outputPath],
      inputCode,
    );
  }

  override async buildObjectFileFromC(
    input: string,
    output: string,
    includeDirs?: string[],
  ): Promise<void> {
    const clangArgs = ["-O3", "-c", "-x", "c"];

    if (includeDirs)
      for (const dir of includeDirs) {
        clangArgs.push(`-I${dir}`);
      }

    clangArgs.push("-", "-o", output);

    const stdout = await CommandRunner.runWithPipe("clang", clangArgs, input);
  }

  override async buildOutput(
    includedFiles: Set<string>,
    buildRoot: string,
  ): Promise<string[]> {
    const outputFilePath = path.join(buildRoot, "output");

    const stdout = await CommandRunner.run("clang", [
      ...includedFiles,
      "-o",
      outputFilePath,
    ]);

    return [outputFilePath];
  }
}
