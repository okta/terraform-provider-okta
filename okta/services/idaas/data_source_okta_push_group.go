package idaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

type pushGroupDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	AppId         types.String `tfsdk:"app_id"`
	SourceGroupId types.String `tfsdk:"source_group_id"`
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Status        types.String `tfsdk:"status"`
}

type pushGroupDataSource struct {
	config *config.Config
}

func newPushGroupDataSource() datasource.DataSource {
	return &pushGroupDataSource{}
}

func (r *pushGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_group"
}

func (r *pushGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.config = dataSourceConfiguration(req, resp)
}

func (r *pushGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Okta Application.",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Push Group Mapping ID",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("id"),
						path.MatchRoot("source_group_id"),
					),
				},
			},
			"source_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the source group in Okta.",
			},
			"target_group_id": schema.StringAttribute{
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "The ID of the existing target group for the push group mapping. This is used to link to an existing group",
			},
			"status": schema.StringAttribute{
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: "The status of the push group mapping. Valid values: `ACTIVE` and `INACTIVE`",
			},
		},
		Description: "Get a Push Group assignment for an Application in Okta.",
	}
}

func (r *pushGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pushGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupPushMapping *v6okta.GroupPushMapping
	var err error
	if state.ID.ValueString() != "" {
		groupPushMapping, _, err = r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.GetGroupPushMapping(ctx, state.AppId.ValueString(), state.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError("failed to read push group mapping: ", err.Error())
			return
		}
	} else {
		err := retry.RetryContext(ctx, 3*time.Second, func() *retry.RetryError {
			groupPushMappings, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.ListGroupPushMappings(ctx, state.AppId.ValueString()).SourceGroupId(state.SourceGroupId.ValueString()).Execute()
			if err != nil {
				resp.Diagnostics.AddError("failed to list push group mappings: ", err.Error())
				return retry.NonRetryableError(err)
			}
			if len(groupPushMappings) == 0 {
				return retry.RetryableError(fmt.Errorf("no push group mapping found with the specified source group id: %s", state.SourceGroupId.ValueString()))
			}
			if len(groupPushMappings) > 1 {
				return retry.NonRetryableError(fmt.Errorf("multiple push group mappings found with the specified source group id: %s", state.SourceGroupId.ValueString()))
			}
			groupPushMapping = &groupPushMappings[0]
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError("failed to read push group mapping: ", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mapPushGroupDataSourceToState(groupPushMapping, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func mapPushGroupDataSourceToState(groupPushMapping *v6okta.GroupPushMapping, state *pushGroupDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringPointerValue(groupPushMapping.Id)
	state.SourceGroupId = types.StringPointerValue(groupPushMapping.SourceGroupId)
	state.TargetGroupId = types.StringPointerValue(groupPushMapping.TargetGroupId)
	state.Status = types.StringPointerValue(groupPushMapping.Status)
	return diags
}
