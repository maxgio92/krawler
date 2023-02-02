package output_test

import (
	"github.com/maxgio92/krawler/pkg/output"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressOptions(t *testing.T) {
	t.Parallel()

	po := output.NewProgressOptions(100)
	assert.NotNil(t, po)
}
