package okta

import (
	"testing"
)

func TestLogoStateFunc(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "../examples/okta_app_basic_auth/terraform_icon.png",
			expected: "188b6050b43d2fbc9be327e39bf5f7849b120bb4529bcd22cde78b02ccce6777", // compare to `shasum -a 256 filepath`
		},
		{
			input:    "invalid/file/path",
			expected: "",
		},
		{
			input:    "",
			expected: "",
		},
	}
	for _, c := range cases {
		result := localFileStateFunc(c.input)
		if result != c.expected {
			t.Errorf("Error matching logo, expected %q, got %q, for file %q", c.expected, result, c.input)
		}
	}
}
