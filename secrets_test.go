// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig_test

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/AlphaOne1/templig"
)

func TestHideSecrets(t *testing.T) {
	tests := []struct {
		in            any
		want          any
		hideStructure bool
	}{
		{ // 0
			in:            "hello",
			want:          "hello",
			hideStructure: true,
		},
		{ // 1
			in:            []string{"a", "b", "c"},
			want:          []string{"a", "b", "c"},
			hideStructure: true,
		},
		{ // 2
			in: map[string]any{
				"Hello": "World",
			},
			want: map[string]any{
				"Hello": "World",
			},
			hideStructure: true,
		},
		{ // 3
			in:            "secret",
			want:          "secret",
			hideStructure: true,
		},
		{ // 4
			in: map[string]any{
				"secret": "World",
			},
			want: map[string]any{
				"secret": "*****",
			},
			hideStructure: true,
		},
		{ // 5
			in: map[string]any{
				"connections": []any{
					map[string]any{
						"user": "us",
						"pass": "pa",
					},
				},
			},
			want: map[string]any{
				"connections": []any{
					map[string]any{
						"user": "us",
						"pass": "**",
					},
				},
			},
			hideStructure: true,
		},
		{ // 6
			in: map[any]any{
				1: []any{
					map[string]any{
						"user": "us",
						"pass": "pa",
					},
				},
			},
			want: map[any]any{
				1: []any{
					map[string]any{
						"user": "us",
						"pass": "**",
					},
				},
			},
			hideStructure: true,
		},
		{ // 7
			in: map[string]any{
				"secrets": map[string]any{
					"user": "us",
					"pass": "pa",
				},
			},
			want: map[any]any{
				"secrets": "*",
			},
			hideStructure: true,
		},
		{ // 8
			in: map[string]any{
				"connections": map[string]any{
					"user":    "us",
					"secrets": []string{"a", "b", "c"},
				},
			},
			want: map[string]any{
				"connections": map[string]any{
					"user":    "us",
					"secrets": "*",
				},
			},
			hideStructure: true,
		},
		{ // 9
			in: map[string]any{
				"connections": map[string]any{
					"user":    "us",
					"secrets": []string{"a", "bb", "ccc"},
				},
			},
			want: map[string]any{
				"connections": map[string]any{
					"user":    "us",
					"secrets": []string{"*", "**", "***"},
				},
			},
			hideStructure: false,
		},
	}

	gotBuf := bytes.Buffer{}
	wantBuf := bytes.Buffer{}

	for k, test := range tests {
		node := yaml.Node{}
		encodeErr := node.Encode(test.in)

		if encodeErr != nil {
			t.Errorf("%v: could not encode value", k)

			continue
		}

		templig.HideSecrets(&node, test.hideStructure)

		if err := yaml.NewEncoder(&gotBuf).Encode(&node); err != nil {
			t.Errorf("%v: Got error serializing got", k)
		}
		if err := yaml.NewEncoder(&wantBuf).Encode(test.want); err != nil {
			t.Errorf("%v: Got error serializing want", k)
		}

		if gotBuf.String() != wantBuf.String() {
			t.Errorf("%v: got %v\nbut wanted %v", k, gotBuf.String(), wantBuf.String())
		}

		gotBuf.Reset()
		wantBuf.Reset()
	}
}

func TestHideSecretsNil(t *testing.T) {
	var a *yaml.Node = nil

	templig.HideSecrets(a, true)
}

func TestHideSecretAlias(t *testing.T) {
	in := `
open: &ref |-
    value
pass: *ref`
	want := `open: &ref |-
    *****
pass: *ref
`

	node := &yaml.Node{}

	if decodeErr := yaml.NewDecoder(bytes.NewBufferString(in)).Decode(node); decodeErr != nil {
		t.Errorf("unexpted encode error: %v", decodeErr)

		return
	}

	buf := bytes.Buffer{}

	templig.HideSecrets(node, true)

	if encodeErr := yaml.NewEncoder(&buf).Encode(node); encodeErr != nil {
		t.Errorf("could not encode node: %v", encodeErr)
	}

	if buf.String() != want {
		t.Errorf("unexpected output:\n%v\nwanted:\n%v", buf.String(), want)
	}
}
