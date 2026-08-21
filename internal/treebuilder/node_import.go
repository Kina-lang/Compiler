package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/lexer"
	"martinpetr.dev/kina/compiler/internal/performance"
)

func ParseImport(scanner *Scanner) (importNode, bool) {
	importKw, ok := scanner.Expect(lexer.KwImportToken)
	if !ok {
		return importNode{}, false
	}

	_, ok = scanner.Expect(lexer.BraceOpenToken)
	if !ok {
		return importNode{}, false
	}

	members, ok := ParseImportMembers(scanner)
	if !ok {
		return importNode{}, false
	}

	_, ok = scanner.Expect(lexer.BraceCloseToken)
	if !ok {
		return importNode{}, false
	}

	_, ok = scanner.Expect(lexer.KwFromToken)
	if !ok {
		return importNode{}, false
	}

	modulePathToken, ok := scanner.Expect(lexer.StringLiteralToken)
	if !ok {
		return importNode{}, false
	}

	// Optional semicolon
	scanner.Match(lexer.SemicolonToken)

	return NewImportNode(Span{
		Start: importKw.Span.Start,
		End: modulePathToken.Span.End,
	}, modulePathToken.Value, members), true
}

func ParseImportMembers(scanner *Scanner) ([]importMemberNode, bool) {
	var members = performance.NewFastArray[importMemberNode](8)

	for !scanner.IsAtEOF() {
		currentToken := scanner.Peek()
		if currentToken.Kind == lexer.BraceCloseToken {
			break
		}

		memberIdentifier, ok := scanner.Expect(lexer.IdentifierToken)
		if !ok {
			return nil, false
		}

		_, hasAlias := scanner.Match(lexer.KwAsToken)
		var alias string = ""

		if hasAlias {
			aliasIdentifier, ok := scanner.Expect(lexer.IdentifierToken)
			if !ok {
				return nil, false
			}

			alias = aliasIdentifier.Value
		}

		maybeCommaToken, hasComma := scanner.Match(lexer.CommaToken)
		var end int = memberIdentifier.Span.End

		if hasComma {
			end = maybeCommaToken.Span.End
		}

		members.Append(NewImportMemberNode(Span{
			Start: memberIdentifier.Span.Start,
			End:   end,
		}, memberIdentifier.Value, alias))

		if !hasComma {
			break
		}
	}

	return members.Items(), true
}
