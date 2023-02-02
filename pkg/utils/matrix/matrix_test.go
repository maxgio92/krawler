package matrix_test

import (
	"testing"

	"github.com/maxgio92/krawler/pkg/utils/matrix"

	"github.com/stretchr/testify/assert"
)

func TestGetColumnOrderedCombinationRows(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		given []matrix.Column

		// The number of combinations.
		want int
	}{
		"single column should return a number of combinations which is equal to the number of the element of the column": {
			given: []matrix.Column{
				{0,
					[]string{
						"a/",
						"b/",
						"c/",
					},
				},
			},
			want: 3,
		},
		"two columns should return a number of combinations equal to the multiplication between the numer of the elements in each column": {
			given: []matrix.Column{
				{0,
					[]string{
						"a/",
						"b/",
						"c/",
					},
				},
				{0,
					[]string{
						"1/",
						"2/",
						"3/",
					},
				},
			},
			want: 9,
		},
		"three columns should return a number of combinations equal to the multiplication between the numer of the elements in each column": {
			given: []matrix.Column{
				{0,
					[]string{
						"a/",
						"b/",
						"c/",
					},
				},
				{0,
					[]string{
						"1/",
						"2/",
						"3/",
					},
				},
				{0,
					[]string{
						"x/",
						"y/",
						"z/",
					},
				},
			},
			want: 27,
		},
	}

	for _, v := range tests {
		combinations, err := matrix.GetColumnOrderedCombinationRows(v.given)
		assert.ErrorIs(t, err, nil)
		assert.Len(t, combinations, v.want)
	}
}
