// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"errors"
	"io"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// templigFunctions gives all the functions that are enabled for the templating engine.
func templigFunctions() template.FuncMap {
	result := sprig.TxtFuncMap()
	result["arg"] = argumentValue
	result["hasArg"] = argumentPresent
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

func argumentValue(name string) (any, error) {
	index := slices.IndexFunc(os.Args, func(s string) bool {
		tmp := strings.TrimLeft(s, "-")
		return len(tmp) != len(s) && strings.HasPrefix(tmp, name)
	})

	// handle arguments that give the value using assignment
	if index >= 0 && strings.Contains(os.Args[index], "=") {
		argument := strings.SplitN(os.Args[index], "=", 2)

		if len(argument) == 2 {
			return argument[1], nil
		}
	}

	// handle arguments, with the value in the next argument
	// (that then may not start with a dash)
	if index >= 0 &&
		len(os.Args) > index+1 &&
		!strings.HasPrefix(os.Args[index+1], "-") {

		return os.Args[index+1], nil
	}

	// no argument value given
	return "", nil
}

func argumentPresent(name string) (any, error) {
	index := slices.IndexFunc(os.Args, func(s string) bool {
		tmp := strings.TrimLeft(s, "-")
		return tmp != s && strings.HasPrefix(tmp, name)
	})

	if index >= 0 {
		return true, nil
	}

	return false, nil
}
