package lexer

import "martinpetr.dev/kina/compiler/internal/performance"

func (r *LexerResult) InsertSemicolons() LexerResult {
	resultTokens := performance.NewFastArray[Token](len(r.Tokens))

	for i, token := range r.Tokens {
		resultTokens.Append(token)

		switch token.Kind {
		case KwReturnToken:
			// Ignore if EOF
			if i >= len(r.Tokens)-1 {
				continue
			}

			// Insert a semicolon after the return keyword, if the next token is a newline
			nextToken := r.Tokens[i+1]

			if nextToken.Kind == NewlineToken {
				resultTokens.Append(Token{
					Kind:  SemicolonToken,
					Value: ";",
					Span: Span{
						Start: token.Span.End,
						End:   token.Span.End,
					}, // Zero-width
				})
			}
		}
	}

	return LexerResult{
		Tokens: resultTokens.Items(),
	}
}
