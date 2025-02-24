// Copyright the templig contributors.
// SPDX-License-Identifier: MPL-2.0

package templig

import (
	"bytes"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
)

type TestConn struct {
	URL    string   `yaml:"url"`
	Passes []string `yaml:"passes"`
}

type TestConfig struct {
	ID   int       `yaml:"id"`
	Name string    `yaml:"name"`
	Conn *TestConn `yaml:"conn,omitempty"`
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		in      string
		inFile  string
		env     map[string]string
		want    TestConfig
		wantErr bool
	}{
		{ // 0
			inFile:  "testData/test_empty.yaml",
			want:    TestConfig{},
			wantErr: true,
		},
		{ // 1
			in: `
                name: "Name0"`,
			want: TestConfig{
				Name: "Name0",
			},
			wantErr: false,
		},
		{ // 2
			in: `
                id: 23`,
			want: TestConfig{
				ID: 23,
			},
			wantErr: false,
		},
		{ // 3
			in: `
                id:   23
                name: Name0`,
			want: TestConfig{
				ID:   23,
				Name: "Name0",
			},
			wantErr: false,
		},
		{ // 4
			in: `
                id:   23
                name: Name0
                conn:
                  url: http://www.tests.to
                  passes:
                    - password0
                    - password1`,
			want: TestConfig{
				ID:   23,
				Name: "Name0",
				Conn: &TestConn{
					URL: "http://www.tests.to",
					Passes: []string{
						"password0",
						"password1",
					},
				},
			},
			wantErr: false,
		},
		{ // 5
			in: `
                id:   23
                name: {{ required "has to be set" "Name0" | quote }}
                conn:
                  url: http://www.tests.to
                  passes:
                    - password0
                    - password1`,
			want: TestConfig{
				ID:   23,
				Name: "Name0",
				Conn: &TestConn{
					URL: "http://www.tests.to",
					Passes: []string{
						"password0",
						"password1",
					},
				},
			},
			wantErr: false,
		},
		{ // 6
			in: `
                id:   23
                name: {{ required "has to be set" "" | quote }}`,
			want: TestConfig{
				ID:   23,
				Name: "Name0",
			},
			wantErr: true,
		},
		{ // 7
			in: `
                id:   23
                name: {{ required "has to be set" nil | quote }}`,
			want: TestConfig{
				ID:   23,
				Name: "Name0",
			},
			wantErr: true,
		},
		{ // 8
			in: `
                id:   23
                name: {{ required "has to be set" 9 | quote }}`,
			want: TestConfig{
				ID:   23,
				Name: "9",
			},
			wantErr: false,
		},
		{ // 9
			in: `
                id:   23
                name: {{ required "has to be set" 9 | quote`,
			want: TestConfig{
				ID:   23,
				Name: "9",
			},
			wantErr: true,
		},
		{ // 10
			inFile: "testData/test_config_0.yaml",
			want: TestConfig{
				ID:   9,
				Name: "Name0",
				Conn: &TestConn{
					URL: "http://www.tests.to",
					Passes: []string{
						"pass0",
						"pass1",
					},
				},
			},
			wantErr: false,
		},
		{ // 11
			inFile: "testData/test_config_1.yaml",
			env: map[string]string{
				"PASS1": "pass1",
			},
			want: TestConfig{
				ID:   9,
				Name: "Name1",
				Conn: &TestConn{
					URL: "http://www.tests.to",
					Passes: []string{
						"pass0",
						"pass1",
					},
				},
			},
			wantErr: false,
		},
	}

	testBuf := bytes.Buffer{}

	for k, test := range tests {
		if len(test.in) > 0 && len(test.inFile) > 0 {
			t.Errorf("%v: input data and file given at the same time", k)
		}

		testBuf.Reset()
		testBuf.WriteString(test.in)

		if test.env != nil {
			for ei, ev := range test.env {
				_ = os.Setenv(ei, ev)
			}
		}

		var c *Config[TestConfig]
		var fromErr error

		if len(test.in) > 0 {
			c, fromErr = From[TestConfig](&testBuf)
		} else if len(test.inFile) > 0 {
			c, fromErr = FromFile[TestConfig](test.inFile)
		} else {
			t.Errorf("%v: neither input data nor input file given", k)
		}

		if test.wantErr && fromErr == nil {
			t.Errorf("%v: wanted error but got nil", k)
		}
		if !test.wantErr && fromErr != nil {
			t.Errorf("%v: did not want error but got %v", k, fromErr)
		}

		if c != nil {
			if c.Get().ID != test.want.ID {
				t.Errorf("%v: wanted ID %v but got %v", k, test.want.ID, c.Get().ID)
			}
			if c.Get().Name != test.want.Name {
				t.Errorf("%v: wanted Name %v but got %v", k, test.want.Name, c.Get().Name)
			}
			if (c.Get().Conn != nil) != (test.want.Conn != nil) {
				t.Errorf("%v: wanted Conn == nil -> %v but got %v", k,
					test.want.Conn != nil,
					c.Get().Conn != nil)
			}
			if c.Get().Conn != nil && test.want.Conn != nil {
				if c.Get().Conn.URL != test.want.Conn.URL {
					t.Errorf("%v: wanted URL %v but got %v", k, test.want.Conn.URL, c.Get().Conn.URL)
				}
				for _, p := range test.want.Conn.Passes {
					if !slices.Contains(c.Get().Conn.Passes, p) {
						t.Errorf("%v: wanted passes to containt %v but was not there", k, p)
					}
				}
				for _, p := range c.Get().Conn.Passes {
					if !slices.Contains(test.want.Conn.Passes, p) {
						t.Errorf("%v: found pass %v but should not there", k, p)
					}
				}
			}
		}

		if test.env != nil {
			for ei := range test.env {
				_ = os.Unsetenv(ei)
			}
		}
	}
}

