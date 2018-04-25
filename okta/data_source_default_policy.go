package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// data source to retrieve information on a Default Policy

func dataSourceDefaultPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDefaultPolicyRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"OKTA_SIGN_ON", "PASSWORD", "MFA_ENROLL", "OAUTH_AUTHORIZATION_POLICY"}, false),
				Description:  "Policy type: OKTA_SIGN_ON, PASSWORD, MFA_ENROLL, or OAUTH_AUTHORIZATION_POLICY",
				Required:     true,
			},
		},
	}
}

func dataSourceDefaultPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Data Source Default Policy Read %v", d.Get("type").(string))
	client := m.(*Config).oktaClient

	set := false
	currentPolicies, _, err := client.Policies.GetPoliciesByType(d.Get("type").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policies in Okta: %v", err)
	}
	if currentPolicies != nil {
		for _, policy := range currentPolicies.Policies {
			if policy.Name == "Default Policy" {
				d.SetId(policy.ID)
				set = true
				break
			}
		}
		if set == false {
			return fmt.Errorf("[ERROR] Unable to retrieve Default Policy for type %v", d.Get("type").(string))
		}
	} else {
		return fmt.Errorf("[ERROR] No policies retrieved for policy type %v", d.Get("type").(string))
	}
	return nil
}
