# Configuration reference

## Languages

Configuration can be expressed in:
- `json`
- `yaml`

## The structure

```yaml
distros:
  <Distro name>:
    versions: [""]
    archs: [""]
    mirrors: [{url: "", repositories: {name: "", packagesUriTemplate: ""}}]
```

> All `versions`, `archs`, `mirrors` are optional fields of the distro configuration.

### Distros

`distros` is a map of well-known supported distro structures. Supported keys are:
- *centos*
 
`distro` structure is a map of `versions`, `archs`, `mirrors`.

##### Example

```
distros:
  centos:
    versions: []
    archs: []
    mirrors: []
```

### Distro.Versions

`versions` is an array of well-known distribution versions, as named under package repository trees (e.g. [*8-stream*](http://mirrors.edge.kernel.org/centos/8-stream/)).

### Distro.Archs

`archs` is an array of supported architecture IDs. The name follows the one provided by package repository trees.
For example *x86_64*, *aarch86*, *ppc64le*.
 
### Distro.Mirrors

`mirrors` is an array of `mirror` structure, which in turn is a map of:
- `url`
- `repositories`

`url` is the root URL of the mirror (e.g. *https://mirrors.kernel.org/centos*).

`repositories` is an array of `repository` structures.

##### Example

```
centos:
  mirrors:
  - url: https://mirrors.kernel.org/centos
    repositories: []
```

### Mirror.Repositories

`repositoties` is an array of `repository` structure, which in turn is a map of `name` and `packagesUriTemplate`.

- `name` is a string label for the name of the repository (e.g. [*AppStream*](http://mirrors.edge.kernel.org/centos/8-stream/AppStream/) for Centos). Please note that this is a label, the value does not have side effects in the crawling flow.
- `packagesUriTemplate` is a string that contains the URI path to the packages folder, starting from the root URL of the mirror (as defined in `mirror.url`). Note that the URI format should start with a "/".

##### Example

```
mirrors:
- url: https://mirrors.kernel.org/centos/
  repositories:
  - name: AppStream
    packagesUriTemplate: /AppStream/x86_64/os/Packages/
```

### Repositories Templating

`packagesUriTemplate` field supports templates in the Go template format for annotations that refer to elements of the related distro's data structure (e.g. `distros.centos`). These elements can be both system-declared and user-declared data structures.
The supported element types are:
- array of strings

**Example**

For example, to configure both old and new Centos repositories, given both the archive and current kernel.org mirrors, you can template the repository URLs like below:

```yaml
distros:
  centos:
    archs: ["aarch64", "x86_64"]
    new_repos: ["BaseOS", "AppStream"]
    old_repos: ["os", "updates"]
    packages_folder: ["Packages"]
    mirrors:
    - url: https://archive.kernel.org/centos-vault/
      repositories:
      - name: old
        packagesUriTemplate: "/{{ .old_repos }}/{{ .archs }}/{{ .packages_folder }}/"
    - url: https://mirrors.edge.kernel.org/centos/
      repositories:
      - name: new
        packagesUriTemplate: "/{{ .new_repos }}/{{ .archs }}/os/{{ .packages_folder }}/"
```

As you can see both system-declared (e.g. `archs`) and user-declared (e.g. `new_repos`) data structure can be referenced in the template string.
