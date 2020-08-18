package weather

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProvider_WS(t *testing.T) {
	provider, err := GetProvider("weather-stack")
	if err != nil {
		t.Error("unexpected error")
	}

	assert.IsType(t, &weatherStackCli{}, provider)
}

func TestGetProvider_WB(t *testing.T) {
	provider, err := GetProvider("weather-bit")
	if err != nil {
		t.Error("unexpected error")
	}

	assert.IsType(t, &weatherBitCli{}, provider)
}

func TestGetProvider_NotValid(t *testing.T) {
	provider, err := GetProvider("weather-any-other")
	if err == nil {
		t.Error("expected error")
	}

	assert.Nil(t, provider)
}
