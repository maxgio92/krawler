package output

type ProgressOptions struct {
	InitFunc     func()
	ProgressFunc func()
}

func NewProgressOptions(InitFunc func(), ProgressFunc func()) *ProgressOptions {
	return &ProgressOptions{
		InitFunc:     InitFunc,
		ProgressFunc: ProgressFunc,
	}
}

func (o *ProgressOptions) Init() {
	o.InitFunc()
}

func (o *ProgressOptions) Progress() {
	o.ProgressFunc()
}
