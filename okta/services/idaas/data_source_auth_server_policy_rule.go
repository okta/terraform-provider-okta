package idaas

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &authorizationServerPolicyRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &authorizationServerPolicyRulesDataSource{}
)

// AuthorizationServerPolicyRulesDataSource defines the data source implementation.
type authorizationServerPolicyRulesDataSource struct {
	Config *config.Config
}

// AuthorizationServerPolicyRulesDataSourceModel describes the data source data model.
type authorizationServerPolicyRulesDataSourceModel struct {
	ID           types.String                                                  `tfsdk:"id"`
	AuthServerId types.String                                                  `tfsdk:"auth_server_id"`
	PolicyId     types.String                                                  `tfsdk:"policy_id"`
	Actions      *AuthorizationServerPolicyRulesDataSourceModelActionsModel    `tfsdk:"actions"`
	Conditions   *AuthorizationServerPolicyRulesDataSourceModelConditionsModel `tfsdk:"conditions"`
	Created      types.String                                                  `tfsdk:"created"`
	LastUpdated  types.String                                                  `tfsdk:"last_updated"`
	Name         types.String                                                  `tfsdk:"name"`
	Priority     types.Int64                                                   `tfsdk:"priority"`
	Status       types.String                                                  `tfsdk:"status"`
	System       types.Bool                                                    `tfsdk:"system"`
	Type         types.String                                                  `tfsdk:"type"`
}

// AuthorizationServerPolicyRulesDataSourceModelActionsModel is the nested model for actions.
type AuthorizationServerPolicyRulesDataSourceModelActionsModel struct {
	Token *AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModel `tfsdk:"token"`
}

// AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModel is the nested model for token.
type AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModel struct {
	AccessTokenLifetimeMinutes  types.Int64                                                                         `tfsdk:"access_token_lifetime_minutes"`
	InlineHook                  *AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModelInlineHookModel `tfsdk:"inline_hook"`
	RefreshTokenLifetimeMinutes types.Int64                                                                         `tfsdk:"refresh_token_lifetime_minutes"`
	RefreshTokenWindowMinutes   types.Int64                                                                         `tfsdk:"refresh_token_window_minutes"`
}

// AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModelInlineHookModel is the nested model for inline_hook.
type AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModelInlineHookModel struct {
	Id types.String `tfsdk:"id"`
}

// AuthorizationServerPolicyRulesDataSourceModelConditionsModel is the nested model for conditions.
type AuthorizationServerPolicyRulesDataSourceModelConditionsModel struct {
	GrantTypes *AuthorizationServerPolicyRulesDataSourceModelConditionsModelGrantTypesModel `tfsdk:"grant_types"`
	People     *AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModel     `tfsdk:"people"`
	Scopes     *AuthorizationServerPolicyRulesDataSourceModelConditionsModelScopesModel     `tfsdk:"scopes"`
}

// AuthorizationServerPolicyRulesDataSourceModelConditionsModelGrantTypesModel is the nested model for grant_types.
type AuthorizationServerPolicyRulesDataSourceModelConditionsModelGrantTypesModel struct {
	Include types.List `tfsdk:"include"`
}

// AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModel is the nested model for people.
type AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModel struct {
	Groups *AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelGroupsModel `tfsdk:"groups"`
	Users  *AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelUsersModel  `tfsdk:"users"`
}

// AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelGroupsModel is the nested model for groups.
type AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelGroupsModel struct {
	Include types.List `tfsdk:"include"`
}

// AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelUsersModel is the nested model for users.
type AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelUsersModel struct {
	Include types.List `tfsdk:"include"`
}

// AuthorizationServerPolicyRulesDataSourceModelConditionsModelScopesModel is the nested model for scopes.
type AuthorizationServerPolicyRulesDataSourceModelConditionsModelScopesModel struct {
	Include types.List `tfsdk:"include"`
}

func newAuthorizationServerPolicyRulesDataSource() datasource.DataSource {
	return &authorizationServerPolicyRulesDataSource{}
}

