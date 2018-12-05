package parcom

// State is a parser state.
type State struct {
	source                              []rune
	sourceIndex, lineIndex, columnIndex int
}

// NewState creates a parser state.
func NewState(s string) *State {
	return &State{[]rune(s), 0, 0, 0}
}

func (s State) exhausted() bool {
	return s.sourceIndex >= len(s.source)
}

func (s State) currentRune() rune {
	if s.exhausted() {
		return '\x00'
	}

	return s.source[s.sourceIndex]
}

func (s *State) readRune() {
	if s.currentRune() == '\n' {
		s.lineIndex++
		s.columnIndex = 0
	} else {
		s.columnIndex++
	}

	s.sourceIndex++
}

// Line returns a current line number.
func (s State) Line() int {
	return s.lineIndex + 1
}

// Column returns a position in a current line.
func (s State) Column() int {
	return s.columnIndex + 1
}
