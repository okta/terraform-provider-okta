package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
					sdk.OauthAuthorizationPolicyType,
				}),
				Description: fmt.Sprintf("Policy type: %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.OauthAuthorizationPolicyType),
				Required:    true,
			},
		},
	}
}

func dataSourceDefaultPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return setPolicyByName(ctx, d, m, "Default Policy")
}
