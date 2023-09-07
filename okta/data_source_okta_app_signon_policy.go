package okta

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceAppSignOnPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppSignOnPolicyRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy name",
			},
		},
	}
}

func dataSourceAppSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return datasourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	app := sdk.NewApplication()
	_, _, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, d.Get("app_id").(string), app, nil)
	if err != nil {
		return diag.Errorf("failed get app by ID: %v", err)
	}
	accessPolicy := linksValue(app.Links, "accessPolicy", "href")
	if accessPolicy == "" {
		return diag.Errorf("app does not support sign-on policy or this feature is not available")
	}
	policy := &sdk.Policy{}
	_policy, _, err := getOktaClientFromMetadata(m).Policy.GetPolicy(ctx, path.Base(accessPolicy), policy, nil)
	if err != nil {
		return diag.Errorf("failed get policy by ID: %v", err)
	}
	policy = _policy.(*sdk.Policy)
	d.SetId(policy.Id)
	_ = d.Set("name", policy.Name)
	return nil
}
