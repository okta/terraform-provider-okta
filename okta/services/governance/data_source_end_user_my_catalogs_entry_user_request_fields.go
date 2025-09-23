package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &EndUserMyCatalogsEntryUserRequestFieldsDataSource{}

func newEndUserMyCatalogsEntryUserRequestFieldsDataSource() datasource.DataSource {
	return &EndUserMyCatalogsEntryUserRequestFieldsDataSource{}
}

type EndUserMyCatalogsEntryUserRequestFieldsDataSource struct {
	EndUserMyCatalogsEntryRequestFieldsDataSource
	*config.Config
}

type EndUserMyCatalogsEntryUserRequestFieldsDataSourceModel struct {
	UserId types.String `tfsdk:"user_id"`
	EndUserMyCatalogsEntryRequestFieldsDataSourceModel
}

func (r *EndUserMyCatalogsEntryUserRequestFieldsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_end_user_my_catalogs_entry_user_request_fields"
}

func (r *EndUserMyCatalogsEntryUserRequestFieldsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateData EndUserMyCatalogsEntryUserRequestFieldsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : Retrieve An Entry From My Catalog
	getMyCatalogEntryUserRequestFieldsV2Request := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyCatalogsAPI.GetMyCatalogEntryUserRequestFieldsV2(ctx, stateData.EntryId.ValueString(), stateData.UserId.ValueString())
	endUserMyCatalogEntryUserRequestFields, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyCatalogsAPI.GetMyCatalogEntryUserRequestFieldsV2Execute(getMyCatalogEntryUserRequestFieldsV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving catalog entry request fields", "Could not retrieve catalog entry request fields, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	dataItems := []Data{}
	for _, dataItem := range endUserMyCatalogEntryUserRequestFields.GetData() {

		choices := []types.String{}
		for _, choice := range dataItem.Choices {
			choices = append(choices, types.StringValue(choice))
		}

		dataItems = append(dataItems, Data{
			Id:           types.StringValue(dataItem.GetId()),
			Required:     types.BoolValue(dataItem.Required),
			Type:         types.StringValue(string(dataItem.Type)),
			Choices:      choices,
			Label:        types.StringValue(dataItem.GetLabel()),
			MaximumValue: types.StringValue(*dataItem.MaximumValue),
			ReadOnly:     types.BoolValue(dataItem.GetReadOnly()),
			Value:        types.StringValue(*dataItem.Value),
		})
	}

	stateData.Data = dataItems
	stateData.Metadata.RiskAssessment.RequestSubmissionType = types.StringValue(string(endUserMyCatalogEntryUserRequestFields.GetMetadata().RiskAssessment.RequestSubmissionType))
	for _, riskRule := range endUserMyCatalogEntryUserRequestFields.GetMetadata().RiskAssessment.RiskRules {
		stateData.Metadata.RiskAssessment.RiskRules = append(stateData.Metadata.RiskAssessment.RiskRules, RiskRule{
			Name:         types.StringValue(riskRule.Name),
			Description:  types.StringValue(*riskRule.Description),
			ResourceName: types.StringValue(*riskRule.ResourceName),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}
