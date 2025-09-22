package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &EndUserMyCatalogsEntryDataSource{}

func newEndUserMyCatalogsEntryDataSource() datasource.DataSource {
	return &EndUserMyCatalogsEntryDataSource{}
}

type EndUserMyCatalogsEntryDataSource struct {
	*config.Config
}

type EndUserMyCatalogsDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	EntryId     types.String `tfsdk:"entry_id"`
	Name        types.String `tfsdk:"name"`
	Requestable types.Bool   `tfsdk:"requestable"`
	Description types.String `tfsdk:"description"`
	Label       types.String `tfsdk:"label"`
	Parent      types.String `tfsdk:"parent"`
	Counts      struct {
		ResourceCounts struct {
			Applications types.Int32 `tfsdk:"applications"`
		} `tfsdk:"resource_counts"`
	} `tfsdk:"counts"`
}

func (r *EndUserMyCatalogsEntryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_end_user_my_requests_entry"
}

func (d *EndUserMyCatalogsEntryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (r *EndUserMyCatalogsEntryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"entry_id": schema.StringAttribute{
				Description: "The ID of the catalog entry",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "Unique identifier for the entry",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the entry",
				Computed:    true,
			},
			"requestable": schema.BoolAttribute{
				Description: "Is the resource requestable",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the entry",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "Label of the entry. Currently either Application or Resource Collection",
				Computed:    true,
			},
			"parent": schema.StringAttribute{
				Description: "Parent of the entry",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"counts": schema.SingleNestedBlock{
				Description: "Entry count metadata",
				Blocks: map[string]schema.Block{
					"resource_counts": schema.SingleNestedBlock{
						Description: "Collection resource counts",
						Attributes: map[string]schema.Attribute{
							"applications": schema.Int32Attribute{
								Description: "Number of app resources in a collection",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}
func (r *EndUserMyCatalogsEntryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateData EndUserMyCatalogsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// step 1 : Retrieve An Entry From My Catalog
	getMyCatalogEntryV2Request := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyCatalogsAPI.GetMyEntryV2(ctx, stateData.EntryId.ValueString())
	endUserMyCatalogEntry, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().MyCatalogsAPI.GetMyEntryV2Execute(getMyCatalogEntryV2Request)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving entry from my catalog", "Could not entry from my catalog, unexpected error: "+err.Error())
		return
	}

	// step 2 : "Convert" API Compatible Type Back To Terraform Type.
	stateData.Id = types.StringValue(endUserMyCatalogEntry.GetId())
	stateData.Name = types.StringValue(endUserMyCatalogEntry.GetName())
	stateData.Requestable = types.BoolValue(endUserMyCatalogEntry.GetRequestable())
	stateData.Description = types.StringValue(endUserMyCatalogEntry.GetDescription())
	stateData.Label = types.StringValue(endUserMyCatalogEntry.GetLabel())
	stateData.Parent = types.StringValue(endUserMyCatalogEntry.GetParent())
	stateData.Counts.ResourceCounts.Applications = types.Int32Value(endUserMyCatalogEntry.GetCounts().ResourceCounts.GetApplications())
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}
