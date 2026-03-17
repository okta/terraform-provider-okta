package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
)

func resourceUserRisk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRiskCreate,
		ReadContext:   resourceUserRiskRead,
		UpdateContext: resourceUserRiskUpdate,
		DeleteContext: resourceUserRiskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserRiskImport,
		},
		Description: "Manages a user's risk level in Okta. This resource allows you to set and manage the risk level for a specific user.",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the user.",
			},
			"risk_level": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Risk level of the user. Valid values: `HIGH`, `LOW`.",
				ValidateFunc: validation.StringInSlice([]string{"HIGH", "LOW"}, false),
			},
		},
	}
}

func resourceUserRiskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	riskLevel := d.Get("risk_level").(string)

	logger(meta).Info("setting user risk", "user_id", userId, "risk_level", riskLevel)

	req := v6okta.NewUserRiskRequest()
	req.SetRiskLevel(riskLevel)

	_, _, err := getOktaV6ClientFromMetadata(meta).UserRiskAPI.UpsertUserRisk(ctx, userId).UserRiskRequest(*req).Execute()
	if err != nil {
		return diag.Errorf("failed to set user risk: %v", err)
	}

	d.SetId(userId)
	return resourceUserRiskRead(ctx, d, meta)
}

func resourceUserRiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Id()

	logger(meta).Info("reading user risk", "user_id", userId)

	resp, _, err := getOktaV6ClientFromMetadata(meta).UserRiskAPI.GetUserRisk(ctx, userId).Execute()
	if err != nil {
		return diag.Errorf("failed to get user risk: %v", err)
	}

	var riskLevel string
	if resp.UserRiskLevelExists != nil {
		riskLevel = resp.UserRiskLevelExists.GetRiskLevel()
	} else if resp.UserRiskLevelNone != nil {
		// User has no risk level set - remove from state
		d.SetId("")
		return nil
	}

	_ = d.Set("user_id", userId)
	_ = d.Set("risk_level", riskLevel)
	return nil
}

func resourceUserRiskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Id()
	riskLevel := d.Get("risk_level").(string)

	logger(meta).Info("updating user risk", "user_id", userId, "risk_level", riskLevel)

	req := v6okta.NewUserRiskRequest()
	req.SetRiskLevel(riskLevel)

	_, _, err := getOktaV6ClientFromMetadata(meta).UserRiskAPI.UpsertUserRisk(ctx, userId).UserRiskRequest(*req).Execute()
	if err != nil {
		return diag.Errorf("failed to update user risk: %v", err)
	}

	return resourceUserRiskRead(ctx, d, meta)
}

func resourceUserRiskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// On destroy, we simply remove the resource from Terraform state.
	// The user's risk level will remain at whatever it was last set to.
	logger(meta).Info("removing user risk from state", "user_id", d.Id())
	return nil
}

func resourceUserRiskImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	userId := d.Id()

	if strings.TrimSpace(userId) == "" {
		return nil, &ImportError{message: "user_id is required for import"}
	}

	// Check if the user has a risk level set before importing
	resp, _, err := getOktaV6ClientFromMetadata(meta).UserRiskAPI.GetUserRisk(ctx, userId).Execute()
	if err != nil {
		return nil, &ImportError{message: "failed to get user risk: " + err.Error()}
	}

	// If the user has no risk level set (NONE), we cannot import
	if resp.UserRiskLevelNone != nil {
		return nil, &ImportError{message: "user has no risk level set (NONE). Set a risk level (HIGH or LOW) before importing, or create a new okta_user_risk resource instead."}
	}

	d.SetId(userId)
	_ = d.Set("user_id", userId)

	return []*schema.ResourceData{d}, nil
}

type ImportError struct {
	message string
}

func (e *ImportError) Error() string {
	return e.message
}
