// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"regexp"
	"strings"
)

// SecretDefaultRE is the default regular expression used to identify secret values automatically.
const SecretDefaultRE = "(key)|(secret)|(pass)|(password)|(cert)|(certificate)"

// SecretRE is the regular expression used to identify secret values automatically.
// In case there are different properties to identify secrets, extend it.
var SecretRE = regexp.MustCompile(SecretDefaultRE)

// hideSecretsStringMap hides secrets inside a map of string keys. It is the only place secrets can
// be detected due to their naming inside the map.
func hideSecretsStringMap(m map[string]any) {
	for k, v := range m {
		if SecretRE.MatchString(strings.ToLower(k)) {
			switch vt := v.(type) {
			case string:
				m[k] = strings.Repeat("*", len(vt))
			default:
				m[k] = "*"
			}
		} else {
			hideSecrets(v)
		}
	}
}

// hideSecretsAnyMap hides secrets inside a non-string keyed map. That means, that secrets can
// only be inside the substructure, or they would be hidden by the hideSecretsStringMap function.
func hideSecretsAnyMap(m map[any]any) {
	for _, v := range m {
		hideSecrets(v)
	}
}

// hideSecrets hides secrets inside the value given. It analyses the basic structure and decides which
// one of the specialized hiding functions to apply.
func hideSecrets(a any) {
	switch ta := a.(type) {
	case map[string]any:
		hideSecretsStringMap(ta)
	case map[any]any:
		hideSecretsAnyMap(ta)
	case []any:
		hideSecretsSlice(ta)
	}
}

// hideSecretsSlice hides secrets inside the given slice. It applies the hideSecrets function on each value.
func hideSecretsSlice(s []any) {
	for _, v := range s {
		hideSecrets(v)
	}
}
