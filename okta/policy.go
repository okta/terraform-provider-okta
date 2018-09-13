package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// Basis of policy schema
var basePolicySchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Policy Name",
	},
	"description": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Policy Description",
	},
	// Need to add an enhancement here, maybe we can provide a hierarchy system, priorityBefore = <refPolicy>
	"priority": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Policy Priority",
	},
	"status": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "ACTIVE",
		ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
		Description:  "Policy Status: ACTIVE or INACTIVE.",
	},
	"groups_excluded": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of Group IDs to Include",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"groups_included": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of Group IDs to Include",
		Elem: &schema.Schema{
			Type: schema.TypeString,
			// Suppress diff if config is empty, the API will apply the default.
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == ""
			},
		},
		// Suppress diff if config is empty, the API will apply the default.
		DiffSuppressFunc: suppressDefaultedArrayDiff,
	},
}

func buildPolicySchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	for k := range basePolicySchema {
		target[k] = basePolicySchema[k]
	}

	return target
}

func syncPolicyFromUpstream(d *schema.ResourceData, policy *articulateOkta.Policy) error {
	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("status", policy.Status)
	d.Set("priority", policy.Priority)

	return setNonPrimitives(d, map[string]interface{}{
		"groups_excluded": convertStringArrToInterface(policy.Conditions.People.Groups.Exclude),
		"groups_included": convertStringArrToInterface(policy.Conditions.People.Groups.Include),
	})
}

// Grabs policy from upstream, if the resource does not exist the returned policy will be nil which is not considered an error
func getPolicy(d *schema.ResourceData, m interface{}) (*articulateOkta.Policy, error) {
	client := m.(*Config).articulateOktaClient
	policy, _, err := client.Policies.GetPolicy(d.Id())

	return policy, err
}

func getGroups(d *schema.ResourceData) *articulateOkta.People {
	var people *articulateOkta.People
	include := d.Get("groups_included")
	exclude := d.Get("groups_excluded")

	if include != nil && len(include.([]interface{})) > 1 {
		people = &articulateOkta.People{
			Groups: &articulateOkta.Groups{
				Include: convertInterfaceToStringArr(include),
			},
		}
	} else if exclude != nil && len(exclude.([]interface{})) > 1 {
		people = &articulateOkta.People{
			Groups: &articulateOkta.Groups{
				Exclude: convertInterfaceToStringArr(exclude),
			},
		}
	}

	return people
}

// activate or deactivate a policy according to the terraform schema status field
func policyActivate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).articulateOktaClient

	if d.Get("status").(string) == "ACTIVE" {
		_, err := client.Policies.ActivatePolicy(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Activating Policy: %v", err)
		}
	}
	if d.Get("status").(string) == "INACTIVE" {
		_, err := client.Policies.DeactivatePolicy(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
		}
	}
	return nil
}

func updatePolicy(d *schema.ResourceData, meta interface{}, template articulateOkta.Policy) error {
	client := getClientFromMetadata(meta)
	policy, _, err := client.Policies.UpdatePolicy(d.Id(), template)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Policy: %v", err)
	}
	log.Printf("[INFO] Okta Policy Updated: %+v", policy)

	return policyActivate(d, meta)
}

func createPolicy(d *schema.ResourceData, meta interface{}, template articulateOkta.Policy) error {
	client := getClientFromMetadata(meta)
	policy, _, err := client.Policies.CreatePolicy(template)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating Policy: %v", err)
	}
	log.Printf("[INFO] Okta Policy Created: %+v. Adding Policy to Terraform.", policy)
	d.SetId(policy.ID)

	return policyActivate(d, meta)
}
