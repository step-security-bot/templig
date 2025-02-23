package templig

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		want    string
		wantErr bool
	}{
		{ // 0
			name:    "sequence concat",
			a:       `["a", 2, 3, 4]`,
			b:       `["b", 5, 6]`,
			want:    `["a", 2, 3, 4, "b", 5, 6]`,
			wantErr: false,
		},
		{ // 1
			name:    "mapping concat",
			a:       `{"a": 2, "b": 3, "c": 4}`,
			b:       `{"d": 5, "e": 6}`,
			want:    `{"a": 2, "b": 3, "c": 4, "d": 5, "e": 6}`,
			wantErr: false,
		},
		{ // 2
			name:    "mapping concat with overlap",
			a:       `{"a": 2, "b": 3, "c": 4}`,
			b:       `{"b": 5, "c": 6, "d": 7}`,
			want:    `{"a": 2, "b": 5, "c": 6, "d": 7}`,
			wantErr: false,
		},
		{ // 3
			name:    "mapping recursive",
			a:       `{"a": {"a": 1, "b": 2}}`,
			b:       `{"a": {"a": 2, "c": 3}}`,
			want:    `{"a": {"a": 2, "b": 2, "c": 3}}`,
			wantErr: false,
		},
		{ // 4
			name:    "mapping recursive",
			a:       `{"a": {"a": 1, "b": 2}}`,
			b:       `{"a": ["a", 2]}`,
			want:    ``,
			wantErr: true,
		},
		{ // 5
			name: "alias contained",
			a: `
a: &ref
    a: 3
b: *ref`,
			b: `
a:
    a: 4`,
			want: `a: &ref
    a: 4
b: *ref`,
			wantErr: false,
		},
		{ // 6
			name: "alias modified",
			a: `
a: &ref
    a: 3
b: *ref`,
			b: `
b:
    a: 5`,
			want: `a: &ref
    a: 3
b:
    a: 5`,
			wantErr: false,
		},
		{ // 7
			name: "alias in merge",
			a: `
a: &ref
    a: 3
b: *ref`,
			b: `
x-blah: &blah
    a: 5
b: *blah`,
			want: `a: &ref
    a: 3
b:
    a: 5
x-blah: &blah
    a: 5`,
			wantErr: false,
		},
		{ // 8
			name: "anchor mismatch",
			a: `
a: &ref0
    a: 3`,
			b: `
a: &ref1
    a: 3`,
			want:    ``,
			wantErr: true,
		},
	}

	for k, v := range tests {
		var an yaml.Node
		var bn yaml.Node

		v.want = v.want + "\n"

		if aErr := yaml.Unmarshal([]byte(v.a), &an); aErr != nil {
			t.Errorf("%v - %v: could not unmarshal a: %v", k, v.name, aErr)
		}

		if bErr := yaml.Unmarshal([]byte(v.b), &bn); bErr != nil {
			t.Errorf("%v - %v: could not unmarshal b: %v", k, v.name, bErr)
		}

		result, resultErr := MergeYAMLNodes(&an, &bn)

		if !v.wantErr && resultErr != nil {
			t.Errorf("%v - %v: could not merge: %v", k, v.name, resultErr)
			continue
		}

		if v.wantErr && resultErr == nil {
			t.Errorf("%v - %v: should not have been able to merge", k, v.name)
			continue
		}

		if resultErr != nil {
			continue
		}

		buf := bytes.Buffer{}

		buf.Reset()

		if err := yaml.NewEncoder(&buf).Encode(result.Content[0]); err != nil {
			t.Errorf("%v - %v: could not marshal document to yaml: %v",
				k, v.name, err)
		}

		if buf.String() != v.want {
			t.Errorf("%v - %v: result mismatch, wanted:\n%v\nbut got:\n%v",
				k, v.name, v.want, buf.String())
		}
	}
}

func TestNullArgs(t *testing.T) {
	mergeFuncs := []func(*yaml.Node, *yaml.Node) (*yaml.Node, error){
		MergeYAMLNodes,
		mergeAliasNodes,
		mergeDocumentNodes,
		mergeMappingNodes,
		mergeScalarNodes,
		mergeSequenceNodes,
	}

	for k, v := range mergeFuncs {
		res, resErr := v(nil, nil)

		if res != nil {
			t.Errorf("%v: merging nil nodes must produce nil", k)
		}

		if resErr == nil {
			t.Errorf("%v: expected an error merging nil nodes", k)
		}
	}
}

func TestMismatchArgs(t *testing.T) {
	mergeFuncs := []func(*yaml.Node, *yaml.Node) (*yaml.Node, error){
		MergeYAMLNodes,
		mergeAliasNodes,
		mergeDocumentNodes,
		mergeMappingNodes,
		mergeScalarNodes,
		mergeSequenceNodes,
	}

	a := yaml.Node{
		Kind: yaml.MappingNode,
	}
	b := yaml.Node{
		Kind: yaml.DocumentNode,
	}

	for k, v := range mergeFuncs {
		res, resErr := v(&a, &b)

		if res != nil {
			t.Errorf("%v: merging kind-mismatched nodes must produce nil", k)
		}

		if resErr == nil {
			t.Errorf("%v: expected an error merging kind-mismatched nodes", k)
		}
	}
}

func TestStrangeDocumentNodes(t *testing.T) {
	a := yaml.Node{
		Kind: yaml.DocumentNode,
	}
	b := yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			&yaml.Node{
				Kind: yaml.ScalarNode,
			},
		},
	}

	res, resErr := MergeYAMLNodes(&a, &b)

	if res != nil {
		t.Errorf("merging strange document nodes must produce nil")
	}

	if resErr == nil {
		t.Errorf("expected an error merging strange document nodes")
	}
}

func TestUnknownNodeKind(t *testing.T) {
	a := yaml.Node{
		Kind: 0,
	}
	b := yaml.Node{
		Kind: 0,
	}

	res, resErr := MergeYAMLNodes(&a, &b)

	if res != nil {
		t.Errorf("merging unknown node kind must produce nil")
	}

	if resErr == nil {
		t.Errorf("expected an error merging unknown node kind")
	}
}
