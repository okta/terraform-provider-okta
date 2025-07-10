package governance

import (
	"context"
	"fmt"
	"github.com/okta/terraform-provider-okta/okta/config"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &campaignDataSource{}
	//_ datasource.DataSourceWithConfigure = &campaignDataSource{}
)

func newCampaignDataSource() datasource.DataSource {
	return &campaignDataSource{}
}

type campaignDataSource struct {
	*config.Config
}

type campaignDataSourceModel struct {
	Id                   types.String                      `tfsdk:"id"`
	Created              types.String                      `tfsdk:"created"`
	CreatedBy            types.String                      `tfsdk:"created_by"`
	LastUpdated          types.String                      `tfsdk:"last_updated"`
	LastUpdatedBy        types.String                      `tfsdk:"last_updated_by"`
	Name                 types.String                      `tfsdk:"name"`
	Status               types.String                      `tfsdk:"status"`
	CampaignType         types.String                      `tfsdk:"campaign_type"`
	Description          types.String                      `tfsdk:"description"`
	RecurringCampaignId  types.String                      `tfsdk:"recurring_campaign_id"`
	RemediationSettings  *campaignRemediationSettingsModel `tfsdk:"remediation_settings"`
	AutoRemediation      *autoRemediationSettingsModel     `tfsdk:"auto_remediation_settings"`
	ResourceSettings     *resourceSettingsModel            `tfsdk:"resource_settings"`
	ReviewerSettings     *reviewerSettingsModel            `tfsdk:"reviewer_settings"`
	ScheduleSettings     *scheduleSettingsModel            `tfsdk:"schedule_settings"`
	NotificationSettings *notificationSettingsModel        `tfsdk:"notification_settings"`
	PrincipalScope       *principalScopeSettingsModel      `tfsdk:"principal_scope_settings"`
}

type autoRemediationSettingsModel struct {
	IncludeAllIndirectAssignments types.Bool                    `tfsdk:"include_all_indirect_assignments"`
	IncludeOnly                   []campaignTargetResourceModel `tfsdk:"include_only"`
}

type resourceSettingsModel struct {
	Type                               types.String                  `tfsdk:"type"`
	IncludeAdminRoles                  types.Bool                    `tfsdk:"includeAdminRoles"`
	IncludeEntitlements                types.Bool                    `tfsdk:"includeEntitlements"`
	IndividuallyAssignedAppsOnly       types.Bool                    `tfsdk:"individuallyAssignedAppsOnly"`
	IndividuallyAssignedGroupsOnly     types.Bool                    `tfsdk:"individuallyAssignedGroupsOnly"`
	OnlyIncludeOutOfPolicyEntitlements types.Bool                    `tfsdk:"onlyIncludeOutOfPolicyEntitlements"`
	ExcludedResources                  []campaignTargetResourceModel `tfsdk:"excluded_resources"`
	TargetResources                    []targetResourceModel         `tfsdk:"target_resources"`
}

type targetResourceModel struct {
	ResourceId                       types.String             `tfsdk:"resource_id"`
	IncludeAllEntitlementsAndBundles types.Bool               `tfsdk:"includeAllEntitlementsAndBundles"`
	ResourceType                     types.String             `tfsdk:"resource_type"`
	EntitlementBundles               []entitlementBundleModel `tfsdk:"entitlement_bundles"`
	Entitlements                     []entitlementModel       `tfsdk:"entitlements"`
}

type entitlementBundleModel struct {
	Id types.String `tfsdk:"id"`
}

type entitlementModel struct {
	Id               types.String            `tfsdk:"id"`
	IncludeAllValues types.Bool              `tfsdk:"include_all_values"`
	Values           []entitlementValueModel `tfsdk:"values"`
}

type entitlementValueModel struct {
	Id types.String `tfsdk:"id"`
}

type reviewerSettingsModel struct {
	Type                    types.String         `tfsdk:"type"`
	BulkDecisionDisabled    types.Bool           `tfsdk:"bulk_decision_disabled"`
	FallbackReviewerId      types.String         `tfsdk:"fallback_reviewer_id"`
	JustificationRequired   types.Bool           `tfsdk:"justification_required"`
	ReassignmentDisabled    types.Bool           `tfsdk:"reassignment_disabled"`
	ReviewerGroupId         types.String         `tfsdk:"reviewer_group_id"`
	ReviewerId              types.String         `tfsdk:"reviewer_id"`
	ReviewerScopeExpression types.String         `tfsdk:"reviewer_scope_expression"`
	SelfReviewDisabled      types.Bool           `tfsdk:"self_review_disabled"`
	ReviewerLevels          []reviewerLevelModel `tfsdk:"reviewer_levels"`
}

