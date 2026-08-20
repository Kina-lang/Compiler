package treebuilder

import (
	"encoding/json"

	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
)

type Tree struct {
	Node fileNode
}

func BuildTree(filePath string, tokens []lexer.Token, reporter *diagnostics.Reporter) Tree {
	scanner := NewScanner(filePath, tokens, reporter)

	fileNode := parseFile(scanner)

	return Tree{
		Node: fileNode,
	}
}

func (t *Tree) String() (string, error) {
	json, err := json.MarshalIndent(t, "", "  ")
	return string(json), err
}
