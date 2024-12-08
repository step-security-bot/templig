// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

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