type reviewerLevelModel struct {
	Type                    types.String       `tfsdk:"type"`
	FallBackReviewerId      types.String       `tfsdk:"fallBackReviewerId"`
	ReviewerGroupId         types.String       `tfsdk:"reviewerGroupId"`
	ReviewerId              types.String       `tfsdk:"reviewerId"`
	ReviewerScopeExpression types.String       `tfsdk:"reviewerScopeExpression"`
	SelfReviewDisabled      types.Bool         `tfsdk:"selfReviewDisabled"`
	StartReview             []startReviewModel `tfsdk:"start_review"`
}

type startReviewModel struct {
	OnDay types.Int64  `tfsdk:"on_day"`
	When  types.String `tfsdk:"when"`
}

type scheduleSettingsModel struct {
	DurationInDays types.Int64       `tfsdk:"duration_in_days"`
	EndDate        types.String      `tfsdk:"end_date"`
	StartDate      types.String      `tfsdk:"start_date"`
	TimeZone       types.String      `tfsdk:"time_zone"`
	Type           types.String      `tfsdk:"type"`
	Recurrence     []recurrenceModel `tfsdk:"recurrence"`
}

type recurrenceModel struct {
	Interval     types.String `tfsdk:"interval"`
	Ends         types.String `tfsdk:"ends"`
	RepeatOnType types.String `tfsdk:"repeat_on_type"`
}

type notificationSettingsModel struct {
	NotifyReviewerAtCampaignEnd                types.Bool `tfsdk:"notifyReviewerAtCampaignEnd"`
	NotifyReviewerDuringMidpointOfReview       types.Bool `tfsdk:"notifyReviewerDuringMidpointOfReview"`
	NotifyReviewerWhenOverdue                  types.Bool `tfsdk:"notifyReviewerWhenOverdue"`
	NotifyReviewerWhenReviewAssigned           types.Bool `tfsdk:"notifyReviewerWhenReviewAssigned"`
	NotifyReviewPeriodEnd                      types.Bool `tfsdk:"notifyReviewPeriodEnd"`
	RemindersReviewerBeforeCampaignCloseInSecs types.List `tfsdk:"remindersReviewerBeforeCampaignCloseInSecs"`
}

type principalScopeSettingsModel struct {
	Type                             types.String              `tfsdk:"type"`
	ExcludedUserIds                  types.String              `tfsdk:"excluded_user_ids"`
	GroupIds                         types.List                `tfsdk:"group_ids"`
	IncludeOnlyActiveUsers           types.Bool                `tfsdk:"include_only_active_users"`
	OnlyIncludeUsersWithSODConflicts types.Bool                `tfsdk:"only_include_users_with_sod_conflicts"`
	UserIds                          types.List                `tfsdk:"user_ids"`
	UserScopeExpression              types.String              `tfsdk:"user_scope_expression"`
	PredefinedInactiveUsersScope     []inactiveUsersScopeModel `tfsdk:"predefined_inactive_users_scope"`
}

type inactiveUsersScopeModel struct {
	InactiveDays types.Int64 `tfsdk:"inactive_days"`
}

func (d *campaignDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaign"
}

