package templig

import (
	"errors"
	"testing"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		text       string
		err        error
		wantErr    bool
		wantErrMsg string
	}{
		{ // 0
			text:       "wrap this: %v",
			err:        errors.New("original error"),
			wantErr:    true,
			wantErrMsg: "wrap this: original error",
		},
		{ // 1
			text:    "wrap this: %v",
			err:     nil,
			wantErr: false,
		},
		{ // 2
			text:       "",
			err:        errors.New("error case"),
			wantErr:    true,
			wantErrMsg: "error case",
		},
		{
			text:    "",
			err:     nil,
			wantErr: false,
		},
	}

	for k, v := range tests {
		got := wrapError(v.text, v.err)

		if (got != nil) != v.wantErr {
			t.Errorf(`%v: got error "%v", but wanted "%v"`, k, got, v.wantErr)
		}

		if v.wantErr && got.Error() != v.wantErrMsg {
			t.Errorf(`%v: got error "%v" but wanted "%v"`, k, got.Error(), v.wantErrMsg)
		}
	}
}
