package idaas_test

import (
	"testing"

	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestDomainValidationString(t *testing.T) {
	tests := []struct {
		element  string
		expected bool
	}{
		{"VERIFIED", true},
		{"COMPLETED", true},
		{"NOT_STARTED", false},
		{"IN_PROGRESS", false},
		{"verified", false},
		{"completed", false},
		{"invalid", false},
	}

	for _, test := range tests {
		actual := idaas.IsDomainValidated(test.element)

		if actual != test.expected {
			t.Errorf("domain validation status failed for status = \"%s\" - Expected: %t, Actual: %t", test.element, test.expected, actual)
		}
	}
}
