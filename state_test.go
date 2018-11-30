package parcom_test

import (
	"testing"

	"github.com/raviqqe/parcom"
	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	parcom.NewState("foo")
}

func TestStateLineNumber(t *testing.T) {
	assert.Equal(t, 1, parcom.NewState("").LineNumber())
}

func TestStateCharacterPosition(t *testing.T) {
	assert.Equal(t, 1, parcom.NewState("").CharacterPosition())
}

func TestStateWithNewLine(t *testing.T) {
	s := parcom.NewState("\n")
	s.Char('\n')()
	assert.Equal(t, 2, s.LineNumber())
	assert.Equal(t, 1, s.CharacterPosition())
}
