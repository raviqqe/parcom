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
		"",
		"foo",
		"foo\nfoo",
	} {
		s := newState(ss)
		_, err := s.Exhaust(s.Block(s.stripRight(s.Str("foo"))))()

		assert.Nil(t, err)
	}
}

func TestPositionalStateBlockError(t *testing.T) {
	s := newState("foo\n foo")
	_, err := s.Exhaust(s.Block(s.stripRight(s.Str("foo"))))()

	assert.Error(t, err)
}

func TestPositionalStateBlock1(t *testing.T) {
	for _, ss := range []string{
		"foo",
		"foo\nfoo",
	} {
		s := newState(ss)
		_, err := s.Exhaust(s.Block(s.stripRight(s.Str("foo"))))()

		assert.Nil(t, err)
	}
}

func TestPositionalStateBlock1Error(t *testing.T) {
	for _, ss := range []string{
		"",
		"foo\n foo",
	} {
		s := newState(ss)
		_, err := s.Exhaust(s.Block1(s.stripRight(s.Str("foo"))))()

		assert.Error(t, err)
	}
}

func TestPositionalStateWithBlock(t *testing.T) {
	for _, ss := range []string{
		"foo\n  bar\n  bar",
		"foo",
	} {
		s := newState(ss)
		_, err := s.Exhaust(s.WithBlock(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

		assert.Nil(t, err)
	}
}

func TestPositionalStateWithBlockErrorWithInvalidBlockIndent(t *testing.T) {
	s := newState("foo\nbar")
	_, err := s.Exhaust(s.WithBlock(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

	assert.Error(t, err)
}

func TestPositionalStateWithBlock1(t *testing.T) {
	s := newState("foo\n  bar\n  bar")
	_, err := s.Exhaust(s.WithBlock1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

	assert.Nil(t, err)
}

func TestPositionalStateWithBlock1WithNestedBlocks(t *testing.T) {
	s := newState("foo\n  bar\n  foo\n   bar\n  bar")
	_, err := s.Exhaust(
		s.WithBlock1(
			s.stripRight(s.Str("foo")),
			s.Or(
				s.stripRight(s.Str("bar")),
				s.WithBlock1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))),
			),
		),
	)()

	assert.Nil(t, err)
}

func TestPositionalStateWithBlock1Error(t *testing.T) {
	s := newState("foo")
	_, err := s.Exhaust(s.WithBlock1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

	assert.Error(t, err)
}

func TestPositionalStateWithBlock1ErrorWithInvalidBlockIndent(t *testing.T) {
	s := newState("foo\nbar")
	_, err := s.Exhaust(s.WithBlock1(s.stripRight(s.Str("foo")), s.stripRight(s.Str("bar"))))()

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

func TestPositionalStateSameLine(t *testing.T) {
	s := newState("foo foo")
	_, err := s.WithPosition(s.And(s.stripRight(s.Str("foo")), s.SameLine(s.Str("foo"))))()

	assert.Nil(t, err)
}

func TestPositionalStateSameLineError(t *testing.T) {
	s := newState("foo\n foo")
	_, err := s.WithPosition(s.And(s.stripRight(s.Str("foo")), s.SameLine(s.Str("foo"))))()

	assert.Error(t, err)
}

func TestPositionalStateSameLineOrIndent(t *testing.T) {
	for _, ss := range []string{"foo foo", "foo\n foo"} {
		s := newState(ss)
		_, err := s.WithPosition(
			s.And(s.stripRight(s.Str("foo")), s.SameLineOrIndent(s.Str("foo"))),
		)()

		assert.Nil(t, err)
	}

}

func TestPositionalStateSameLineOrIndentError(t *testing.T) {
	s := newState("foo\nfoo")
	_, err := s.WithPosition(s.And(s.stripRight(s.Str("foo")), s.SameLineOrIndent(s.Str("foo"))))()

	assert.Error(t, err)
}

func TestPositionalStateSameColumn(t *testing.T) {
	s := newState("foo\nfoo")
	_, err := s.WithPosition(s.And(s.stripRight(s.Str("foo")), s.SameColumn(s.Str("foo"))))()

	assert.Nil(t, err)
}

func TestPositionalStateSameColumnError(t *testing.T) {
	s := newState("foo\n foo")
	_, err := s.WithPosition(s.And(s.stripRight(s.Str("foo")), s.SameColumn(s.Str("foo"))))()

	assert.Error(t, err)
}
