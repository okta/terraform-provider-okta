package idaas

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

var _ datasource.DataSource = &resourceSetDataSource{}

func newResourceSetDataSource() datasource.DataSource {
	return &resourceSetDataSource{}
}

type resourceSetDataSource struct {
	config *config.Config
}

type resourceSetDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	Resources   types.List   `tfsdk:"resources"`
}

func (d *resourceSetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_set"
}

func (d *resourceSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *resourceSetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve a Resource Set from Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Resource Set.",
			},
			"label": schema.StringAttribute{
				Computed:    true,
				Description: "Unique name given to the Resource Set.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A description of the Resource Set.",
			},
			"resources": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The resources included in the Resource Set. Each entry is the self-link href of the resource.",
			},
		},
	}
}

func (d *resourceSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data resourceSetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.config.OktaIDaaSClient.OktaSDKSupplementClient()

	rs, _, err := client.GetResourceSet(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Resource Set",
			"Could not read Resource Set, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(rs.Id)
	data.Label = types.StringValue(rs.Label)
	data.Description = types.StringValue(rs.Description)

	// Fetch the resources in the set with pagination
	var resLinks []string
	resources, _, err := client.ListResourceSetResources(ctx, rs.Id, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Resource Set resources",
			"Could not read Resource Set resources, unexpected error: "+err.Error(),
		)
		return
	}
	for _, r := range resources.Resources {
		if r.Links != nil {
			resLinks = append(resLinks, utils.LinksValue(r.Links, "self", "href"))
		}
	}
	for {
		if nextURL := utils.LinksValue(resources.Links, "next", "href"); nextURL != "" {
			u, parseErr := url.Parse(nextURL)
			if parseErr != nil {
				break
			}
			after := u.Query().Get("after")
			resources, _, err = client.ListResourceSetResources(ctx, rs.Id, &query.Params{After: after, Limit: utils.DefaultPaginationLimit})
			if err != nil {
				resp.Diagnostics.AddError(
					"Error reading Resource Set resources",
					"Could not read Resource Set resources during pagination, unexpected error: "+err.Error(),
				)
				return
			}
			for _, r := range resources.Resources {
				if r.Links != nil {
					resLinks = append(resLinks, utils.LinksValue(r.Links, "self", "href"))
				}
			}
		} else {
			break
		}
	}

	resourcesList, diags := types.ListValueFrom(ctx, types.StringType, resLinks)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Resources = resourcesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
