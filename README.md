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
    <a href="https://www.bestpractices.dev/projects/9789"
       rel="external"
       target="_blank">
        <img src="https://www.bestpractices.dev/projects/9789/badge"
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
    <a href="https://godoc.org/github.com/AlphaOne1/templig"
       rel="external"
       target="_blank">
        <img src="https://godoc.org/github.com/AlphaOne1/templig?status.svg"
             alt="GoDoc Reference">
    </a>
</p>

templig
=======

*templig* (pronounced [ˈtɛmplɪç]) is a non-intrusive configuration library utilizing the text templating engine and the
functions best known from [helm](https://github.com/helm/helm) charts, that originally stem from
[Masterminds/sprig](https://github.com/Masterminds/sprig).
Its primary goal is to enable dynamic configuration files, that have access to the system environment to fill
information using function like `env` and and `read`. To facilitate different environments, overlays can be defined,
that amend a base configuration with environment specific attributes and changes.
Configurations providing the methods for the `Validator` interface also have automated checking on load enabled.

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

The `Get` method gives a pointer to the internally held Config structure that the user supplied. The pointer is always
non-nil, so additional nil-checks are not necessary. Running that program would give:

```text
read errors: <nil>
ID:   23
Name: Interesting Name
```

### Reading with Overlays

Having a base configuration file `my_config.yaml` like the following:

```yaml
id:   23
name: Interesting DevName
```

and a file that contains specific configuration for e.g. the production environment `my_prod_overlay.yaml`:

```yaml
name: Important ProdName
```

The code to read those files would look like this:

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
	c, confErr := templig.FromFile[Config](
		"my_config.yaml",
		"my_prod_overlay.yaml",
	)

	fmt.Printf("read errors: %v", confErr)

	if confErr == nil {
		fmt.Printf("ID:   %v\n", c.Get().ID)
		fmt.Printf("Name: %v\n", c.Get().Name)
	}
}
```

That way, the different configuration files are read in order, with he first one as the base. Every additional file
gives changes to all the ones read before. In this example, changing the name. Running this program would give:

```text
read errors: <nil>
ID:   23
Name: Important ProdName
```

As expected, the value of `Name` was replaced by the one provided in overlay configuration.

### Template functionality
#### Overview

*templig* supports templating the configuration files. On top of the few templating functions provided by the Go
`text/template` library, the functions of [sprig](http://github.com/Masterminds/sprig), that are maybe best known for
their use in [helm](https://github.com/helm/helm) charts. On top of that, the following functions are provided for
convenience:

| Function | Description                                                         | Example                            |
|----------|---------------------------------------------------------------------|------------------------------------|
| arg      | reads the value of the command line argument with the given name    | [Link](examples/templating/arg)    |
| hasArg   | true if an argument with the given name is present, false otherwise | [Link](examples/templating/hasArg) |
| required | checks that its second argument is not zero length or nil           | [Link](examples/templating/env)    |
| read     | reads the content of a file                                         | [Link](examples/templating/read)   |

The expansion of the templated parts is done __before__ overlaying takes place. Any errors of templating will thus be
displayed in their respective source locations.

#### Reading environment

Having a templated configuration file like this one:

```yaml
id:   23
name: Interesting Name
pass: {{ env "PASSWORD" | required "password required" | quote }}
```

or this one;

```yaml
id:   23
name: Interesting Name
pass: {{ read "pass.txt" | required "password required" | quote }}
```

One can see the templating code between the double curly braces `{{` and `}}`. 
The following programm is essentially the same as in the [Simple Case](#simple-case).
It just adds the `pass` field to the configuration: 

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

### Validation

The templating facilities allow also for a wide range of tests, but depend on the configuration file read. As it is
most likely user supplied, possible consistency checks are not reliable in the form of template code.
For this purpose, *templig* also allows for the configuration structure to implement the `Validator` interface.
Implementing types provide a function `Validate` that allows *templig* to check __after__ the configuration was read, if
its structure should be considered valid and report errors accordingly.

```go
package main

import (
    "errors"
    "fmt"

	"github.com/AlphaOne1/templig"
)

type Config struct {
    ID   int    `yaml:"id"`
    Name string `yaml:"name"`
}

// Validate fulfills the Validator interface provided by templig.
// This method is called, if it is defined. It influences the outcome of the configuration reading.
func (c *Config) Validate() error {
    result := make([]error, 0)

	if len(c.Name) == 0 {
		result = append(result, errors.New("name is required"))
	}

	if c.ID < 0 {
		result = append(result, errors.New("id greater than zero is required"))
	}

	return errors.Join(result...)
}

func main() {
	c, confErr := templig.FromFile[Config]("my_config_good.yaml")

	if confErr == nil {
		fmt.Printf("ID:   %v\n", c.Get().ID)
		fmt.Printf("Name: %v\n", c.Get().Name)
	}
}
```

Validation functionality can be as simple as in this example. But as the complexity of the configuration grows,
automated tools to generate the configuration structure and basic consistency checks could be employed. These use
e.g. JSON Schema or its embedded form in OpenAPI 2 or 3.

A non-exhaustive list of these:

* https://github.com/omissis/go-jsonschema (JSON Schema)
* https://github.com/ogen-go/ogen (OpenAPI 3.x)
* https://github.com/go-swagger/go-swagger (OpenAPI 2.0 / Swagger 2.0)
