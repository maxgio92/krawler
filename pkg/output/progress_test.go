package output

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProgressOptions(t *testing.T) {
	po := NewProgressOptions(
		func() {},
		func() {},
	)
	assert.NotNil(t, po)
	assert.NotNil(t, po.InitFunc)
	assert.NotNil(t, po.ProgressFunc)
}
