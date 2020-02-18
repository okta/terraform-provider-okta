package okta

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceAppGroupAssignment() *schema.Resource {
	return &schema.Resource{
		// No point in having an exist function, since only the group has to exist
		Create: resourceAppGroupAssignmentCreate,
		Exists: resourceAppGroupAssignmentExists,
		Read:   resourceAppGroupAssignmentRead,
		Delete: resourceAppGroupAssignmentDelete,
		Update: resourceAppGroupAssignmentUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("Invalid resource import specifier. Use: terraform import <app_id>/<group_id>")
				}

				d.Set("app_id", parts[0])
				d.Set("group_id", parts[1])

				assignment, _, err := getOktaClientFromMetadata(m).Application.
					GetApplicationGroupAssignment(parts[0], parts[1], nil)

				if err != nil {
					return nil, err
				}

				d.SetId(assignment.Id)

				return []*schema.ResourceData{d}, nil
			},
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
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"profile": &schema.Schema{
				Type:      schema.TypeString,
				StateFunc: normalizeDataJSON,
				Optional:  true,
				Default:   "{}",
			},
		},
	}
}

func resourceAppGroupAssignmentExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getOktaClientFromMetadata(m)
	_, resp, err := client.Application.GetApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)

	if is404(resp.StatusCode) {
		return false, nil
	}

	return err == nil, err
}

func getAppGroupAssignment(d *schema.ResourceData) okta.ApplicationGroupAssignment {
	var profile interface{}

	rawProfile := d.Get("profile").(string)
	// JSON is already validated
	json.Unmarshal([]byte(rawProfile), &profile)
	priority := d.Get("priority").(int)

	return okta.ApplicationGroupAssignment{
		Profile:  profile,
		Priority: int64(priority),
	}
}

func resourceAppGroupAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	assignment, _, err := getOktaClientFromMetadata(m).Application.CreateApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		getAppGroupAssignment(d),
	)

	if err != nil {
		return err
	}

	d.SetId(assignment.Id)

	return resourceAppGroupAssignmentRead(d, m)
}

func resourceAppGroupAssignmentUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	// Create actually does a PUT
	_, _, err := client.Application.CreateApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		getAppGroupAssignment(d),
	)

	if err != nil {
		return err
	}

	return resourceAppGroupAssignmentRead(d, m)
}

func resourceAppGroupAssignmentRead(d *schema.ResourceData, m interface{}) error {
	g, resp, err := getOktaClientFromMetadata(m).Application.GetApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	jsonProfile, err := json.Marshal(g.Profile)
	if err != nil {
		return fmt.Errorf("Failed to marshal app user profile to JSON, error: %s", err)
	}

	d.Set("profile", string(jsonProfile))
	d.Set("priority", g.Priority)

	return nil
}

func resourceAppGroupAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationGroupAssignment(
		d.Get("app_id").(string),
		d.Get("group_id").(string),
	)
	return err
}
