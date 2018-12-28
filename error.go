package parcom

import "fmt"

// Error is an parsing error.
type Error struct {
	message      string
	line, column int
}

// NewError creates an error.
func NewError(m string, s *State) Error {
	return Error{m, s.Line(), s.Column()}
}

func (e Error) Error() string {
	return e.message
}

// Line returns a line number.
func (e Error) Line() int {
	return e.line
}

// Column returns a column number.
func (e Error) Column() int {
	return e.column
}

func newInvalidCharacterError(s *State) Error {
	if s.currentRune() == 0 {
		return NewError("unexpected end of source", s)
	}

	return NewError(fmt.Sprintf("invalid character '%c'", s.currentRune()), s)
}
