package parcom_test

import (
	"testing"

	"github.com/raviqqe/parcom"
	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	parcom.NewState("foo")
}

func TestStateLine(t *testing.T) {
	assert.Equal(t, 1, parcom.NewState("").Line())
}

func TestStateColumn(t *testing.T) {
	assert.Equal(t, 1, parcom.NewState("").Column())
}

func TestStateWithNewLine(t *testing.T) {
	s := parcom.NewState("\n")
	_, err := s.Char('\n')()

	assert.Nil(t, err)
	assert.Equal(t, 2, s.Line())
	assert.Equal(t, 1, s.Column())
}

func TestStateColumnIncrements(t *testing.T) {
	s := parcom.NewState("foo")
	_, err := s.Str("fo")()

	assert.Nil(t, err)
	assert.Equal(t, 3, s.Column())
}
