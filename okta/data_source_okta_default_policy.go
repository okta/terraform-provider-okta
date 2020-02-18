package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// data source to retrieve information on a Default Policy

func dataSourceDefaultPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDefaultPolicyRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{signOnPolicyType, passwordPolicyType, "MFA_ENROLL", "OAUTH_AUTHORIZATION_POLICY"}, false),
				Description:  fmt.Sprintf("Policy type: %s, %s, MFA_ENROLL, or OAUTH_AUTHORIZATION_POLICY", signOnPolicyType, passwordPolicyType),
				Required:     true,
			},
		},
	}
}

func dataSourceDefaultPolicyRead(d *schema.ResourceData, m interface{}) error {
	return setPolicyByName(d, m, "Default Policy")
}
