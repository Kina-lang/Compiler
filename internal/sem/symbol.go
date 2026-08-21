package sem

import "encoding/json"

type Symbol struct {
	Name string
	Signature *Signature
}

type SymbolTable struct {
	Symbols map[string]*Symbol
	Parent *SymbolTable
}

type FileContext struct {
	FilePath string
	SymbolTable *SymbolTable
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Symbols: make(map[string]*Symbol),
		Parent: parent,
	}
}

func NewSymbol(name string) *Symbol {
	return &Symbol{
		Name: name,
	}
}

// Tries to find the symbol from the inner-most to the outer-most scope
func (t *SymbolTable) Lookup(name string) (*Symbol, bool) {
	// Check if the current symbol has the symbol
	symbol, found := t.Symbols[name]
	if found {
		return symbol, true
	}

	// Lookup in parent table
	symbol, found = t.Parent.Lookup(name)
	return symbol, found
}

// Defines a new symbol in the current symbol table
func (t *SymbolTable) Define(symbol *Symbol) {
	t.Symbols[symbol.Name] = symbol
}

func (t *SymbolTable) String() (string, error) {
	str, err := json.MarshalIndent(t, "", "  ")
	return string(str), err
}

func (s *Symbol) SetSignature(signature *Signature) {
	s.Signature = signature
}
