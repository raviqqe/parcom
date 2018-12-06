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
	return s.block(s.Many, p, pp)
}

// Block1 is the same as the Block but blocks must have at least one element.
func (s *PositionalState) Block1(p, pp Parser) Parser {
	return s.block(s.Many1, p, pp)
}

func (s *PositionalState) block(m func(Parser) Parser, p, pp Parser) Parser {
	return s.WithPosition(s.Prefix(p, s.SameLineOrIndent(s.WithPosition(m(s.SameColumn(pp))))))
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

// SameLine creates a parser which parses something in the same line.
func (s *PositionalState) SameLine(p Parser) Parser {
	return func() (interface{}, error) {
		if s.position.columnIndex >= 0 && s.lineIndex != s.position.lineIndex {
			return nil, errors.New("should be in the same line")
		}

		return p()
	}
}

// SameLineOrIndent creates a parser which parses something in the same line or indented.
func (s *PositionalState) SameLineOrIndent(p Parser) Parser {
	return s.Or(s.SameLine(p), s.Indent(p))
}

// SameColumn creates a parser which parses something in the same column.
func (s *PositionalState) SameColumn(p Parser) Parser {
	return func() (interface{}, error) {
		if s.columnIndex != s.position.columnIndex {
			return nil, errors.New("invalid indent")
		}

		return p()
	}
}
