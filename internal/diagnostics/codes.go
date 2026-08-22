package diagnostics

var DiagnosticCode string

const (
	InvalidTokenDiagnosticCode  = "E0001"
	InvalidSyntaxDiagnosticCode = "E0002"
	IllegalNodeInTopLevelDiagnosticCode = "E0003"
	SymbolNotFoundDiagnosticCode = "E0004"
	SymbolAlreadyDefinedDiagnosticCode = "E0005"
	TypeMismatchDiagnosticCode = "E0006"
	MissingMainFunctionDiagnosticCode = "E0007"
)
