package governance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *EndUserMyCatalogsEntryUserRequestFieldsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (r *EndUserMyCatalogsEntryUserRequestFieldsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"entry_id": schema.StringAttribute{
				Description: "The ID of the catalog entry",
				Required:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed: true, // the response doesn't return any ID that we can rely upon.
			},
		},
		Blocks: map[string]schema.Block{
			"data": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Useful for specifying requesterFieldValues when adding a request.",
							Computed:    true,
						},
						"required": schema.BoolAttribute{
							Description: "Useful for specifying requesterFieldValues when adding a request.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of value for the requester field.",
							Computed:    true,
						},
						"choices": schema.ListAttribute{
							Description: "Valid choices when type is SELECT or MULTISELECT.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "label of the requester field",
							Computed:    true,
						},
						"maximum_value": schema.StringAttribute{
							Description: "The maximum value allowed for this field. Only applies to DURATION fields.",
							Computed:    true,
						},
						"read_only": schema.BoolAttribute{
							Description: "Indicates this field is immutable.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "An admin configured value for this field. Only applies to DURATION fields.",
							Computed:    true,
						},
					},
				},
			},
			"metadata": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"risk_assessment": schema.SingleNestedBlock{
						Description: "A risk assessment indicates whether request submission is allowed or restricted and contains the risk rules that lead to possible conflicts for the requested resource.",
						Attributes: map[string]schema.Attribute{
							"request_submission_type": schema.StringAttribute{
								Computed:    true,
								Description: "Whether request submission is allowed or restricted in the risk settings.",
							},
						},
						Blocks: map[string]schema.Block{
							"risk_rules": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "The name of a resource rule causing a conflict",
											Computed:    true,
										},
										"description": schema.StringAttribute{
											Description: "The human readable description",
											Computed:    true,
										},
										"resource_name": schema.StringAttribute{
											Description: "Human readable name of the resource",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
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
	if endUserMyCatalogEntryUserRequestFields.GetMetadata().RiskAssessment != nil {
		stateData.Metadata = &Metadata{}
		stateData.Metadata.RiskAssessment = &RiskAssessment{}
		stateData.Metadata.RiskAssessment.RequestSubmissionType = types.StringValue(string(endUserMyCatalogEntryUserRequestFields.GetMetadata().RiskAssessment.RequestSubmissionType))
		for _, riskRule := range endUserMyCatalogEntryUserRequestFields.GetMetadata().RiskAssessment.RiskRules {
			stateData.Metadata.RiskAssessment.RiskRules = append(stateData.Metadata.RiskAssessment.RiskRules, RiskRule{
				Name:         types.StringValue(riskRule.Name),
				Description:  types.StringValue(*riskRule.Description),
				ResourceName: types.StringValue(*riskRule.ResourceName),
			})
		}
	}
	stateData.Id = types.StringValue(fmt.Sprintf("%v-%v", stateData.EntryId.ValueString(), stateData.UserId.ValueString())) // the response doesn't return any ID that we can rely upon.
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}
