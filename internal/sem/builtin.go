package sem

var builtinTypes = map[string]TypeSignature{
	"int":  NewPrimitiveTypeSignature("int"),
	"bool": NewPrimitiveTypeSignature("bool"),
	"string":  NewPrimitiveTypeSignature("string"),
}

var GlobalSymbolTable = NewSymbolTable(nil)

func init() {
	for name, sig := range builtinTypes {
		symbol := NewSymbol(name, GlobalSymbolTable)
		symbol.SetSignature(sig)

		GlobalSymbolTable.Define(symbol)
	}
}
