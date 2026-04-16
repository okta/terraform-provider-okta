package idaas

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

func dataSourcePolicyRulePassword() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRulePasswordRead,
		Description: "Get a Password Policy Rule from Okta.",
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Policy owning this rule.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the rule.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the rule: `ACTIVE` or `INACTIVE`.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Priority of the rule.",
			},
			"network_connection": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network selection mode: `ANYWHERE`, `ZONE`, `ON_NETWORK`, or `OFF_NETWORK`.",
			},
			"network_includes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Network zones to include (when `network_connection` = `ZONE`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"network_excludes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Network zones to exclude (when `network_connection` = `ZONE`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"users_excluded": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "User IDs to exclude from this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"users_included": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "User IDs to include in this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"groups_excluded": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Group IDs to exclude from this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"groups_included": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Group IDs to include in this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"password_change": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allow or deny a user to change their password: `ALLOW` or `DENY`.",
			},
			"password_reset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allow or deny a user to reset their password: `ALLOW` or `DENY`.",
			},
			"password_unlock": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allow or deny a user to unlock: `ALLOW` or `DENY`.",
			},
			"password_reset_access_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether SSPR access is governed by an authentication policy or legacy behavior. Options: `LEGACY`, `AUTH_POLICY`.",
			},
			"password_reset_requirement": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Self-service password reset (SSPR) requirement settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method_constraints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Constraints on the values specified in `primary_methods`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The method to constrain (e.g. `otp`).",
									},
									"allowed_authenticators": {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Keys of the authenticators allowed for this method (e.g. `google_otp`).",
									},
								},
							},
						},
						"primary_methods": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Primary authentication methods for SSPR.",
						},
						"step_up_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether step-up authentication is required for SSPR.",
						}, "step_up_methods": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Authenticator methods required for the secondary authentication step of password recovery. Items value: `security_question`.",
						}},
				},
			},
		},
	}
}

func dataSourcePolicyRulePasswordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policyID := d.Get("policy_id").(string)
	ruleID := d.Get("id").(string)

	resp, httpResp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.GetPolicyRule(ctx, policyID, ruleID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.Response != nil && httpResp.StatusCode == http.StatusNotFound {
			return diag.Errorf("password policy rule with id '%s' not found in policy '%s'", ruleID, policyID)
		}
		return diag.Errorf("failed to get password policy rule: %v", err)
	}

	if resp == nil || resp.PasswordPolicyRule == nil {
		return diag.Errorf("password policy rule response is empty")
	}

	return flattenPolicyRulePassword(d, resp.PasswordPolicyRule)
}

func flattenPolicyRulePassword(d *schema.ResourceData, rule *v6okta.PasswordPolicyRule) diag.Diagnostics {
	d.SetId(rule.GetId())
	_ = d.Set("name", rule.GetName())
	_ = d.Set("status", rule.GetStatus())
	_ = d.Set("priority", int(rule.GetPriority()))

	if conds := rule.Conditions; conds != nil {
		if net := conds.Network; net != nil {
			_ = d.Set("network_connection", net.GetConnection())
			if len(net.Include) > 0 {
				_ = d.Set("network_includes", net.Include)
			}
			if len(net.Exclude) > 0 {
				_ = d.Set("network_excludes", net.Exclude)
			}
		}
		if people := conds.People; people != nil {
			if users := people.Users; users != nil {
				_ = d.Set("users_excluded", users.Exclude)
				_ = d.Set("users_included", users.Include)
			}
			if groups := people.Groups; groups != nil {
				_ = d.Set("groups_excluded", groups.Exclude)
				_ = d.Set("groups_included", groups.Include)
			}
		}
	}

	if actions := rule.Actions; actions != nil {
		if pc := actions.PasswordChange; pc != nil {
			_ = d.Set("password_change", pc.GetAccess())
		}
		if su := actions.SelfServiceUnlock; su != nil {
			_ = d.Set("password_unlock", su.GetAccess())
		}
		if sspr := actions.SelfServicePasswordReset; sspr != nil {
			_ = d.Set("password_reset", sspr.GetAccess())
			if req := sspr.Requirement; req != nil {
				if ac := req.AccessControl; ac != nil {
					_ = d.Set("password_reset_access_control", ac)
				}
				reqMap := map[string]interface{}{
					"primary_methods":    []string{},
					"step_up_enabled":    false,
					"step_up_methods":    []string{},
					"method_constraints": []interface{}{},
				}
				if primary := req.Primary; primary != nil {
					reqMap["primary_methods"] = primary.Methods
					if len(primary.MethodConstraints) > 0 {
						constraints := make([]interface{}, 0, len(primary.MethodConstraints))
						for _, mc := range primary.MethodConstraints {
							mcMap := map[string]interface{}{
								"method":                 mc.GetMethod(),
								"allowed_authenticators": flattenAuthenticatorKeys(mc.GetAllowedAuthenticators()),
							}
							constraints = append(constraints, mcMap)
						}
						reqMap["method_constraints"] = constraints
					}
				}
				if stepUp := req.StepUp; stepUp != nil {
					reqMap["step_up_enabled"] = stepUp.GetRequired()
					reqMap["step_up_methods"] = stepUp.Methods
				}
				_ = d.Set("password_reset_requirement", []interface{}{reqMap})
			}
		}
	}

	return nil
}

// flattenAuthenticatorKeys extracts the key string from each AuthenticatorIdentity.
func flattenAuthenticatorKeys(identities []v6okta.AuthenticatorIdentity) []string {
	keys := make([]string, 0, len(identities))
	for _, identity := range identities {
		if k := identity.GetKey(); k != "" {
			keys = append(keys, k)
		}
	}
	return keys
}
