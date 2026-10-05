package idaas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSetProfileEnrollmentTargetGroupID(t *testing.T) {
	tests := []struct {
		name  string
		prior string
		ids   []string
		want  string
	}{
		{name: "group assigned in Okta", prior: "", ids: []string{"00g1"}, want: "00g1"},
		{name: "group replaced in Okta", prior: "00g1", ids: []string{"00g2"}, want: "00g2"},
		{name: "group removed in Okta clears state", prior: "00g1", ids: nil, want: ""},
		{name: "empty list clears state", prior: "00g1", ids: []string{}, want: ""},
		{name: "never had a group", prior: "", ids: nil, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourcePolicyProfileEnrollmentRule().Schema, map[string]interface{}{
				"unknown_user_action": "REGISTER",
				"target_group_id":     tt.prior,
			})
			setProfileEnrollmentTargetGroupID(d, tt.ids)
			if got := d.Get("target_group_id").(string); got != tt.want {
				t.Fatalf("target_group_id = %q, want %q", got, tt.want)
			}
		})
	}
}
