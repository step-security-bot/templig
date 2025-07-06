// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// SecretDefaultRE is the default regular expression used to identify secret values automatically. All inputs to the
// regular expressions are preprocessed using the strings.ToLower function, to not unnecessarily complicate the regexp.
const SecretDefaultRE = "(key)|(secret)|(pass)|(password)|(cert)|(certificate)"

// SecretRE is the regular expression used to identify secret values automatically.
// In case there are different properties to identify secrets, extend it.
var SecretRE = regexp.MustCompile(SecretDefaultRE)

// HideSecrets hides secrets in the given YAML node structure. Secrets are identified using the [SecretRE].
// Depending on the parameter `hideStructure`, the structure of the secret is hidden too (`true`) or visible (`false`).
func HideSecrets(node *yaml.Node, hideStructure bool) {
	if node == nil {
		return
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if SecretRE.MatchString(strings.ToLower(node.Content[i].Value)) {
				hideAll(node.Content[i+1], hideStructure)
			} else {
				HideSecrets(node.Content[i+1], hideStructure)
			}
		}
	} else {
		for _, v := range node.Content {
			HideSecrets(v, hideStructure)
		}
	}
}

func hideAll(node *yaml.Node, hideStructure bool) {
	switch node.Kind {
	case yaml.ScalarNode:
		node.Tag = "!!str"
		node.Value = strings.Repeat("*", len(node.Value))
	case yaml.AliasNode:
		if node.Alias != nil {
			hideAll(node.Alias, hideStructure)
		}
	default:
		if hideStructure {
			node.Kind = yaml.ScalarNode
			node.Tag = "!!str"
			node.Value = "*"
			node.Content = nil
		} else {
			for _, v := range node.Content {
				hideAll(v, hideStructure)
			}
		}
	}
}
