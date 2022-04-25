package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

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
				Type:        schema.TypeList,
				Required:    true,
				Description: "A group to assign to this application",
				MinItems:    1,
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
							Required:         true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == ""
							},
							DefaultFunc: func() (interface{}, error) {
								return "{}", nil
							},
						},
					},
				},
			},
		},
	}
}

func resourceAppGroupAssignmentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	assignments := tfGroupsToGroupAssignments(d)

	// run through all groups in the set and create an assignment
	for i := range assignments {
		_, _, err := client.Application.CreateApplicationGroupAssignment(
			ctx,
			d.Get("app_id").(string),
			assignments[i].Id,
			*assignments[i],
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
	assignments, resp, err := listApplicationGroupAssignments(
		ctx,
		client,
		d.Get("app_id").(string),
	)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to fetch group assignments: %v", err)
	}
	if assignments == nil {
		d.SetId("")
		return nil
	}
	g, ok := d.GetOk("group")
	if ok {
		err := setNonPrimitives(d, map[string]interface{}{"group": syncGroups(d, g.([]interface{}), assignments)})
		if err != nil {
			return diag.Errorf("failed to set OAuth application properties: %v", err)
		}
	} else {
		arr := make([]map[string]interface{}, len(assignments))
		for i := range assignments {
			arr[i] = groupAssignmentToTFGroup(assignments[i])
		}
		err := setNonPrimitives(d, map[string]interface{}{"group": arr})
		if err != nil {
			return diag.Errorf("failed to set OAuth application properties: %v", err)
		}
	}
	return nil
}

func resourceAppGroupAssignmentsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	appID := d.Get("app_id").(string)
	assignments, _, err := listApplicationGroupAssignments(
		ctx,
		client,
		d.Get("app_id").(string),
	)
	if err != nil {
		return diag.Errorf("failed to fetch group assignments: %v", err)
	}
	toAssign, toRemove := splitAssignmentsTargets(tfGroupsToGroupAssignments(d), assignments)
	err = deleteGroupAssignments(
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
		toAssign,
	)
	if err != nil {
		return diag.Errorf("failed to add/update group assignment: %v", err)
	}
	return resourceAppGroupAssignmentsRead(ctx, d, m)
}

func resourceAppGroupAssignmentsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	for _, rawGroup := range d.Get("group").([]interface{}) {
		group := rawGroup.(map[string]interface{})
		resp, err := client.Application.DeleteApplicationGroupAssignment(
			ctx,
			d.Get("app_id").(string),
			group["id"].(string),
		)
		if err := suppressErrorOn404(resp, err); err != nil {
			return diag.Errorf("failed to delete application group assignment: %v", err)
		}
	}
	return nil
}

func syncGroups(d *schema.ResourceData, groups []interface{}, assignments []*okta.ApplicationGroupAssignment) []interface{} {
	var newGroups []interface{}
	for i := range groups {
		present := false
		for _, assignment := range assignments {
			if assignment.Id == d.Get(fmt.Sprintf("group.%d.id", i)).(string) {
				present = true
				if assignment.Priority >= 0 {
					groups[i].(map[string]interface{})["priority"] = assignment.Priority
				}
				groups[i].(map[string]interface{})["profile"] = buildProfile(d, i, assignment)
			}
		}
		if present {
			newGroups = append(newGroups, groups[i])
		}
	}
	return newGroups
}

func buildProfile(d *schema.ResourceData, i int, assignment *okta.ApplicationGroupAssignment) string {
	oldProfile := d.Get(fmt.Sprintf("group.%d.profile", i)).(string)
	opm := make(map[string]interface{})
	_ = json.Unmarshal([]byte(oldProfile), &opm)
	ap, ok := assignment.Profile.(map[string]interface{})
	if ok {
		for k, v := range ap {
			if _, ok := opm[k]; ok {
				opm[k] = v
				continue
			}
			if v != nil {
				opm[k] = v
			}
		}
	}
	jsonProfile, _ := json.Marshal(&opm)
	return string(jsonProfile)
}

func splitAssignmentsTargets(expectedAssignments, existingAssignments []*okta.ApplicationGroupAssignment) (toAssign, toRemove []*okta.ApplicationGroupAssignment) {
	for i := range expectedAssignments {
		if !containsEqualAssignment(existingAssignments, expectedAssignments[i]) {
			toAssign = append(toAssign, expectedAssignments[i])
		}
	}
	for i := range existingAssignments {
		if !containsAssignment(expectedAssignments, existingAssignments[i]) {
			toRemove = append(toRemove, existingAssignments[i])
		}
	}
	return
}

func containsAssignment(assignments []*okta.ApplicationGroupAssignment, assignment *okta.ApplicationGroupAssignment) bool {
	for i := range assignments {
		if assignments[i].Id == assignment.Id {
			return true
		}
	}
	return false
}

func containsEqualAssignment(assignments []*okta.ApplicationGroupAssignment, assignment *okta.ApplicationGroupAssignment) bool {
	for i := range assignments {
		if assignments[i].Id == assignment.Id && reflect.DeepEqual(assignments[i].Profile, assignment.Profile) {
			if assignment.Priority >= 0 {
				return reflect.DeepEqual(assignments[i].Priority, assignment.Priority)
			}
			return true
		}
	}
	return false
}

func groupAssignmentToTFGroup(assignment *okta.ApplicationGroupAssignment) map[string]interface{} {
	jsonProfile, _ := json.Marshal(assignment.Profile)
	profile := "{}"
	if string(jsonProfile) != "" {
		profile = string(jsonProfile)
	}
	return map[string]interface{}{
		"id":       assignment.Id,
		"priority": assignment.Priority,
		"profile":  profile,
	}
}

func tfGroupsToGroupAssignments(d *schema.ResourceData) []*okta.ApplicationGroupAssignment {
	assignments := make([]*okta.ApplicationGroupAssignment, len(d.Get("group").([]interface{})))
	for i := range d.Get("group").([]interface{}) {
		rawProfile := d.Get(fmt.Sprintf("group.%d.profile", i))
		var profile interface{}
		_ = json.Unmarshal([]byte(rawProfile.(string)), &profile)
		a := &okta.ApplicationGroupAssignment{
			Id:      d.Get(fmt.Sprintf("group.%d.id", i)).(string),
			Profile: profile,
		}
		priority, ok := d.GetOk(fmt.Sprintf("group.%d.priority", i))
		if ok {
			a.Priority = int64(priority.(int))
		}
		assignments[i] = a
	}
	return assignments
}

// addGroupAssignments adds all group assignments
func addGroupAssignments(
	add func(context.Context, string, string, okta.ApplicationGroupAssignment) (*okta.ApplicationGroupAssignment, *okta.Response, error),
	ctx context.Context,
	appID string,
	assignments []*okta.ApplicationGroupAssignment,
) error {
	for _, assignment := range assignments {
		_, _, err := add(ctx, appID, assignment.Id, *assignment)
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
	assignments []*okta.ApplicationGroupAssignment,
) error {
	for i := range assignments {
		_, err := delete(ctx, appID, assignments[i].Id)
		if err != nil {
			return fmt.Errorf("could not delete assignment for group %s, to application %s: %w", assignments[i].Id, appID, err)
		}
	}
	return nil
}
