package treebuilder

type NodeKind string

const (
	FileNodeKind                NodeKind = "FILE"
	FunctionDeclarationNodeKind NodeKind = "FUNCTION_DECLARATION"
	FunctionParameterNodeKind   NodeKind = "FUNCTION_PARAMETER"
	TypeAnnotationNodeKind      NodeKind = "TYPE_ANNOTATION"
	BasicBlockNodeKind          NodeKind = "BASIC_BLOCK"
)

type Span struct {
	Start int
	End   int
}

type baseNode struct {
	Kind NodeKind `json:"kind"`
	Span Span     `json:"span"`
}

type Node interface{ Base() baseNode }

func (n baseNode) Base() baseNode { return n }

type fileNode struct {
	baseNode
	Children []Node `json:"children"`
}

func NewFileNode(span Span, children []Node) fileNode {
	return fileNode{
		baseNode: baseNode{
			Kind: FileNodeKind,
			Span: span,
		},
		Children: children,
	}
}

type functionDeclarationNode struct {
	baseNode
	Name       string                  `json:"name"`
	Parameters []functionParameterNode `json:"parameters"`
	ReturnType typeAnnotationNode      `json:"returnType"`
	Body       basicBlockNode          `json:"body"`
}

func NewFunctionDeclarationNode(span Span, name string, parameters []functionParameterNode, returnType typeAnnotationNode, body basicBlockNode) functionDeclarationNode {
	return functionDeclarationNode{
		baseNode: baseNode{
			Kind: FunctionDeclarationNodeKind,
			Span: span,
		},
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
}

type functionParameterNode struct {
	baseNode
	Name string             `json:"name"`
	Type typeAnnotationNode `json:"type"`
}

func NewFunctionParameterNode(span Span, name string, typeAnnotation typeAnnotationNode) functionParameterNode {
	return functionParameterNode{
		baseNode: baseNode{
			Kind: FunctionParameterNodeKind,
			Span: span,
		},
		Name: name,
		Type: typeAnnotation,
	}
}

type typeAnnotationNode struct {
	baseNode
	TypeName string `json:"typeName"`
}

func NewTypeAnnotationNode(span Span, typeName string) typeAnnotationNode {
	return typeAnnotationNode{
		baseNode: baseNode{
			Kind: TypeAnnotationNodeKind,
			Span: span,
		},
		TypeName: typeName,
	}
}

type basicBlockNode struct {
	baseNode
	Statements []Node `json:"statements"`
}

func NewBasicBlockNode(span Span, statements []Node) basicBlockNode {
	return basicBlockNode{
		baseNode: baseNode{
			Kind: BasicBlockNodeKind,
			Span: span,
		},
		Statements: statements,
	}
}
