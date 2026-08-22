package sem

import (
	"fmt"
	"strings"
)

func BeautifySignature(signature Signature) string {
	switch sig := signature.(type) {
		case PrimitiveTypeSignature:
			return sig.Name
		case FunctionSignature:
			var parameters = make([]string, 0, len(sig.Parameters))

			for _, param := range sig.Parameters {
				parameters = append(parameters, BeautifySignature(param))
			}

			return fmt.Sprintf("(%s) -> %s", strings.Join(parameters, ", "), BeautifySignature(sig.ReturnType))
		default:
			return "<unknown>"
	}
}
