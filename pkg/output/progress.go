package output

import (
	"github.com/schollz/progressbar/v3"
)

type ProgressOptions struct {
	bar *progressbar.ProgressBar
}

func NewProgressOptions(total int, message ...string) *ProgressOptions {
	return &ProgressOptions{
		bar: progressbar.Default(int64(total), message...),
	}
}

func (b *ProgressOptions) Progress(n int) {
	b.bar.Add(n)
}
