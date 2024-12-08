<p align="center">
    <img src="templig_logo.svg" width="25%" alt="Logo"><br>
    <a href="https://github.com/AlphaOne1/templig/actions/workflows/test.yml"
       rel="external"
       target="_blank">
        <img src="https://github.com/AlphaOne1/templig/actions/workflows/test.yml/badge.svg"
             alt="Test Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/templig/actions/workflows/codeql.yml"
       rel="external"
       target="_blank">
        <img src="https://github.com/AlphaOne1/templig/actions/workflows/codeql.yml/badge.svg"
             alt="CodeQL Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/templig/actions/workflows/security.yml"
       rel="external"
       target="_blank">
        <img src="https://github.com/AlphaOne1/templig/actions/workflows/security.yml/badge.svg"
             alt="Security Pipeline Result">
    </a>
    <a href="https://goreportcard.com/report/github.com/AlphaOne1/templig"
       rel="external"
       target="_blank">
        <img src="https://goreportcard.com/badge/github.com/AlphaOne1/templig"
             alt="Go Report Card">
    </a>
    <a href="https://codecov.io/github/AlphaOne1/templig"
       rel="external"
       target="_blank">
        <img src="https://codecov.io/github/AlphaOne1/templig/graph/badge.svg?token=P18EOCUPU8"
             alt="Code Coverage">
    </a>
    <a href="https://www.bestpractices.dev/projects/9251"
       rel="external"
       target="_blank">
        <img src="https://www.bestpractices.dev/projects/9251/badge"
             alt="OpenSSF Best Practises">
    </a>
    <a href="https://scorecard.dev/viewer/?uri=github.com/AlphaOne1/templig"
       rel="external"
       target="_blank">
        <img src="https://api.scorecard.dev/projects/github.com/AlphaOne1/templig/badge"
             alt="OpenSSF Scorecard">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Ftemplig?ref=badge_shield&issueType=license"
       rel="external"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Ftemplig.svg?type=shield&issueType=license"
            alt="FOSSA Status">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Ftemplig?ref=badge_shield&issueType=security" 
       rel="external"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Ftemplig.svg?type=shield&issueType=security"
             alt="FOSSA Status">
    </a>
    <a href="http://godoc.org/github.com/AlphaOne1/templig"
       rel="external"
       target="_blank">
        <img src="https://godoc.org/github.com/AlphaOne1/templig?status.svg"
             alt="GoDoc Reference">
    </a>
</p>

```
 ._^_.
_{{_}}_
```

templig
=======

*templig* is configuration library utilizing the text templating engine and the functions best known from helm charts,
that originally stem from [Masterminds/sprig](http://github.com/Masterminds/sprig/v3).
Its primary goal is to enable access to the system environment to fill information using the `env` function. It also
enables to include verifications inside the configuration.

Usage
-----

### Simple Case

Having a configuration file like the following:

```yaml
id:   23
name: Interesting Name
```

The code to read that file would look like this:

```go
package main

import (
	"fmt"
	"github.com/AlphaOne1/templig"
)

type Config struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
}

func main() {
	c, confErr := templig.FromFile[Config]("my_config.yaml")

	fmt.Printf("read errors: %v", confErr)

	if confErr == nil {
		fmt.Printf("ID:   %v\n", c.Get().ID)
		fmt.Printf("Name: %v\n", c.Get().Name)
	}
}
```

The `Get` method gives a pointer to the internally held Config strucutre that the use supplied. The pinter is always
non-nil, so additional nil-checks are not necessary.

### Advanced Case

Having a templated configuration file like this one:

```yaml
id:   23
name: Interesting Name
pass: {{ env "PASSWORD" | required "password required" | quote }}
```

As demonstrated, one can use the templating functionality that is best known from helm charts. The functions provided
come from the aforementionened [strig](http://github.com/Masterminds/strig/v3)-library.

```go
package main

import (
	"fmt"
	"strings"

	"github.com/AlphaOne1/templig"
)

type Config struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
	Pass string `yaml:"pass"`
}

func main() {
	c, confErr := templig.FromFile[Config]("my_config.yaml")

	fmt.Printf("read errors: %v", confErr)

	if confErr == nil {
		fmt.Printf("ID:   %v\n", c.Get().ID)
		fmt.Printf("Name: %v\n", c.Get().Name)
		fmt.Printf("Pass: %v\n", strings.Repeat("*", len(c.Get().Pass)))
	}
}
```