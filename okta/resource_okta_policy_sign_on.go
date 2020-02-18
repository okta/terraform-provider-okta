package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePolicySignon() *schema.Resource {
	return &schema.Resource{
		Exists: resourcePolicyExists,
		Create: resourcePolicySignonCreate,
		Read:   resourcePolicySignonRead,
		Update: resourcePolicySignonUpdate,
		Delete: resourcePolicySignonDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: basePolicySchema,
	}
}

func resourcePolicySignonCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	log.Printf("[INFO] Creating Policy %v", d.Get("name").(string))

	template := buildSignOnPolicy(d, m)
	err := createPolicy(d, m, template)
	if err != nil {
		return err
	}

	return resourcePolicySignonRead(d, m)
}

func resourcePolicySignonRead(d *schema.ResourceData, m interface{}) error {
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

func resourcePolicySignonUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	log.Printf("[INFO] Update Policy %v", d.Get("name").(string))
	d.Partial(true)
	template := buildSignOnPolicy(d, m)
	err := updatePolicy(d, m, template)
	if err != nil {
		return err
	}
	d.Partial(false)

	return resourcePolicySignonRead(d, m)
}

func resourcePolicySignonDelete(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	log.Printf("[INFO] Delete Policy %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient
	_, err := client.Policies.DeletePolicy(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Policy from Okta: %v", err)
	}
	// remove the policy resource from terraform
	d.SetId("")

	return nil
}

// create or update a signon policy
func buildSignOnPolicy(d *schema.ResourceData, m interface{}) *articulateOkta.Policy {
	client := m.(*Config).articulateOktaClient

	template := client.Policies.SignOnPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	template.Type = signOnPolicyType
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = priority.(int)
	}

	template.Conditions = &articulateOkta.PolicyConditions{
		People: getGroups(d),
	}

	return &template
}
