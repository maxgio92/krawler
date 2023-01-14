package output

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProgressOptions(t *testing.T) {
	po := NewProgressOptions(100)
	assert.NotNil(t, po)
}
