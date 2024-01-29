package okta

import (
	"context"
	"fmt"
	"strings"

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
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Policy type: %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.IdpDiscoveryType),
				Required:    true,
			},
		},
		Description: "Get a Default policy from Okta.",
	}
}

func dataSourceDefaultPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policyType := d.Get("type").(string)
	policies, err := findSystemPolicyByType(ctx, m, policyType)
	if err != nil {
		return diag.FromErr(err)
	}
	var policy *sdk.Policy
	for _, p := range policies {
		if strings.Contains(p.Name, "Default") || strings.Contains(p.Description, "default") {
			policy = p
			break
		}
	}
	if policy == nil {
		return diag.FromErr(fmt.Errorf("cannot find default %v policy", policyType))
	}
	d.SetId(policy.Id)
	return nil
}
