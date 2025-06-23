package templig

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// MergeYAMLNodes merges the content of node `b` into node `a`.
// If `a` contains already an element with the same name and of the same kind as `b`,
// they are merged recursively.
func MergeYAMLNodes(a, b *yaml.Node) (*yaml.Node, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if a.Kind != b.Kind && a.Kind != yaml.AliasNode && b.Kind != yaml.AliasNode {
		return nil, fmt.Errorf("node kind mismatch")
	}

	for b.Kind == yaml.AliasNode {
		b = b.Alias
	}

	var res *yaml.Node
	var resErr error

	switch a.Kind {
	case yaml.DocumentNode:
		res, resErr = mergeDocumentNodes(a, b)
	case yaml.SequenceNode:
		res, resErr = mergeSequenceNodes(a, b)
	case yaml.MappingNode:
		res, resErr = mergeMappingNodes(a, b)
	case yaml.ScalarNode:
		res, resErr = mergeScalarNodes(a, b)
	case yaml.AliasNode:
		res, resErr = mergeAliasNodes(a, b)
	default:
		resErr = fmt.Errorf("unhandled node type %v", a.Kind)
	}

	if len(a.Anchor) > 0 {
		if len(b.Anchor) > 0 && a.Anchor != b.Anchor {
			res = nil
			resErr = fmt.Errorf("unequal named anchors not yet supported (source %v, merge %v)", a.Anchor, b.Anchor)
		}
	} else {
		a.Anchor = b.Anchor
	}

	return res, resErr
}

func mergeAliasNodes(a, b *yaml.Node) (*yaml.Node, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if a.Kind != yaml.AliasNode {
		return nil, fmt.Errorf("source node is not AliasNode")
	}

	tmp := *a.Alias
	tmp.Anchor = ""

	return MergeYAMLNodes(&tmp, b)
}

func mergeDocumentNodes(a, b *yaml.Node) (*yaml.Node, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if a.Kind != yaml.DocumentNode || b.Kind != yaml.DocumentNode {
		return nil, fmt.Errorf("nodes have incompatible kind")
	}

	// this is the top level, the yaml library does not support multiple documents in one file
	if len(a.Content) == 1 && len(b.Content) == 1 {
		ret := *a
		ret.Content = make([]*yaml.Node, 1)

		merged, mergedErr := MergeYAMLNodes(a.Content[0], b.Content[0])

		if mergedErr != nil {
			return nil, mergedErr
		}

		ret.Content[0] = merged

		return &ret, nil
	} else {
		return nil, fmt.Errorf("document node in strange configuration")
	}
}

func mergeScalarNodes(a, b *yaml.Node) (*yaml.Node, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if a.Kind != yaml.ScalarNode || b.Kind != yaml.ScalarNode {
		return nil, fmt.Errorf("nodes have incompatible kind")
	}

	ret := *b

	return &ret, nil
}

func mergeSequenceNodes(a, b *yaml.Node) (*yaml.Node, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if a.Kind != yaml.SequenceNode || b.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("nodes have incompatible kind")
	}

	ret := *a
	ret.Content = make([]*yaml.Node, 0, len(a.Content)+len(b.Content))
	ret.Content = append(ret.Content, a.Content...)
	ret.Content = append(ret.Content, b.Content...)

	return &ret, nil
}

func mergeAddValue(a, key, value *yaml.Node) error {
	var k int
	var v int

	for i := 0; i+1 < len(a.Content); i += 2 {
		k = i
		v = i + 1

		if a.Content[k].Kind == yaml.ScalarNode &&
			a.Content[k].Kind == key.Kind &&
			a.Content[k].Value == key.Value {

			merged, mergedErr := MergeYAMLNodes(a.Content[v], value)

			if mergedErr == nil {
				a.Content[v] = merged
			}

			return mergedErr
		}
	}

	a.Content = append(a.Content, key, value)

	return nil
}

func mergeMappingNodes(a, b *yaml.Node) (*yaml.Node, error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if a.Kind != yaml.MappingNode || b.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("nodes have incompatible kind")
	}

	var keyNode *yaml.Node
	var valueNode *yaml.Node

	ret := *a
	ret.Content = make([]*yaml.Node, 0, len(a.Content)+len(b.Content))
	ret.Content = append(ret.Content, a.Content...)

	for i := 0; i+1 < len(b.Content); i += 2 {
		keyNode = b.Content[i]
		valueNode = b.Content[i+1]

		if mergeErr := mergeAddValue(&ret, keyNode, valueNode); mergeErr != nil {
			return nil, mergeErr
		}
	}

	return &ret, nil
}
