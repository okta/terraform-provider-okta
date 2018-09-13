package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSignOnPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSignOnPolicyCreate,
		Read:   resourceSignOnPolicyRead,
		Update: resourceSignOnPolicyUpdate,
		Delete: resourceSignOnPolicyDelete,

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			// user cannot edit a default policy
			// editing default Password Policies not supported in the Okta api
			// please upvote the support request here: https://support.okta.com/help/ideas/viewIdea.apexp?id=0870Z000000SS6mQAG
			if d.Get("name").(string) == "Default Policy" {
				return fmt.Errorf("You cannot edit a default Policy")
			}

			return nil
		},

		Schema: basePolicySchema,
	}
}

func resourceSignOnPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Policy %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	var policyID string
	exists := false
	currentPolicies, _, err := client.Policies.GetPoliciesByType(signOnPolicyType)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if currentPolicies != nil {
		for _, policy := range currentPolicies.Policies {
			if policy.Name == d.Get("name").(string) {
				policyID = policy.ID
				exists = true
				break
			}
		}
	}

	if exists == true {
		log.Printf("[INFO] Policy %v already exists in Okta. Adding to Terraform.", d.Get("name").(string))
		d.SetId(policyID)
	} else {
		template := buildSignOnPolicy(d, m)
		err = createPolicy(d, m, template)
		if err != nil {
			return err
		}
	}

	return resourceSignOnPolicyRead(d, m)
}

func resourceSignOnPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy %v", d.Get("name").(string))

	policy, err := getPolicy(d, m)
	if err != nil {
		return err
	}
	if policy == nil {
		// if the policy does not exist in Okta, delete from terraform state
		d.SetId("")
		return nil
	}

	return syncPolicyFromUpstream(d, policy)
}

func resourceSignOnPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Policy %v", d.Get("name").(string))
	d.Partial(true)

	policy, err := getPolicy(d, m)
	if err != nil {
		return err
	}
	if policy.ID != "" {
		template := buildSignOnPolicy(d, m)
		err = updatePolicy(d, m, template)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}
	d.Partial(false)

	return resourceSignOnPolicyRead(d, m)
}

func resourceSignOnPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Policy %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	policy, err := getPolicy(d, m)
	if err != nil {
		return err
	}
	if policy.ID != "" {
		if policy.System == true {
			log.Printf("[INFO] Policy %v is a System Policy, cannot delete from Okta", d.Get("name").(string))
		} else {
			_, err = client.Policies.DeletePolicy(d.Id())
			if err != nil {
				return fmt.Errorf("[ERROR] Error Deleting Policy from Okta: %v", err)
			}
		}
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}
	// remove the policy resource from terraform
	d.SetId("")

	return nil
}

// create or update a signon policy
func buildSignOnPolicy(d *schema.ResourceData, m interface{}) articulateOkta.Policy {
	client := m.(*Config).articulateOktaClient

	template := client.Policies.SignOnPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	template.Type = signOnPolicyType
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}

	template.Conditions = &articulateOkta.PolicyConditions{
		People: getGroups(d),
	}

	return template
}
