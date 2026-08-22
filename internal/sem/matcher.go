package sem

func MatchSignature(signature Signature, wantedSignature Signature) bool {
	if signature.Base().Kind != wantedSignature.Base().Kind {
		return false
	}

	switch sig := signature.(type) {
	case PrimitiveTypeSignature:
		wantedSig, ok := wantedSignature.(PrimitiveTypeSignature)
		if !ok {
			return false
		}

		return sig.Name == wantedSig.Name
	case FunctionSignature:
		wantedSig, ok := wantedSignature.(FunctionSignature)
		if !ok {
			return false
		}

		if len(sig.Parameters) != len(wantedSig.Parameters) {
			return false
		}

		for i := range sig.Parameters {
			if !MatchSignature(sig.Parameters[i], wantedSig.Parameters[i]) {
				return false
			}
		}

		return MatchSignature(sig.ReturnType, wantedSig.ReturnType)
	default:
		return false
	}
}
