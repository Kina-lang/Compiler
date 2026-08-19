package lexer

type TokenType string

const (
	LineCommentToken  TokenType = "COMMENT_LINE"
	BlockCommentToken TokenType = "COMMENT_BLOCK"

	IdentifierToken TokenType = "IDENTIFIER"

	KwFuncToken   TokenType = "KEYWORD_FUNC"
	KwReturnToken TokenType = "KEYWORD_RETURN"

	KwIntToken    TokenType = "KEYWORD_INT"
	KwBoolToken   TokenType = "KEYWORD_BOOL"
	KwStringToken TokenType = "KEYWORD_STRING"
	KwFloatToken  TokenType = "KEYWORD_FLOAT"

	KwTrueToken  TokenType = "KEYWORD_TRUE"
	KwFalseToken TokenType = "KEYWORD_FALSE"

	ParenOpenToken    TokenType = "PAREN_OPEN"
	ParenCloseToken   TokenType = "PAREN_CLOSE"
	BraceOpenToken    TokenType = "BRACE_OPEN"
	BraceCloseToken   TokenType = "BRACE_CLOSE"
	BracketOpenToken  TokenType = "BRACKET_OPEN"
	BracketCloseToken TokenType = "BRACKET_CLOSE"
	CommaToken        TokenType = "COMMA"
	DotToken          TokenType = "DOT"
	SemicolonToken    TokenType = "SEMICOLON"
	ColonToken        TokenType = "COLON"

	IntLiteralToken    TokenType = "LITERAL_INT"
	FloatLiteralToken  TokenType = "LITERAL_FLOAT"
	StringLiteralToken TokenType = "LITERAL_STRING"

	NewlineToken TokenType = "NEWLINE"
	EOFToken     TokenType = "EOF"
)

type Span struct {
	Start int
	End   int
}

type Token struct {
	Kind  TokenType `json:"kind"`
	Value string    `json:"value"`
	Span  Span      `json:"span"`
}
