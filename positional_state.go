package parcom

import (
	"errors"
)

type position struct {
	lineIndex, columnIndex int // -1 indicates invalid position.
}

// PositionalState is a position-aware parser state.
type PositionalState struct {
	State
	position position
}

// NewPositionalState creates a parser state.
func NewPositionalState(s string) *PositionalState {
	return &PositionalState{*NewState(s), position{-1, -1}}
}

// WithPosition creates a parser saving a current position.
func (s *PositionalState) WithPosition(p Parser) Parser {
	return func() (interface{}, error) {
		pp := s.position
		s.position = position{s.lineIndex, s.columnIndex}
		defer func() { s.position = pp }()

		return p()
	}
}

// Block creates a parser which parses a block of the second parsers prefixed
// by the first parser.
func (s *PositionalState) Block(p, pp Parser) Parser {
	return s.Prefix(p, s.WithPosition(s.Many1(s.atColumn(pp))))
}

// Indent creates a parser which parses an indent before running a given parser.
// It is equivalent to a given parser and parses no indent if no position is
// saved beforehand.
func (s *PositionalState) Indent(p Parser) Parser {
	return func() (interface{}, error) {
		if s.position.columnIndex >= 0 && s.columnIndex <= s.position.columnIndex {
			return nil, errors.New("invalid indent")
		}

		return p()
	}
}

func (s *PositionalState) atColumn(p Parser) Parser {
	return func() (interface{}, error) {
		if s.columnIndex != s.position.columnIndex {
			return nil, errors.New("invalid indent")
		}

		return p()
	}
}
