//go:build archlinux

package alpm

import (
	"github.com/pkg/errors"
)

var (
	ErrDirEmpty = errors.New("directory is empty")
)
