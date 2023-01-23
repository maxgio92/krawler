package output

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressOptions struct {
	bar *progressbar.ProgressBar
}

const (
	progressBarThrottleMilliseconds = 65
	progressBarWidth                = 10
	progressBarSpinnerType          = 14
)

func NewProgressOptions(total int, message ...string) *ProgressOptions {
	desc := ""
	if len(message) > 0 {
		desc = message[0]
	}

	bar := progressbar.NewOptions64(
		int64(total),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(progressBarWidth),
		progressbar.OptionThrottle(progressBarThrottleMilliseconds*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(progressBarSpinnerType),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(false),
	)

	return &ProgressOptions{
		bar: bar,
	}
}

func (b *ProgressOptions) Progress(n int) {
	if b.bar != nil {
		//nolint:errcheck
		b.bar.Add(n)
	}
}
