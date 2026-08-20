package treebuilder

import (
	"martinpetr.dev/kina/compiler/internal/lexer"
	"martinpetr.dev/kina/compiler/internal/performance"
)

func ParseFunctionDeclaration(scanner *Scanner) (functionDeclarationNode, bool) {
	funcKw, ok := scanner.Expect(lexer.KwFuncToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	nameIdentifier, ok := scanner.Expect(lexer.IdentifierToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	_, ok = scanner.Expect(lexer.ParenOpenToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	functionParameters, ok := ParseFunctionParameters(scanner)
	if !ok {
		return functionDeclarationNode{}, false
	}

	_, ok = scanner.Expect(lexer.ParenCloseToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	_, ok = scanner.Expect(lexer.ColonToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	returnTypeNode, ok := ParseTypeAnnotation(scanner)
	if !ok {
		return functionDeclarationNode{}, false
	}

	_, ok = scanner.Expect(lexer.BraceOpenToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	functionBody, ok := ParseBasicBlock(scanner)
	if !ok {
		return functionDeclarationNode{}, false
	}

	closeBraceToken, ok := scanner.Expect(lexer.BraceCloseToken)
	if !ok {
		return functionDeclarationNode{}, false
	}

	return NewFunctionDeclarationNode(Span{
		Start: funcKw.Span.Start,
		End:   closeBraceToken.Span.End,
	}, nameIdentifier.Value, functionParameters, returnTypeNode, functionBody), true
}

func ParseFunctionParameters(scanner *Scanner) ([]functionParameterNode, bool) {
	var parameters = performance.NewFastArray[functionParameterNode](2)

	for !scanner.IsAtEOF() {
		currentToken := scanner.Peek()
		if currentToken.Kind == lexer.ParenCloseToken {
			break
		}

		paramNameToken, ok := scanner.Expect(lexer.IdentifierToken)
		if !ok {
			return nil, false
		}

		_, ok = scanner.Expect(lexer.ColonToken)
		if !ok {
			return nil, false
		}

		paramTypeNode, ok := ParseTypeAnnotation(scanner)
		if !ok {
			return nil, false
		}

		maybeCommaToken := scanner.Peek()
		var hasComma bool = maybeCommaToken.Kind == lexer.CommaToken
		var end int = paramTypeNode.Base().Span.End

		if hasComma {
			scanner.Advance() // Consume the comma token
			end = maybeCommaToken.Span.End
		}

		parameters.Append(NewFunctionParameterNode(Span{
			Start: paramNameToken.Span.Start,
			End:   end,
		}, paramNameToken.Value, paramTypeNode))

		if !hasComma {
			break
		}
	}

	return parameters.Items(), true
}
