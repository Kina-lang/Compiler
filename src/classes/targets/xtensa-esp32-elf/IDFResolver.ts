import path from "path";
import { readdir, copyFile, writeFile, mkdir } from "fs/promises";
import { existsSync } from "fs";
import { homedir } from "os";

export class IDFResolver {
  private _llcPath: string | null = null;
  private _gccBinDir: string | null = null;

  public readonly _espressifRoot = path.join(homedir(), ".espressif");

  async resolveLLC(): Promise<string> {
    if (this._llcPath) return this._llcPath;

    const espClangDir = path.join(this._espressifRoot, "tools", "esp-clang");
    const version = await this.findLatestVersion(espClangDir, "esp-clang");

    this._llcPath = path.join(espClangDir, version, "esp-clang", "bin", "llc");

    return this._llcPath;
  }

  async resolveGCC(): Promise<string> {
    return path.join(await this.resolveGCCBinDir(), "xtensa-esp32-elf-gcc");
  }

  public async resolveGCCBinDir(): Promise<string> {
    if (this._gccBinDir) return this._gccBinDir;

    const xtensaDir = path.join(this._espressifRoot, "tools", "xtensa-esp-elf");
    const version = await this.findLatestVersion(xtensaDir, "xtensa-esp-elf");
    this._gccBinDir = path.join(xtensaDir, version, "xtensa-esp-elf", "bin");

    return this._gccBinDir;
  }

  async resolveIDFPy(): Promise<string> {
    const idfPath = await this.resolveIDFPath();

    if (idfPath) {
      const idfPy = path.join(idfPath, "tools", "idf.py");

      if (existsSync(idfPy)) return idfPy;
    }

    // Fallback - maybe it's on PATH
    return "idf.py";
  }

  async resolvePythonBinDir(): Promise<string | null> {
    const pythonDir = path.join(this._espressifRoot, "tools", "python");
    const versions = await readdir(pythonDir).catch(() => []);

    for (const ver of versions.sort().reverse()) {
      const binDir = path.join(pythonDir, ver, "venv", "bin");

      if (existsSync(binDir)) return binDir;
    }

    return null;
  }

  async resolvePython(): Promise<string> {
    const binDir = await this.resolvePythonBinDir();

    if (binDir) {
      const pythonPath = path.join(binDir, "python");
      if (existsSync(pythonPath)) return pythonPath;
    }

    return "python";
  }

  async resolveIDFPath(): Promise<string | null> {
    let entries: string[];
    try {
      entries = await readdir(this._espressifRoot);
    } catch {
      return null;
    }

    const candidates = entries
      .filter((e) => e.startsWith("v") || e === "master")
      .sort()
      .reverse();

    for (const ver of candidates) {
      const idfPath = path.join(this._espressifRoot, ver, "esp-idf");
      if (existsSync(idfPath)) return idfPath;
    }

    return null;
  }

  async buildIDFEnv(): Promise<NodeJS.ProcessEnv> {
    const env = { ...process.env };

    const idfPath = await this.resolveIDFPath();
    if (idfPath) env.IDF_PATH = idfPath;

    const binPaths: string[] = [];

    // Add python venv bin
    const pythonBinDir = await this.resolvePythonBinDir();
    if (pythonBinDir) binPaths.push(pythonBinDir);

    // Add compiler bin dir
    const gccBinDir = await this.resolveGCCBinDir();
    if (gccBinDir) binPaths.push(gccBinDir);

    const toolsRoot = path.join(this._espressifRoot, "tools");
    try {
      const toolNames = await readdir(toolsRoot);
      for (const toolName of toolNames) {
        if (toolName === "python" || toolName === "xtensa-esp-elf") continue;

        const toolDir = path.join(toolsRoot, toolName);
        const versions = await readdir(toolDir).catch(() => []);

        for (const ver of versions) {
          const versionDir = path.join(toolDir, ver);
          const binDir = path.join(versionDir, "bin");

          if (existsSync(binDir)) binPaths.push(binDir);
          else if (existsSync(versionDir)) binPaths.push(versionDir);
        }
      }
    } catch {}

    if (binPaths.length > 0) {
      const separator = process.platform === "win32" ? ";" : ":";
      env.PATH = `${binPaths.join(separator)}${env.PATH ? separator + env.PATH : ""}`;
    }

    return env;
  }

  private async findLatestVersion(
    dir: string,
    toolName: string,
  ): Promise<string> {
    let entries: string[];

    try {
      entries = await readdir(dir);
    } catch {
      throw new Error(
        `${toolName} not found at ${dir}.\n` +
          "Install ESP-IDF using the Espressif installer.",
      );
    }

    const versions = entries
      .filter((e) => !e.startsWith("."))
      .sort()
      .reverse();
    if (versions.length === 0)
      throw new Error(`No ${toolName} versions found in: ${dir}`);

    return versions[0]!;
  }
}

export class IDFProjectBridge {
  private readonly _projectDir: string;

  constructor(buildRoot: string) {
    this._projectDir = path.join(buildRoot, "idf_project");
  }

  get projectPath(): string {
    return this._projectDir;
  }

  get mainDir(): string {
    return path.join(this._projectDir, "main");
  }

  get buildDir(): string {
    return path.join(this._projectDir, "build");
  }

  async setup(objectFiles: Set<string>): Promise<void> {
    await mkdir(this.mainDir, { recursive: true });

    const copiedNames: string[] = [];

    for (const objFile of objectFiles) {
      const basename = path.basename(objFile);
      const dest = path.join(this.mainDir, basename);

      await copyFile(objFile, dest);

      copiedNames.push(basename);
    }

    await this.writeTopCMakeLists();
    await this.writeMainCMakeLists(copiedNames);
    await this.writeEntryPoint();
  }

  private async writeTopCMakeLists(): Promise<void> {
    const content = `cmake_minimum_required(VERSION 3.22)
include($ENV{IDF_PATH}/tools/cmake/project.cmake)
project(kina_app)
`;

    await writeFile(path.join(this._projectDir, "CMakeLists.txt"), content);
  }

  private async writeMainCMakeLists(objectFiles: string[]): Promise<void> {
    const objLines = objectFiles
      .map((f) => `  "\${CMAKE_CURRENT_SOURCE_DIR}/${f}"`)
      .join("\n");

    const content = `idf_component_register(SRCS "kina_entry.c"
                    INCLUDE_DIRS ".")

target_link_libraries(\${COMPONENT_LIB} PRIVATE
${objLines}
)
`;

    await writeFile(path.join(this.mainDir, "CMakeLists.txt"), content);
  }

  private async writeEntryPoint(): Promise<void> {
    const content = `#include <stdio.h>

extern int main(int argc, char **argv);

void app_main(void)
{
  char *argv[] = { "kina_app", NULL };
  
  // Call the Kina runtime entrypoint
  main(1, argv);
}
`;

    await writeFile(path.join(this.mainDir, "kina_entry.c"), content);
  }

  getOutputBinPath(): string {
    return path.join(this.buildDir, "kina_app.bin");
  }
}
