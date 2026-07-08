import path from "path";
import type { IKinaCompilerOptions } from "../../types/compiler";
import type { KinaCompiler } from "../KinaCompiler";
import { CompilationStep } from "./_base";
import { existsSync } from "fs";
import { mkdir, rm } from "fs/promises";

export class PrepareFSStep extends CompilationStep<{
  buildRoot: string;
  objDir: string;
}> {
  constructor(compiler: KinaCompiler) {
    super(compiler);
  }

  override async execute(opts: IKinaCompilerOptions) {
    const buildRoot = path.join(
      opts.buildDir,
      opts.name,
      opts.version,
      opts.target,
    );
    const objDir = path.join(buildRoot, "obj");

    // TODO: Add support for incremental builds
    if (existsSync(buildRoot))
      await rm(buildRoot, { recursive: true, force: true });
    await mkdir(buildRoot, { recursive: true });
    await mkdir(objDir, { recursive: true });

    return { buildRoot, objDir };
  }
}
