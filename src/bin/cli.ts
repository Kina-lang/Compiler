#!/usr/bin/env tsx

import path from "path";
import { KinaCompiler } from "../compiler";
import { KinaProjectConfig } from "../project_config";
import { readFile } from "fs/promises";

const rootDir = process.cwd();
const buildDir = path.join(rootDir, "build");
const projectConfigPath = path.join(rootDir, "kina.toml");

const compiler = new KinaCompiler();

const projectConfig = KinaProjectConfig.parse(
  await readFile(projectConfigPath, "utf-8"),
);

await compiler.compile({
  name: projectConfig.package.name,
  version: projectConfig.package.version,
  entry: projectConfig.package.entry,
  rootDir: rootDir,
  buildDir: buildDir,
});
