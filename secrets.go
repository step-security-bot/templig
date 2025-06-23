// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// SecretDefaultRE is the default regular expression used to identify secret values automatically. All input to the
// regular expression is preprocessed using the strings.ToLower function, to not unnecessarily complicate the regexp.
const SecretDefaultRE = "(key)|(secret)|(pass)|(password)|(cert)|(certificate)"

// SecretRE is the regular expression used to identify secret values automatically.
// In case there are different properties to identify secrets, extend it.
var SecretRE = regexp.MustCompile(SecretDefaultRE)

// HideSecrets hides secrets in the given YAML node structure. Secrets are identified using the [SecretRE].
// Depending on the parameter `hideStructure` the structure of the secret is hidden too (`true`) or visible (`false`).
func HideSecrets(n *yaml.Node, hideStructure bool) {
	if n == nil {
		return
	}

	if n.Kind == yaml.MappingNode {
		for i := 0; i < len(n.Content); i += 2 {
			if SecretRE.MatchString(strings.ToLower(n.Content[i].Value)) {
				hideAll(n.Content[i+1], hideStructure)
			} else {
				HideSecrets(n.Content[i+1], hideStructure)
			}
		}
	} else {
		for _, v := range n.Content {
			HideSecrets(v, hideStructure)
		}
	}
}

func hideAll(n *yaml.Node, hideStructure bool) {
	switch n.Kind {
	case yaml.ScalarNode:
		n.Tag = "!!str"
		n.Value = strings.Repeat("*", len(n.Value))
	case yaml.AliasNode:
		if n.Alias != nil {
			hideAll(n.Alias, hideStructure)
		}
	default:
		if hideStructure {
			n.Kind = yaml.ScalarNode
			n.Tag = "!!str"
			n.Value = "*"
			n.Content = nil
		} else {
			for _, v := range n.Content {
				hideAll(v, hideStructure)
			}
		}
	}
}
