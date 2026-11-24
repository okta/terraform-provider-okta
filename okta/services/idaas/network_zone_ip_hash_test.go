package idaas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCanonicalIPAddress(t *testing.T) {
	tests := map[string]string{
		"203.0.113.5":             "203.0.113.5",
		"203.0.113.5-203.0.113.5": "203.0.113.5",
		"203.0.113.5-203.0.113.9": "203.0.113.5-203.0.113.9",
		"203.0.113.0/24":          "203.0.113.0/24",
		"2001:db8::1-2001:db8::1": "2001:db8::1",
		"":                        "",
	}
	for in, want := range tests {
		if got := canonicalIPAddress(in); got != want {
			t.Errorf("canonicalIPAddress(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestNetworkZoneAddressDiff runs the real schema diff for gateways and
// proxies, with state as written by Read and a given config. Regression test
// for https://github.com/okta/terraform-provider-okta/issues/2928.
func TestNetworkZoneAddressDiff(t *testing.T) {
	tests := []struct {
		name     string
		state    []interface{}
		config   []interface{}
		wantDiff bool
	}{
		{"single IP in range form matches bare IP in state", []interface{}{"203.0.113.5"}, []interface{}{"203.0.113.5-203.0.113.5"}, false},
		{"bare IP matches range form in state", []interface{}{"203.0.113.5-203.0.113.5"}, []interface{}{"203.0.113.5"}, false},
		{"identical values", []interface{}{"203.0.113.5"}, []interface{}{"203.0.113.5"}, false},
		{"mixed list with range form", []interface{}{"203.0.113.5", "198.51.100.0/24"}, []interface{}{"198.51.100.0/24", "203.0.113.5-203.0.113.5"}, false},
		{"different IP is a change", []interface{}{"203.0.113.5"}, []interface{}{"203.0.113.6"}, true},
		{"real range is a change", []interface{}{"203.0.113.5"}, []interface{}{"203.0.113.5-203.0.113.9"}, true},
		{"added address is a change", []interface{}{"203.0.113.5"}, []interface{}{"203.0.113.5-203.0.113.5", "203.0.113.6"}, true},
	}

	for _, attr := range []string{"gateways", "proxies"} {
		for _, tt := range tests {
			t.Run(attr+"/"+tt.name, func(t *testing.T) {
				r := resourceNetworkZone()
				base := map[string]interface{}{"name": "zone", "type": "IP", "status": "ACTIVE", "usage": "POLICY"}

				state := r.Data(nil)
				state.SetId("nzo123")
				_ = state.Set("name", base["name"])
				_ = state.Set("type", base["type"])
				_ = state.Set("status", base["status"])
				_ = state.Set("usage", base["usage"])
				// lintignore:R001 // the same cases run for gateways and proxies
				if err := state.Set(attr, tt.state); err != nil {
					t.Fatalf("setting %s: %v", attr, err)
				}

				config := map[string]interface{}{attr: tt.config}
				for k, v := range base {
					config[k] = v
				}

				diff, err := r.Diff(context.Background(), state.State(), terraform.NewResourceConfigRaw(config), nil)
				if err != nil {
					t.Fatalf("Diff: %v", err)
				}
				gotDiff := diff != nil && len(diff.Attributes) > 0
				if gotDiff != tt.wantDiff {
					t.Errorf("diff present = %v, want %v (diff: %v)", gotDiff, tt.wantDiff, diff)
				}
			})
		}
	}
}
