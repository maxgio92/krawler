# Quickstart

Below an example `main.go`:

```
package main

import (
    "fmt"
    "github.com/maxgio92/krawler/pkg/matrix"
)

var (
	columns = []matrix.Column{
		{OrdinateIndex: 0, Points: []string{"A", "B"}},                // part
		{OrdinateIndex: 0, Points: []string{"1", "2", "3", "4", "5"}}, // part
		{OrdinateIndex: 0, Points: []string{"w", "x", "y", "z"}},      // part
		{OrdinateIndex: 0, Points: []string{"E", "F", "G", "H"}},      // part
		{OrdinateIndex: 0, Points: []string{"A", "B"}},                // part
	}
)

func main() {
	for _, v := range columns {
		fmt.Println(v.Points)
	}
	combinations, err := matrix.GetColumnOrderedCombinationRows(columns)
	if err != nil {
		panic(err)
	}
	fmt.Println(combinations)
}
```