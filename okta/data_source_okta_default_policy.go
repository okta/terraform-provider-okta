package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

// data source to retrieve information on a Default Policy

func dataSourceDefaultPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDefaultPolicyRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.OauthAuthorizationPolicyType}, false),
				Description:  fmt.Sprintf("Policy type: %s, %s, %s, or %s", sdk.SignOnPolicyType, sdk.PasswordPolicyType, sdk.MfaPolicyType, sdk.OauthAuthorizationPolicyType),
				Required:     true,
			},
		},
	}
}

func dataSourceDefaultPolicyRead(d *schema.ResourceData, m interface{}) error {
	return setPolicyByName(d, m, "Default Policy")
}
