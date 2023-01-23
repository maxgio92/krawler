package kernelrelease_test

import (
	"testing"

	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/deb"

	"gotest.tools/assert"
)

//nolint:funlen
func TestBuildFromPackage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pkg  packages.Package
		want KernelRelease
	}{
		"just kernel version": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "5.5.2",
				Release: "",
				Arch:    "",
			},
			want: KernelRelease{
				Fullversion:      "5.5.2",
				Version:          5,
				PatchLevel:       5,
				Sublevel:         2,
				Extraversion:     "",
				FullExtraversion: "",
				PackageName:      "linux-headers",
				Architecture:     Arch(""),
			},
		},
		"an empty string": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "",
				Release: "",
				Arch:    "",
			},
			want: KernelRelease{
				Fullversion:      "",
				Version:          0,
				PatchLevel:       0,
				Sublevel:         0,
				Extraversion:     "",
				FullExtraversion: "",
				PackageName:      "linux-headers",
				Architecture:     Arch(""),
			},
		},
		"Architecture Linux version": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "6.1.5",
				Release: "arch2-1",
				Arch:    "x86_64",
			},
			want: KernelRelease{
				Fullversion:      "6.1.5",
				Version:          6,
				PatchLevel:       1,
				Sublevel:         5,
				Extraversion:     "arch2-1",
				FullExtraversion: "-arch2-1.x86_64",
				PackageName:      "linux-headers",
				Architecture:     Arch("x86_64"),
			},
		},
		"Debian Jessie version": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "3.16.0",
				Release: "10",
				Arch:    "amd64",
			},
			want: KernelRelease{
				Fullversion:      "3.16.0",
				Version:          3,
				PatchLevel:       16,
				Sublevel:         0,
				Extraversion:     "10",
				FullExtraversion: "-10.amd64",
				PackageName:      "linux-headers",
				Architecture:     Arch("amd64"),
			},
		},
		"Debian Buster version": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "4.19.0",
				Release: "6",
				Arch:    "amd64",
			},
			want: KernelRelease{
				Fullversion:      "4.19.0",
				Version:          4,
				PatchLevel:       19,
				Sublevel:         0,
				Extraversion:     "6",
				FullExtraversion: "-6.amd64",
				PackageName:      "linux-headers",
				Architecture:     Arch("amd64"),
			},
		},
		"Debian version with tilde separator": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "4.9.65",
				Release: "2+grsecunoff1~bpo9+1",
				Arch:    "amd64",
			},
			want: KernelRelease{
				Fullversion:      "4.9.65",
				Version:          4,
				PatchLevel:       9,
				Sublevel:         65,
				Extraversion:     "2",
				FullExtraversion: "-2+grsecunoff1~bpo9+1.amd64",
				PackageName:      "linux-headers",
				Architecture:     Arch("amd64"),
			},
		},
		"Debian version with plus separator": {
			pkg: &deb.Package{
				Name:    "linux-headers",
				Version: "4.19+105",
				Release: "deb10u4~bpo9+1",
				Arch:    "amd64",
			},
			want: KernelRelease{
				Fullversion:      "4.19+105",
				Version:          4,
				PatchLevel:       19,
				Sublevel:         105,
				Extraversion:     "deb10u4",
				FullExtraversion: "-deb10u4~bpo9+1.amd64",
				PackageName:      "linux-headers",
				Architecture:     Arch("amd64"),
			},
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := KernelRelease{}
			err := got.BuildFromPackage(tt.pkg)

			assert.NilError(t, err)
			assert.DeepEqual(t, tt.want, got)
		})
	}
}
