package kernelrelease

import "fmt"

var (
	ErrKernelCompilerVersionNotFound = fmt.Errorf("compiler version not found")
	ErrKernelConfigValueNotFound     = fmt.Errorf("the line does not contain the config value")
)
