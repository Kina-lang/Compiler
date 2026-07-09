import { copyFile } from "fs/promises";
import { CommandRunner } from "../../CommandRunner";
import { CompilationTarget } from "../_base";
import { IDFResolver, IDFProjectBridge } from "./IDFResolver";
import path from "path";

export class XtensaESP32ELFTarget extends CompilationTarget {
  private readonly _idf = new IDFResolver();

  constructor() {
    super();
  }

  override async buildObjectFileFromLLVM(
    inputCode: string,
    outputPath: string,
  ): Promise<void> {
    const llcPath = await this._idf.resolveLLC();

    await CommandRunner.runWithPipe(
      llcPath,
      ["-mtriple=xtensa-esp32-elf", "-O3", "-filetype=obj", "-o", outputPath],
      inputCode,
    );
  }

  override async buildObjectFileFromC(
    input: string,
    output: string,
    includeDirs?: string[],
  ): Promise<void> {
    const gccPath = await this._idf.resolveGCC();
    const gccArgs = ["-mlongcalls", "-O3", "-c", "-x", "c"];

    if (includeDirs)
      for (const dir of includeDirs) {
        gccArgs.push(`-I${dir}`);
      }

    gccArgs.push("-", "-o", output);

    await CommandRunner.runWithPipe(gccPath, gccArgs, input);
  }

  override async buildOutput(
    includedFiles: Set<string>,
    buildRoot: string,
  ): Promise<string[]> {
    const bridge = new IDFProjectBridge(buildRoot);

    // Scaffold IDF project and inject .o files
    await bridge.setup(includedFiles);

    // Let IDF handle linking, linker scripts, esptool, ...
    const python = await this._idf.resolvePython();
    const idfPy = await this._idf.resolveIDFPy();
    const env = await this._idf.buildIDFEnv();

    await CommandRunner.run(python, [idfPy, "build"], {
      cwd: bridge.projectPath,
      env,
    });

    await copyFile(
      bridge.getOutputBinPath(),
      path.join(buildRoot, "output.bin"),
    );

    return [
      "echo",
      `Build complete. Output binary is located at: ${path.join(buildRoot, "output.bin")}. Use esptool.py to flash the binary to your ESP32 device.`,
    ];
  }
}
