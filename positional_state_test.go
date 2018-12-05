package parcom_test

import (
	"testing"

	"github.com/raviqqe/parcom"
	"github.com/stretchr/testify/assert"
)

type state struct {
	*parcom.PositionalState
}

func newState(s string) *state {
	return &state{parcom.NewPositionalState(s)}
}

func (s *state) blanks() parcom.Parser {
	return s.Many(s.Chars(" \t\n"))
}

func (s *state) stripRight(p parcom.Parser) parcom.Parser {
	return s.Suffix(p, s.blanks())
}

func TestPositionalStateBlock(t *testing.T) {
	for _, ss := range []string{
		"foo\n  bar\n  bar",
		"foo",
	} {
		s := newState(ss)
		_, err := s.Exhaust(s.Block(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

		assert.Nil(t, err)
	}
}

func TestPositionalStateBlock1(t *testing.T) {
	s := newState("foo\n  bar\n  bar")
	_, err := s.Exhaust(s.Block1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

	assert.Nil(t, err)
}

func TestPositionalStateBlock1WithNestedBlock1s(t *testing.T) {
	s := newState("foo\n  bar\n  bar\n  bar foo\n       bar\n  bar")
	_, err := s.Exhaust(
		s.Block1(
			s.stripRight(s.Str("foo")),
			s.And(
				s.stripRight(s.Str("bar")),
				s.Maybe(s.Block1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar")))),
			),
		),
	)()

	assert.Nil(t, err)
}

func TestPositionalStateBlock1Error(t *testing.T) {
	s := newState("foo")
	_, err := s.Exhaust(s.Block1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

	assert.Error(t, err)
}

func TestPositionalStateIndent(t *testing.T) {
	s := newState(" foo")
	_, err := s.WithPosition(s.And(s.blanks(), s.Indent(s.Str("foo"))))()

	assert.Nil(t, err)
}

func TestPositionalStateIndentWithoutPosition(t *testing.T) {
	s := newState("foo")
	_, err := s.Indent(s.Str("foo"))()

	assert.Nil(t, err)
}

func TestPositionalStateIndentError(t *testing.T) {
	s := newState("foo")
	_, err := s.WithPosition(s.And(s.blanks(), s.Indent(s.Str("foo"))))()

	assert.Error(t, err)
}

func TestPositionalStateWithOldAndNewPositions(t *testing.T) {
	s := newState("foo foo")
	p := s.Indent(s.stripRight(s.Str("foo")))
	_, err := s.WithOldAndNewPositions(p, p)()

	assert.Nil(t, err)
}

func TestPositionalStateWithOldAndNewPositionsError(t *testing.T) {
	s := newState("bar foo")
	p := s.Indent(s.stripRight(s.Str("foo")))
	_, err := s.WithOldAndNewPositions(p, p)()

	assert.Error(t, err)
}

func TestPositionalStateWithPosition(t *testing.T) {
	s := newState(" foo")
	x, err := s.WithPosition(s.Prefix(s.blanks(), s.Indent(s.Str("foo"))))()

	assert.Equal(t, "foo", x)
	assert.Nil(t, err)
}