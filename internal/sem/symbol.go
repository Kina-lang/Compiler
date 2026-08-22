package sem

import (
	"encoding/json"

	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
)

type Symbol struct {
	Name string
	Signature *Signature
	Table *SymbolTable
	Node treebuilder.Node `json:"-"`
}

type SymbolTable struct {
	Symbols map[string]*Symbol
	Parent *SymbolTable `json:"-"`
}

type FileContext struct {
	FilePath string
	SymbolTable *SymbolTable
	Reporter *diagnostics.Reporter
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Symbols: make(map[string]*Symbol),
		Parent: parent,
	}
}

func NewSymbol(name string, table *SymbolTable, node treebuilder.Node) *Symbol {
	return &Symbol{
		Name: name,
		Table: table,
		Node: node,
	}
}

// Tries to find the symbol from the inner-most to the outer-most scope
func (t *SymbolTable) Lookup(name string) (*Symbol, bool) {
	// Check if the current symbol has the symbol
	symbol, found := t.Symbols[name]
	if found {
		return symbol, true
	}

	hasParent := t.Parent != nil
	if !hasParent {
		return nil, false
	}

	// Lookup in parent table
	symbol, found = t.Parent.Lookup(name)
	return symbol, found
}

// Tries to find the symbol only in the current scope
func (t *SymbolTable) DirectLookup(name string) (*Symbol, bool) {
	symbol, found := t.Symbols[name]
	return symbol, found
}

// Defines a new symbol in the current symbol table
func (t *SymbolTable) Define(symbol *Symbol) bool {
	if _, exists := t.Symbols[symbol.Name]; exists {
		return false
	}

	t.Symbols[symbol.Name] = symbol

	return true
}

func (t *SymbolTable) String() (string, error) {
	str, err := json.MarshalIndent(t, "", "  ")
	return string(str), err
}

func (s *Symbol) SetSignature(signature Signature) {
	s.Signature = &signature
}
