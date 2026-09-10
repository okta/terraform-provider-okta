package idaas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func mfaPolicySchemasForTest() []struct {
	name   string
	schema map[string]*schema.Schema
} {
	return []struct {
		name   string
		schema map[string]*schema.Schema
	}{
		{"okta_policy_mfa", resourcePolicyMfa().Schema},
		{"okta_policy_mfa_default", resourcePolicyMfaDefault().Schema},
	}
}

func TestBuildSettings_OktaVerifyBreakdownKeys(t *testing.T) {
	for _, tt := range mfaPolicySchemasForTest() {
		t.Run(tt.name, func(t *testing.T) {
			raw := map[string]interface{}{
				"is_oie":                     true,
				sdk.OktaVerifyFastPassFactor: map[string]interface{}{"enroll": "REQUIRED"},
				sdk.OktaVerifyPushFactor:     map[string]interface{}{"enroll": "OPTIONAL"},
				sdk.OktaVerifyTotpFactor:     map[string]interface{}{"enroll": "NOT_ALLOWED"},
			}
			d := schema.TestResourceDataRaw(t, tt.schema, raw)

			settings := buildSettings(d)
			if settings == nil {
				t.Fatalf("buildSettings returned nil")
			}
			if settings.Type != "AUTHENTICATORS" {
				t.Errorf("expected Type AUTHENTICATORS, got %q", settings.Type)
			}

			enrollByKey := map[string]string{}
			for _, a := range settings.Authenticators {
				enrollByKey[a.Key] = a.Enroll.Self
			}

			want := map[string]string{
				sdk.OktaVerifyFastPassFactor: "REQUIRED",
				sdk.OktaVerifyPushFactor:     "OPTIONAL",
				sdk.OktaVerifyTotpFactor:     "NOT_ALLOWED",
			}
			for key, enroll := range want {
				if got := enrollByKey[key]; got != enroll {
					t.Errorf("expected %s enroll %q, got %q", key, enroll, got)
				}
			}
			if _, ok := enrollByKey[sdk.OktaVerifyFactor]; ok {
				t.Errorf("did not expect combined %s key in authenticators, got %v", sdk.OktaVerifyFactor, enrollByKey)
			}
		})
	}
}

func TestSyncSettings_OktaVerifyBreakdownKeys(t *testing.T) {
	for _, tt := range mfaPolicySchemasForTest() {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tt.schema, map[string]interface{}{})

			settings := &sdk.SdkPolicySettings{
				Type: "AUTHENTICATORS",
				Authenticators: []*sdk.PolicyAuthenticator{
					{Key: sdk.OktaVerifyFastPassFactor, Enroll: &sdk.Enroll{Self: "REQUIRED"}},
					{Key: sdk.OktaVerifyPushFactor, Enroll: &sdk.Enroll{Self: "OPTIONAL"}},
					{Key: sdk.OktaVerifyTotpFactor, Enroll: &sdk.Enroll{Self: "NOT_ALLOWED"}},
				},
			}
			syncSettings(d, settings)

			checks := map[string]string{
				sdk.OktaVerifyFastPassFactor: "REQUIRED",
				sdk.OktaVerifyPushFactor:     "OPTIONAL",
				sdk.OktaVerifyTotpFactor:     "NOT_ALLOWED",
			}
			for key, want := range checks {
				got, _ := d.Get(key).(map[string]interface{})
				if got["enroll"] != want {
					t.Errorf("expected %s enroll %q, got %q", key, want, got["enroll"])
				}
			}

			combined, _ := d.Get(sdk.OktaVerifyFactor).(map[string]interface{})
			if len(combined) != 0 {
				t.Errorf("expected combined %s to be empty, got %v", sdk.OktaVerifyFactor, combined)
			}
		})
	}
}

func TestSyncSettings_CombinedOktaVerifyKeyStillRoundTrips(t *testing.T) {
	for _, tt := range mfaPolicySchemasForTest() {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tt.schema, map[string]interface{}{})

			settings := &sdk.SdkPolicySettings{
				Type: "AUTHENTICATORS",
				Authenticators: []*sdk.PolicyAuthenticator{
					{Key: sdk.OktaVerifyFactor, Enroll: &sdk.Enroll{Self: "OPTIONAL"}},
				},
			}
			syncSettings(d, settings)

			combined, _ := d.Get(sdk.OktaVerifyFactor).(map[string]interface{})
			if combined["enroll"] != "OPTIONAL" {
				t.Errorf("expected combined %s enroll OPTIONAL, got %q", sdk.OktaVerifyFactor, combined["enroll"])
			}

			for _, key := range []string{sdk.OktaVerifyFastPassFactor, sdk.OktaVerifyPushFactor, sdk.OktaVerifyTotpFactor} {
				got, _ := d.Get(key).(map[string]interface{})
				if len(got) != 0 {
					t.Errorf("expected %s to be empty, got %v", key, got)
				}
			}
		})
	}
}
