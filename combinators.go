package parcom

import (
	"fmt"
	"strings"
)

// Char creates a parser to parse a character.
func (s *State) Char(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() != r {
			return nil, fmt.Errorf("invalid character, '%c'", s.currentRune())
		}

		s.readRune()

		return r, nil
	}
}

// NotChar creates a parser to parse a character which is not the one of an argument.
func (s *State) NotChar(r rune) Parser {
	return s.NotChars(string(r))
}

// Chars creates a parser to parse one of given characters.
func (s *State) Chars(cs string) Parser {
	rs := stringToRuneSet(cs)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; ok {
			defer s.readRune()
			return s.currentRune(), nil
		}

		return nil, fmt.Errorf("invalid character, '%c'", s.currentRune())
	}
}

// NotChars creates a parser to parse a character not in a given string.
func (s *State) NotChars(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; !ok {
			defer s.readRune()
			return s.currentRune(), nil
		}

		return nil, fmt.Errorf("invalid character, '%c'", s.currentRune())
	}
}

// Str creates a parser to parse a string.
func (s *State) Str(str string) Parser {
	rs := []rune(str)
	ps := make([]Parser, 0, len(rs))

	for _, r := range rs {
		ps = append(ps, s.Char(r))
	}

	return s.Stringify(s.And(ps...))
}

// Wrap wraps a parser with parsers which parse something before and after.
// Resulting parsers' parsing results are ones of the middle parsers.
func (s *State) Wrap(l, m, r Parser) Parser {
	return second(s.And(l, m, r))
}

// Prefix creates a parser with 2 parsers which returns the second one's result.
func (s *State) Prefix(pre, p Parser) Parser {
	return second(s.And(pre, p))
}

// Many creates a parser of more than or equal to 0 repetition of a given parser.
func (s *State) Many(p Parser) Parser {
	return func() (interface{}, error) {
		xs := []interface{}{}

		for {
			ss := *s
			x, err := p()

			if err != nil {
				*s = ss
				break
			}

			xs = append(xs, x)
		}

		return xs, nil
	}
}

// Many1 creates a parser of more than 0 repetition of a given parser.
func (s *State) Many1(p Parser) Parser {
	pp := s.Many(p)

	return func() (interface{}, error) {
		x, err := p()

		if err != nil {
			return nil, err
		}

		y, err := pp()

		if err != nil {
			return nil, err
		}

		return append([]interface{}{x}, y.([]interface{})...), nil
	}
}

// Or creates a selectional parser from given parsers.
func (s *State) Or(ps ...Parser) Parser {
	return func() (interface{}, error) {
		err := error(nil)
		ss := *s

		for _, p := range ps {
			x := interface{}(nil)
			x, err = p()

			if err == nil {
				return x, nil
			}

			*s = ss
		}

		return nil, err
	}
}

// And creates a parser which combines given parsers sequentially.
func (s *State) And(ps ...Parser) Parser {
	return func() (interface{}, error) {
		xs := make([]interface{}, 0, len(ps))

		for _, p := range ps {
			x, err := p()

			if err != nil {
				return nil, err
			}

			xs = append(xs, x)
		}

		return xs, nil
	}
}

// Lazy creates a parser which runs a parser created by a given constructor.
// This combinator is useful to define recursive parsers.
func (s *State) Lazy(f func() Parser) Parser {
	p := Parser(nil)

	return func() (interface{}, error) {
		if p == nil {
			p = f()
		}

		return p()
	}
}

// Void creates a parser whose result is always nil but parses something from
// a given parser.
func (State) Void(p Parser) Parser {
	return func() (interface{}, error) {
		_, err := p()
		return nil, err
	}
}

// Exhaust creates a parser which fails when a source string is not exhausted
// after running a given parser. This combinator takes a custom error
// constructor because "not exhausted" errors are always useless.
func (s *State) Exhaust(p Parser, f func(State) error) Parser {
	return func() (interface{}, error) {
		x, err := p()

		if err != nil {
			return nil, err
		} else if !s.exhausted() {
			return nil, f(*s)
		}

		return x, nil
	}
}

// App creates a parser which applies a function to results of a given parser.
func (s *State) App(f func(interface{}) (interface{}, error), p Parser) Parser {
	return func() (interface{}, error) {
		x, err := p()

		if err != nil {
			return nil, err
		}

		return f(x)
	}
}

// None creates a parser which parses nothing and succeeds always.
func (s *State) None() Parser {
	return func() (interface{}, error) {
		return nil, nil
	}
}

// Maybe creates a parser which runs a given parser or parses nothing when it
// fails.
func (s *State) Maybe(p Parser) Parser {
	return s.Or(p, s.None())
}

// Stringify creates a parser which returns a string converted from a result of
// a given parser. The result of a given parser must be a rune, a string or a
// sequence of them in []interface{}.
func (s *State) Stringify(p Parser) Parser {
	return s.App(func(x interface{}) (interface{}, error) { return stringify(x), nil }, p)
}

func stringify(x interface{}) string {
	switch x := x.(type) {
	case nil:
		return ""
	case string:
		return x
	case rune:
		return string(x)
	case []interface{}:
		ss := make([]string, 0, len(x))

		for _, s := range x {
			ss = append(ss, stringify(s))
		}

		return strings.Join(ss, "")
	}

	panic("invalid result type for stringify combinator")
}

func stringToRuneSet(s string) map[rune]bool {
	rs := make(map[rune]bool)

	for _, r := range s {
		rs[r] = true
	}

	return rs
}

func second(p Parser) Parser {
	return func() (interface{}, error) {
		xs, err := p()

		if err != nil {
			return nil, err
		}

		return xs.([]interface{})[1], nil
	}
}
