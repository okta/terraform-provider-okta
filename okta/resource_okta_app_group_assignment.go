package okta

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceAppGroupAssignment() *schema.Resource {
	return &schema.Resource{
		// No point in having an exist function, since only the group has to exist
		Create: resourceAppGroupAssignmentCreate,
		Exists: resourceAppGroupAssignmentExists,
		Read:   resourceAppGroupAssignmentRead,
		Delete: resourceAppGroupAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "App to associate group with",
				ForceNew:    true,
			},
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group associated with the application",
				ForceNew:    true,
			},
			"profile": &schema.Schema{
				Type:      schema.TypeString,
				StateFunc: normalizeDataJSON,
				Computed:  true,
			},
		},
	}
}

func resourceAppGroupAssignmentExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getOktaClientFromMetadata(m)
	g, _, err := client.Application.GetApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)

	return g != nil, err
}

func resourceAppGroupAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	assignment, _, err := getOktaClientFromMetadata(m).Application.CreateApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		okta.ApplicationGroupAssignment{},
	)

	if err != nil {
		return err
	}

	d.SetId(assignment.Id)

	return resourceAppGroupAssignmentRead(d, m)
}

func resourceAppGroupAssignmentRead(d *schema.ResourceData, m interface{}) error {
	g, _, err := getOktaClientFromMetadata(m).Application.GetApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)

	if err != nil {
		return err
	}

	jsonProfile, err := json.Marshal(g.Profile)
	if err != nil {
		return fmt.Errorf("Failed to marshal app user profile to JSON, error: %s", err)
	}

	d.Set("profile", string(jsonProfile))

	return nil
}

func resourceAppGroupAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
	)
	return err
}
