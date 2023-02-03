package output_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxgio92/krawler/pkg/output"
)

func TestNewProgressOptions(t *testing.T) {
	t.Parallel()

	po := output.NewProgressOptions(100)
	assert.NotNil(t, po)
}
