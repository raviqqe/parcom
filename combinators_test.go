package parcom_test

import (
	"fmt"
	"testing"

	"github.com/raviqqe/parcom"
	"github.com/stretchr/testify/assert"
)

func TestChars(t *testing.T) {
	s := parcom.NewState("b")
	x, err := s.Chars("abc")()
	assert.Equal(t, 'b', x)
	assert.Nil(t, err)
}

func TestCharsError(t *testing.T) {
	s := parcom.NewState("d")
	x, err := s.Chars("abc")()

	assert.Nil(t, x)
	assert.Error(t, err)
	assert.Equal(t, 1, err.(parcom.Error).Line())
	assert.Equal(t, 1, err.(parcom.Error).Column())
}

func TestCharsErrorWithEmptySource(t *testing.T) {
	s := parcom.NewState("")
	_, err := s.Chars("a")()

	assert.Error(t, err)
	assert.Equal(t, "unexpected end of source", err.Error())
}

func TestNotChar(t *testing.T) {
	s := parcom.NewState("a")
	x, err := s.NotChar(' ')()
	assert.Equal(t, 'a', x)
	assert.Nil(t, err)
}

func TestNotCharError(t *testing.T) {
	s := parcom.NewState(" ")
	x, err := s.NotChar(' ')()
	assert.Nil(t, x)
	assert.Error(t, err)
}

func TestStrError(t *testing.T) {
	s := parcom.NewState("foe")
	_, err := s.Str("foo")()

	assert.Error(t, err)
	assert.Equal(t, 1, err.(parcom.Error).Line())
	assert.Equal(t, 3, err.(parcom.Error).Column())
}

func TestWrap(t *testing.T) {
	s := parcom.NewState("abc")
	x, err := s.Wrap(s.Str("a"), s.Str("b"), s.Str("c"))()
	assert.Equal(t, "b", x)
	assert.Nil(t, err)
}

func TestPrefix(t *testing.T) {
	s := parcom.NewState("abc")
	x, err := s.Prefix(s.Str("ab"), s.Str("c"))()
	assert.Equal(t, "c", x)
	assert.Nil(t, err)
}

func TestPrefixError(t *testing.T) {
	s := parcom.NewState("abc")
	x, err := s.Prefix(s.Str("ad"), s.Str("c"))()
	assert.Nil(t, x)
	assert.Error(t, err)
}

func TestSuffix(t *testing.T) {
	s := parcom.NewState("abc")
	x, err := s.Suffix(s.Str("ab"), s.Str("c"))()
	assert.Equal(t, "ab", x)
	assert.Nil(t, err)
}

func TestMany(t *testing.T) {
	for _, str := range []string{"", "  "} {
		s := parcom.NewState(str)
		x, err := s.Many(s.Char(' '))()

		t.Logf("%#v", x)

		assert.NotNil(t, x)
		assert.Nil(t, err)
	}
}

func TestManyError(t *testing.T) {
	for _, str := range []string{"="} {
		s := parcom.NewState(str)
		x, err := s.Exhaust(s.Many(func() (interface{}, error) {
			x, err := s.Str("=")()

			if err != nil {
				return nil, err
			}

			if x.(string) == "=" {
				return nil, fmt.Errorf("Invalid word")
			}

			return x, nil
		}))()

		t.Logf("%#v", x)

		assert.Nil(t, x)
		assert.Error(t, err)
	}
}

func testMany1Space(str string) (interface{}, error) {
	s := parcom.NewState(str)
	return s.Many1(s.Char(' '))()
}

func TestMany1(t *testing.T) {
	x, err := testMany1Space(" ")

	t.Logf("%#v", x)

	assert.NotNil(t, x)
	assert.Nil(t, err)
}

func TestMany1Error(t *testing.T) {
	x, err := testMany1Space("")

	t.Log(err)

	assert.Nil(t, x)
	assert.Error(t, err)
}

