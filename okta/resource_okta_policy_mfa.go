package okta

import (
	"context"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
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
	policy := buildMFAPolicy(d)
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
	syncFactor(d, sdk.DuoFactor, policy.Settings.Factors.Duo)
	syncFactor(d, sdk.FidoU2fFactor, policy.Settings.Factors.FidoU2f)
	syncFactor(d, sdk.FidoWebauthnFactor, policy.Settings.Factors.FidoWebauthn)
	syncFactor(d, sdk.GoogleOtpFactor, policy.Settings.Factors.GoogleOtp)
	syncFactor(d, sdk.OktaCallFactor, policy.Settings.Factors.OktaCall)
	syncFactor(d, sdk.OktaOtpFactor, policy.Settings.Factors.OktaOtp)
	syncFactor(d, sdk.OktaPasswordFactor, policy.Settings.Factors.OktaPassword)
	syncFactor(d, sdk.OktaPushFactor, policy.Settings.Factors.OktaPush)
	syncFactor(d, sdk.OktaQuestionFactor, policy.Settings.Factors.OktaQuestion)
	syncFactor(d, sdk.OktaSmsFactor, policy.Settings.Factors.OktaSms)
	syncFactor(d, sdk.OktaEmailFactor, policy.Settings.Factors.OktaEmail)
	syncFactor(d, sdk.RsaTokenFactor, policy.Settings.Factors.RsaToken)
	syncFactor(d, sdk.SymantecVipFactor, policy.Settings.Factors.SymantecVip)
	syncFactor(d, sdk.YubikeyTokenFactor, policy.Settings.Factors.YubikeyToken)
	syncFactor(d, sdk.HotpFactor, policy.Settings.Factors.YubikeyToken)
	err = syncPolicyFromUpstream(d, policy)
	if err != nil {
		return diag.Errorf("failed to sync policy: %v", err)
	}
	return nil
}

func resourcePolicyMfaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy := buildMFAPolicy(d)
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

// create or update a MFA policy
func buildMFAPolicy(d *schema.ResourceData) sdk.Policy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	if priority, ok := d.GetOk("priority"); ok {
		policy.Priority = int64(priority.(int))
	}
	policy.Settings = &sdk.PolicySettings{
		Factors: &sdk.PolicyFactorsSettings{
			Duo:          buildFactorProvider(d, sdk.DuoFactor),
			FidoU2f:      buildFactorProvider(d, sdk.FidoU2fFactor),
			FidoWebauthn: buildFactorProvider(d, sdk.FidoWebauthnFactor),
			GoogleOtp:    buildFactorProvider(d, sdk.GoogleOtpFactor),
			OktaCall:     buildFactorProvider(d, sdk.OktaCallFactor),
			OktaOtp:      buildFactorProvider(d, sdk.OktaOtpFactor),
			OktaPassword: buildFactorProvider(d, sdk.OktaPasswordFactor),
			OktaPush:     buildFactorProvider(d, sdk.OktaPushFactor),
			OktaQuestion: buildFactorProvider(d, sdk.OktaQuestionFactor),
			OktaSms:      buildFactorProvider(d, sdk.OktaSmsFactor),
			OktaEmail:    buildFactorProvider(d, sdk.OktaEmailFactor),
			RsaToken:     buildFactorProvider(d, sdk.RsaTokenFactor),
			SymantecVip:  buildFactorProvider(d, sdk.SymantecVipFactor),
			YubikeyToken: buildFactorProvider(d, sdk.YubikeyTokenFactor),
			Hotp:         buildFactorProvider(d, sdk.HotpFactor),
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
	sdk.DuoFactor,
	sdk.FidoU2fFactor,
	sdk.FidoWebauthnFactor,
	sdk.GoogleOtpFactor,
	sdk.OktaCallFactor,
	sdk.OktaOtpFactor,
	sdk.OktaPasswordFactor,
	sdk.OktaPushFactor,
	sdk.OktaQuestionFactor,
	sdk.OktaSmsFactor,
	sdk.OktaEmailFactor,
	sdk.RsaTokenFactor,
	sdk.SymantecVipFactor,
	sdk.YubikeyTokenFactor,
	sdk.HotpFactor,
}

// List of factor provider above, they all follow the same schema
func buildFactorProviders() map[string]*schema.Schema {
	res := make(map[string]*schema.Schema)
	for _, key := range factorProviders {
		res[key] = &schema.Schema{
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
		}
	}
	return res
}
