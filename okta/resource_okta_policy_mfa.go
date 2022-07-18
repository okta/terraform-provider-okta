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

// resourcePolicyMfa requires Org Feature Flag OKTA_MFA_POLICY
func resourcePolicyMfa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaCreate,
		ReadContext:   resourcePolicyMfaRead,
		UpdateContext: resourcePolicyMfaUpdate,
		DeleteContext: resourcePolicyMfaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: buildMfaPolicySchema(buildFactorSchemaProviders()),
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

	syncSettings(d, policy.Settings)

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

// create or update a MFA policy
func buildMFAPolicy(d *schema.ResourceData) sdk.Policy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	if priority, ok := d.GetOk("priority"); ok {
		policy.Priority = int64(priority.(int))
	}
	policy.Settings = buildSettings(d)
	policy.Conditions = &okta.PolicyRuleConditions{
		People: getGroups(d),
	}
	return policy
}

// Opposite of syncSettings(): Build the corresponding sdk.PolicySettings based on the schema.ResourceData
func buildSettings(d *schema.ResourceData) *sdk.PolicySettings {
	if d.Get("is_oie") == true {
		authenticators := []*sdk.PolicyAuthenticator{}

		for _, key := range sdk.AuthenticatorProviders {
			rawFactor := d.Get(key).(map[string]interface{})
			enroll := rawFactor["enroll"]
			if enroll == nil {
				continue
			}

			authenticator := &sdk.PolicyAuthenticator{}
			authenticator.Key = key
			if enroll != nil {
				authenticator.Enroll = &sdk.Enroll{Self: enroll.(string)}
			}
			authenticators = append(authenticators, authenticator)
		}

		return &sdk.PolicySettings{
			Type:           "AUTHENTICATORS",
			Authenticators: authenticators,
		}
	}

	return &sdk.PolicySettings{
		Type: "FACTORS",
		Factors: &sdk.PolicyFactorsSettings{
			Duo:          buildFactorProvider(d, sdk.DuoFactor),
			FidoU2f:      buildFactorProvider(d, sdk.FidoU2fFactor),
			FidoWebauthn: buildFactorProvider(d, sdk.FidoWebauthnFactor),
			GoogleOtp:    buildFactorProvider(d, sdk.GoogleOtpFactor),
			Hotp:         buildFactorProvider(d, sdk.HotpFactor),
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
		},
	}
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

// Syncs either classic factors or OIE authenticators into the resource data.
func syncSettings(d *schema.ResourceData, settings *sdk.PolicySettings) {
	_ = d.Set("is_oie", settings.Type == "AUTHENTICATORS")

	if settings.Type == "AUTHENTICATORS" {
		for _, key := range sdk.AuthenticatorProviders {
			syncAuthenticator(d, key, settings.Authenticators)
		}
	} else {
		syncFactor(d, sdk.DuoFactor, settings.Factors.Duo)
		syncFactor(d, sdk.HotpFactor, settings.Factors.YubikeyToken)
		syncFactor(d, sdk.FidoU2fFactor, settings.Factors.FidoU2f)
		syncFactor(d, sdk.FidoWebauthnFactor, settings.Factors.FidoWebauthn)
		syncFactor(d, sdk.GoogleOtpFactor, settings.Factors.GoogleOtp)
		syncFactor(d, sdk.OktaCallFactor, settings.Factors.OktaCall)
		syncFactor(d, sdk.OktaOtpFactor, settings.Factors.OktaOtp)
		syncFactor(d, sdk.OktaPasswordFactor, settings.Factors.OktaPassword)
		syncFactor(d, sdk.OktaPushFactor, settings.Factors.OktaPush)
		syncFactor(d, sdk.OktaQuestionFactor, settings.Factors.OktaQuestion)
		syncFactor(d, sdk.OktaSmsFactor, settings.Factors.OktaSms)
		syncFactor(d, sdk.OktaEmailFactor, settings.Factors.OktaEmail)
		syncFactor(d, sdk.RsaTokenFactor, settings.Factors.RsaToken)
		syncFactor(d, sdk.SymantecVipFactor, settings.Factors.SymantecVip)
		syncFactor(d, sdk.YubikeyTokenFactor, settings.Factors.YubikeyToken)
	}
}

func syncFactor(d *schema.ResourceData, k string, f *sdk.PolicyFactor) {
	if f != nil {
		_ = d.Set(k, map[string]interface{}{
			"consent_type": f.Consent.Type,
			"enroll":       f.Enroll.Self,
		})
	}
}

func syncAuthenticator(d *schema.ResourceData, k string, authenticators []*sdk.PolicyAuthenticator) {
	for _, authenticator := range authenticators {
		if authenticator.Key == k {
			_ = d.Set(k, map[string]interface{}{
				"enroll": authenticator.Enroll.Self,
			})
			return
		}
	}
}

// List of factor provider above, they all follow the same schema
func buildFactorSchemaProviders() map[string]*schema.Schema {
	res := make(map[string]*schema.Schema)
	// Note: It's okay to append and have duplicates as we're setting back into a map here
	for _, key := range append(sdk.FactorProviders, sdk.AuthenticatorProviders...) {
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
					dErr := elemInSlice([]string{"NOT_ALLOWED", "OPTIONAL", "REQUIRED"})(enroll, cty.GetAttrPath("enroll"))
					if dErr != nil {
						errs = append(errs, dErr...)
					}
				}
				if consentType, ok := m["consent_type"]; ok {
					dErr := elemInSlice([]string{"NONE", "TERMS_OF_SERVICE"})(consentType, cty.GetAttrPath("consent_type"))
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