func (d *campaignDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *campaignDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"campaign_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"recurring_campaign_id": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"remediation_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"access_approved": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the action by default if the reviewer approves access. NO_ACTION indicates there is no remediation action and the user retains access.",
					},
					"access_revoked": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the action if the reviewer revokes access. NO_ACTION indicates the user retains the same access. DENY indicates the user will have their access revoked as long as they are not assigned to a group through Group Rules.",
					},
					"no_response": schema.StringAttribute{
						Computed: true,
					},
				},
				Description: "Specify the action to be taken after a reviewer makes a decision to APPROVE or REVOKE the access, or if the campaign was CLOSED and there was no response from the reviewer.",
			},
			"auto_remediation_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"include_all_indirect_assignments": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"include_only": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Required: true,
								},
								"resource_type": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			"resource_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
					"includeAdminRoles": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"includeEntitlements": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"individuallyAssignedAppsOnly": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"individuallyAssignedGroupsOnly": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"onlyIncludeOutOfPolicyEntitlements": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"excluded_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Required: true,
									Computed: true,
								},
								"resource_type": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
					"target_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Required: true,
								},
								"includeAllEntitlementsAndBundles": schema.BoolAttribute{
									Computed: true,
								},
								"resource_type": schema.StringAttribute{
									Computed: true,
								},
							},
							Blocks: map[string]schema.Block{
								"entitlement_bundles": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
								"entitlements": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required: true,
											},
											"include_all_values": schema.BoolAttribute{
												Computed: true,
											},
										},
										Blocks: map[string]schema.Block{
											"values": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"reviewer_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
					"bulk_decision_disabled": schema.BoolAttribute{
						Computed: true,
					},
					"fallback_reviewer_id": schema.StringAttribute{
						Computed: true,
					},
					"justification_required": schema.BoolAttribute{
						Computed: true,
					},
					"reassignment_disabled": schema.BoolAttribute{
						Computed: true,
					},
					"reviewer_group_id": schema.StringAttribute{
						Computed: true,
					},
					"reviewer_id": schema.StringAttribute{
						Computed: true,
					},
					"reviewer_scope_expression": schema.StringAttribute{
						Computed: true,
					},
					"self_review_disabled": schema.BoolAttribute{
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"reviewer_levels": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
								},
								"fallBackReviewerId": schema.StringAttribute{
									Computed: true,
								},
								"reviewerGroupId": schema.StringAttribute{
									Computed: true,
								},
								"reviewerId": schema.StringAttribute{
									Computed: true,
								},
								"reviewerScopeExpression": schema.StringAttribute{
									Computed: true,
								},
								"selfReviewDisabled": schema.BoolAttribute{
									Computed: true,
								},
							},
							Blocks: map[string]schema.Block{
								"start_review": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"on_day": schema.Int64Attribute{
												Computed: true,
											},
											"when": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"schedule_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"duration_in_days": schema.Int64Attribute{
						Computed: true,
					},
					"end_date": schema.StringAttribute{
						Computed: true,
					},
					"start_date": schema.StringAttribute{
						Computed: true,
					},
					"time_zone": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"recurrence": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"interval": schema.StringAttribute{
									Computed: true,
								},
								"ends": schema.StringAttribute{
									Computed: true,
								},
								"repeat_on_type": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"notification_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"notifyReviewerAtCampaignEnd": schema.BoolAttribute{
						Computed: true,
					},
					"notifyReviewerDuringMidpointOfReview": schema.BoolAttribute{
						Computed: true,
					},
					"notifyReviewerWhenOverdue": schema.BoolAttribute{
						Computed: true,
					},
					"notifyReviewerWhenReviewAssigned": schema.BoolAttribute{
						Computed: true,
					},
					"notifyReviewPeriodEnd": schema.BoolAttribute{
						Computed: true,
					},
					"remindersReviewerBeforeCampaignCloseInSecs": schema.ListAttribute{
						Computed: true,
					},
				},
			},
			"principal_scope_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
					"excluded_user_ids": schema.StringAttribute{
						Computed: true,
					},
					"group_ids": schema.ListAttribute{
						Computed: true,
					},
					"include_only_active_users": schema.BoolAttribute{
						Computed: true,
					},
					"only_include_users_with_sod_conflicts": schema.BoolAttribute{
						Computed: true,
					},
					"user_ids": schema.ListAttribute{
						Computed: true,
					},
					"user_scope_expression": schema.StringAttribute{
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"predefined_inactive_users_scope": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"inactive_days": schema.Int64Attribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *campaignDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data campaignDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	campaign, _, err := d.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.GetCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading campaign",
			fmt.Sprintf("Could not retrieve campaign with ID %s: %s", data.Id.ValueString(), err.Error()),
		)
		return
	}

	data.Id = types.StringValue(campaign.Id)
	data.Name = types.StringValue(campaign.Name)
	data.Status = types.StringValue(string(campaign.Status))
	data.Description = types.StringValue(*campaign.Description)
	//data.CampaignType = campaign.CampaignType
	data.Created = types.StringValue(campaign.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(campaign.CreatedBy)
	data.LastUpdated = types.StringValue(campaign.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(campaign.LastUpdatedBy)
	data.RecurringCampaignId = types.StringValue(*campaign.RecurringCampaignId.Get())

	//if rs := campaign.RemediationSettings; rs != nil {
	data.RemediationSettings = &campaignRemediationSettingsModel{
		AccessApproved: types.StringValue(string(campaign.RemediationSettings.AccessApproved)),
		AccessRevoked:  types.StringValue(string(campaign.RemediationSettings.AccessRevoked)),
		NoResponse:     types.StringValue(string(campaign.RemediationSettings.NoResponse)),
	}
	//}

	//todo handle RemediationSettings
	//if ars := campaign.RemediationSettings; ars != nil {
	//	includeOnly := make([]campaignTargetResourceModel, len(ars.Inc))
	//	for i, r := range ars.IncludeOnly {
	//		includeOnly[i] = campaignTargetResourceModel{
	//			ResourceId:   types.StringValue(r.),
	//			ResourceType: types.StringValue(string(r.ResourceType)),
	//		}
	//	}
	//	data.AutoRemediation = &autoRemediationSettingsModel{
	//		IncludeAllIndirectAssignments: types.BoolValue(ars.IncludeAllIndirectAssignments),
	//		IncludeOnly:                   includeOnly,
	//	}
	//}

	//if rs := campaign.ResourceSettings; rs != nil {
	excluded := make([]campaignTargetResourceModel, len(campaign.ResourceSettings.ExcludedResources))
	for i, ex := range campaign.ResourceSettings.ExcludedResources {
		excluded[i] = campaignTargetResourceModel{
			ResourceId: types.StringValue(*ex.ResourceId),
			//todo handle ResourceType
			//ResourceType: ex.ResourceType,
		}
		//}

		targets := make([]targetResourceModel, len(campaign.ResourceSettings.TargetResources))
		for i, tr := range campaign.ResourceSettings.TargetResources {
			entitlements := make([]entitlementModel, len(tr.Entitlements))
			for j, ent := range tr.Entitlements {
				values := make([]entitlementValueModel, len(ent.Values))
				for k, v := range ent.Values {
					values[k] = entitlementValueModel{Id: types.StringValue(v.Id)}
				}
				entitlements[j] = entitlementModel{
					Id:               types.StringValue(ent.Id),
					IncludeAllValues: types.BoolValue(*ent.IncludeAllValues),
					Values:           values,
				}
			}
			bundles := make([]entitlementBundleModel, len(tr.EntitlementBundles))
			for j, b := range tr.EntitlementBundles {
				bundles[j] = entitlementBundleModel{Id: types.StringValue(b.Id)}
			}
			targets[i] = targetResourceModel{
				ResourceId:                       types.StringValue(tr.ResourceId),
				IncludeAllEntitlementsAndBundles: types.BoolValue(*tr.IncludeAllEntitlementsAndBundles),
				//todo handle ResourceType
				//ResourceType:                     types.StringValue(string(tr.ResourceType)),
				EntitlementBundles: bundles,
				Entitlements:       entitlements,
			}
		}
		data.ResourceSettings = &resourceSettingsModel{
			Type:                               types.StringValue(string(campaign.ResourceSettings.Type)),
			IncludeAdminRoles:                  types.BoolValue(*campaign.ResourceSettings.IncludeAdminRoles),
			IncludeEntitlements:                types.BoolValue(*campaign.ResourceSettings.IncludeEntitlements),
			IndividuallyAssignedAppsOnly:       types.BoolValue(*campaign.ResourceSettings.IndividuallyAssignedAppsOnly),
			IndividuallyAssignedGroupsOnly:     types.BoolValue(*campaign.ResourceSettings.IndividuallyAssignedGroupsOnly),
			OnlyIncludeOutOfPolicyEntitlements: types.BoolValue(*campaign.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements),
			ExcludedResources:                  excluded,
			TargetResources:                    targets,
		}
	}

	// Similarly map other nested blocks: ReviewerSettings, ScheduleSettings, NotificationSettings, PrincipalScope...
	// (To be added depending on your API structure)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
