package governance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &requestConditionsDataSource{}

func newRequestConditionsDataSource() datasource.DataSource {
	return &requestConditionsDataSource{}
}

func (d *requestConditionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_conditions"
}

func (d *requestConditionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *requestConditionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all request conditions for a resource. Request conditions define what resources and access levels requesters can request from their resource catalog.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this data source, required by Terraform to track state. This field does not exist in the Okta API response.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format or ORN format.",
			},
		},
		Blocks: map[string]schema.Block{
			"conditions": schema.ListNestedBlock{
				Description: "List of request conditions for the resource.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the request condition.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the request condition.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the request condition.",
						},
						"priority": schema.Int32Attribute{
							Computed:    true,
							Description: "The priority of the request condition. Lower numbers indicate higher priority.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status indicates if this condition is active or not. Possible values: ACTIVE, INACTIVE, DELETED, INVALID.",
						},
						"approval_sequence_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the approval sequence.",
						},
						"created": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 formatted date and time when the resource was created.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the Okta user who created the resource.",
						},
						"last_updated": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 formatted date and time when the object was last updated.",
						},
						"last_updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the Okta user who last updated the object.",
						},
					},
					Blocks: map[string]schema.Block{
						"access_scope_settings": schema.SingleNestedBlock{
							Description: "Settings for the access request scope (such as groups, entitlement bundles, or default resources).",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of access scope. Possible values: RESOURCE_DEFAULT, GROUPS, ENTITLEMENT_BUNDLES.",
								},
							},
							Blocks: map[string]schema.Block{
								"ids": schema.ListNestedBlock{
									Description: "Block list of groups/entitlement bundles ids.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed:    true,
												Description: "The group/entitlement bundle ID.",
											},
										},
									},
								},
							},
						},
						"requester_settings": schema.SingleNestedBlock{
							Description: "Requester settings define who may submit an access request for the related resource and access scopes.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of requester. Possible values: EVERYONE, TEAMS, GROUPS.",
								},
							},
							Blocks: map[string]schema.Block{
								"ids": schema.ListNestedBlock{
									Description: "Block list of teams/groups ids.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed:    true,
												Description: "The group/team ID.",
											},
										},
									},
								},
							},
						},
						"access_duration_settings": schema.SingleNestedBlock{
							Description: "Settings that control who may specify the access duration allowed by this request condition, as well as what duration may be requested.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of access duration setting. Possible values: ADMIN_FIXED_DURATION, REQUESTER_SPECIFIED_DURATION.",
								},
								"duration": schema.StringAttribute{
									Computed:    true,
									Description: "The duration of access in ISO 8601 duration format. Only applicable for ADMIN_FIXED_DURATION type.",
								},
								"maximum_duration": schema.StringAttribute{
									Computed:    true,
									Description: "The maximum duration of access in ISO 8601 duration format. Only applicable for REQUESTER_SPECIFIED_DURATION type.",
								},
							},
						},
					},
				},
			},
		},
	}
}

type requestConditionsDataSource struct {
	*config.Config
}

type requestConditionsListModel struct {
	Id         types.String                `tfsdk:"id"`
	ResourceId types.String                `tfsdk:"resource_id"`
	Conditions []requestConditionItemModel `tfsdk:"conditions"`
}

type requestConditionItemModel struct {
	Id                     types.String            `tfsdk:"id"`
	Name                   types.String            `tfsdk:"name"`
	Description            types.String            `tfsdk:"description"`
	Priority               types.Int32             `tfsdk:"priority"`
	Status                 types.String            `tfsdk:"status"`
	ApprovalSequenceId     types.String            `tfsdk:"approval_sequence_id"`
	Created                types.String            `tfsdk:"created"`
	CreatedBy              types.String            `tfsdk:"created_by"`
	LastUpdated            types.String            `tfsdk:"last_updated"`
	LastUpdatedBy          types.String            `tfsdk:"last_updated_by"`
	AccessScopeSettings    *Settings               `tfsdk:"access_scope_settings"`
	RequesterSettings      *Settings               `tfsdk:"requester_settings"`
	AccessDurationSettings *AccessDurationSettings `tfsdk:"access_duration_settings"`
}

