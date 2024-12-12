// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"bytes"
	"errors"
	"io"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Validator is the interface to facility validity checks on configuration types.
type Validator interface {
	// Validate is used to validate a configuration.
	Validate() error
}

// Config is the generic structure holding the configuration information for the specified type.
type Config[T any] struct {
	content T
}

// Get gives a pointer to the deserialized configuration.
func (c *Config[T]) Get() *T {
	return &c.content
}

// From reads a configuration from the given io.Reader.
func From[T any](r io.Reader) (*Config[T], error) {
	var c Config[T]
	fileContent, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var t *template.Template

	if t, err = template.
		New("config").
		Funcs(templigFuncs()).
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

	switch v := any(&c.content).(type) {
	case Validator:
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// To writes a configuration to the given io.Writer.
func (c *Config[T]) To(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(&c.content)
}

// ToSecretsHidden writes the configuration to the given io.Writer and hides secret values using the SecretRE.
// Strings are replaced with the number of * corresponding to their length.
// Substructures (e.g. `"secrets": [ "a", "b", "c"]`) containing secrets, are replaced with a single *.
func (c *Config[T]) ToSecretsHidden(w io.Writer) error {
	tmpMap := make(map[string]any)
	data := bytes.Buffer{}

	err0 := yaml.NewEncoder(&data).Encode(&c.content)
	err1 := yaml.NewDecoder(&data).Decode(tmpMap)

	hideSecrets(tmpMap)

	return errors.Join(err0, err1, yaml.NewEncoder(w).Encode(tmpMap))
}

// FromFile loads a configuration from a file with the given name.
func FromFile[T any](path string) (*Config[T], error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() { _ = f.Close() }()

	return From[T](f)
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
