package sem

type SignatureKind string

const (
	PrimitiveTypeSignatureKind SignatureKind = "PRIMITIVE_TYPE"
	FunctionSignatureKind      SignatureKind = "FUNCTION"
)

type BaseSignature struct {
	Kind SignatureKind
}

type Signature interface{ Base() BaseSignature }

func (s BaseSignature) Base() BaseSignature { return s }

type BaseTypeSignature struct {
	BaseSignature
}

// TypeSignature is a Signature usable in a type position
type TypeSignature interface {
	Signature
	isTypeSignature()
}

func (BaseTypeSignature) isTypeSignature() {}

type PrimitiveTypeSignature struct {
	BaseTypeSignature
	Name string
}

func NewPrimitiveTypeSignature(name string) PrimitiveTypeSignature {
	return PrimitiveTypeSignature{
		BaseTypeSignature: BaseTypeSignature{
			BaseSignature: BaseSignature{
				Kind: PrimitiveTypeSignatureKind,
			},
		},
		Name: name,
	}
}

type FunctionSignature struct {
	BaseSignature
	Parameters []TypeSignature
	ReturnType TypeSignature
}

func NewFunctionSignature(parameters []TypeSignature, returnType TypeSignature) FunctionSignature {
	return FunctionSignature{
		BaseSignature: BaseSignature{
			Kind: FunctionSignatureKind,
		},
		Parameters: parameters,
		ReturnType: returnType,
	}
}
