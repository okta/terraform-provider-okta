package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyMfaDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		ReadContext:   resourcePolicyMfaDefaultRead,
		UpdateContext: resourcePolicyMfaDefaultCreateOrUpdate,
		DeleteContext: resourcePolicyMfaDefaultDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_, err := setDefaultPolicy(ctx, d, m, sdk.MfaPolicyType)
				if err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: buildDefaultPolicySchema(buildFactorProviders()),
	}
}

func resourcePolicyMfaDefaultCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	if id == "" {
		policy, err := setDefaultPolicy(ctx, d, m, sdk.MfaPolicyType)
		if err != nil {
			return diag.FromErr(err)
		}
		id = policy.Id
	}
	_, _, err := getSupplementFromMetadata(m).UpdatePolicy(ctx, id, buildDefaultMFAPolicy(d))
	if err != nil {
		return diag.Errorf("failed to update default MFA policy: %v", err)
	}
	return resourcePolicyMfaDefaultRead(ctx, d, m)
}

func resourcePolicyMfaDefaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get default MFA policy: %v", err)
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
	return nil
}

// Default policy can not be removed
func resourcePolicyMfaDefaultDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func buildDefaultMFAPolicy(d *schema.ResourceData) sdk.Policy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	policy.Priority = int64(d.Get("priority").(int))
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
		People: &okta.PolicyPeopleCondition{
			Groups: &okta.GroupCondition{
				Include: []string{d.Get("default_included_group_id").(string)},
			},
		},
	}
	return policy
}
