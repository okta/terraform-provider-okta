package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
				ValidateFunc: validation.StringInSlice([]string{sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.IdpDiscoveryType, sdk.OauthAuthorizationPolicyType}, false),
				Description:  fmt.Sprintf("Policy type: %s, %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.IdpDiscoveryType, sdk.OauthAuthorizationPolicyType),
				Required:     true,
			},
		},
	}
}

func dataSourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	return setPolicyByName(d, m, d.Get("name").(string))
}

func setPolicyByName(d *schema.ResourceData, m interface{}, name string) error {
	ptype := d.Get("type").(string)
	policies, _, err := getOktaClientFromMetadata(m).Policy.ListPolicies(context.Background(), &query.Params{Type: ptype})
	if err != nil {
		return fmt.Errorf("failed to list policies: %v", err)
	}
	for _, policy := range policies {
		if policy.Name == name {
			d.SetId(policy.Id)
			return nil
		}
	}
	return fmt.Errorf("no policies retrieved for policy type %v, name %s", ptype, name)
}
