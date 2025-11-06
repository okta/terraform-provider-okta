package idaas

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/okta/terraform-provider-okta/okta/utils"
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

func resourcePolicyMfaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy := buildMFAPolicy(d)
	err := createPolicy(ctx, d, meta, policy)
	if err != nil {
		return diag.Errorf("failed to create MFA policy: %v", err)
	}
	return resourcePolicyMfaRead(ctx, d, meta)
}

func resourcePolicyMfaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, meta)
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

func resourcePolicyMfaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy := buildMFAPolicy(d)
	err := updatePolicy(ctx, d, meta, policy)
	if err != nil {
		return diag.Errorf("failed to update MFA policy: %v", err)
	}
	return resourcePolicyMfaRead(ctx, d, meta)
}

func resourcePolicyMfaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deletePolicy(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete MFA policy: %v", err)
	}
	return nil
}

// create or update a MFA policy
func buildMFAPolicy(d *schema.ResourceData) sdk.SdkPolicy {
	policy := sdk.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	if priority, ok := d.GetOk("priority"); ok {
		policy.PriorityPtr = utils.Int64Ptr(priority.(int))
	}
	policy.Settings = buildSettings(d)
	policy.Conditions = &sdk.PolicyRuleConditions{
		People: getGroups(d),
	}
	return policy
}

// Opposite of syncSettings(): Build the corresponding sdk.PolicySettings based on the schema.ResourceData
func buildSettings(d *schema.ResourceData) *sdk.SdkPolicySettings {
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
			authenticator.Enroll = &sdk.Enroll{Self: enroll.(string)}
			constraints := rawFactor["constraints"]
			if constraints != nil {
				c, ok := constraints.(string)
				if ok {
					// NOTE: we should consider using diff suppress func if we start seeing updating issue
					slice := strings.Split(c, ",")
					sort.Strings(slice)
					authenticator.Constraints = &sdk.PolicyAuthenticatorConstraints{AaguidGroups: slice}
				}
			}
			authenticators = append(authenticators, authenticator)
		}
		_, ok := d.GetOk("external_idps")
		if ok {
			rawExternalIDPs := d.Get("external_idps").(*schema.Set).List()
			for _, r := range rawExternalIDPs {
				rawExternalIDP := r.(map[string]interface{})
				enroll := rawExternalIDP["enroll"]
				if enroll == nil {
					continue
				}
				id := rawExternalIDP["id"]
				if id == nil {
					continue
				}
				authenticator := &sdk.PolicyAuthenticator{
					Key: "external_idp",
					ID:  id.(string),
					Enroll: &sdk.Enroll{
						Self: enroll.(string),
					},
				}
				constraints := rawExternalIDP["constraints"]
				if constraints != nil {
					c, ok := constraints.(string)
					if ok {
						slice := strings.Split(c, ",")
						sort.Strings(slice)
						authenticator.Constraints = &sdk.PolicyAuthenticatorConstraints{AaguidGroups: slice}
					}
				}
				authenticators = append(authenticators, authenticator)
			}
		}

		return &sdk.SdkPolicySettings{
			Type:           "AUTHENTICATORS",
			Authenticators: authenticators,
		}
	}

	return &sdk.SdkPolicySettings{
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
func syncSettings(d *schema.ResourceData, settings *sdk.SdkPolicySettings) {
	if settings == nil {
		// NOTE when sdk/policy.go is gone we probably won't need this guard
		return
	}

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
		// lintignore:R001
		_ = d.Set(k, map[string]interface{}{
			"consent_type": f.Consent.Type,
			"enroll":       f.Enroll.Self,
		})
	}
}

func syncAuthenticator(d *schema.ResourceData, k string, authenticators []*sdk.PolicyAuthenticator) {
	externalIdps := make([]interface{}, 0)
	for _, authenticator := range authenticators {
		if authenticator.Key == k {
			if k != "external_idp" {
				if authenticator.Constraints != nil {
					slice := authenticator.Constraints.AaguidGroups
					sort.Strings(slice)
					// lintignore:R001
					_ = d.Set(k, map[string]interface{}{
						"enroll":      authenticator.Enroll.Self,
						"constraints": strings.Join(slice, ","),
					})
				} else {
					// lintignore:R001
					_ = d.Set(k, map[string]interface{}{
						"enroll": authenticator.Enroll.Self,
					})
				}
				return
			} else {
				if idp, ok := d.GetOk("external_idp"); ok && idp != nil {
					if authenticator.Constraints != nil {
						slice := authenticator.Constraints.AaguidGroups
						sort.Strings(slice)
						// lintignore:R001
						_ = d.Set(k, map[string]interface{}{
							"enroll":      authenticator.Enroll.Self,
							"constraints": strings.Join(slice, ","),
						})
					} else {
						// lintignore:R001
						_ = d.Set(k, map[string]interface{}{
							"enroll": authenticator.Enroll.Self,
						})
					}
				} else {
					m := make(map[string]interface{})
					m["enroll"] = authenticator.Enroll.Self
					m["id"] = authenticator.ID
					if authenticator.Constraints != nil {
						slice := authenticator.Constraints.AaguidGroups
						sort.Strings(slice)
						m["constraints"] = strings.Join(slice, ",")
					}
				}
			}
		}
	}
	if len(externalIdps) > 0 {
		_ = d.Set("external_idps", externalIdps)
	}
}

// List of factor provider above, they all follow the same schema
func buildFactorSchemaProviders() map[string]*schema.Schema {
	res := make(map[string]*schema.Schema)
	// Note: It's okay to append and have duplicates as we're setting back into a map here
	for _, key := range append(sdk.FactorProviders, sdk.AuthenticatorProviders...) {
		if key == "external_idp" {
			res[key] = &schema.Schema{
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Deprecated: "Since okta now support multiple external_idps, this will be deprecated. Please use `external_idps` instead",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.HasSuffix(k, ".%") || new == ""
				},
				ConflictsWith: []string{"external_idps"},
			}
		} else {
			res[key] = &schema.Schema{
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.HasSuffix(k, ".%") || new == ""
				},
			}
		}
	}
	res["external_idps"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		DiffSuppressFunc: structure.SuppressJsonDiff,
		ConflictsWith:    []string{"external_idp"},
	}
	return res
}
