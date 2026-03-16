package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserRisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRiskRead,
		Description: "Gets a user's risk level in Okta.",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user.",
			},
			"risk_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Risk level of the user. Possible values: `HIGH`, `LOW`.",
			},
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason for the risk level, if available.",
			},
		},
	}
}

func dataSourceUserRiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)

	logger(meta).Info("reading user risk data source", "user_id", userId)

	resp, _, err := getOktaV6ClientFromMetadata(meta).UserRiskAPI.GetUserRisk(ctx, userId).Execute()
	if err != nil {
		return diag.Errorf("failed to get user risk: %v", err)
	}

	// Extract risk level and reason from union type response
	var riskLevel, reason string
	if resp.UserRiskLevelExists != nil {
		riskLevel = resp.UserRiskLevelExists.GetRiskLevel()
		reason = resp.UserRiskLevelExists.GetReason()
	} else if resp.UserRiskLevelNone != nil {
		riskLevel = "NONE"
		reason = ""
	}

	d.SetId(userId)
	_ = d.Set("user_id", userId)
	_ = d.Set("risk_level", riskLevel)
	_ = d.Set("reason", reason)

	return nil
}
