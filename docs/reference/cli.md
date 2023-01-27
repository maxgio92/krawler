# CLI reference

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

An example of a `yaml`-formatted result entry:

```yml
https://github.com/maxgio92/krawler/blob/main/pkg/kernelrelease/kernelrelease.go#L16
fullversion: 4.18.0
version: 4
patchlevel: 18
sublevel: 0
extraversion: "326"
fullextraversion: -326.el8.x86_64
architecture: x86_64
packagename: kernel-devel
packageurl: https://mirrors.edge.kernel.org/centos/8-stream/BaseOS/x86_64/os/Packages/kernel-devel-4.18.0-326.el8.x86_64.rpm
compilerversion: "80500"
```
