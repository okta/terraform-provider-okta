package governance

import (
	"context"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*catalogEntryDefaultDataSource)(nil)

func newCatalogEntryDefaultDataSource() datasource.DataSource {
	return &catalogEntryDefaultDataSource{}
}

func (d *catalogEntryDefaultDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type catalogEntryDefaultDataSource struct {
	*config.Config
}

type resourceCounts struct {
	Applications types.Int32 `tfsdk:"applications"`
}

type counts struct {
	ResourceCounts resourceCounts `tfsdk:"resource_counts"`
}

type self struct {
	Href types.String `tfsdk:"href"`
}

type links struct {
	Self self `tfsdk:"self"`
}

type catalogEntryDefaultDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	EntryId     types.String `tfsdk:"entry_id"`
	Name        types.String `tfsdk:"name"`
	Requestable types.Bool   `tfsdk:"requestable"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Parent      types.String `tfsdk:"parent"`
	Counts      *counts      `tfsdk:"counts"`
	Links       *links       `tfsdk:"links"`
}

func (d *catalogEntryDefaultDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog_entry_default"
}

func (d *catalogEntryDefaultDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"entry_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"requestable": schema.BoolAttribute{
				Computed: true,
			},
			"label": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"parent": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"counts": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"resource_counts": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"applications": schema.Int32Attribute{
								Computed: true,
							},
						},
					},
				},
			},
			"links": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"self": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"href": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *catalogEntryDefaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogEntryDefaultDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getCatalogEntryResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().CatalogsAPI.GetCatalogEntryV2(ctx, data.EntryId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Catalog entry",
			"Could not read Catalog entry, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(getCatalogEntryResp.GetId())
	data.Name = types.StringValue(getCatalogEntryResp.GetName())
	data.Requestable = types.BoolValue(getCatalogEntryResp.GetRequestable())
	data.Description = types.StringValue(getCatalogEntryResp.GetDescription())
	data.Label = types.StringValue(getCatalogEntryResp.GetLabel())
	data.Parent = types.StringValue(getCatalogEntryResp.GetParent())
	respLinks := getCatalogEntryResp.GetLinks()
	respSelf := respLinks.GetSelf()
	linksSelfHref := respSelf.GetHref()
	data.Links = &links{Self: self{Href: types.StringValue(linksSelfHref)}}
	if getCatalogEntryResp.Counts != nil {
		if getCatalogEntryResp.Counts.ResourceCounts != nil {
			data.Counts.ResourceCounts.Applications = types.Int32Value(getCatalogEntryResp.Counts.ResourceCounts.GetApplications())
		}
	}

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
