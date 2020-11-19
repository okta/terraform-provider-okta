package okta

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicySignOn() *schema.Resource {
	return &schema.Resource{
		Exists: resourcePolicyExists,
		Create: resourcePolicySignOnCreate,
		Read:   resourcePolicySignOnRead,
		Update: resourcePolicySignOnUpdate,
		Delete: resourcePolicySignOnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: basePolicySchema,
	}
}

func resourcePolicySignOnCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	log.Printf("[INFO] Creating Policy %v", d.Get("name").(string))
	template := buildSignOnPolicy(d)
	err := createPolicy(d, m, template)
	if err != nil {
		return err
	}
	return resourcePolicySignOnRead(d, m)
}

func resourcePolicySignOnRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy %v", d.Get("name").(string))
	policy, err := getPolicy(d, m)
	if policy == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	return syncPolicyFromUpstream(d, policy)
}

func resourcePolicySignOnUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	log.Printf("[INFO] Update Policy %v", d.Get("name").(string))
	template := buildSignOnPolicy(d)
	err := updatePolicy(d, m, template)
	if err != nil {
		return err
	}
	return resourcePolicySignOnRead(d, m)
}

func resourcePolicySignOnDelete(d *schema.ResourceData, m interface{}) error {
	return deletePolicy(d, m)
}

// create or update a sign on policy
func buildSignOnPolicy(d *schema.ResourceData) sdk.Policy {
	template := sdk.SignOnPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &okta.PolicyRuleConditions{
		People: getGroups(d),
	}
	return template
}
