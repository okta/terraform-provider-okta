package okta

import (
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-okta/sdk"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	postBindingAlias     = "HTTP-POST"
	redirectBindingAlias = "HTTP-REDIRECT"
)

var (
	baseIdpSchema = map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of idp",
		},
		"status": statusSchema,
		"account_link_action": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "AUTO",
		},
		"account_link_group_include": &schema.Schema{
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"provisioning_action": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"AUTO", "DISABLED", ""}, false),
			Default:      "AUTO",
		},
		"deprovisioned_action": actionSchema,
		"suspended_action":     actionSchema,
		"groups_action": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "NONE",
			ValidateFunc: validation.StringInSlice([]string{"NONE", "SYNC", "APPEND", "ASSIGN"}, false),
		},
		"groups_attribute": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"groups_assignment": &schema.Schema{
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"groups_filter": &schema.Schema{
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"username_template": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "idpuser.email",
		},
		"subject_match_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "USERNAME",
		},
		"subject_match_attribute": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"profile_master": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"request_signature_algorithm": algorithmSchema,
		"request_signature_scope": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "algorithm to use to sign response",
			ValidateFunc: validation.StringInSlice([]string{"REQUEST", ""}, false),
		},
		"response_signature_algorithm": algorithmSchema,
		"response_signature_scope": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "algorithm to use to sign response",
			ValidateFunc: validation.StringInSlice([]string{"RESPONSE", "ANY", ""}, false),
		},
	}

	actionSchema = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "NONE",
	}

	algorithmSchema = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "algorithm to use to sign requests",
		ValidateFunc: validation.StringInSlice([]string{"SHA-256", "SHA-1"}, false),
		Default:      "SHA-256",
	}

	optBindingSchema = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	optUrlSchema = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	optionalUrlSchema = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	bindingSchema = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}, false),
	}

	optionalBindingSchema = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}, false),
	}

	issuerMode = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
		ValidateFunc: validation.StringInSlice([]string{"ORG_URL", "CUSTOM_URL_DOMAIN"}, false),
		Default:      "ORG_URL",
		Optional:     true,
	}

	urlSchema = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
)

func buildIdpSchema(idpSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseIdpSchema, idpSchema)
}

func resourceIdpDelete(d *schema.ResourceData, m interface{}) error {
	return resourceDeleteAnyIdp(d, m, d.Get("status").(string) == "ACTIVE")
}

func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {
	return resourceDeleteAnyIdp(d, m, d.Get("active").(bool))
}

func resourceDeleteAnyIdp(d *schema.ResourceData, m interface{}, active bool) error {
	client := getSupplementFromMetadata(m)

	if active {
		if resp, err := client.DeactivateIdentityProvider(d.Id()); err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return err
			}
		}
	}

	if resp, err := client.DeleteIdentityProvider(d.Id()); err != nil {
		return suppressErrorOn404(resp, err)
	}

	return nil
}

func fetchIdp(id string, m interface{}, idp sdk.IdentityProvider) error {
	client := getSupplementFromMetadata(m)
	_, response, err := client.GetIdentityProvider(id, idp)
	if response.StatusCode == http.StatusNotFound {
		idp = nil
		return nil
	}

	return responseErr(response, err)
}

func updateIdp(id string, m interface{}, idp sdk.IdentityProvider) error {
	client := getSupplementFromMetadata(m)
	_, response, err := client.UpdateIdentityProvider(id, idp, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		idp = nil
		return nil
	}

	return responseErr(response, err)
}

func createIdp(m interface{}, idp sdk.IdentityProvider) error {
	client := getSupplementFromMetadata(m)
	_, response, err := client.CreateIdentityProvider(idp, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		idp = nil
		return nil
	}

	return responseErr(response, err)
}

func setIdpStatus(id, status, desiredStatus string, m interface{}) error {
	if status != desiredStatus {
		c := getSupplementFromMetadata(m)

		if desiredStatus == "INACTIVE" {
			return responseErr(c.DeactivateIdentityProvider(id))
		} else if desiredStatus == "ACTIVE" {
			return responseErr(c.ActivateIdentityProvider(id))
		}
	}

	return nil
}

func syncGroupActions(d *schema.ResourceData, groups *sdk.IDPGroupsAction) error {
	if groups != nil {
		d.Set("groups_action", groups.Action)
		d.Set("groups_attribute", groups.SourceAttributeName)

		return setNonPrimitives(d, map[string]interface{}{
			"groups_assignment": groups.Assignments,
			"groups_filter":     groups.Filter,
		})
	}

	return nil
}

func getIdentityProviderExists(idp sdk.IdentityProvider) schema.ExistsFunc {
	return func(d *schema.ResourceData, m interface{}) (bool, error) {
		_, resp, err := getSupplementFromMetadata(m).GetIdentityProvider(d.Id(), idp)

		return resp.StatusCode == 200, err
	}
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

			d.Set("request_algorithm", reqSign.Algorithm)
			d.Set("request_scope", reqSign.Scope)
		}

		if alg.Response != nil && alg.Response.Signature != nil {
			resSign := alg.Response.Signature

			d.Set("response_algorithm", resSign.Algorithm)
			d.Set("response_scope", resSign.Scope)
		}
	}

}
