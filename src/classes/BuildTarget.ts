import type { CompilationTarget } from "./targets/_base";
import { XtensaESP32ELFTarget } from "./targets/xtensa-esp32-elf/xtensa-esp32-elf";
import { X86_64UnknownLinuxGnuTarget } from "./targets/x86_64-unknown-linux-gnu";

export class BuildTarget {
  static readonly DEFAULT_TARGET = "x86_64-unknown-linux-gnu" as const;
  static readonly SUPPORTED_TARGETS = [
    "x86_64-unknown-linux-gnu",
    "xtensa-esp32-elf",
  ] as const;

  public static isSupported(target: string): target is BuildTargetType {
    return this.SUPPORTED_TARGETS.includes(target as BuildTargetType);
  }

  public static getTarget(target: BuildTargetType): CompilationTarget {
    switch (target) {
      case "x86_64-unknown-linux-gnu":
        return new X86_64UnknownLinuxGnuTarget();
      case "xtensa-esp32-elf":
        return new XtensaESP32ELFTarget();
    }
  }
}

export type BuildTargetType = (typeof BuildTarget.SUPPORTED_TARGETS)[number];
