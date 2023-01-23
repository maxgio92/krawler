package output_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressOptions(t *testing.T) {
	t.Parallel()

	po := NewProgressOptions(100)
	assert.NotNil(t, po)
}
