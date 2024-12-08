package templig

import (
	"errors"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// templigFuncs gives all the functions that are enabled for the templating engine.
func templigFuncs() template.FuncMap {
	result := sprig.TxtFuncMap()
	result["required"] = required
	return result
}

// required is a template function to indicate that the second argument cannot be empty or nil.
func required(warn string, val any) (any, error) {
	if s, ok := val.(string); val == nil || (ok && s == "") {
		return val, errors.New(warn)
	}

	return val, nil
}
