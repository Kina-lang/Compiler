import { KinaAssertionError } from "@kina-lang/utils";
import { existsSync, readFileSync } from "fs";
import path from "path";
import type { KinaCompiler } from "./KinaCompiler";
import { KinaProjectConfig } from "./KinaProjectConfig";

export class PathResolver {
  private readonly _compiler: KinaCompiler;

  constructor(compiler: KinaCompiler) {
    this._compiler = compiler;
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
      this._compiler.config.rootDir,
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
}
