package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var pushGroupsDataSourceObjectAttrs = map[string]attr.Type{
	"app_id":          types.StringType,
	"id":              types.StringType,
	"source_group_id": types.StringType,
	"status":          types.StringType,
}

type pushGroupsDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	AppId    types.String `tfsdk:"app_id"`
	Mappings types.List   `tfsdk:"mappings"`
}

type pushGroupsDataSource struct {
	config *config.Config
}

func newPushGroupsDataSource() datasource.DataSource {
	return &pushGroupsDataSource{}
}

func (r *pushGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_groups"
}

func (r *pushGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.config = dataSourceConfiguration(req, resp)
}

func (r *pushGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    false,
				Computed:    true,
				Description: "The ID of the Okta Application.",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Computed:    false,
				Description: "The ID of the Okta Application.",
			},
			"mappings": schema.ListAttribute{
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "List of Push Group mappings for the Application.",
				ElementType: types.ObjectType{
					AttrTypes: pushGroupsDataSourceObjectAttrs,
				},
			},
		},
		Description: "Gets a list of Push Group mappings for an Application in Okta.",
	}
}

func (r *pushGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pushGroupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupPushMappings []v6okta.GroupPushMapping
	var cursor string
	for {
		mappings, response, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.ListGroupPushMappings(ctx, state.AppId.ValueString()).After(cursor).Execute()
		if err != nil {
			resp.Diagnostics.AddError("failed to read push group mapping: ", err.Error())
			return
		}
		groupPushMappings = append(groupPushMappings, mappings...)

		if response.HasNextPage() {
			cursor = *mappings[len(mappings)-1].Id
		} else {
			break
		}
	}

	resp.Diagnostics.Append(mapPushGroupsDataSourceToState(groupPushMappings, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = state.AppId

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func mapPushGroupsDataSourceToState(groupPushMapping []v6okta.GroupPushMapping, state *pushGroupsDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(groupPushMapping) == 0 {
		return diags
	}

	state.Mappings = types.ListValueMust(
		types.ObjectType{
			AttrTypes: pushGroupsDataSourceObjectAttrs,
		},
		func() []attr.Value {
			var values []attr.Value
			for _, mapping := range groupPushMapping {
				objAttrs := map[string]attr.Value{
					"id":              types.StringPointerValue(mapping.Id),
					"app_id":          types.StringPointerValue(state.AppId.ValueStringPointer()),
					"source_group_id": types.StringPointerValue(mapping.SourceGroupId),
					"status":          types.StringPointerValue(mapping.Status),
				}
				values = append(values, types.ObjectValueMust(
					pushGroupsDataSourceObjectAttrs,
					objAttrs,
				))
			}
			return values
		}(),
	)
	return diags
}
