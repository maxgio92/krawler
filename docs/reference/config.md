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
    mirrors: [{name: "", url: ""}]
    repositories: [{name: "", uri: ""}]
    vars: []
output:
  verbosity: [0-6]
```

> All `versions`, `archs`, `mirrors` are optional fields of the distro configuration.

### Distros

`distros` is a map of well-known supported distro structures.

#### Supported distros

As of now, the supported Linux distributions are:
- *amazonlinux1*
- *amazonlinux2*
- *amazonlinux2022*
- *centos*
- *debian*
- *ubuntu*
- *oracle*
 
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

`archs` is an array of supported architecture IDs.

The name follows the one provided by package repository trees. For example *x86_64*/*amd64*, *aarch86*, *ppc64le*.

> If omitted, **all supported** CPU architectures are selected.
 
### Distro.Mirrors

`mirrors` is an array of `mirror` structure, which is a map of:
- `name` (optional)
- `url`

`name` is a string label for the name of the mirror (e.g. [*Edge*](http://mirrors.edge.kernel.org)). Please note that this is a label, the value does not have side effects in the crawling flow.

`url` is the root URL of the mirror (e.g. *https://mirrors.kernel.org/centos*).

##### Example

```
centos:
  mirrors:
  - url: https://mirrors.kernel.org/centos
```

### Distro.Repositories

`repositories` is an array of `repository` structure, which in turn is a map of:
- `name` (optional)
- `uri`

`name` is a string label for the name of the repository (e.g. [*AppStream*](http://mirrors.edge.kernel.org/centos/8-stream/AppStream/) for Centos). Please note that this is a label, the value does not have side effects in the crawling flow.

`uri` is a string that contains the uri path to the repository root folder, starting from the root URL of the mirror. Note that the uri format should start with a "/".

##### Example

```
centos:
  repositories:
  - name: AppStream
    uri: /AppStream/x86_64/os/
```

### Repositories Templating

`uri` field supports templates in the Go template format for annotations that refer to elements of the related distro's data structure (e.g. `distros.centos`). These elements can be both system-declared and user-declared data structures.

##### Supported data types

The supported element types are:

- array of strings

#### System declared variables

- `Distro.Archs`
- `Distro.Versions`

#### Distro.Vars: User declared variables

You can define your declared variables in `Distro.Vars` structure, which is expected at `distros.<distro>.vars` path.

**Example**

For example, to configure both old and new Centos repositories, given both the archive and current kernel.org mirrors, you can template the repository URLs like below:

```yaml
distros:
  centos:
    archs: ["aarch64", "x86_64"]
    mirrors:
    - name: archive
      url: https://archive.kernel.org/centos-vault/
    - name: edge
    url: https://mirrors.edge.kernel.org/centos/
    repositories:
    - name: old
      uri: "/{{ .old_repos }}/{{ .archs }}/"
    - name: new
      uri: "/{{ .new_repos }}/{{ .archs }}/os/"
    vars:
      new_repos: ["BaseOS", "AppStream"]
      old_repos: ["os", "updates"]
```

As you can see both system-declared (e.g. `archs`) and user-declared (e.g. `new_repos`) data structure can be referenced in the template string.

### Output

`output` is a map of settings for visual output of the commands:
- `verbosity`

#### Verbosity

`verbosity` allows to set the verbosity of the visual output of the commands, through a decimal number from 0 to 6.

It can be set either globally (as [above](#the-structure)), and per `distro`. For example:

```
distros:
  ubuntu:
    mirrors:
    - url: "https://mirrors.edge.kernel.org/ubuntu"
      name: Edge
    - url: "http://security.ubuntu.com/ubuntu"
      name: Security
    output:
      verbosity: 6
```

