package notifications

import (
	"reflect"
	"testing"
)

func Test_sanitizeShellInput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			"Valid value",
			"abcdez123",
			"abcdez123",
		},
		{
			"With special chars from other languages",
			"абвгдabcdez123äöü",
			"абвгдabcdez123äöü",
		},
		{
			"With stripping of some disallowed chars",
			"a\"bcdez'123",
			"abcdez123",
		},
		{
			"With some allowed delimeters of some disallowed chars",
			"abc_de-z123",
			"abc_de-z123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := sanitizeShellInput(tt.input)

			if !reflect.DeepEqual(res, tt.output) {
				t.Errorf("sanitizeShellInput = \n%#v, want \n%#v", res, tt.output)
			}
		})
	}
}
