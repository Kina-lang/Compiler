import { spawn } from "child_process";

export class CommandRunner {
  constructor() {}

  static async runWithPipe(
    command: string,
    args: string[],
    input: string,
  ): Promise<string> {
    return await new Promise<string>((resolve, reject) => {
      const proc = spawn(command, args, {
        stdio: ["pipe", "pipe", "inherit"],
        cwd: process.cwd(),
        env: process.env,
      });

      let stdout = "";
      proc.stdout.on("data", (data) => {
        stdout += data.toString();
      });

      proc.on("spawn", () => {
        proc.stdin.write(input);
        proc.stdin.end();
      });

      proc.on("close", (code) => {
        if (code === 0) resolve(stdout);
        else reject(new Error(`Command failed with exit code ${code}`));
      });
    });
  }

  static async run(
    command: string,
    args: string[],
    opts?: { cwd?: string; env?: NodeJS.ProcessEnv },
  ): Promise<void> {
    return await new Promise<void>((resolve, reject) => {
      const proc = spawn(command, args, {
        stdio: "inherit",
        cwd: opts?.cwd || process.cwd(),
        env: opts?.env || process.env,
      });

      proc.on("close", (code) => {
        if (code === 0) resolve();
        else reject(new Error(`Command failed with exit code ${code}`));
      });
    });
  }
}
