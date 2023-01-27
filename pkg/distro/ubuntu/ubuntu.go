package ubuntu

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/distro/debian"
)

type Ubuntu struct {
	debian.Debian
}

func (u *Ubuntu) Configure(config distro.Config) error {
	c, err := u.BuildConfig(DefaultConfig, config)
	if err != nil {
		return err
	}

	u.Config = c

	return nil
}
