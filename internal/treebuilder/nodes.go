package treebuilder

type NodeKind string

const (
	FileNodeKind NodeKind = "FILE"

	FunctionDeclarationNodeKind NodeKind = "FUNCTION_DECLARATION"
	FunctionParameterNodeKind   NodeKind = "FUNCTION_PARAMETER"

	ImportNodeKind NodeKind = "IMPORT"
	ImportMemberNodeKind NodeKind = "IMPORT_MEMBER"

	TypeAnnotationNodeKind NodeKind = "TYPE_ANNOTATION"

	BasicBlockNodeKind NodeKind = "BASIC_BLOCK"

	StatementNodeKind       NodeKind = "STATEMENT"
	ReturnStatementNodeKind NodeKind = "RETURN_STATEMENT"

	ExpressionNodeKind        NodeKind = "EXPRESSION"
	LiteralExpressionNodeKind NodeKind = "LITERAL_EXPRESSION"
)

type Span struct {
	Start int
	End   int
}

type BaseNode struct {
	Kind NodeKind `json:"kind"`
	Span Span     `json:"span"`
}

type Node interface{ Base() BaseNode }

func (n BaseNode) Base() BaseNode { return n }

type FileNode struct {
	BaseNode
	Children []Node `json:"children"`
}

func NewFileNode(span Span, children []Node) FileNode {
	return FileNode{
		BaseNode: BaseNode{
			Kind: FileNodeKind,
			Span: span,
		},
		Children: children,
	}
}

type FunctionDeclarationNode struct {
	BaseNode
	Name       string                  `json:"name"`
	Parameters []FunctionParameterNode `json:"parameters"`
	ReturnType TypeAnnotationNode      `json:"returnType"`
	Body       BasicBlockNode          `json:"body"`
}

func NewFunctionDeclarationNode(span Span, name string, parameters []FunctionParameterNode, returnType TypeAnnotationNode, body BasicBlockNode) FunctionDeclarationNode {
	return FunctionDeclarationNode{
		BaseNode: BaseNode{
			Kind: FunctionDeclarationNodeKind,
			Span: span,
		},
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
}

type FunctionParameterNode struct {
	BaseNode
	Name string             `json:"name"`
	Type TypeAnnotationNode `json:"type"`
}

func NewFunctionParameterNode(span Span, name string, typeAnnotation TypeAnnotationNode) FunctionParameterNode {
	return FunctionParameterNode{
		BaseNode: BaseNode{
			Kind: FunctionParameterNodeKind,
			Span: span,
		},
		Name: name,
		Type: typeAnnotation,
	}
}

type TypeAnnotationNode struct {
	BaseNode
	TypeName string `json:"typeName"`
}

func NewTypeAnnotationNode(span Span, typeName string) TypeAnnotationNode {
	return TypeAnnotationNode{
		BaseNode: BaseNode{
			Kind: TypeAnnotationNodeKind,
			Span: span,
		},
		TypeName: typeName,
	}
}

type BasicBlockNode struct {
	BaseNode
	Statements []StatementNode `json:"statements"`
}

func NewBasicBlockNode(span Span, statements []StatementNode) BasicBlockNode {
	return BasicBlockNode{
		BaseNode: BaseNode{
			Kind: BasicBlockNodeKind,
			Span: span,
		},
		Statements: statements,
	}
}

type ImportNode struct {
	BaseNode
	ModuleName string `json:"moduleName"`
	Members   []ImportMemberNode `json:"members"`
}

func NewImportNode(span Span, moduleName string, members []ImportMemberNode) ImportNode {
	return ImportNode{
		BaseNode: BaseNode{
			Kind: ImportNodeKind,
			Span: span,
		},
		ModuleName: moduleName,
		Members:    members,
	}
}

type ImportMemberNode struct {
	BaseNode
	Name string `json:"name"`
	Alias string `json:"alias"`
}

func NewImportMemberNode(span Span, name string, alias string) ImportMemberNode {
	return ImportMemberNode{
		BaseNode: BaseNode{
			Kind: ImportMemberNodeKind,
			Span: span,
		},
		Name:  name,
		Alias: alias,
	}
}

type BaseStatementNode struct {
	BaseNode
}

type StatementNode interface{ Base() BaseStatementNode }

func (n BaseStatementNode) Base() BaseStatementNode { return n }

type ReturnStatementNode struct {
	BaseStatementNode
	Expression ExpressionNode `json:"expression"`
}

func NewReturnStatementNode(span Span, expression ExpressionNode) ReturnStatementNode {
	return ReturnStatementNode{
		BaseStatementNode: BaseStatementNode{
			BaseNode: BaseNode{
				Kind: ReturnStatementNodeKind,
				Span: span,
			},
		},
		Expression: expression,
	}
}

type BaseExpressionNode struct {
	BaseNode
}

type ExpressionNode interface{ Base() BaseExpressionNode }

func (n BaseExpressionNode) Base() BaseExpressionNode { return n }

type LiteralExpressionNode struct {
	BaseExpressionNode
	LiteralType LiteralType `json:"literalType"`
	Value       string      `json:"value"`
}

func NewLiteralExpressionNode(span Span, literalType LiteralType, value string) LiteralExpressionNode {
	return LiteralExpressionNode{
		BaseExpressionNode: BaseExpressionNode{
			BaseNode: BaseNode{
				Kind: LiteralExpressionNodeKind,
				Span: span,
			},
		},
		LiteralType: literalType,
		Value:       value,
	}
}
