package sem

type SignatureKind string

const (
	PrimitiveTypeSignatureKind SignatureKind = "PRIMITIVE_TYPE"
	FunctionSignatureKind      SignatureKind = "FUNCTION"
)

type BaseSignature struct {
	Kind SignatureKind `json:"kind"`
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
	Name string `json:"name"`
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

func (s *PrimitiveTypeSignature) String() string {
	return s.Name
}

type FunctionSignature struct {
	BaseSignature
	Parameters []TypeSignature `json:"parameters"`
	ReturnType TypeSignature `json:"return"`
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
