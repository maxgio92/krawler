package v2022

import (
	"github.com/maxgio92/krawler/pkg/distro"
	common "github.com/maxgio92/krawler/pkg/distro/amazonlinux"
)

type AmazonLinux struct {
	common.AmazonLinux
}

func (a *AmazonLinux) Configure(config distro.Config) error {
	return a.ConfigureCommon(DefaultConfig, config)
}
