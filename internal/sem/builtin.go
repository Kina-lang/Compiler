package sem

type PrimitiveType string

const (
	IntType    PrimitiveType = "int"
	BoolType   PrimitiveType = "bool"
	StringType PrimitiveType = "string"
	VoidType   PrimitiveType = "void"
	NullType   PrimitiveType = "null"
	AnyType    PrimitiveType = "any"
)

var builtinTypes = map[PrimitiveType]TypeSignature{
	IntType:    NewPrimitiveTypeSignature(string(IntType)),
	BoolType:   NewPrimitiveTypeSignature(string(BoolType)),
	StringType: NewPrimitiveTypeSignature(string(StringType)),
	VoidType:   NewPrimitiveTypeSignature(string(VoidType)),
	NullType:   NewPrimitiveTypeSignature(string(NullType)),
	AnyType:    NewPrimitiveTypeSignature(string(AnyType)),
}

var GlobalSymbolTable = NewSymbolTable(nil)

func init() {
	for name, sig := range builtinTypes {
		symbol := NewSymbol(string(name), GlobalSymbolTable, nil)
		symbol.SetSignature(sig)

		GlobalSymbolTable.Define(symbol)
	}
}
