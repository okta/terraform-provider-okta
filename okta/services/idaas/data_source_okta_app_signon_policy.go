package idaas

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
		Description: "Get a sign-on policy for the application.",
	}
}

func dataSourceAppSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return datasourceOIEOnlyFeatureError(resources.OktaIDaaSAppSignOnPolicy)
	}

	app := sdk.NewApplication()
	_, _, err := getOktaClientFromMetadata(meta).Application.GetApplication(ctx, d.Get("app_id").(string), app, nil)
	if err != nil {
		return diag.Errorf("failed get app by ID: %v", err)
	}
	accessPolicy := utils.LinksValue(app.Links, "accessPolicy", "href")
	if accessPolicy == "" {
		return diag.Errorf("app does not support sign-on policy or this feature is not available")
	}
	policy := &sdk.Policy{}
	_policy, _, err := getOktaClientFromMetadata(meta).Policy.GetPolicy(ctx, path.Base(accessPolicy), policy, nil)
	if err != nil {
		return diag.Errorf("failed get policy by ID: %v", err)
	}
	policy = _policy.(*sdk.Policy)
	d.SetId(policy.Id)
	_ = d.Set("name", policy.Name)
	return nil
}
