# Configuration reference

## Languages

Configuration can be expressed in:
- `json`
- `yaml`

## The structure

```yaml
distros:
  <Distro>:
    versions: <Versions>
    archs: <Archs>
    mirrors:
    - url: "<Mirror root URL>"
      repositories:
        name: "<Package repository name label>"
        packagesUriFormat: "<Packages URI Go string-format>"
```

### Distros

`distros` is a map of well-known supported `distro` structures.

Supported `distro`s are:
- *"centos"*
 
`distro` structure regardless of the distro follows this structure: is a map of `versions`, `archs`, `mirrors`.

### Versions

`versions` is an array of well-known distribution versions, as named under package repository trees (e.g. ["*8-stream*"](http://mirrors.edge.kernel.org/centos/8-stream/)).

### Archs

`archs` is an array of supported architecture IDs. The name follows the one provided by package repository trees.
For example `x86_64`, `aarch86`, `ppc64le`.
 
### Mirrors

`mirrors` is an array of `mirror` structures. In turn `mirror` is a map of:
- `url`
- `repositories`

`url` is the root URL of the mirror (e.g. `https://mirrors.kernel.org/centos`).

`repositories` is an array of `repository` structures.

### Repositories

> Note: You need to know how the repository tree is structured before configuring this.

`repository` is a map of `name` and `packagesUriFormat`.

`name` is a string label for the name of the repository (e.e. [*"AppStream"*](http://mirrors.edge.kernel.org/centos/8-stream/AppStream/) for Centos). Please note that this is a label, the value does not have side effects in the crawling flow.

`packagesUriFormat` is a string that contains a Go string-format of the URI path to the packages folder, starting from the root URL of the mirror (as defined in `mirror.url`).
The format should containe a Go string flag to output the `arch` in the URI tree.

> Note: the URI format should start with a "/".

#### Example

For example to configure the *AppStream* `repository`, given the *"https://mirrors.kernel.org/centos"* `distros.centos.mirrors[0].url`, *"8-stream"* Centos `distros.centos.version`, you can set it *"8-stream/AppStream/%s/os/Packages/"*.

```
distros:
  centos:
    archs: ['aarch64', 'x86_64']
    versions: ['8-stream']
    mirrors:
    - url: https://mirrors.kernel.org
      repositories:
      - name: AppStream
        packageUriFormat: "/8-stream/AppStream/%s/os/Packages/"
```

The string flag (`%s`) will be replaced with each of the configured architectures, as for `distros.centos.archs` value.