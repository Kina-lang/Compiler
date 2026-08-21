package utils

import "strconv"

// Removes the quotes and unescapes a string literal. Returns an error if the literal is not valid.
func GetStringLiteralValue(literal string) (string, error) {
	return strconv.Unquote(literal)
}
