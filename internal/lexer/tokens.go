package lexer

type TokenType string

const (
	LineCommentToken TokenType = "COMMENT_LINE"
	BlockCommentToken TokenType = "COMMENT_BLOCK"
)

type Span struct {
	Start int
	End int
}

type Token struct {
	Kind TokenType `json:"kind"`
	Value string `json:"value"`
	Span Span `json:"span"`
}
