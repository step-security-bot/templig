// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Validator is the interface to facility validity checks on configuration types.
type Validator interface {
	// Validate is used to Validate a configuration.
	Validate() error
}

// Config is the generic structure holding the configuration information for the specified type.
type Config[T any] struct {
	node    *yaml.Node
	content T
}

// Get gives a pointer to the deserialized configuration.
func (c *Config[T]) Get() *T {
	return &c.content
}

// overlay is called repeatedly and overlays the current intermediate configuration
// with the content of the given io.Reader.
func (c *Config[T]) overlay(r io.Reader) error {
	a, aErr := fromSingle[yaml.Node](r)

	if aErr != nil {
		return aErr
	}

	if c.node == nil {
		c.node = a.Get()
	} else {
		merged, mergeErr := MergeYAMLNodes(c.node, a.Get())

		if mergeErr != nil {
			return mergeErr
		}

		c.node = merged
	}

	return nil
}

// overlayFile opens a given configuration file and loads it as an intermediate using the overlay function.
func (c *Config[T]) overlayFile(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	return c.overlay(f)
}

// fromSingle reads a configuration from the single given io.Reader and
// runs - if necessary - the contained template functions.
func fromSingle[T any](r io.Reader) (*Config[T], error) {
	var c Config[T]
	fileContent, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var t *template.Template

	if t, err = template.
		New("config").
		Funcs(templigFunctions()).
		Parse(string(fileContent)); err != nil {
		return nil, err
	}

	var b bytes.Buffer

	if err = t.Execute(&b, nil); err != nil {
		return nil, err
	}

	if decodeErr := yaml.NewDecoder(&b).Decode(&c.content); decodeErr != nil {
		return nil, decodeErr
	}

	return &c, nil
}

// Validate checks if the configuration is valid, if the content fulfills the Validator interface.
func (c *Config[T]) Validate() error {
	switch v := any(&c.content).(type) {
	case Validator:
		return v.Validate()
	}

	return nil
}

// From reads a configuration from the given set of io.Reader.
func From[T any](r ...io.Reader) (*Config[T], error) {
	if len(r) == 0 {
		return nil, fmt.Errorf("no configuration readers given")
	}

	var c *Config[T]
	var decodeErr error
	var validateErr error

	if len(r) == 1 {
		// optimize the most common case of a single reader, we do not need to
		// go over the yaml.Node structure first.
		c, decodeErr = fromSingle[T](r[0])
	} else {
		c = new(Config[T])

		for _, v := range r {
			if err := c.overlay(v); err != nil {
				return nil, err
			}
		}

		decodeErr = c.node.Decode(&c.content)

		// cleanup
		c.node = nil
	}

	if decodeErr == nil {
		validateErr = c.Validate()
	}

	if resultErr := errors.Join(decodeErr, validateErr); resultErr != nil {
		return nil, resultErr
	}

	return c, nil
}

// To writes a configuration to the given io.Writer.
func (c *Config[T]) To(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(&c.content)
}

// ToSecretsHidden writes the configuration to the given io.Writer and hides secret values using the [SecretRE].
// Strings are replaced with the number of * corresponding to their length.
// Substructures containing secrets, are replaced with a single '*'.
// The following example
//
//	id: id0
//	secrets:
//	  - secret0
//	  - secret1
//
// thus will be replaced by
//
//	id: id0
//	secrets: *
func (c *Config[T]) ToSecretsHidden(w io.Writer) error {
	var writeErr error = nil
	node := yaml.Node{}

	encodeErr := node.Decode(c.content)

	if encodeErr == nil {
		HideSecrets(&node, true)
		writeErr = yaml.NewEncoder(w).Encode(node)
	}

	return errors.Join(encodeErr, writeErr)
}

// ToSecretsHiddenStructured writes the configuration to the given io.Writer
// and hides secret values using the [SecretRE].
// Strings are replaced with the number of * corresponding to their length.
// Substructures containing secrets, are replaced with a corresponding structure of '*'.
// The following example
//
//	id: id0
//	secrets:
//	  - secret0
//	  - secret1
//
// thus will be replaced by
//
//	id: id0
//	secrets:
//	  - *******
//	  - *******
func (c *Config[T]) ToSecretsHiddenStructured(w io.Writer) error {
	var writeErr error = nil
	node := yaml.Node{}

	encodeErr := node.Decode(c.content)

	if encodeErr == nil {
		HideSecrets(&node, false)
		writeErr = yaml.NewEncoder(w).Encode(node)
	}

	return errors.Join(encodeErr, writeErr)
}

// FromFile loads a series of configuration files. The first file is considered the base, all others are
// loaded on top of that one using the [MergeYAMLNodes] functionality.
func FromFile[T any](paths ...string) (*Config[T], error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no configuration paths given")
	}

	c := new(Config[T])
	var decodeErr error
	var validateErr error

	if len(paths) == 1 {
		// optimize the most common case of a single file, we do not need to
		// go over the yaml.Node structure first.
		f, err := os.Open(paths[0])

		if err != nil {
			return nil, err
		}

		defer func() { _ = f.Close() }()

		c, decodeErr = fromSingle[T](f)
	} else {
		for _, addOn := range paths[0:] {
			aErr := c.overlayFile(addOn)

			if aErr != nil {
				return nil, aErr
			}
		}

		decodeErr = c.node.Decode(&c.content)

		// cleanup
		c.node = nil
	}

	if decodeErr == nil {
		validateErr = c.Validate()
	}

	if resultErr := errors.Join(decodeErr, validateErr); resultErr != nil {
		return nil, resultErr
	}

	return c, nil
}

// FromFiles loads a series of configuration files. The first file is considered the base, all others are
// loaded on top of that one using the [MergeYAMLNodes] functionality.
//
// Deprecated: As of version 0.6.0 this function is deprecated and will be removed in the next major release.
func FromFiles[T any](paths []string) (*Config[T], error) {
	return FromFile[T](paths...)
}

// ToFile saves a configuration to a file with the given name, replacing it in case.
func (c *Config[T]) ToFile(path string) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	return c.To(f)
}
