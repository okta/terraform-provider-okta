package governance

import (
	"fmt"

	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/okta/terraform-provider-okta/okta/config"
)

func FWProviderResources() []func() resource.Resource {
	return []func() resource.Resource{
		newCampaignResource,
		newReviewResource,
		newEntitlementResource,
		newEntitlementBundleResource,
		newGrantResource,
		newRiskRuleResource,
		newCollectionResource,
		newMyRequestsResource,
	}
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newCampaignDataSource,
		newReviewDataSource,
		newEntitlementDataSource,
		newEntitlementBundlesDataSource,
		newPrincipalEntitlementsDataSource,
		newPrincipalAccessDataSource,
		newGrantDataSource,
		newRiskRulesDataSource,
		newCollectionDataSource,
	}
}

func ProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func ProviderDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{

		// resources.OktaInternalGovernanceCampaigns:    dataSourceCampaigns(),
	}
}

func stringIsJSON(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %q JSON to not be empty, got %v", k, i)
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		return diag.Errorf("%q contains an invalid JSON: %s", k, err)
	}
	return nil
}

// doNotRetry helper function to flag if provider should be using backoff.Retry
func doNotRetry(m interface{}, err error) bool {
	return m.(*config.Config).TimeOperations.DoNotRetry(err)
}

func datasourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("data-sources", name)
}

func oieOnlyFeatureError(kind, name string) diag.Diagnostics {
	url := fmt.Sprintf("https://registry.terraform.io/providers/okta/okta/latest/docs/%s/%s", kind, string(name[5:]))
	if kind == "resources" {
		kind = "resource"
	}
	if kind == "data-sources" {
		kind = "datasource"
	}
	return diag.Errorf("%q is a %s for OIE Orgs only, see %s", name, kind, url)
}

func resourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("resources", name)
}

func dataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

func resourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	p, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return p
}
