import type { KinaTypeTokenKind } from "@kina-lang/ast/src/types/types";
import { TokenKind } from "@kina-lang/lexer";
import { KinaAssertionError } from "@kina-lang/utils";

export class TypeTranslator {
  static cToKina(cType: string): KinaTypeTokenKind {
    const normalized = cType.replace(/\s+/g, " ").trim();

    switch (normalized) {
      case "int":
        return TokenKind.TypeInt;
      case "bool":
        return TokenKind.TypeBool;
      case "void":
        return TokenKind.TypeVoid;
      default:
        if (normalized.endsWith("*") || normalized.includes("*"))
          return TokenKind.TypePtr;

        throw new KinaAssertionError(`Unknown C type: ${cType}`);
    }
  }
}
