package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"net/http"
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
		"status":                       statusSchema,
		"request_signature_algorithm":  algorithmSchema,
		"response_signature_algorithm": algorithmSchema,
		"request_signature_scope": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "algorithm to use to sign response",
		},
		"response_signature_scope": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "algorithm to use to sign response",
		},
		"acs_binding": bindingSchema,
		"acs_type": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "INSTANCE",
			ValidateFunc: validation.StringInSlice([]string{"INSTANCE"}, false),
		},
		"account_link_action": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "AUTO",
			ValidateFunc: validation.StringInSlice([]string{"AUTO"}, false),
		},
		"account_link_filter": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"deprovisioned_action": actionSchema,
		"suspended_action":     actionSchema,
		"groups_action":        actionSchema,
		"username_template": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"subject_match_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"profile_master": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
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
		Default:      "SHA-256",
		ValidateFunc: validation.StringInSlice([]string{"SHA-256"}, false),
	}

	bindingSchema = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}, false),
	}

	urlSchema = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
)

func buildIdpSchema(idpSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseIdpSchema, idpSchema)
}

func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	if d.Get("status").(string) == "ACTIVE" {
		if resp, err := client.DeactivateIdentityProvider(d.Id()); err != nil {
			return suppressErrorOn404(resp, err)
		}
	}

	if resp, err := client.DeleteIdentityProvider(d.Id()); err != nil {
		return suppressErrorOn404(resp, err)
	}

	return nil
}

func fetchIdp(id string, m interface{}, idp IdentityProvider) error {
	client := getSupplementFromMetadata(m)
	_, response, err := client.GetIdentityProvider(id, idp)
	if response.StatusCode == http.StatusNotFound {
		idp = nil
		return nil
	}

	return responseErr(response, err)
}

func updateIdp(id string, m interface{}, idp IdentityProvider) error {
	client := getSupplementFromMetadata(m)
	_, response, err := client.UpdateIdentityProvider(id, idp, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		idp = nil
		return nil
	}

	return responseErr(response, err)
}

func createIdp(m interface{}, idp IdentityProvider) error {
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

func getIdentityProviderExists(idp IdentityProvider) schema.ExistsFunc {
	return func(d *schema.ResourceData, m interface{}) (bool, error) {
		_, resp, err := getSupplementFromMetadata(m).GetIdentityProvider(d.Id(), idp)

		return err == nil && !is404(resp.StatusCode), err
	}
}
