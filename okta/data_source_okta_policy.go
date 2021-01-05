package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func dataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of policy",
				Required:    true,
			},
			"type": {
				Type: schema.TypeString,
				ValidateDiagFunc: stringInSlice([]string{
					sdk.SignOnPolicyType,
					sdk.PasswordPolicyType,
					sdk.MfaPolicyType,
					sdk.IdpDiscoveryType,
					sdk.OauthAuthorizationPolicyType,
				}),
				Description: fmt.Sprintf("Policy type: %s, %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.IdpDiscoveryType, sdk.OauthAuthorizationPolicyType),
				Required:    true,
			},
		},
	}
}

func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return setPolicyByName(ctx, d, m, d.Get("name").(string))
}

func setPolicyByName(ctx context.Context, d *schema.ResourceData, m interface{}, name string) diag.Diagnostics {
	policyType := d.Get("type").(string)
	policies, _, err := getOktaClientFromMetadata(m).Policy.ListPolicies(ctx, &query.Params{Type: policyType})
	if err != nil {
		return diag.Errorf("failed to list policies: %v", err)
	}
	for _, policy := range policies {
		if policy.Name == name {
			d.SetId(policy.Id)
			return nil
		}
	}
	return diag.Errorf("no policies retrieved for policy type %s, name %s", policyType, name)
}
