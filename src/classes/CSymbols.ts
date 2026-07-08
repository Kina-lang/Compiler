import type { KinaTypeTokenKind } from "@kina-lang/ast/src/types/types";
import type { KinaCompiler } from "./KinaCompiler";
import Parser from "tree-sitter";
import C from "tree-sitter-c";
import { readFileSync } from "fs";
import { KinaAssertionError } from "@kina-lang/utils";

export class CSymbols {
  private readonly _compiler: KinaCompiler;

  constructor(compiler: KinaCompiler) {
    this._compiler = compiler;
  }

  // TODO: Refactor
  public get(
    filePath: string,
  ): Record<
    string,
    { returnType: KinaTypeTokenKind; parameterTypes: KinaTypeTokenKind[] }
  > {
    const parser = new Parser();
    // @ts-ignore Module 'tree-sitter-c' has bad typing
    parser.setLanguage(C as any);

    const fileContents = readFileSync(filePath, "utf-8");
    const tree = parser.parse(fileContents);

    const symbols: Record<
      string,
      { returnType: KinaTypeTokenKind; parameterTypes: KinaTypeTokenKind[] }
    > = {};

    for (const node of tree.rootNode.children) {
      if (node.type !== "function_definition") continue;

      const declaratorNode = node.childForFieldName("declarator");
      if (!declaratorNode) continue;

      let currentDecl = declaratorNode;
      while (currentDecl.childForFieldName("declarator")) {
        currentDecl = currentDecl.childForFieldName("declarator")!;
      }

      const functionNameNode = currentDecl;
      if (!functionNameNode || functionNameNode.type !== "identifier") continue;

      const returnTypeNode = node.childForFieldName("type");
      if (!returnTypeNode) continue;

      let functionDeclaratorNode = declaratorNode;
      while (
        functionDeclaratorNode &&
        functionDeclaratorNode.type !== "function_declarator"
      ) {
        functionDeclaratorNode =
          functionDeclaratorNode.childForFieldName("declarator")!;
      }

      const parameterTypeNodes = functionDeclaratorNode
        ?.childForFieldName("parameters")
        ?.namedChildren.filter((n) => n.type === "parameter_declaration");
      if (!parameterTypeNodes) continue;

      const functionName = functionNameNode.text;

      let returnType = returnTypeNode.text;
      if (declaratorNode.type === "pointer_declarator") {
        returnType = `${returnType}*`;
      }

      const parameterTypes = parameterTypeNodes.map((n) => {
        const typeNode = n.childForFieldName("type");
        if (!typeNode)
          throw new KinaAssertionError(
            `Parameter declaration missing type: ${n.text}`,
          );

        const paramDecl = n.childForFieldName("declarator");
        const isPointer = paramDecl && paramDecl.text.startsWith("*");

        return isPointer ? `${typeNode.text}*` : typeNode.text;
      });

      symbols[functionName] = {
        returnType: this._compiler.typeTranslator.cToKina(returnType),
        parameterTypes: parameterTypes.map((t) =>
          this._compiler.typeTranslator.cToKina(t),
        ),
      };
    }

    return symbols;
  }
}
