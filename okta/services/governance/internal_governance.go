package governance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
		newRequestTypeResource,
		newRequestConditionResource,
		newRequestSequenceResource,
		newRequestSettingOrganizationResource,
		newRequestSettingResourceResource,
		newRequestV2Resource,
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
		newRequestConditionDataSource,
		newRequestSequencesDataSource,
		newRequestSettingOrganizationDataSource,
		newRequestSettingResourceDataSource,
		newRequestV2DataSource,
	}
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
