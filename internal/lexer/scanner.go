package lexer

type Scanner struct {
	bytes []byte
	filePath string

	cursor int
}

func NewScanner(filePath string, src []byte) *Scanner {
	return &Scanner{
		bytes: src,
		filePath: filePath,
	}
}

func (s *Scanner) isEOF(pos int) bool {
	return pos >= len(s.bytes)
}

func (s *Scanner) IsAtEOF() bool {
	return s.isEOF(s.cursor)
}

func (s *Scanner) Peek() byte {
	if s.isEOF(s.cursor) {
		return 0
	}

	return s.bytes[s.cursor]
}

func (s *Scanner) Advance() byte {
	current := s.Peek()

	if !s.isEOF(s.cursor) {
		s.cursor++
	}

	return current
}

func (s *Scanner) AdvanceUntil(predicate func(byte) bool) []byte {
	start := s.cursor

	for !s.isEOF(s.cursor) && !predicate(s.Peek()) {
		s.cursor++
	}

	return s.bytes[start:s.cursor]
}

func (s *Scanner) AdvanceWhile(predicate func(byte) bool) []byte {
	start := s.cursor

	for !s.isEOF(s.cursor) && predicate(s.Peek()) {
		s.cursor++
	}

	return s.bytes[start:s.cursor]
}

func (s *Scanner) PeekAhead(offset int) byte {
	pos := s.cursor + offset

	if s.isEOF(pos) {
		return 0
	}

	return s.bytes[pos]
}

func (s *Scanner) Revert() {
	s.cursor--
}
