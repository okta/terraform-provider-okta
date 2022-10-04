package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

// data source to retrieve information on a Default Policy
func dataSourceDefaultPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDefaultPolicyRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type: schema.TypeString,
				ValidateDiagFunc: elemInSlice([]string{
					sdk.SignOnPolicyType,
					sdk.PasswordPolicyType,
					sdk.MfaPolicyType,
					sdk.IdpDiscoveryType,
					sdk.AccessPolicyType,
					sdk.ProfileEnrollmentPolicyType,
				}),
				Description: fmt.Sprintf("Policy type: %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.IdpDiscoveryType),
				Required:    true,
			},
		},
	}
}

func dataSourceDefaultPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policyType := d.Get("type").(string)
	policy, err := findSystemPolicyByType(ctx, m, policyType)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policy.Id)
	return nil
}
