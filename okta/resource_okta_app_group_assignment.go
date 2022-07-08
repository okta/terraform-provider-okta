package okta

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppGroupAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppGroupAssignmentCreate,
		ReadContext:   resourceAppGroupAssignmentRead,
		DeleteContext: resourceAppGroupAssignmentDelete,
		UpdateContext: resourceAppGroupAssignmentUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, errors.New("invalid resource import specifier. Use: terraform import <app_id>/<group_id>")
				}
				_ = d.Set("app_id", parts[0])
				_ = d.Set("group_id", parts[1])
				_ = d.Set("retain_assignment", false)
				assignment, _, err := getOktaClientFromMetadata(m).Application.
					GetApplicationGroupAssignment(ctx, parts[0], parts[1], nil)
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
			"retain_assignment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Retain the group assignment on destroy. If set to true, the resource will be removed from state but not from the Okta app.",
			},
		},
	}
}

func resourceAppGroupAssignmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = context.WithValue(ctx, retryOnStatusCodes, []int{http.StatusInternalServerError})
	assignment, _, err := getOktaClientFromMetadata(m).Application.CreateApplicationGroupAssignment(
		ctx,
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		buildAppGroupAssignment(d),
	)
	if err != nil {
		return diag.Errorf("failed to create application group assignment: %v", err)
	}
	d.SetId(assignment.Id)
	return resourceAppGroupAssignmentRead(ctx, d, m)
}

func resourceAppGroupAssignmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Create actually does a PUT
	_, _, err := getOktaClientFromMetadata(m).Application.CreateApplicationGroupAssignment(
		ctx,
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		buildAppGroupAssignment(d),
	)
	if err != nil {
		return diag.Errorf("failed to update application group assignment: %v", err)
	}
	return resourceAppGroupAssignmentRead(ctx, d, m)
}

func resourceAppGroupAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g, resp, err := getOktaClientFromMetadata(m).Application.GetApplicationGroupAssignment(
		ctx,
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		nil,
	)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get application group assignment: %v", err)
	}
	if g == nil {
		d.SetId("")
		return nil
	}
	jsonProfile, err := json.Marshal(g.Profile)
	if err != nil {
		return diag.Errorf("failed to marshal app user profile to JSON: %v", err)
	}
	_ = d.Set("profile", string(jsonProfile))
	_ = d.Set("priority", g.Priority)
	return nil
}

func resourceAppGroupAssignmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	retain := d.Get("retain_assignment").(bool)
	if retain {
		// The assignment should be retained, bail before DeleteApplicationGroupAssignment is called
		return nil
	}

	_, err := getOktaClientFromMetadata(m).Application.DeleteApplicationGroupAssignment(
		ctx,
		d.Get("app_id").(string),
		d.Get("group_id").(string),
	)
	if err != nil {
		return diag.Errorf("failed to delete application group assignment: %v", err)
	}
	return nil
}

func buildAppGroupAssignment(d *schema.ResourceData) okta.ApplicationGroupAssignment {
	var profile interface{}
	rawProfile := d.Get("profile").(string)
	_ = json.Unmarshal([]byte(rawProfile), &profile)
	assignment := okta.ApplicationGroupAssignment{
		Profile: profile,
	}
	p, ok := d.GetOk("priority")
	if ok {
		assignment.Priority = int64(p.(int))
	}
	return assignment
}
