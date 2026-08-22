package sem

import "martinpetr.dev/kina/compiler/internal/diagnostics"

func validateMainFunction(table *SymbolTable, reporter *diagnostics.Reporter) bool {
	var wantedSignature = NewFunctionSignature(make([]TypeSignature, 0), NewPrimitiveTypeSignature(string(VoidType)))

	symbol, found := table.Lookup("main")
	if !found {
		reporter.Errorf(-1, -1, diagnostics.MissingMainFunctionDiagnosticCode, "Missing main function")
		return false
	}

	matches := MatchSignature(*symbol.Signature, wantedSignature)
	if !matches {
		reportTypeMismatchError(reporter, symbol.Node, wantedSignature, *symbol.Signature)
		return false
	}

	return true
}

func ValidateAllRules(entryPointPath string, ctxs map[string]*FileContext) bool {
	// Main file validation
	mainCtx := ctxs[entryPointPath]
	ok := validateMainFunction(mainCtx.SymbolTable, mainCtx.Reporter)
	if !ok {
		return false
	}

	return true
}
