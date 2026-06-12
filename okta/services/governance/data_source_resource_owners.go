package governance

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*resourceOwnersDataSource)(nil)

func newResourceOwnersDataSource() datasource.DataSource {
	return &resourceOwnersDataSource{}
}

type resourceOwnersDataSource struct {
	*config.Config
}

type resourceOwnersDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Filter         types.String `tfsdk:"filter"`
	ResourceOwners types.List   `tfsdk:"resource_owners"`
}

var principalObjectAttrTypes = map[string]attr.Type{
	"id":   types.StringType,
	"type": types.StringType,
	"orn":  types.StringType,
	"name": types.StringType,
}

var resourceOwnerDSObjectAttrTypes = map[string]attr.Type{
	"resource_id":        types.StringType,
	"resource_type":      types.StringType,
	"resource_orn":       types.StringType,
	"resource_name":      types.StringType,
	"parent_resource_orn": types.StringType,
	"principals": types.ListType{
		ElemType: types.ObjectType{AttrTypes: principalObjectAttrTypes},
	},
}

func (d *resourceOwnersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_owners"
}

func (d *resourceOwnersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *resourceOwnersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists governance resources and their assigned owners.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder ID.",
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "A filter expression for listing resource owners. " +
					"Supports parentResourceOrn (eq) and resource.orn (eq) filters.",
				Required: true,
			},
			"resource_owners": schema.ListAttribute{
				Description: "List of resources with their assigned owners.",
				Computed:    true,
				ElementType: types.ObjectType{AttrTypes: resourceOwnerDSObjectAttrTypes},
			},
		},
	}
}

func (d *resourceOwnersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data resourceOwnersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := data.Filter.ValueString()
	if filter == "" {
		resp.Diagnostics.AddError("Missing filter", "The 'filter' attribute must be set.")
		return
	}

	// Collect all pages. Initialize to empty (not nil) to avoid null vs empty list in state.
	allElements := make([]attr.Value, 0)
	var after string

	for {
		listReq := d.OktaGovernanceClient.OktaGovernanceSDKClient().ResourceOwnersAPI.
			ListResourceOwners(ctx).
			Filter(filter)
		if after != "" {
			listReq = listReq.After(after)
		}

		listResp, _, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing resource owners",
				fmt.Sprintf("Could not list resource owners: %s", err.Error()),
			)
			return
		}

		for _, ro := range listResp.GetData() {
			res := ro.GetResource()

			// Build principals list
			principalElements := make([]attr.Value, 0, len(ro.GetPrincipals()))
			for _, p := range ro.GetPrincipals() {
				profile := p.GetProfile()
				principalElements = append(principalElements, types.ObjectValueMust(
					principalObjectAttrTypes,
					map[string]attr.Value{
						"id":   types.StringValue(p.GetId()),
						"type": types.StringValue(p.GetType()),
						"orn":  types.StringValue(p.GetOrn()),
						"name": types.StringValue(profile.GetName()),
					},
				))
			}
			principalsList, diags := types.ListValue(
				types.ObjectType{AttrTypes: principalObjectAttrTypes},
				principalElements,
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			profile := res.GetProfile()
			allElements = append(allElements, types.ObjectValueMust(
				resourceOwnerDSObjectAttrTypes,
				map[string]attr.Value{
					"resource_id":        types.StringValue(res.GetId()),
					"resource_type":      types.StringValue(res.GetType()),
					"resource_orn":       types.StringValue(res.GetOrn()),
					"resource_name":      types.StringValue(profile.GetName()),
					"parent_resource_orn": types.StringValue(ro.GetParentResourceOrn()),
					"principals":         principalsList,
				},
			))
		}

		// Handle pagination
		links, hasLinks := listResp.GetLinksOk()
		if !hasLinks || links == nil {
			break
		}
		nextLink, hasNext := links.GetNextOk()
		if !hasNext || nextLink == nil {
			break
		}
		nextHref := nextLink.GetHref()
		if nextHref == "" {
			break
		}
		u, err := url.Parse(nextHref)
		if err != nil {
			break
		}
		afterParam := u.Query().Get("after")
		if afterParam == "" {
			break
		}
		after = afterParam
	}

	resourceOwnersList, diags := types.ListValue(
		types.ObjectType{AttrTypes: resourceOwnerDSObjectAttrTypes},
		allElements,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ResourceOwners = resourceOwnersList
	data.ID = types.StringValue("resource_owners")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