func TestMany1Nest(t *testing.T) {
	s := parcom.NewState("    ")
	x, err := s.Many1(s.Many1(s.Char(' ')))()

	t.Logf("%#v", x)

	assert.NotNil(t, x)
	assert.Nil(t, err)
}

func testOr(str string) (interface{}, error) {
	s := parcom.NewState(str)
	return s.Or(s.Char('a'), s.Char('b'))()
}

func TestOr(t *testing.T) {
	for _, str := range []string{"a", "b"} {
		x, err := testOr(str)

		t.Logf("%#v", x)

		assert.NotNil(t, x)
		assert.Nil(t, err)
	}
}

func TestOrError(t *testing.T) {
	x, err := testOr("c")

	t.Log(err)

	assert.Nil(t, x)
	assert.Error(t, err)
}

func TestMaybeSuccess(t *testing.T) {
	s := parcom.NewState("foo")
	x, err := s.Maybe(s.Str("foo"))()

	t.Log(x)

	assert.Equal(t, "foo", x)
	assert.Nil(t, err)
}

func TestMaybeError(t *testing.T) {
	s := parcom.NewState("bar")
	x, err := s.Maybe(s.Str("foo"))()

	t.Log(x)

	assert.Nil(t, x)
	assert.Nil(t, err)
}

func TestExhaustError(t *testing.T) {
	s := parcom.NewState("foo bar")
	_, err := s.Exhaust(s.Str("foo"))()

	assert.Error(t, err)
	assert.Equal(t, 1, err.(parcom.Error).Line())
	assert.Equal(t, 4, err.(parcom.Error).Column())
}

func TestExhaustErrorWithMany(t *testing.T) {
	s := parcom.NewState("foofoe")
	_, err := s.Exhaust(s.Many(s.Str("foo")))()

	assert.Error(t, err)
	assert.Equal(t, 1, err.(parcom.Error).Line())
	assert.Equal(t, 4, err.(parcom.Error).Column())
}

func TestExhaustiveMany(t *testing.T) {
	s := parcom.NewState("foofoo")
	x, err := s.ExhaustiveMany(s.Str("foo"))()

	assert.Equal(t, []interface{}{"foo", "foo"}, x)
	assert.Nil(t, err)
}

func TestExhaustErrorWithExhaustiveMany(t *testing.T) {
	s := parcom.NewState("foofoe")
	_, err := s.Exhaust(s.ExhaustiveMany(s.Str("foo")))()

	assert.Error(t, err)
	assert.Equal(t, 1, err.(parcom.Error).Line())
	assert.Equal(t, 6, err.(parcom.Error).Column())
}

func TestExhaustWithErroneousParser(t *testing.T) {
	s := parcom.NewState("")
	_, err := s.Exhaust(s.Str("foo"))()

	assert.Error(t, err)
}

func TestStringify(t *testing.T) {
	str := "foo"
	s := parcom.NewState(str)
	x, err := s.Exhaust(s.Stringify(s.And(s.Str(str))))()
	assert.Equal(t, str, x)
	assert.Nil(t, err)
}

func TestStringifyWithNil(t *testing.T) {
	s := parcom.NewState("")
	x, err := s.Stringify(s.None())()
	assert.Equal(t, "", x)
	assert.Nil(t, err)
}

func TestStringifyPanic(t *testing.T) {
	s := parcom.NewState("")
	assert.Panics(
		t,
		func() {
			s.Stringify(
				s.App(func(interface{}) (interface{}, error) { return struct{}{}, nil }, s.None()),
			)()
		},
	)
}

func TestLazy(t *testing.T) {
	s := parcom.NewState("foo")
	x, err := s.Lazy(func() parcom.Parser { return s.Str("foo") })()
	assert.Equal(t, "foo", x)
	assert.Nil(t, err)
}

func TestVoid(t *testing.T) {
	s := parcom.NewState("foo")
	x, err := s.Void(s.Str("foo"))()
	assert.Nil(t, x)
	assert.Nil(t, err)
}
