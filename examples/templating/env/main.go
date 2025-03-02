// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package main

// This example demonstrates the use of the `env` functions in a templated configuration.
// The use of the `required` function is demonstrated in conjunction with `env`.

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlphaOne1/templig"
)

// Config is the configuration structure that is to be filled by templig.
type Config struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
	Pass string `yaml:"pass"`
}

func main() {
	c, confErr := templig.FromFile[Config]("my_config.yaml")

	fmt.Printf("read errors: %v\n", confErr)

	if confErr == nil {
		fmt.Printf("ID:   %v\n", c.Get().ID)
		fmt.Printf("Name: %v\n", c.Get().Name)
		fmt.Printf("Pass: %v\n", strings.Repeat("*", len(c.Get().Pass)))

		fmt.Println("Config printed by templig with hidden secrets:")
		_ = c.ToSecretsHiddenStructured(os.Stdout)
	}
}
