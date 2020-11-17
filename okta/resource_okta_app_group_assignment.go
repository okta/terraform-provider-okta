package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppGroupAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppGroupAssignmentCreate,
		Exists: resourceAppGroupAssignmentExists,
		Read:   resourceAppGroupAssignmentRead,
		Delete: resourceAppGroupAssignmentDelete,
		Update: resourceAppGroupAssignmentUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("invalid resource import specifier. Use: terraform import <app_id>/<group_id>")
				}

				_ = d.Set("app_id", parts[0])
				_ = d.Set("group_id", parts[1])

				assignment, _, err := getOktaClientFromMetadata(m).Application.
					GetApplicationGroupAssignment(context.Background(), parts[0], parts[1], nil)

				if err != nil {
					return nil, err
				}

				d.SetId(assignment.Id)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App to associate group with",
				ForceNew:    true,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group associated with the application",
				ForceNew:    true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"profile": {
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
		context.Background(),
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)

	if resp != nil && is404(resp.StatusCode) {
		return false, nil
	}

	return err == nil, err
}

func getAppGroupAssignment(d *schema.ResourceData) okta.ApplicationGroupAssignment {
	var profile interface{}

	rawProfile := d.Get("profile").(string)
	// JSON is already validated
	_ = json.Unmarshal([]byte(rawProfile), &profile)
	priority := d.Get("priority").(int)

	return okta.ApplicationGroupAssignment{
		Profile:  profile,
		Priority: int64(priority),
	}
}

func resourceAppGroupAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	assignment, _, err := getOktaClientFromMetadata(m).Application.CreateApplicationGroupAssignment(
		context.Background(),
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
		context.Background(),
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
		context.Background(),
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	jsonProfile, err := json.Marshal(g.Profile)
	if err != nil {
		return fmt.Errorf("failed to marshal app user profile to JSON, error: %s", err)
	}

	_ = d.Set("profile", string(jsonProfile))
	_ = d.Set("priority", g.Priority)

	return nil
}

func resourceAppGroupAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationGroupAssignment(
		context.Background(),
		d.Get("app_id").(string),
		d.Get("group_id").(string),
	)
	return err
}