func (d *requestConditionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestConditionsListModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to list request conditions for the resource
	requestConditionsResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().
		RequestConditionsAPI.
		ListResourceRequestConditionsV2(ctx, data.ResourceId.ValueString()).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Conditions",
			fmt.Sprintf("Could not list request conditions for resource %s, unexpected error: %s", data.ResourceId.ValueString(), err.Error()),
		)
		return
	}

	// Map the API response to the Terraform model
	var conditions []requestConditionItemModel
	for _, condition := range requestConditionsResp.Data {
		item := requestConditionItemModel{
			Id:            types.StringValue(condition.GetId()),
			Name:          types.StringValue(condition.GetName()),
			Priority:      types.Int32Value(condition.GetPriority()),
			Status:        types.StringValue(string(condition.GetStatus())),
			Created:       types.StringValue(condition.GetCreated().Format(time.RFC3339)),
			CreatedBy:     types.StringValue(condition.GetCreatedBy()),
			LastUpdated:   types.StringValue(condition.GetLastUpdated().Format(time.RFC3339)),
			LastUpdatedBy: types.StringValue(condition.GetLastUpdatedBy()),
		}

		// Set optional approval sequence ID if present
		if condition.ApprovalSequenceId != nil {
			item.ApprovalSequenceId = types.StringValue(condition.GetApprovalSequenceId())
		} else {
			item.ApprovalSequenceId = types.StringNull()
		}

		// Set optional description if present
		if condition.Description != nil {
			item.Description = types.StringValue(condition.GetDescription())
		} else {
			item.Description = types.StringNull()
		}

		// Set access scope settings
		item.AccessScopeSettings, _ = setAccessScopeSettingsSparse(condition.GetAccessScopeSettings())

		// Set requester settings
		item.RequesterSettings, _ = setRequesterSettingsSparse(condition.GetRequesterSettings())

		// Set access duration settings if present
		if condition.AccessDurationSettings != nil {
			item.AccessDurationSettings = setAccessDurationSettings(*condition.AccessDurationSettings)
		}

		conditions = append(conditions, item)
	}

	// Set the data in the model
	data.Conditions = conditions
	data.Id = types.StringValue(fmt.Sprintf("request-conditions-%s", data.ResourceId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// setAccessScopeSettingsSparse converts AccessScopeSettings (from sparse response) to Settings model
func setAccessScopeSettingsSparse(settings governance.AccessScopeSettings) (*Settings, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := &Settings{}

	// Get the type if present
	if settings.Type != nil {
		result.Type = types.StringValue(string(*settings.Type))
	} else {
		result.Type = types.StringNull()
	}

	// Handle groups or entitlement bundles IDs from AdditionalProperties if present
	if settings.AdditionalProperties != nil {
		if groups, ok := settings.AdditionalProperties["groups"].([]interface{}); ok {
			var ids []IdModel
			for _, g := range groups {
				if groupMap, ok := g.(map[string]interface{}); ok {
					if id, ok := groupMap["id"].(string); ok {
						ids = append(ids, IdModel{
							Id: types.StringValue(id),
						})
					}
				}
			}
			result.Ids = ids
		} else if bundles, ok := settings.AdditionalProperties["entitlementBundles"].([]interface{}); ok {
			var ids []IdModel
			for _, b := range bundles {
				if bundleMap, ok := b.(map[string]interface{}); ok {
					if id, ok := bundleMap["id"].(string); ok {
						ids = append(ids, IdModel{
							Id: types.StringValue(id),
						})
					}
				}
			}
			result.Ids = ids
		}
	}

	return result, diags
}

// setRequesterSettingsSparse converts RequesterSettings (from sparse response) to Settings model
func setRequesterSettingsSparse(settings governance.RequesterSettings) (*Settings, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := &Settings{}

	// Get the type if present
	if settings.Type != nil {
		result.Type = types.StringValue(string(*settings.Type))
	} else {
		result.Type = types.StringNull()
	}

	// Handle groups or teams IDs from AdditionalProperties if present
	if settings.AdditionalProperties != nil {
		if groups, ok := settings.AdditionalProperties["groups"].([]interface{}); ok {
			var ids []IdModel
			for _, g := range groups {
				if groupMap, ok := g.(map[string]interface{}); ok {
					if id, ok := groupMap["id"].(string); ok {
						ids = append(ids, IdModel{
							Id: types.StringValue(id),
						})
					}
				}
			}
			result.Ids = ids
		} else if teams, ok := settings.AdditionalProperties["teams"].([]interface{}); ok {
			var ids []IdModel
			for _, t := range teams {
				if teamMap, ok := t.(map[string]interface{}); ok {
					if id, ok := teamMap["id"].(string); ok {
						ids = append(ids, IdModel{
							Id: types.StringValue(id),
						})
					}
				}
			}
			result.Ids = ids
		}
	}

	return result, diags
}
