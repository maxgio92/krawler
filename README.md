# Krawler: a kernel releases crawler

> Under development

[![Latest](https://img.shields.io/github/v/release/maxgio92/krawler)](https://github.com/maxgio92/krawler/releases/latest)
[![CI](https://github.com/maxgio92/krawler/actions/workflows/ci.yaml/badge.svg)](https://github.com/maxgio92/krawler/actions/workflows/ci.yaml)
[![Release](https://github.com/maxgio92/krawler/actions/workflows/release.yaml/badge.svg)](https://github.com/maxgio92/krawler/actions/workflows/release.yaml)

A crawler for kernel releases distributed by the major Linux distributions.
It works by scraping mirrors for published kernel headers packages on package mirrors.

```
krawler [options] <command>
```

## Options
- `-c, --config file`: (optional) the config file to customize the list of mirrors to scrape for kernel releases (by default it looks at *$HOME/.krawler.yaml*).
- `-v, --verbosity level`: (optional) the verbosity level (*debug*, *info*, *warn*, *error*, *fatal*, *panic*). By (default *warning*).

## Commands

### `list`|`ls`

List available kernel releases with distributed headers, by Linux distribution.
It returns a list of `kernelRelease` objects. The output format can be specified by flag parameter.

```
krawler [options] list|ls <distribution> [-o <format>] 
```

### Parameters
`distribution`: (**required**) The Linux distribution for which the release has been pubished.
Available distributions:

- amazonlinux
- amazonlinux2
- amazonlinux2022
- centos
- debian
- ubuntu

### Options
`-o, --output format`: (optional) the format of the output of the list of kernel releases (one of *text*, *json* or *yaml*). By default *yaml*.

### Output

The `list`|`ls` command prints on standard ouput a is a list of kernel release objects of type [`KernelRelease`](https://github.com/maxgio92/krawler/blob/main/pkg/kernelrelease/kernelrelease.go#L16).

An example of a `json` result entry:

```
{
  "full_version": "4.18.0",
  "version": 4,
  "patch_level": 18,
  "sublevel": 0,
  "extra_version": "331",
  "full_extra_version": "-331.el8.aarch64",
  "architecture": "aarch64",
  "package_name": "kernel-devel",
  "package_url": "https://mirrors.edge.kernel.org/centos/8-stream/BaseOS/aarch64/os/Packages/kernel-devel-4.18.0-331.el8.aarch64.rpm",
  "compiler_version": "80500"
}
```

## Getting started

Let's imagine you want to list the available CentOS kernel releases, scraping default mirrors. You do it by running:

```
krawler ls centos
```

## Configuration

A configuration lets you configure parameters for the crawling, like the mirrors to scrape.

The default configuration file path is `$HOME/.krawler.yaml`. You can specify a custom path with the `--config` option.

When a configuration is not present, a default configurations for repositories are used (for example [this](https://github.com/maxgio92/krawler/blob/main/pkg/distro/centos/constants.go#L20) is the default for Centos).

For a detailed overview see the [**reference**](docs/reference/CONFIG.md).

Moreover, sample configurations are available [here](./config/samples).

## Roadmap

- [ ] Provide GCC versions for all releases
- [ ] Support new distributions