func (d *authorizationServerPolicyRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorization_server_policy_rule"
}

func (d *authorizationServerPolicyRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *authorizationServerPolicyRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a policy rule by `ruleId`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the authorization_server_policy_rules.",
				Optional:            true,
				Computed:            true,
			},
			"auth_server_id": schema.StringAttribute{
				MarkdownDescription: "ID of the authorization server",
				Required:            true,
			},
			"policy_id": schema.StringAttribute{
				MarkdownDescription: "ID of the authorization server policy",
				Required:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the rule was created",
				Computed:            true,
			},
			"last_updated": schema.StringAttribute{
				MarkdownDescription: "Timestamp when the rule was last modified",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the rule",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority of the rule",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the rule",
				Computed:            true,
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "Set to `true` for system rules.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Rule type",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"actions": schema.SingleNestedBlock{
				Description: "Actions",
				Blocks: map[string]schema.Block{
					"token": schema.SingleNestedBlock{
						Description: "Token",
						Attributes: map[string]schema.Attribute{
							"access_token_lifetime_minutes": schema.Int64Attribute{
								Description: "Lifetime of the access token in minutes.",
								Computed:    true,
							},
							"refresh_token_lifetime_minutes": schema.Int64Attribute{
								Description: "Lifetime of the refresh token is the minimum access token lifetime.",
								Computed:    true,
							},
							"refresh_token_window_minutes": schema.Int64Attribute{
								Description: "Timeframe when the refresh token is valid.",
								Computed:    true,
							},
						},
						Blocks: map[string]schema.Block{
							"inline_hook": schema.SingleNestedBlock{
								Description: "InlineHook",
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Id",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"conditions": schema.SingleNestedBlock{
				Description: "Conditions",
				Blocks: map[string]schema.Block{
					"grant_types": schema.SingleNestedBlock{
						Description: "Array of grant types that this condition includes.",
						Attributes: map[string]schema.Attribute{
							"include": schema.ListAttribute{
								Description: "Array of grant types that this condition includes.",
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
					"people": schema.SingleNestedBlock{
						Description: "Identifies Users and Groups that are used together",
						Blocks: map[string]schema.Block{
							"groups": schema.SingleNestedBlock{
								Description: "Specifies a set of Groups whose Users are to be included",
								Attributes: map[string]schema.Attribute{
									"include": schema.ListAttribute{
										Description: "Groups to be included",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
							"users": schema.SingleNestedBlock{
								Description: "Specifies a set of Users to be included",
								Attributes: map[string]schema.Attribute{
									"include": schema.ListAttribute{
										Description: "Users to be included",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
					},
					"scopes": schema.SingleNestedBlock{
						Description: "Array of scopes that the condition includes",
						Attributes: map[string]schema.Attribute{
							"include": schema.ListAttribute{
								Description: "Include",
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *authorizationServerPolicyRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authorizationServerPolicyRulesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Config.OktaIDaaSClient.OktaSDKClientV6()
	if !state.ID.IsNull() && !state.ID.IsUnknown() && state.ID.ValueString() != "" {
		id := state.ID.ValueString()
		authServerId := state.AuthServerId.ValueString()
		policyId := state.PolicyId.ValueString()
		result, httpResp, err := client.AuthorizationServerRulesAPI.GetAuthorizationServerPolicyRule(ctx, authServerId, policyId, id).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddError("Not Found", "authorization_server_policy_rules with the given ID was not found.")
				return
			}
			resp.Diagnostics.AddError("Error reading authorization_server_policy_rules", err.Error())
			return
		}
		state.ID = types.StringValue(result.GetId())
		if t := result.GetCreated(); !t.IsZero() {
			state.Created = types.StringValue(t.Format(time.RFC3339))
		}
		if t := result.GetLastUpdated(); !t.IsZero() {
			state.LastUpdated = types.StringValue(t.Format(time.RFC3339))
		}
		state.Name = types.StringValue(result.GetName())
		state.Priority = types.Int64Value(int64(result.GetPriority()))
		state.Status = types.StringValue(result.GetStatus())
		state.System = types.BoolValue(result.GetSystem())
		state.Type = types.StringValue(result.GetType())
		if actionsRaw, ok := result.GetActionsOk(); ok && actionsRaw != nil {
			actionsModel := &AuthorizationServerPolicyRulesDataSourceModelActionsModel{}
			if tokenRaw, ok := actionsRaw.GetTokenOk(); ok && tokenRaw != nil {
				tokenModel := &AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModel{}
				tokenModel.AccessTokenLifetimeMinutes = types.Int64Value(int64(tokenRaw.GetAccessTokenLifetimeMinutes()))
				if inlineHookRaw, ok := tokenRaw.GetInlineHookOk(); ok && inlineHookRaw != nil {
					inlineHookModel := &AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModelInlineHookModel{}
					inlineHookModel.Id = types.StringValue(inlineHookRaw.GetId())
					tokenModel.InlineHook = inlineHookModel
				}
				tokenModel.RefreshTokenLifetimeMinutes = types.Int64Value(int64(tokenRaw.GetRefreshTokenLifetimeMinutes()))
				tokenModel.RefreshTokenWindowMinutes = types.Int64Value(int64(tokenRaw.GetRefreshTokenWindowMinutes()))
				actionsModel.Token = tokenModel
			}
			state.Actions = actionsModel
		}
		if conditionsRaw, ok := result.GetConditionsOk(); ok && conditionsRaw != nil {
			conditionsModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModel{}
			if grantTypesRaw, ok := conditionsRaw.GetGrantTypesOk(); ok && grantTypesRaw != nil {
				grantTypesModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelGrantTypesModel{}
				{
					listVal, listDiags := types.ListValueFrom(ctx, types.StringType, grantTypesRaw.GetInclude())
					resp.Diagnostics.Append(listDiags...)
					grantTypesModel.Include = listVal
				}
				conditionsModel.GrantTypes = grantTypesModel
			}
			if peopleRaw, ok := conditionsRaw.GetPeopleOk(); ok && peopleRaw != nil {
				peopleModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModel{}
				if groupsRaw, ok := peopleRaw.GetGroupsOk(); ok && groupsRaw != nil {
					groupsModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelGroupsModel{}
					{
						listVal, listDiags := types.ListValueFrom(ctx, types.StringType, groupsRaw.GetInclude())
						resp.Diagnostics.Append(listDiags...)
						groupsModel.Include = listVal
					}
					peopleModel.Groups = groupsModel
				}
				if usersRaw, ok := peopleRaw.GetUsersOk(); ok && usersRaw != nil {
					usersModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelUsersModel{}
					{
						listVal, listDiags := types.ListValueFrom(ctx, types.StringType, usersRaw.GetInclude())
						resp.Diagnostics.Append(listDiags...)
						usersModel.Include = listVal
					}
					peopleModel.Users = usersModel
				}
				conditionsModel.People = peopleModel
			}
			if scopesRaw, ok := conditionsRaw.GetScopesOk(); ok && scopesRaw != nil {
				scopesModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelScopesModel{}
				{
					listVal, listDiags := types.ListValueFrom(ctx, types.StringType, scopesRaw.GetInclude())
					resp.Diagnostics.Append(listDiags...)
					scopesModel.Include = listVal
				}
				conditionsModel.Scopes = scopesModel
			}
			state.Conditions = conditionsModel
		}

	} else {
		authServerId := state.AuthServerId.ValueString()
		policyId := state.PolicyId.ValueString()
		results, httpResp, err := client.AuthorizationServerRulesAPI.ListAuthorizationServerPolicyRules(ctx, authServerId, policyId).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddError("Not Found", "No authorization_server_policy_rules resources were found.")
				return
			}
			resp.Diagnostics.AddError("Error listing authorization_server_policy_rules", err.Error())
			return
		}
		if len(results) == 0 {
			resp.Diagnostics.AddError("Not Found", "No authorization_server_policy_rules resources were found.")
			return
		}
		// TODO: Add filtering logic to select the correct result from the list.
		result := results[0]
		state.ID = types.StringValue(result.GetId())
		if t := result.GetCreated(); !t.IsZero() {
			state.Created = types.StringValue(t.Format(time.RFC3339))
		}
		if t := result.GetLastUpdated(); !t.IsZero() {
			state.LastUpdated = types.StringValue(t.Format(time.RFC3339))
		}
		state.Name = types.StringValue(result.GetName())
		state.Priority = types.Int64Value(int64(result.GetPriority()))
		state.Status = types.StringValue(result.GetStatus())
		state.System = types.BoolValue(result.GetSystem())
		state.Type = types.StringValue(result.GetType())
		if actionsRaw, ok := result.GetActionsOk(); ok && actionsRaw != nil {
			actionsModel := &AuthorizationServerPolicyRulesDataSourceModelActionsModel{}
			if tokenRaw, ok := actionsRaw.GetTokenOk(); ok && tokenRaw != nil {
				tokenModel := &AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModel{}
				tokenModel.AccessTokenLifetimeMinutes = types.Int64Value(int64(tokenRaw.GetAccessTokenLifetimeMinutes()))
				if inlineHookRaw, ok := tokenRaw.GetInlineHookOk(); ok && inlineHookRaw != nil {
					inlineHookModel := &AuthorizationServerPolicyRulesDataSourceModelActionsModelTokenModelInlineHookModel{}
					inlineHookModel.Id = types.StringValue(inlineHookRaw.GetId())
					tokenModel.InlineHook = inlineHookModel
				}
				tokenModel.RefreshTokenLifetimeMinutes = types.Int64Value(int64(tokenRaw.GetRefreshTokenLifetimeMinutes()))
				tokenModel.RefreshTokenWindowMinutes = types.Int64Value(int64(tokenRaw.GetRefreshTokenWindowMinutes()))
				actionsModel.Token = tokenModel
			}
			state.Actions = actionsModel
		}
		if conditionsRaw, ok := result.GetConditionsOk(); ok && conditionsRaw != nil {
			conditionsModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModel{}
			if grantTypesRaw, ok := conditionsRaw.GetGrantTypesOk(); ok && grantTypesRaw != nil {
				grantTypesModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelGrantTypesModel{}
				{
					listVal, listDiags := types.ListValueFrom(ctx, types.StringType, grantTypesRaw.GetInclude())
					resp.Diagnostics.Append(listDiags...)
					grantTypesModel.Include = listVal
				}
				conditionsModel.GrantTypes = grantTypesModel
			}
			if peopleRaw, ok := conditionsRaw.GetPeopleOk(); ok && peopleRaw != nil {
				peopleModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModel{}
				if groupsRaw, ok := peopleRaw.GetGroupsOk(); ok && groupsRaw != nil {
					groupsModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelGroupsModel{}
					{
						listVal, listDiags := types.ListValueFrom(ctx, types.StringType, groupsRaw.GetInclude())
						resp.Diagnostics.Append(listDiags...)
						groupsModel.Include = listVal
					}
					peopleModel.Groups = groupsModel
				}
				if usersRaw, ok := peopleRaw.GetUsersOk(); ok && usersRaw != nil {
					usersModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelPeopleModelUsersModel{}
					{
						listVal, listDiags := types.ListValueFrom(ctx, types.StringType, usersRaw.GetInclude())
						resp.Diagnostics.Append(listDiags...)
						usersModel.Include = listVal
					}
					peopleModel.Users = usersModel
				}
				conditionsModel.People = peopleModel
			}
			if scopesRaw, ok := conditionsRaw.GetScopesOk(); ok && scopesRaw != nil {
				scopesModel := &AuthorizationServerPolicyRulesDataSourceModelConditionsModelScopesModel{}
				{
					listVal, listDiags := types.ListValueFrom(ctx, types.StringType, scopesRaw.GetInclude())
					resp.Diagnostics.Append(listDiags...)
					scopesModel.Include = listVal
				}
				conditionsModel.Scopes = scopesModel
			}
			state.Conditions = conditionsModel
		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
