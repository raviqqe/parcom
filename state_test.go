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

func TestStateLinePosition(t *testing.T) {
	assert.Equal(t, 1, parcom.NewState("").CharacterPosition())
}
