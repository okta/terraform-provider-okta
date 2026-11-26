package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// GH-3044: Okta returns grace period durations in hours ("P7D" comes back as "PT168H").
func TestSameISO8601Duration_GH3044(t *testing.T) {
	cases := []struct {
		a, b string
		same bool
	}{
		{"P7D", "PT168H", true},
		{"P30D", "PT720H", true},
		{"P1W", "P7D", true},
		{"PT1H", "PT60M", true},
		{"P1DT12H", "PT36H", true},
		{"PT90S", "PT1M30S", true},
		{"P7D", "P7D", true},
		{"P7D", "PT167H", false},
		{"P1M", "P30D", false}, // months and years have no fixed length
		{"P1M", "P1M", true},
		{"P1Y", "P1Y", true},
		{"P1Y", "P365D", false},
		{"P7D", "not-a-duration", false},
		{"", "", true},
	}
	for _, c := range cases {
		assert.Equalf(t, c.same, SameISO8601Duration(c.a, c.b), "%q vs %q", c.a, c.b)
	}
}

func TestPreferPriorDuration_GH3044(t *testing.T) {
	// Same duration: keep what the configuration says, so no diff and no inconsistent result.
	assert.Equal(t, types.StringValue("P7D"), PreferPriorDuration(types.StringValue("P7D"), "PT168H"))
	// Changed outside Terraform: report the API value.
	assert.Equal(t, types.StringValue("PT240H"), PreferPriorDuration(types.StringValue("P7D"), "PT240H"))
	// No prior value (import, first read): API value.
	assert.Equal(t, types.StringValue("PT168H"), PreferPriorDuration(types.StringNull(), "PT168H"))
	assert.Equal(t, types.StringValue("PT168H"), PreferPriorDuration(types.StringUnknown(), "PT168H"))
}
