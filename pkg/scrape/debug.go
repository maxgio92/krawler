package scrape

import "fmt"

func debug(format string, a interface{}) {
	fmt.Printf("\n" + format, a)
}
