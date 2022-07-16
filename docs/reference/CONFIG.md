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
    mirrors: [{url: "", }]
    repositories: [{name: "", uri: ""}}]
```

> All `versions`, `archs`, `mirrors` are optional fields of the distro configuration.

### Distros

`distros` is a map of well-known supported distro structures. Supported keys are:
- *centos*
 
`distro` structure is a map of `versions`, `archs`, `mirrors`, `repositories`.

##### Example

```
distros:
  centos:
    versions: []
    archs: []
    mirrors: []
    repositories: []
    vars: []
```

### Distro.Versions

`versions` is an array of well-known distribution versions, as named under package repository trees (e.g. [*8-stream*](http://mirrors.edge.kernel.org/centos/8-stream/)).

### Distro.Archs

`archs` is an array of supported architecture IDs. The name follows the one provided by package repository trees.
For example *x86_64*, *aarch86*, *ppc64le*.
 
### Distro.Mirrors

`mirrors` is an array of `mirror` structure, which is a map of:
- `url`

`url` is the root URL of the mirror (e.g. *https://mirrors.kernel.org/centos*).

##### Example

```
centos:
  mirrors:
  - url: https://mirrors.kernel.org/centos
```

### Distro.Repositories

`repositories` is an array of `repository` structure, which in turn is a map of `name` and `uri`.

- `name` is a string label for the name of the repository (e.g. [*AppStream*](http://mirrors.edge.kernel.org/centos/8-stream/AppStream/) for Centos). Please note that this is a label, the value does not have side effects in the crawling flow.
- `uri` is a string that contains the uri path to the packages folder, starting from the root URL of the mirror. Note that the uri format should start with a "/".

##### Example

```
centos:
  repositories:
  - name: AppStream
    uri: /AppStream/x86_64/os/Packages/
```

### Repositories Templating

`uri` field supports templates in the Go template format for annotations that refer to elements of the related distro's data structure (e.g. `distros.centos`). These elements can be both system-declared and user-declared data structures.

##### Supported data types

The supported element types are:
- array of strings

#### System declared variables

- `Distro.Archs`
- `Distro.Versions`

#### User declared variables

You can define your declared variables in `Distro.Vars` structure, which is expected at `distros.<distro>.vars` path.

**Example**

For example, to configure both old and new Centos repositories, given both the archive and current kernel.org mirrors, you can template the repository URLs like below:

```yaml
distros:
  centos:
    archs: ["aarch64", "x86_64"]
    mirrors:
    - url: https://archive.kernel.org/centos-vault/
    - url: https://mirrors.edge.kernel.org/centos/
    repositories:
    - name: old
      uri: "/{{ .old_repos }}/{{ .archs }}/{{ .packages_folder }}/"
    - name: new
      uri: "/{{ .new_repos }}/{{ .archs }}/os/{{ .packages_folder }}/"
    vars:
      new_repos: ["BaseOS", "AppStream"]
      old_repos: ["os", "updates"]
      packages_folder: ["Packages"]
```

As you can see both system-declared (e.g. `archs`) and user-declared (e.g. `new_repos`) data structure can be referenced in the template string.
