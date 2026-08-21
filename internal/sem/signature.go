package sem

type SignatureKind string

const (
	TypeSignatureKind     SignatureKind = "TYPE"
	FunctionSignatureKind SignatureKind = "FUNCTION"
)

type BaseSignature struct {
	Kind SignatureKind
}

type Signature interface{ Base() BaseSignature
}
func (s BaseSignature) Base() BaseSignature { return s }

type TypeSignature struct {
	BaseSignature
}

func NewTypeSignature() *TypeSignature {
	return &TypeSignature{
		BaseSignature: BaseSignature{
			Kind: TypeSignatureKind,
		},
	}
}

type FunctionSignature struct {
	BaseSignature
	Parameters []TypeSignature
	ReturnType TypeSignature
}

func NewFunctionSignature(parameters []TypeSignature, returnType TypeSignature) *FunctionSignature {
	return &FunctionSignature{
		BaseSignature: BaseSignature{
			Kind: FunctionSignatureKind,
		},
		Parameters: parameters,
		ReturnType: returnType,
	}
}
