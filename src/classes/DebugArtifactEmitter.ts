import { existsSync } from "fs";
import { mkdir, writeFile } from "fs/promises";
import path from "path";

export class DebugArtifactEmitter {
  private readonly _artifacts: {
    name: string;
    type: string;
    content: string;
    fileType: string;
  }[] = [];

  constructor() {}

  public add(
    name: string,
    type: string,
    content: string,
    fileType: string,
  ): void {
    this._artifacts.push({ name, type, content, fileType });
  }

  public async emit(buildDir: string): Promise<void> {
    for (const artifact of this._artifacts) {
      const artifactPath = path.join(
        buildDir,
        "debug",
        artifact.type,
        artifact.name.replaceAll("_", "__").replaceAll("/", "_") +
          artifact.fileType,
      );
      const artifactDir = path.dirname(artifactPath);

      if (!existsSync(artifactDir))
        await mkdir(artifactDir, { recursive: true });

      await writeFile(artifactPath, artifact.content, "utf-8");
    }
  }
}
