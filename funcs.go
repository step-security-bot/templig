// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"errors"
	"io"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// templigFunctions gives all the functions that are enabled for the templating engine.
func templigFunctions() template.FuncMap {
	result := sprig.TxtFuncMap()
	result["required"] = required
	result["read"] = readFile
	return result
}

// required is a template function to indicate that the second argument cannot be empty or nil.
func required(warn string, val any) (any, error) {
	if s, ok := val.(string); val == nil || (ok && s == "") {
		return val, errors.New(warn)
	}

	return val, nil
}

// readFile is a template function to read a file and store its content into a string.
// If the file does not exist, an empty string is generated, facilitating the use of `required` for customized
// user interaction.
func readFile(fileName string) (any, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return "", nil
	}

	defer func() { _ = file.Close() }()

	content, readErr := io.ReadAll(file)

	return string(content), readErr
}
