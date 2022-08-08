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

- centos

### Options
`-o, --output format`: (optional) the format of the output of the list of kernel releases (one of *text*, *json* or *yaml*). By default *yaml*.

### Output

The `list`|`ls` command prints on standard ouput a is a list of kernel release objects of type [`KernelRelease`](https://github.com/falcosecurity/driverkit/blob/master/pkg/kernelrelease/kernelrelease.go#L13).

An example of a `yaml`-formatted result entry:

```
fullversion: 4.18.0
version: "4"
patchlevel: "18"
sublevel: "0"
extraversion: "348"
fullextraversion: -348.2.1.el8_5.x86_64
```

## Getting started

Let's imagine you want to list the available CentOS kernel releases, scraping default mirrors. You do it by running:

```
krawler ls centos -o yaml
```

## Configuration

A configuration lets you configure parameters for the crawling, like the mirrors to scrape.

The default configuration file path is `$HOME/.krawler.yaml`. You can specify a custom path with the `--config` option.

When a configuration is not present, the [default configurations](./pkg/scrape/defaults.go) for repositories are used.

For a detailed overview see the [**reference**](docs/reference/CONFIG.md).

Moreover, sample configurations are available [here](./config/samples).

## Roadmap

- Support new distributions (Debian, Ubuntu, Fedora, Amazon Linux)
- Provide GCC versions

## Thanks

- Falco [driverkit](https://github.com/falcosecurity/driverkit)
