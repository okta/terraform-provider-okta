package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	ResourceSettings     *resourceSettingsModel            `tfsdk:"resource_settings"`
	ReviewerSettings     *reviewerSettingsModel            `tfsdk:"reviewer_settings"`
	ScheduleSettings     *scheduleSettingsModel            `tfsdk:"schedule_settings"`
	NotificationSettings *notificationSettingsModel        `tfsdk:"notification_settings"`
	PrincipalScope       *principalScopeSettingsModel      `tfsdk:"principal_scope_settings"`
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
				Optional: true,
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
				Description: "Specify the action to be taken after a reviewer makes a decision to APPROVE or REVOKE the access, or if the campaign was CLOSED and there was no response from the reviewer.",
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
				Blocks: map[string]schema.Block{
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
				},
			},
			"resource_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
					"include_admin_roles": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"include_entitlements": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"individually_assigned_apps_only": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"individually_assigned_groups_only": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
					"only_include_out_of_policy_entitlements": schema.BoolAttribute{
						Computed: true,
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"excluded_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
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
								"include_all_entitlements_and_bundles": schema.BoolAttribute{
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
								"fallback_reviewer_id": schema.StringAttribute{
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
					"notify_reviewer_at_campaign_end": schema.BoolAttribute{
						Computed: true,
					},
					"notify_reviewer_during_midpoint_of_review": schema.BoolAttribute{
						Computed: true,
					},
					"notify_reviewer_when_overdue": schema.BoolAttribute{
						Computed: true,
					},
					"notify_reviewer_when_review_assigned": schema.BoolAttribute{
						Computed: true,
					},
					"notify_review_period_end": schema.BoolAttribute{
						Computed: true,
					},
					"reminders_reviewer_before_campaign_close_in_secs": schema.ListAttribute{
						Computed:    true,
						ElementType: types.Int64Type,
					},
				},
			},
			"principal_scope_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
					"excluded_user_ids": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"group_ids": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"include_only_active_users": schema.BoolAttribute{
						Computed: true,
					},
					"only_include_users_with_sod_conflicts": schema.BoolAttribute{
						Computed: true,
					},
					"user_ids": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
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

	fmt.Println("INSIDE READ FOR CAMPAIGNS")
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
	data.CampaignType = types.StringValue(string(*campaign.CampaignType))
	data.Created = types.StringValue(campaign.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(campaign.CreatedBy)
	data.LastUpdated = types.StringValue(campaign.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(campaign.LastUpdatedBy)
	if campaign.RecurringCampaignId.Get() != nil {
		data.RecurringCampaignId = types.StringValue(*campaign.RecurringCampaignId.Get())
	} else {
		data.RecurringCampaignId = types.StringNull()
	}

	if campaign.RemediationSettings.AccessApproved != "" ||
		campaign.RemediationSettings.AccessRevoked != "" ||
		campaign.RemediationSettings.NoResponse != "" {

		if campaign.RemediationSettings.AutoRemediationSettings != nil {
			data.RemediationSettings.AutoRemediationSettings = &autoRemediationSettingsModel{
				IncludeAllIndirectAssignments: types.BoolValue(
					derefBool(campaign.RemediationSettings.AutoRemediationSettings.IncludeAllIndirectAssignments),
				),
			}
		}
		data.RemediationSettings = &campaignRemediationSettingsModel{
			AccessApproved: stringOrNull(string(campaign.RemediationSettings.AccessApproved)),
			AccessRevoked:  stringOrNull(string(campaign.RemediationSettings.AccessRevoked)),
			NoResponse:     stringOrNull(string(campaign.RemediationSettings.NoResponse)),
		}
	} else {
		data.RemediationSettings = nil // required to avoid "null to non-nullable" error
	}

	excluded := make([]targetResourceModel, len(campaign.ResourceSettings.ExcludedResources))
	for i, ex := range campaign.ResourceSettings.ExcludedResources {
		excluded[i] = targetResourceModel{
			ResourceId:   types.StringValue(*ex.ResourceId),
			ResourceType: types.StringValue(string(*ex.ResourceType)),
		}
	}

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
			ResourceType:                     types.StringValue(string(*tr.ResourceType)),
			EntitlementBundles:               bundles,
			Entitlements:                     entitlements,
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
	data.ReviewerSettings = d.createSettingModelDS(campaign)

	var endDate types.String
	if campaign.ScheduleSettings.EndDate != nil {
		endDate = types.StringValue(campaign.ScheduleSettings.EndDate.Format(time.RFC3339))
	}
	var recurrence []recurrenceModel
	if campaign.ScheduleSettings.Recurrence != nil {
		recurrenceList := make([]recurrenceModel, 1)
		recurrenceList[0] = recurrenceModel{
			Interval:     types.StringValue(campaign.ScheduleSettings.Recurrence.Interval),
			Ends:         types.StringValue(campaign.ScheduleSettings.Recurrence.Ends.Format(time.RFC3339)),
			RepeatOnType: types.StringValue(string(*campaign.ScheduleSettings.Recurrence.RepeatOnType)),
		}
		recurrence = recurrenceList
	} else {
		recurrence = nil
	}
	data.ScheduleSettings = &scheduleSettingsModel{
		DurationInDays: types.Int32Value(int32(campaign.ScheduleSettings.DurationInDays)),
		EndDate:        endDate,
		StartDate:      types.StringValue(campaign.ScheduleSettings.StartDate.Format(time.RFC3339)),
		TimeZone:       types.StringValue(campaign.ScheduleSettings.TimeZone),
		Type:           types.StringValue(string(campaign.ScheduleSettings.Type)),
		Recurrence:     recurrence,
	}

	var remindersReviewerBeforeCampaignCloseInSecs types.List
	if campaign.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs != nil {
		intValues := make([]attr.Value, 0, len(campaign.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs))
		for _, v := range campaign.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs {
			intValues = append(intValues, types.Int64Value(int64(v)))
		}
		remindersReviewerBeforeCampaignCloseInSecs = types.ListValueMust(
			types.Int32Type,
			intValues,
		)
	} else {
		remindersReviewerBeforeCampaignCloseInSecs = types.ListNull(types.Int32Type)
	}
	data.NotificationSettings = &notificationSettingsModel{
		NotifyReviewerAtCampaignEnd:                types.BoolValue(*campaign.NotificationSettings.NotifyReviewerAtCampaignEnd),
		NotifyReviewerDuringMidpointOfReview:       types.BoolValue(*campaign.NotificationSettings.NotifyReviewerDuringMidpointOfReview.Get()),
		NotifyReviewerWhenOverdue:                  types.BoolValue(*campaign.NotificationSettings.NotifyReviewerWhenOverdue.Get()),
		NotifyReviewerWhenReviewAssigned:           types.BoolValue(*campaign.NotificationSettings.NotifyReviewerWhenReviewAssigned),
		NotifyReviewPeriodEnd:                      types.BoolValue(*campaign.NotificationSettings.NotifyReviewPeriodEnd.Get()),
		RemindersReviewerBeforeCampaignCloseInSecs: remindersReviewerBeforeCampaignCloseInSecs,
	}

	var excludedUsersIds types.List
	if len(campaign.PrincipalScopeSettings.ExcludedUserIds) > 0 {
		excluded := make([]attr.Value, 0, len(campaign.PrincipalScopeSettings.ExcludedUserIds))
		for _, id := range campaign.PrincipalScopeSettings.ExcludedUserIds {
			excluded = append(excluded, types.StringValue(id))
		}
		excludedUsersIds = types.ListValueMust(types.StringType, excluded)
	} else {
		excludedUsersIds = types.ListNull(types.StringType)
	}

	var userIdList types.List
	if len(campaign.PrincipalScopeSettings.UserIds) > 0 {
		userIds := make([]attr.Value, 0, len(campaign.PrincipalScopeSettings.UserIds))
		for _, id := range campaign.PrincipalScopeSettings.UserIds {
			userIds = append(userIds, types.StringValue(id))
		}
		userIdList = types.ListValueMust(types.StringType, userIds)
	} else {
		userIdList = types.ListNull(types.StringType)
	}

	var groupIdList types.List
	if len(campaign.PrincipalScopeSettings.GroupIds) > 0 {
		groupIds := make([]attr.Value, 0, len(campaign.PrincipalScopeSettings.GroupIds))
		for _, id := range campaign.PrincipalScopeSettings.GroupIds {
			groupIds = append(groupIds, types.StringValue(id))
		}
		groupIdList = types.ListValueMust(types.StringType, groupIds)
	} else {
		groupIdList = types.ListNull(types.StringType)
	}

	var inactiveUsersScope []inactiveUsersScopeModel
	if campaign.PrincipalScopeSettings.PredefinedInactiveUsersScope != nil {
		inactiveUsersScope = []inactiveUsersScopeModel{
			{
				InactiveDays: types.Int32Value(*campaign.PrincipalScopeSettings.PredefinedInactiveUsersScope.InactiveDays),
			},
		}
	}

	var includeOnlyActiveUsers types.Bool
	if campaign.PrincipalScopeSettings.IncludeOnlyActiveUsers != nil {
		includeOnlyActiveUsers = types.BoolValue(*campaign.PrincipalScopeSettings.IncludeOnlyActiveUsers)
	} else {
		includeOnlyActiveUsers = types.BoolNull()
	}

	var onlyIncludeUsersWithSODConflicts types.Bool
	if campaign.PrincipalScopeSettings.OnlyIncludeUsersWithSODConflicts != nil {
		onlyIncludeUsersWithSODConflicts = types.BoolValue(*campaign.PrincipalScopeSettings.OnlyIncludeUsersWithSODConflicts)
	} else {
		onlyIncludeUsersWithSODConflicts = types.BoolNull()
	}

	var userScopeExpression types.String
	if campaign.PrincipalScopeSettings.UserScopeExpression != nil {
		userScopeExpression = types.StringValue(*campaign.PrincipalScopeSettings.UserScopeExpression)
	} else {
		userScopeExpression = types.StringNull()
	}

	data.PrincipalScope = &principalScopeSettingsModel{
		Type:                             types.StringValue(string(campaign.PrincipalScopeSettings.Type)),
		ExcludedUserIds:                  excludedUsersIds,
		GroupIds:                         groupIdList,
		IncludeOnlyActiveUsers:           includeOnlyActiveUsers,
		OnlyIncludeUsersWithSODConflicts: onlyIncludeUsersWithSODConflicts,
		UserIds:                          userIdList,
		UserScopeExpression:              userScopeExpression,
		PredefinedInactiveUsersScope:     inactiveUsersScope,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *campaignDataSource) createSettingModelDS(campaign *oktaInternalGovernance.CampaignFull) *reviewerSettingsModel {
	var (
		bulkDecisionDisabled    types.Bool
		fallbackReviewerId      types.String
		justificationRequired   types.Bool
		reassignmentDisabled    types.Bool
		reviewerGroupId         types.String
		reviewerID              types.String
		reviewerScopeExpression types.String
		selfReviewDisabled      types.Bool
	)

	if campaign.ReviewerSettings.BulkDecisionDisabled != nil {
		bulkDecisionDisabled = types.BoolValue(*campaign.ReviewerSettings.BulkDecisionDisabled)
	} else {
		bulkDecisionDisabled = types.BoolNull()
	}

	if campaign.ReviewerSettings.FallBackReviewerId != nil {
		fallbackReviewerId = types.StringValue(*campaign.ReviewerSettings.FallBackReviewerId)
	} else {
		fallbackReviewerId = types.StringNull()
	}

	if campaign.ReviewerSettings.JustificationRequired != nil {
		justificationRequired = types.BoolValue(*campaign.ReviewerSettings.JustificationRequired)
	} else {
		justificationRequired = types.BoolNull()
	}

	if campaign.ReviewerSettings.ReassignmentDisabled != nil {
		reassignmentDisabled = types.BoolValue(*campaign.ReviewerSettings.ReassignmentDisabled)
	} else {
		reassignmentDisabled = types.BoolNull()
	}

	if campaign.ReviewerSettings.ReviewerGroupId != nil {
		reviewerGroupId = types.StringValue(*campaign.ReviewerSettings.ReviewerGroupId)
	} else {
		reviewerGroupId = types.StringNull()
	}

	if campaign.ReviewerSettings.ReviewerId != nil {
		reviewerID = types.StringValue(*campaign.ReviewerSettings.ReviewerId)
	} else {
		reviewerID = types.StringNull()
	}

	if campaign.ReviewerSettings.ReviewerScopeExpression != nil {
		reviewerScopeExpression = types.StringValue(*campaign.ReviewerSettings.ReviewerScopeExpression)
	} else {
		reviewerScopeExpression = types.StringNull()
	}

	if campaign.ReviewerSettings.SelfReviewDisabled != nil {
		selfReviewDisabled = types.BoolValue(*campaign.ReviewerSettings.SelfReviewDisabled)
	} else {
		selfReviewDisabled = types.BoolNull()
	}

	reviewerLevels := make([]reviewerLevelModel, len(campaign.ReviewerSettings.ReviewerLevels))
	if campaign.ReviewerSettings.ReviewerLevels != nil {
		for _, level := range campaign.ReviewerSettings.ReviewerLevels {
			reviewerLevel := reviewerLevelModel{
				Type:                    types.StringValue(string(level.Type)),
				FallBackReviewerId:      types.StringValue(*level.FallBackReviewerId),
				ReviewerGroupId:         types.StringValue(*level.ReviewerGroupId),
				ReviewerId:              types.StringValue(*level.ReviewerId),
				ReviewerScopeExpression: types.StringValue(*level.ReviewerScopeExpression),
				SelfReviewDisabled:      types.BoolValue(*level.SelfReviewDisabled),
			}
			startReviews := make([]startReviewModel, 1)
			startReviews[0].OnDay = types.Int32Value(level.StartReview.OnDay)
			startReviews[0].When = types.StringValue(string(*level.StartReview.When))

			reviewerLevels = append(reviewerLevels, reviewerLevel)
		}
	}

	return &reviewerSettingsModel{
		Type:                    types.StringValue(string(campaign.ReviewerSettings.Type)),
		BulkDecisionDisabled:    bulkDecisionDisabled,
		FallbackReviewerId:      fallbackReviewerId,
		JustificationRequired:   justificationRequired,
		ReassignmentDisabled:    reassignmentDisabled,
		ReviewerGroupId:         reviewerGroupId,
		ReviewerId:              reviewerID,
		ReviewerScopeExpression: reviewerScopeExpression,
		SelfReviewDisabled:      selfReviewDisabled,
		ReviewerLevels:          reviewerLevels,
	}
}

func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func derefBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
