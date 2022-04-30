# [UNDER DEVELOPMENT] Krawler: a kernel releases crawler

![ci workflow](https://github.com/maxgio92/krawler/actions/workflows/ci.yaml/badge.svg)

A crawler for kernel releases distributed by the major Linux distributions.
It works by scraping mirrors for published kernel headers packages on package mirrors.

```
krawler [options] <command>
```

## Options
- `-c, --config file`: (optional) the config file to customize the list of mirrors to scrape for kernel releases (by default it looks at *$HOME/.krawler.yaml*).
- `-v, --verbosity level`: (optional) the verbosity level (*debug*, *info*, *warn*, *error*, *fatal*, *panic*). By (default *warning*).

## Commands

### `list`

List available kernel releases with distributed headers, by Linux distribution.
It returns a list of `kernelRelease` objects. The output format can be specified by flag parameter.

```
krawler [options] list <distribution> [-o <format>] 
```

### Parameters
`distribution`: (**required**) The Linux distribution for which the release has been pubished.
Available distributions:

- centos

### Options
`-o, --output format`: (optional) the format of the output of the list of kernel releases (one of *text*, *json* or *yaml*). By default *yaml*.

### Output

The `list` command prints on standard ouput a is a list of kernel release objects of type [`KernelRelease`](https://github.com/falcosecurity/driverkit/blob/master/pkg/kernelrelease/kernelrelease.go#L13).

An example of a `yaml`-formatted result entry:

```
fullversion: 4.18.0
version: "4"
patchlevel: "18"
sublevel: "0"
extraversion: "348"
fullextraversion: -348.2.1.el8_5.x86_64
```

## Configuration

A configuration lets you configure parameters for the scraping, like the mirrors to scrape for the kernel headers packages. The file is automatically read at the path `$HOME/.krawler.yaml`. Otherwise you can specify the path with the `--config` option.
The configuration follows the structure that you can see below:

```yaml
distros:
  <Distribution name>:
      versions: ["<Distribution version>"]
      archs: ["<Architecture ID>"]
      mirrors:
      - url: "<Mirror root URL>"
        repositories:
          name: "<Package repository name label>"
          packagesUriFormat: "<Packages URI Go string-format>"
```

All fields are optional. When not specified [default configurations](./pkg/scrape/defaults.go) for repositories are used. For the configuration reference see [docs/reference/config.md](docs/reference/CONFIG.md).

Example configurations are available in [config/samples](./config/samples).

## Getting started

Let's imagine you want to list the available CentOS kernel releases, scraping default mirrors. You do it by running:

```
krawler list centos -o yaml
```

### Configuration

#### Mirrors

It is possible to configure the mirrors to scrape for kernel headers.
The supported formats are *json* and *yaml*.

For example, consider if you'd want to scrape only current CentOS kernel releases and to ignore archived ones:

```
cat <<EOF > config.yaml
mirrors:
  centos:
  - https://mirrors.edge.kernel.org/centos/
EOF
krawler list centos -o yaml -c config.yaml
```

## VNEXT

- Support new distributions (Debian, Ubuntu, Fedora, Amazon Linux)
- Support package repositories URI as templates

## Thanks

- Falco [driverkit](https://github.com/falcosecurity/driverkit)
