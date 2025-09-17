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
		newEntitlementResource,
		newReviewResource,
		newRequestConditionResource,
		newEndUserMyRequestsResource,
		newRequestSequenceResource,
		newRequestSettingOrganizationResource,
		newRequestSettingResourceResource,
		newRequestV2Resource,
	}
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newCampaignDataSource,
		newEntitlementDataSource,
		newReviewDataSource,
		newPrincipalEntitlementsDataSource,
		newRequestSequencesDataSource,
		newRequestSettingOrganizationDataSource,
		newRequestSettingResourceDataSource,
		newRequestV2DataSource,
		newRequestConditionDataSource,
		newEndUserMyRequestsDataSource,
	}
}

func dataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	conf, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return conf
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
