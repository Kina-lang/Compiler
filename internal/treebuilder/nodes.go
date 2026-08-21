package treebuilder

type NodeKind string

const (
	FileNodeKind NodeKind = "FILE"

	FunctionDeclarationNodeKind NodeKind = "FUNCTION_DECLARATION"
	FunctionParameterNodeKind   NodeKind = "FUNCTION_PARAMETER"

	ImportNode NodeKind = "IMPORT"
	ImportMember NodeKind = "IMPORT_MEMBER"

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
	Statements []StatementNode `json:"statements"`
}

func NewBasicBlockNode(span Span, statements []StatementNode) basicBlockNode {
	return basicBlockNode{
		baseNode: baseNode{
			Kind: BasicBlockNodeKind,
			Span: span,
		},
		Statements: statements,
	}
}

type importNode struct {
	baseNode
	ModuleName string `json:"moduleName"`
	Members   []importMemberNode `json:"members"`
}

func NewImportNode(span Span, moduleName string, members []importMemberNode) importNode {
	return importNode{
		baseNode: baseNode{
			Kind: ImportNode,
			Span: span,
		},
		ModuleName: moduleName,
		Members:    members,
	}
}

type importMemberNode struct {
	baseNode
	Name string `json:"name"`
	Alias string `json:"alias"`
}

func NewImportMemberNode(span Span, name string, alias string) importMemberNode {
	return importMemberNode{
		baseNode: baseNode{
			Kind: ImportMember,
			Span: span,
		},
		Name:  name,
		Alias: alias,
	}
}

type statementNode struct {
	baseNode
}

type StatementNode interface{ Base() statementNode }

func (n statementNode) Base() statementNode { return n }

type returnStatementNode struct {
	statementNode
	Expression ExpressionNode `json:"expression"`
}

func NewReturnStatementNode(span Span, expression ExpressionNode) returnStatementNode {
	return returnStatementNode{
		statementNode: statementNode{
			baseNode: baseNode{
				Kind: ReturnStatementNodeKind,
				Span: span,
			},
		},
		Expression: expression,
	}
}

type expressionNode struct {
	baseNode
}

type ExpressionNode interface{ Base() expressionNode }

func (n expressionNode) Base() expressionNode { return n }

type literalExpressionNode struct {
	expressionNode
	LiteralType LiteralType `json:"literalType"`
	Value       string      `json:"value"`
}

func NewLiteralExpressionNode(span Span, literalType LiteralType, value string) literalExpressionNode {
	return literalExpressionNode{
		expressionNode: expressionNode{
			baseNode: baseNode{
				Kind: LiteralExpressionNodeKind,
				Span: span,
			},
		},
		LiteralType: literalType,
		Value:       value,
	}
}
