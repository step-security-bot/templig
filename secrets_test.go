// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestHideSecrets(t *testing.T) {
	tests := []struct {
		in   any
		want any
	}{
		{ // 0
			in:   "hello",
			want: "hello",
		},
		{ // 1
			in:   []string{"a", "b", "c"},
			want: []string{"a", "b", "c"},
		},
		{ // 2
			in: map[string]any{
				"Hello": "World",
			},
			want: map[string]any{
				"Hello": "World",
			},
		},
		{ // 3
			in:   "secret",
			want: "secret",
		},
		{ // 4
			in: map[string]any{
				"secret": "World",
			},
			want: map[string]any{
				"secret": "*****",
			},
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
		},
	}

	gotBuf := bytes.Buffer{}
	wantBuf := bytes.Buffer{}

	for k, test := range tests {
		hideSecrets(test.in)

		if err := yaml.NewEncoder(&gotBuf).Encode(test.in); err != nil {
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
	var a *int = nil

	hideSecrets(&a)

	if a != nil {
		t.Errorf("secret hiding changed a nil pointer")
	}
}
