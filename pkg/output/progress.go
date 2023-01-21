package output

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

type ProgressOptions struct {
	bar *progressbar.ProgressBar
}

func NewProgressOptions(total int, message ...string) *ProgressOptions {
	desc := ""
	if len(message) > 0 {
		desc = message[0]
	}
	bar := progressbar.NewOptions64(
		int64(total),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(false),
	)

	return &ProgressOptions{
		bar: bar,
	}
}

func (b *ProgressOptions) Progress(n int) {
	if b.bar != nil {
		b.bar.Add(n)
	}
}
