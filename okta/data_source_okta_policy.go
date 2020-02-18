package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of policy",
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{signOnPolicyType, passwordPolicyType, "MFA_ENROLL", "OAUTH_AUTHORIZATION_POLICY", "IDP_DISCOVERY"}, false),
				Description:  fmt.Sprintf("Policy type: %s, %s, MFA_ENROLL, IDP_DISCOVERY, or OAUTH_AUTHORIZATION_POLICY", signOnPolicyType, passwordPolicyType),
				Required:     true,
			},
		},
	}
}

func dataSourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	return setPolicyByName(d, m, d.Get("name").(string))
}

func setPolicyByName(d *schema.ResourceData, m interface{}, name string) error {
	client := getClientFromMetadata(m)
	ptype := d.Get("type").(string)

	currentPolicies, _, err := client.Policies.GetPoliciesByType(ptype)
	if err != nil {
		return fmt.Errorf("Error Listing Policies in Okta: %v", err)
	}
	if currentPolicies != nil {
		for _, policy := range currentPolicies.Policies {
			if policy.Name == name {
				d.SetId(policy.ID)
				return nil
			}
		}
	}
	return fmt.Errorf("No policies retrieved for policy type %v, name %s", ptype, name)
}
