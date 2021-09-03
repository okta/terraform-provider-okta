package okta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppGroupAssignments() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppGroupAssignmentsCreate,
		ReadContext:   resourceAppGroupAssignmentsRead,
		DeleteContext: resourceAppGroupAssignmentsDelete,
		UpdateContext: resourceAppGroupAssignmentsUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("app_id", d.Id())

				client := getOktaClientFromMetadata(m)

				assignments, err := listApplicationGroupAssignments(
					ctx,
					client,
					d.Get("app_id").(string),
				)
				if err != nil {
					return nil, err
				}

				tfFlattenedAssignments := make([]interface{}, len(assignments))
				for i, assignment := range assignments {
					tfAssignment, err := groupAssignmentToTFGroup(assignment)
					if err != nil {
						return nil, err
					}
					tfFlattenedAssignments[i] = tfAssignment
				}

				err = d.Set("group", tfFlattenedAssignments)
				if err != nil {
					return nil, err
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "A group to assign to this application",
				MinItems:    1,
				Set: func(v interface{}) int {
					buf := bytes.NewBuffer(nil)
					group := v.(map[string]interface{})

					buf.WriteString(fmt.Sprintf("%s-", group["id"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", normalizeDataJSON(group["profile"])))
					buf.WriteString(fmt.Sprintf("%d-", group["priority"].(int)))

					return schema.HashString(buf.String())
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A group to associate with the application",
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								p, n := d.GetChange("priority")
								return p == n && new == "0"
							},
						},
						"profile": {
							Type:             schema.TypeString,
							ValidateDiagFunc: stringIsJSON,
							StateFunc:        normalizeDataJSON,
							Optional:         true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
						},
					},
				},
			},
		},
	}
}

func resourceAppGroupAssignmentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groups := d.Get("group").(*schema.Set).List()
	client := getOktaClientFromMetadata(m)

	assignments := tfGroupsToGroupAssignments(groups...)

	// run through all groups in the set and create an assignment
	for groupID, assignment := range assignments {
		_, _, err := client.Application.CreateApplicationGroupAssignment(
			ctx,
			d.Get("app_id").(string),
			groupID,
			assignment,
		)
		if err != nil {
			return diag.Errorf("failed to create application group assignment: %v", err)
		}
	}

	// okta_app_group_assignments completely control all assignments for an application
	d.SetId(d.Get("app_id").(string))
	return resourceAppGroupAssignmentsRead(ctx, d, m)
}

func resourceAppGroupAssignmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)

	assignments, err := listApplicationGroupAssignments(
		ctx,
		client,
		d.Get("app_id").(string),
	)
	if err != nil {
		return diag.Errorf("failed to fetch group assignments: %v", err)
	}

	tfFlattenedAssignments := make([]interface{}, len(assignments))
	for i, assignment := range assignments {
		tfAssignment, err := groupAssignmentToTFGroup(assignment)
		if err != nil {
			return diag.Errorf("failed to marshal group profile: %v", err)
		}
		tfFlattenedAssignments[i] = tfAssignment
	}

	err = d.Set("group", tfFlattenedAssignments)
	if err != nil {
		return diag.Errorf("failed to set groups in tf state: %v", err)
	}
	return nil
}

func resourceAppGroupAssignmentsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)

	for _, rawGroup := range d.Get("group").(*schema.Set).List() {
		group := rawGroup.(map[string]interface{})

		_, err := client.Application.DeleteApplicationGroupAssignment(
			ctx,
			d.Get("app_id").(string),
			group["id"].(string),
		)
		if err != nil {
			return diag.Errorf("failed to delete application group assignment: %v", err)
		}
	}
	return nil
}

func resourceAppGroupAssignmentsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	appID := d.Get("app_id").(string)

	old, new := d.GetChange("group")
	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

	toAdd := tfGroupsToGroupAssignments(
		newSet.Difference(oldSet).List()...,
	)
	toRemove := tfGroupsToGroupAssignments(
		oldSet.Difference(newSet).List()...,
	)

	err := deleteGroupAssignments(
		client.Application.DeleteApplicationGroupAssignment,
		ctx,
		appID,
		toRemove,
	)
	if err != nil {
		return diag.Errorf("failed to delete group assignment: %v", err)
	}

	err = addGroupAssignments(
		client.Application.CreateApplicationGroupAssignment,
		ctx,
		appID,
		toAdd,
	)
	if err != nil {
		return diag.Errorf("failed to add group assignment: %v", err)
	}
	return resourceAppGroupAssignmentsRead(ctx, d, m)
}

// groupAssignmentToTFGroup
func groupAssignmentToTFGroup(assignment *okta.ApplicationGroupAssignment) (map[string]interface{}, error) {
	profile := "{}"

	jsonProfile, err := json.Marshal(assignment.Profile)
	if err != nil {
		return nil, err
	}
	if string(jsonProfile) != "" {
		profile = string(jsonProfile)
	}

	tfAssignment := map[string]interface{}{
		"id":       assignment.Id,
		"priority": assignment.Priority,
		"profile":  profile,
	}
	return tfAssignment, nil
}

func tfGroupsToGroupAssignments(groups ...interface{}) map[string]okta.ApplicationGroupAssignment {
	assignments := map[string]okta.ApplicationGroupAssignment{}
	// run through all groups in the set and create an assignment
	for _, untypedGroup := range groups {
		group := untypedGroup.(map[string]interface{})

		id := group["id"].(string)
		// skip empty groups with no id
		if id == "" {
			continue
		}
		priority := group["priority"].(int)

		rawProfile := group["profile"]
		var profile interface{}
		_ = json.Unmarshal([]byte(rawProfile.(string)), &profile)

		assignments[id] = okta.ApplicationGroupAssignment{
			Profile:  profile,
			Priority: int64(priority),
			Id:       id,
		}
	}
	return assignments
}

// addGroupAssignments adds all group assignments
func addGroupAssignments(
	add func(context.Context, string, string, okta.ApplicationGroupAssignment) (*okta.ApplicationGroupAssignment, *okta.Response, error),
	ctx context.Context,
	appID string,
	assignments map[string]okta.ApplicationGroupAssignment,
) error {
	for groupID, assignment := range assignments {
		_, _, err := add(ctx, appID, groupID, assignment)
		if err != nil {
			return err
		}
	}
	return nil
}

// deleteGroupAssignments deletes all group assignments
func deleteGroupAssignments(
	delete func(context.Context, string, string) (*okta.Response, error),
	ctx context.Context,
	appID string,
	assignments map[string]okta.ApplicationGroupAssignment,
) error {
	for groupID := range assignments {
		_, err := delete(ctx, appID, groupID)
		if err != nil {
			return fmt.Errorf(
				"could not delete assignment for group %s, to application %s: %w",
				groupID,
				appID,
				err,
			)
		}
	}
	return nil
}
