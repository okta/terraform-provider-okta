package okta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

const (
	postBindingAlias     = "HTTP-POST"
	redirectBindingAlias = "HTTP-REDIRECT"
)

var (
	baseIdpSchema = map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of idp",
		},
		"status": statusSchema,
		"account_link_action": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "AUTO",
		},
		"account_link_group_include": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"provisioning_action": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringInSlice([]string{"AUTO", "DISABLED", ""}),
			Default:          "AUTO",
		},
		"deprovisioned_action": actionSchema,
		"suspended_action":     actionSchema,
		"groups_action": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "NONE",
			ValidateDiagFunc: stringInSlice([]string{"NONE", "SYNC", "APPEND", "ASSIGN"}),
		},
		"groups_attribute": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"groups_assignment": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"groups_filter": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"username_template": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "idpuser.email",
		},
		"subject_match_type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "USERNAME",
		},
		"subject_match_attribute": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"profile_master": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"request_signature_algorithm": algorithmSchema,
		"request_signature_scope": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "algorithm to use to sign response",
			ValidateDiagFunc: stringInSlice([]string{"REQUEST", ""}),
		},
		"response_signature_algorithm": algorithmSchema,
		"response_signature_scope": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "algorithm to use to sign response",
			ValidateDiagFunc: stringInSlice([]string{"RESPONSE", "ANY", ""}),
		},
	}

	actionSchema = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "NONE",
	}

	algorithmSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "algorithm to use to sign requests",
		ValidateDiagFunc: stringInSlice([]string{"SHA-256", "SHA-1"}),
		Default:          "SHA-256",
	}

	optBindingSchema = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	optURLSchema = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	optionalURLSchema = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	bindingSchema = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: stringInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}),
	}

	optionalBindingSchema = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: stringInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}),
	}

	issuerMode = &schema.Schema{
		Type:             schema.TypeString,
		Description:      "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
		ValidateDiagFunc: stringInSlice([]string{"ORG_URL", "CUSTOM_URL_DOMAIN"}),
		Default:          "ORG_URL",
		Optional:         true,
	}

	urlSchema = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
)

func buildIdpSchema(idpSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseIdpSchema, idpSchema)
}

func resourceIdpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDeleteAnyIdp(ctx, d, m, d.Get("status").(string) == statusActive)
}

func resourceDeleteAnyIdp(ctx context.Context, d *schema.ResourceData, m interface{}, active bool) diag.Diagnostics {
	client := getSupplementFromMetadata(m)

	if active {
		resp, err := client.DeactivateIdentityProvider(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to deactivate identity provider: %v", err)
		}
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
	}
	_, err := client.DeleteIdentityProvider(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete identity provider: %v", err)
	}
	return nil
}

func setIdpStatus(ctx context.Context, d *schema.ResourceData, client *okta.Client, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	var err error
	if desiredStatus == statusInactive {
		_, _, err = client.IdentityProvider.DeactivateIdentityProvider(ctx, d.Id())
	} else {
		_, _, err = client.IdentityProvider.ActivateIdentityProvider(ctx, d.Id())
	}
	return err
}

func syncGroupActions(d *schema.ResourceData, groups *sdk.IDPGroupsAction) error {
	if groups == nil {
		return nil
	}
	_ = d.Set("groups_action", groups.Action)
	_ = d.Set("groups_attribute", groups.SourceAttributeName)
	return setNonPrimitives(d, map[string]interface{}{
		"groups_assignment": groups.Assignments,
		"groups_filter":     groups.Filter,
	})
}

func NewIdpProvisioning(d *schema.ResourceData) *sdk.IDPProvisioning {
	return &sdk.IDPProvisioning{
		Action:        d.Get("provisioning_action").(string),
		ProfileMaster: d.Get("profile_master").(bool),
		Conditions: &sdk.IDPConditions{
			Deprovisioned: &sdk.IDPAction{
				Action: d.Get("deprovisioned_action").(string),
			},
			Suspended: &sdk.IDPAction{
				Action: d.Get("suspended_action").(string),
			},
		},
		Groups: &sdk.IDPGroupsAction{
			Action:              d.Get("groups_action").(string),
			Assignments:         convertInterfaceToStringSetNullable(d.Get("groups_assignment")),
			Filter:              convertInterfaceToStringSetNullable(d.Get("groups_filter")),
			SourceAttributeName: d.Get("groups_attribute").(string),
		},
	}
}

func NewAccountLink(d *schema.ResourceData) *sdk.AccountLink {
	link := convertInterfaceToStringSet(d.Get("account_link_group_include"))
	var filter *sdk.Filter

	if len(link) > 0 {
		filter = &sdk.Filter{
			Groups: &sdk.Included{
				Include: link,
			},
		}
	}

	return &sdk.AccountLink{
		Action: d.Get("account_link_action").(string),
		Filter: filter,
	}
}

func NewAlgorithms(d *schema.ResourceData) *sdk.Algorithms {
	return &sdk.Algorithms{
		Request:  NewSignature(d, "request"),
		Response: NewSignature(d, "response"),
	}
}

func NewSignature(d *schema.ResourceData, key string) *sdk.IDPSignature {
	scopeKey := fmt.Sprintf("%s_signature_scope", key)
	scope := d.Get(scopeKey).(string)

	if scope == "" {
		return nil
	}

	return &sdk.IDPSignature{
		Signature: &sdk.Signature{
			Algorithm: d.Get(fmt.Sprintf("%s_signature_algorithm", key)).(string),
			Scope:     scope,
		},
	}
}

func NewAcs(d *schema.ResourceData) *sdk.ACSSSO {
	return &sdk.ACSSSO{
		Binding: d.Get("acs_binding").(string),
		Type:    d.Get("acs_type").(string),
	}
}

func NewEndpoints(d *schema.ResourceData) *sdk.OIDCEndpoints {
	return &sdk.OIDCEndpoints{
		Acs:           NewAcs(d),
		Authorization: sdk.GetEndpoint(d, "authorization"),
		Token:         sdk.GetEndpoint(d, "token"),
		UserInfo:      sdk.GetEndpoint(d, "user_info"),
		Jwks:          sdk.GetEndpoint(d, "jwks"),
	}
}

func syncAlgo(d *schema.ResourceData, alg *sdk.Algorithms) {
	if alg != nil {
		if alg.Request != nil && alg.Request.Signature != nil {
			reqSign := alg.Request.Signature

			_ = d.Set("request_signature_algorithm", reqSign.Algorithm)
			_ = d.Set("request_signature_scope", reqSign.Scope)
		}
		if alg.Response != nil && alg.Response.Signature != nil {
			resSign := alg.Response.Signature
			_ = d.Set("response_signature_algorithm", resSign.Algorithm)
			_ = d.Set("response_signature_scope", resSign.Scope)
		}
	}
}
