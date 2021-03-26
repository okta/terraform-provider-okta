package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

// data source to retrieve information on a Default Policy
func dataSourceDefaultPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDefaultPolicyRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type: schema.TypeString,
				ValidateDiagFunc: stringInSlice([]string{
					sdk.SignOnPolicyType,
					sdk.PasswordPolicyType,
					sdk.MfaPolicyType,
					sdk.IdpDiscoveryType,
				}),
				Description: fmt.Sprintf("Policy type: %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.IdpDiscoveryType),
				Required:    true,
			},
		},
	}
}

func dataSourceDefaultPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policyType := d.Get("type").(string)
	var name string
	if policyType == sdk.IdpDiscoveryType {
		name = "Idp Discovery Policy"
	} else {
		name = "Default Policy"
	}
	policy, err := findPolicy(ctx, m, name, policyType)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policy.Id)
	return nil
}
