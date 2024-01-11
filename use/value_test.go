package use

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestIsZero(t *testing.T) {
	assert.EqualValues(t, true, IsZero(time.Time{}))
	assert.EqualValues(t, true, IsZero(""))
	assert.EqualValues(t, false, IsZero(time.Millisecond))
}

func TestIf(t *testing.T) {
	textValue := "value-string"
	assert.EqualValues(t, "value-string", If(strings.HasPrefix(textValue, "value"), textValue, "fallback-string"))

	executionTime := time.Hour * 2
	assert.EqualValues(t, time.Second*30, If(executionTime*2 < time.Now().Add(executionTime).Sub(time.Now()), executionTime, time.Second*30))
}

func TestGetOrDefault(t *testing.T) {
	assert.EqualValues(t, "value-string", GetOrDefault("", "value-string"))
	assert.EqualValues(t, 250, GetOrDefaultFunc(250, func() int {
		return 256
	}))
	assert.EqualValues(t, time.Minute, GetOrDefaultFunc(Zero[time.Duration](), func() time.Duration {
		return time.Minute
	}))
}
