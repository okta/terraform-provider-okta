package okta

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

// Basis of policy schema
var (
	basePolicySchema = map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Policy Name",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Policy Description",
		},
		"priority": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Policy Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.",
			// Suppress diff if config is empty.
			DiffSuppressFunc: createValueDiffSuppression("0"),
		},
		"status": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      statusActive,
			ValidateFunc: validation.StringInSlice([]string{statusActive, statusInactive}, false),
			Description:  "Policy Status: ACTIVE or INACTIVE.",
		},
		"groups_included": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "List of Group IDs to Include",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	// Pattern used in a few spots, whitelisting/blacklisting users and groups
	peopleSchema = map[string]*schema.Schema{
		"user_whitelist": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"user_blacklist": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"group_whitelist": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"group_blacklist": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
	}

	statusSchema = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      statusActive,
		ValidateFunc: validation.StringInSlice([]string{statusActive, statusInactive}, false),
	}
)

func addPeopleAssignments(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(peopleSchema, target)
}

func setPeopleAssignments(d *schema.ResourceData, c *okta.GroupRulePeopleCondition) error {
	// Don't think the API omits these when they are empty thus the unguarded accessing
	return setNonPrimitives(d, map[string]interface{}{
		"group_whitelist": convertStringSetToInterface(c.Groups.Include),
		"group_blacklist": convertStringSetToInterface(c.Groups.Exclude),
		"user_whitelist":  convertStringSetToInterface(c.Users.Include),
		"user_blacklist":  convertStringSetToInterface(c.Users.Exclude),
	})
}

func getPeopleConditions(d *schema.ResourceData) *okta.GroupRulePeopleCondition {
	return &okta.GroupRulePeopleCondition{
		Groups: &okta.GroupRuleGroupCondition{
			Include: convertInterfaceToStringSet(d.Get("group_whitelist")),
			Exclude: convertInterfaceToStringSet(d.Get("group_blacklist")),
		},
		Users: &okta.GroupRuleUserCondition{
			Include: convertInterfaceToStringSet(d.Get("user_whitelist")),
			Exclude: convertInterfaceToStringSet(d.Get("user_blacklist")),
		},
	}
}

func buildPolicySchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(basePolicySchema, target)
}

func createPolicy(d *schema.ResourceData, m interface{}, template sdk.Policy) error {
	client := getSupplementFromMetadata(m)
	policy, _, err := client.CreatePolicy(context.Background(), template)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating Policy: %v", err)
	}
	log.Printf("[INFO] Okta Policy Created: %+v. Adding Policy to Terraform.", policy)
	d.SetId(policy.Id)

	// Even if priority is invalid we want to add the policy to Terraform to reflect upstream.
	err = validatePriority(template.Priority, policy.Priority)
	if err != nil {
		return err
	}

	return policyActivate(d, m)
}

func ensureNotDefaultPolicy(d *schema.ResourceData) error {
	return ensureNotDefault(d, "Policy")
}

func getGroups(d *schema.ResourceData) *okta.PolicyPeopleCondition {
	var people *okta.PolicyPeopleCondition
	if include, ok := d.GetOk("groups_included"); ok {
		people = &okta.PolicyPeopleCondition{
			Groups: &okta.GroupCondition{
				Include: convertInterfaceToStringSet(include),
			},
		}
	}
	return people
}

// Grabs policy from upstream, if the resource does not exist the returned policy will be nil which is not considered an error
func getPolicy(d *schema.ResourceData, m interface{}) (*sdk.Policy, error) {
	policy, resp, err := getSupplementFromMetadata(m).GetPolicy(context.Background(), d.Id())
	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil, nil
	}
	return policy, err
}

// activate or deactivate a policy according to the terraform schema status field
func policyActivate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)

	if d.Get("status").(string) == statusActive {
		_, err := client.Policy.ActivatePolicy(context.Background(), d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Activating Policy: %v", err)
		}
	}
	if d.Get("status").(string) == statusInactive {
		_, err := client.Policy.DeactivatePolicy(context.Background(), d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
		}
	}
	return nil
}

func updatePolicy(d *schema.ResourceData, m interface{}, template sdk.Policy) error {
	client := getSupplementFromMetadata(m)
	policy, _, err := client.UpdatePolicy(context.Background(), d.Id(), template)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Policy: %v", err)
	}
	// avoiding perpetual diffs by erroring when the configured priority is not valid and the API defaults it.
	err = validatePriority(template.Priority, policy.Priority)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Okta Policy Updated: %+v", policy)

	return policyActivate(d, m)
}

func deletePolicy(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	log.Printf("[INFO] Delete Policy %v", d.Get("name").(string))
	client := getOktaClientFromMetadata(m)
	_, err := client.Policy.DeletePolicy(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Policy from Okta: %v", err)
	}
	// remove the policy resource from terraform
	d.SetId("")
	return nil
}

func resourcePolicyExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	policy, err := getPolicy(d, m)

	if err != nil || policy == nil {
		return false, err
	}

	return true, nil
}

func syncPolicyFromUpstream(d *schema.ResourceData, policy *sdk.Policy) error {
	_ = d.Set("name", policy.Name)
	_ = d.Set("description", policy.Description)
	_ = d.Set("status", policy.Status)
	_ = d.Set("priority", policy.Priority)

	return setNonPrimitives(d, map[string]interface{}{
		"groups_included": convertStringSetToInterface(policy.Conditions.People.Groups.Include),
	})
}
