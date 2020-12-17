package okta

import (
	"context"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicyMfa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaCreate,
		ReadContext:   resourcePolicyMfaRead,
		UpdateContext: resourcePolicyMfaUpdate,
		DeleteContext: resourcePolicyMfaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: buildPolicySchema(buildFactorProviders()),
	}
}

func resourcePolicyMfaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy := buildMfaPolicy(d)
	err := createPolicy(ctx, d, m, policy)
	if err != nil {
		return diag.Errorf("failed to create MFA policy: %v", err)
	}
	return resourcePolicyMfaRead(ctx, d, m)
}

func resourcePolicyMfaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get MFA policy: %v", err)
	}
	if policy == nil {
		return nil
	}
	syncFactor(d, "duo", policy.Settings.Factors.Duo)
	syncFactor(d, "fido_u2f", policy.Settings.Factors.FidoU2f)
	syncFactor(d, "fido_webauthn", policy.Settings.Factors.FidoWebauthn)
	syncFactor(d, "google_otp", policy.Settings.Factors.GoogleOtp)
	syncFactor(d, "okta_call", policy.Settings.Factors.OktaCall)
	syncFactor(d, "okta_otp", policy.Settings.Factors.OktaOtp)           //
	syncFactor(d, "okta_password", policy.Settings.Factors.OktaPassword) //
	syncFactor(d, "okta_push", policy.Settings.Factors.OktaPush)         //
	syncFactor(d, "okta_question", policy.Settings.Factors.OktaQuestion)
	syncFactor(d, "okta_sms", policy.Settings.Factors.OktaSms)
	syncFactor(d, "rsa_token", policy.Settings.Factors.RsaToken)
	syncFactor(d, "symantec_vip", policy.Settings.Factors.SymantecVip)
	syncFactor(d, "yubikey_token", policy.Settings.Factors.YubikeyToken)
	err = syncPolicyFromUpstream(d, policy)
	if err != nil {
		return diag.Errorf("failed to sync policy: %v", err)
	}
	return nil
}

func resourcePolicyMfaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy := buildMfaPolicy(d)
	err := updatePolicy(ctx, d, m, policy)
	if err != nil {
		return diag.Errorf("failed to update MFA policy: %v", err)
	}
	return resourcePolicyMfaRead(ctx, d, m)
}

func resourcePolicyMfaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deletePolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete MFA policy: %v", err)
	}
	return nil
}

func buildFactorProvider(d *schema.ResourceData, key string) *sdk.PolicyFactor {
	rawFactor := d.Get(key).(map[string]interface{})
	consent := rawFactor["consent_type"]
	enroll := rawFactor["enroll"]
	if consent == nil && enroll == nil {
		return nil
	}
	f := &sdk.PolicyFactor{}
	if consent != nil {
		f.Consent = &sdk.Consent{Type: consent.(string)}
	}
	if enroll != nil {
		f.Enroll = &sdk.Enroll{Self: enroll.(string)}
	}
	return f
}

// create or update a password policy
func buildMfaPolicy(d *schema.ResourceData) sdk.Policy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	if priority, ok := d.GetOk("priority"); ok {
		policy.Priority = int64(priority.(int))
	}
	policy.Settings = &sdk.PolicySettings{
		Factors: &sdk.PolicyFactorsSettings{
			Duo:          buildFactorProvider(d, "duo"),
			FidoU2f:      buildFactorProvider(d, "fido_u2f"),
			FidoWebauthn: buildFactorProvider(d, "fido_webauthn"),
			GoogleOtp:    buildFactorProvider(d, "google_otp"),
			OktaCall:     buildFactorProvider(d, "okta_call"),
			OktaOtp:      buildFactorProvider(d, "okta_otp"),
			OktaPassword: buildFactorProvider(d, "okta_password"),
			OktaPush:     buildFactorProvider(d, "okta_push"),
			OktaQuestion: buildFactorProvider(d, "okta_question"),
			OktaSms:      buildFactorProvider(d, "okta_sms"),
			RsaToken:     buildFactorProvider(d, "rsa_token"),
			SymantecVip:  buildFactorProvider(d, "symantec_vip"),
			YubikeyToken: buildFactorProvider(d, "yubikey_token"),
		},
	}
	policy.Conditions = &okta.PolicyRuleConditions{
		People: getGroups(d),
	}
	return policy
}

func syncFactor(d *schema.ResourceData, k string, f *sdk.PolicyFactor) {
	if f == nil {
		return
	}
	_ = d.Set(k, map[string]interface{}{
		"consent_type": f.Consent.Type,
		"enroll":       f.Enroll.Self,
	})
}

var factorProviders = []string{
	"duo",
	"fido_u2f",
	"fido_webauthn",
	"google_otp",
	"okta_call",
	"okta_otp",
	"okta_password",
	"okta_push",
	"okta_question",
	"okta_sms",
	"rsa_token",
	"symantec_vip",
	"yubikey_token",
}

// List of factor provider above, they all follow the same schema
func buildFactorProviders() map[string]*schema.Schema {
	res := make(map[string]*schema.Schema)
	for _, key := range factorProviders {
		sMap := getPolicyFactorSchema(key)
		for nestedKey, nestedVal := range sMap {
			res[nestedKey] = nestedVal
		}
	}
	return res
}

func getPolicyFactorSchema(key string) map[string]*schema.Schema {
	// These are primitives to allow defaulting. Terraform still does not support aggregate defaults.
	return map[string]*schema.Schema{
		key: {
			Optional: true,
			Type:     schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.HasSuffix(k, ".%") || new == ""
			},
			ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
				var errs diag.Diagnostics
				m := i.(map[string]interface{})
				if enroll, ok := m["enroll"]; ok {
					dErr := stringInSlice([]string{"NOT_ALLOWED", "OPTIONAL", "REQUIRED"})(enroll, cty.GetAttrPath("enroll"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
				}
				if consentType, ok := m["consent_type"]; ok {
					dErr := stringInSlice([]string{"NONE", "TERMS_OF_SERVICE"})(consentType, cty.GetAttrPath("consent_type"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
				}
				return errs
			},
		},
	}
}