type BrokenIO struct{}

func (b *BrokenIO) Read(_ []byte) (n int, err error) {
	return 0, errors.New("broken reader")
}

func (b *BrokenIO) Write(_ []byte) (n int, err error) {
	return 0, errors.New("broken writer")
}

func TestBrokenReader(t *testing.T) {
	c, fromErr := From[TestConfig](&BrokenIO{})

	if fromErr == nil {
		t.Errorf("reading from broken reader should have returned an error")
	}

	if c != nil {
		t.Errorf("reading from broken reader should have returned nil")
	}
}

func TestNonexistingFile(t *testing.T) {
	c, fromErr := FromFile[TestConfig]("testData/test_does_not_exist.yaml")

	if fromErr == nil {
		t.Errorf("reading from broken reader should have returned an error")
	}

	if c != nil {
		t.Errorf("reading from broken reader should have returned nil")
	}
}

func TestBrokenWriter(t *testing.T) {
	c, _ := FromFile[TestConfig]("testData/test_config_0.yaml")

	toErr := c.To(&BrokenIO{})

	if toErr == nil {
		t.Errorf("reading from broken reader should have returned an error")
	}
}

func TestWriteFile(t *testing.T) {
	c, _ := FromFile[TestConfig]("testData/test_config_0.yaml")

	err := c.ToFile("testData/test_config_written.yaml")

	if err != nil {
		t.Errorf("writing to file should work")
	}
	defer func() { _ = os.Remove("testData/test_config_written.yaml") }()

	bufOrig := bytes.Buffer{}
	bufCopy := bytes.Buffer{}

	_ = c.To(&bufOrig)

	cp, _ := FromFile[TestConfig]("testData/test_config_written.yaml")
	_ = cp.To(&bufCopy)

	if bufOrig.String() != bufCopy.String() {
		t.Errorf("written file does not match original file")
	}
}

func TestWriteProtectedFile(t *testing.T) {
	c, _ := FromFile[TestConfig]("testData/test_config_0.yaml")

	if chmodErr := os.Chmod("testData/test_writeprotected.yaml", 0400); chmodErr != nil {
		t.Errorf("could not writeprotect file for test: %v", chmodErr)
	}

	err := c.ToFile("testData/test_writeprotected.yaml")

	if err == nil {
		t.Errorf("writing to file should not work")
	}
}

func TestSecretsHidden(t *testing.T) {
	c, _ := FromFile[TestConfig]("testData/test_config_0.yaml")

	buf := bytes.Buffer{}

	if err := c.ToSecretsHidden(&buf); err != nil {
		t.Errorf("could not generate secrets-hidden config")
	}

	if strings.Contains(buf.String(), "pass0") || strings.Contains(buf.String(), "pass1") {
		t.Errorf("found secrets in normally secrets-hidden output")
	}
}
