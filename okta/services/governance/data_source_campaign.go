package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &campaignDataSource{}

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
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier for the object.",
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
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the campaign.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the campaign.",
			},
			"campaign_type": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Identifies if it is a resource campaign or a user campaign. By default it is RESOURCE.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Human readable description.",
			},
			"recurring_campaign_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the recurring campaign if this campaign was created as part of a recurring schedule.",
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
						Computed:    true,
						Description: "Specifies the action if the reviewer doesn't respond to the request.",
					},
				},
				Blocks: map[string]schema.Block{
					"auto_remediation_settings": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"include_all_indirect_assignments": schema.BoolAttribute{
								Computed:    true,
								Optional:    true,
								Description: "When a group is selected to be automatically remediated.",
							},
						},
						Blocks: map[string]schema.Block{
							"include_only": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"resource_id": schema.StringAttribute{
											Required:    true,
											Description: "The resource ID of the target resource When type = GROUP, it will point to the group ID.",
										},
										"resource_type": schema.StringAttribute{
											Required:    true,
											Description: "The type of the resource to be automatically remediated. Only GROUP is supported.",
										},
									},
								},
								Description: "An array of resources to be automatically remediated.",
							},
						},
					},
				},
			},
			"resource_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of Okta resource.",
					},
					"include_admin_roles": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Include admin roles.",
					},
					"include_entitlements": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Include entitlements for this application.",
					},
					"individually_assigned_apps_only": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Only include individually assigned groups.",
					},
					"individually_assigned_groups_only": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Only include individually assigned groups.",
					},
					"only_include_out_of_policy_entitlements": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Description: "Only include out-of-policy entitlements.",
					},
				},
				Blocks: map[string]schema.Block{
					"excluded_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Computed:    true,
									Description: "Okta specific resource ID.",
								},
								"resource_type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of Okta resource.",
								},
							},
						},
						Description: "An array of resources that are excluded from the review.",
					},
					"target_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Required:    true,
									Description: "The resource ID that is being reviewed.",
								},
								"include_all_entitlements_and_bundles": schema.BoolAttribute{
									Computed:    true,
									Description: "Include all entitlements and entitlement bundles for this application.",
								},
								"resource_type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of Okta resource.",
								},
							},
							Blocks: map[string]schema.Block{
								"entitlement_bundles": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required:    true,
												Description: "The id of the entitlement bundle.",
											},
										},
									},
									Description: "An array of entitlement bundles associated with resourceId that should be chosen as target when creating reviews.",
								},
								"entitlements": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required:    true,
												Description: "The id of the entitlement.",
											},
											"include_all_values": schema.BoolAttribute{
												Computed:    true,
												Description: "Whether to include all values for this entitlement.",
											},
										},
										Blocks: map[string]schema.Block{
											"values": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Required:    true,
															Description: "The entitlement value id.",
														},
													},
												},
												Description: "Entitlement value ids",
											},
										},
									},
								},
							},
						},
						Description: "Represents a resource that will be part of Access certifications.",
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
	var data campaignDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	campaign, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().CampaignsAPI.GetCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading campaign",
			"Could not read Campaign, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(campaign.GetId())
	data.Name = types.StringValue(campaign.GetName())
	data.Status = types.StringValue(string(campaign.GetStatus()))
	data.Description = types.StringValue(campaign.GetDescription())
	data.CampaignType = types.StringValue(string(campaign.GetCampaignType()))
	data.Created = types.StringValue(campaign.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(campaign.GetCreatedBy())
	data.LastUpdated = types.StringValue(campaign.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(campaign.GetLastUpdatedBy())
	if campaign.RecurringCampaignId.Get() != nil {
		data.RecurringCampaignId = types.StringValue(campaign.GetRecurringCampaignId())
	}

	if campaign.RemediationSettings.AccessApproved != "" ||
		campaign.RemediationSettings.AccessRevoked != "" ||
		campaign.RemediationSettings.NoResponse != "" {

		if campaign.RemediationSettings.AutoRemediationSettings != nil {
			data.RemediationSettings.AutoRemediationSettings = &autoRemediationSettingsModel{
				IncludeAllIndirectAssignments: types.BoolValue(
					campaign.RemediationSettings.AutoRemediationSettings.GetIncludeAllIndirectAssignments(),
				),
			}
		}
		data.RemediationSettings = &campaignRemediationSettingsModel{
			AccessApproved: types.StringValue(string(campaign.RemediationSettings.GetAccessApproved())),
			AccessRevoked:  types.StringValue(string(campaign.RemediationSettings.GetAccessRevoked())),
			NoResponse:     types.StringValue(string(campaign.RemediationSettings.GetNoResponse())),
		}
	} else {
		data.RemediationSettings = nil
	}

	excluded := make([]excludedResourceModel, len(campaign.ResourceSettings.ExcludedResources))

	for _, ex := range campaign.ResourceSettings.ExcludedResources {
		excludedRes := excludedResourceModel{
			ResourceId: types.StringValue(ex.GetResourceId()),
		}
		if ex.ResourceType != nil {
			excludedRes.ResourceType = types.StringValue(string(ex.GetResourceType()))
		}
		excluded = append(excluded, excludedRes)
	}

	targets := make([]targetResourceModel, len(campaign.ResourceSettings.GetTargetResources()))
	for i, tr := range campaign.ResourceSettings.GetTargetResources() {
		entitlements := make([]entitlementModel, len(tr.GetEntitlements()))
		for j, ent := range tr.GetEntitlements() {
			values := make([]entitlementValueModel, len(ent.GetValues()))
			for k, v := range ent.GetValues() {
				values[k] = entitlementValueModel{Id: types.StringValue(v.GetId())}
			}
			entitlements[j] = entitlementModel{
				Id:               types.StringValue(ent.GetId()),
				IncludeAllValues: types.BoolValue(ent.GetIncludeAllValues()),
				Values:           values,
			}
		}
		bundles := make([]entitlementBundleModel, len(tr.GetEntitlementBundles()))
		for j, b := range tr.GetEntitlementBundles() {
			bundles[j] = entitlementBundleModel{Id: types.StringValue(b.GetId())}
		}
		targets[i] = targetResourceModel{
			ResourceId:                       types.StringValue(tr.GetResourceId()),
			IncludeAllEntitlementsAndBundles: types.BoolValue(tr.GetIncludeAllEntitlementsAndBundles()),
			ResourceType:                     types.StringValue(string(tr.GetResourceType())),
			EntitlementBundles:               bundles,
			Entitlements:                     entitlements,
		}
	}

	data.ResourceSettings = &resourceSettingsModel{
		Type:                               types.StringValue(string(campaign.ResourceSettings.GetType())),
		IncludeAdminRoles:                  types.BoolValue(campaign.ResourceSettings.GetIncludeAdminRoles()),
		IncludeEntitlements:                types.BoolValue(campaign.ResourceSettings.GetIncludeEntitlements()),
		IndividuallyAssignedAppsOnly:       types.BoolValue(campaign.ResourceSettings.GetIndividuallyAssignedAppsOnly()),
		IndividuallyAssignedGroupsOnly:     types.BoolValue(campaign.ResourceSettings.GetIndividuallyAssignedGroupsOnly()),
		OnlyIncludeOutOfPolicyEntitlements: types.BoolValue(campaign.ResourceSettings.GetOnlyIncludeOutOfPolicyEntitlements()),
		ExcludedResources:                  excluded,
		TargetResources:                    targets,
	}
	data.ReviewerSettings = d.createSettingModelDS(campaign)

	var endDate types.String
	if campaign.ScheduleSettings.EndDate != nil {
		endDate = types.StringValue(campaign.ScheduleSettings.GetEndDate().Format(time.RFC3339))
	}
	var recurrence []recurrenceModel
	if campaign.ScheduleSettings.Recurrence != nil {
		recurrenceList := make([]recurrenceModel, 1)
		recurrenceList[0] = recurrenceModel{
			Interval:     types.StringValue(campaign.ScheduleSettings.Recurrence.GetInterval()),
			Ends:         types.StringValue(campaign.ScheduleSettings.Recurrence.GetEnds().Format(time.RFC3339)),
			RepeatOnType: types.StringValue(string(campaign.ScheduleSettings.Recurrence.GetRepeatOnType())),
		}
		recurrence = recurrenceList
	} else {
		recurrence = nil
	}
	data.ScheduleSettings = &scheduleSettingsModel{
		DurationInDays: types.Int32Value(int32(campaign.ScheduleSettings.GetDurationInDays())),
		EndDate:        endDate,
		StartDate:      types.StringValue(campaign.ScheduleSettings.GetStartDate().Format(time.RFC3339)),
		TimeZone:       types.StringValue(campaign.ScheduleSettings.GetTimeZone()),
		Type:           types.StringValue(string(campaign.ScheduleSettings.GetType())),
		Recurrence:     recurrence,
	}

	var remindersReviewerBeforeCampaignCloseInSecs types.List
	if campaign.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs != nil {
		intValues := make([]attr.Value, 0, len(campaign.NotificationSettings.GetRemindersReviewerBeforeCampaignCloseInSecs()))
		for _, v := range campaign.NotificationSettings.GetRemindersReviewerBeforeCampaignCloseInSecs() {
			intValues = append(intValues, types.Int64Value(int64(v)))
		}
		remindersReviewerBeforeCampaignCloseInSecs = types.ListValueMust(
			types.Int64Type,
			intValues,
		)
	} else {
		remindersReviewerBeforeCampaignCloseInSecs = types.ListNull(types.Int64Type)
	}
	data.NotificationSettings = &notificationSettingsModel{
		NotifyReviewerAtCampaignEnd:                types.BoolValue(campaign.NotificationSettings.GetNotifyReviewerAtCampaignEnd()),
		NotifyReviewerDuringMidpointOfReview:       types.BoolValue(campaign.NotificationSettings.GetNotifyReviewerDuringMidpointOfReview()),
		NotifyReviewerWhenOverdue:                  types.BoolValue(campaign.NotificationSettings.GetNotifyReviewerWhenOverdue()),
		NotifyReviewerWhenReviewAssigned:           types.BoolValue(campaign.NotificationSettings.GetNotifyReviewerWhenReviewAssigned()),
		NotifyReviewPeriodEnd:                      types.BoolValue(campaign.NotificationSettings.GetNotifyReviewPeriodEnd()),
		RemindersReviewerBeforeCampaignCloseInSecs: remindersReviewerBeforeCampaignCloseInSecs,
	}

	var excludedUsersIds types.List
	if len(campaign.PrincipalScopeSettings.GetExcludedUserIds()) > 0 {
		excluded := make([]attr.Value, 0, len(campaign.PrincipalScopeSettings.GetExcludedUserIds()))
		for _, id := range campaign.PrincipalScopeSettings.GetExcludedUserIds() {
			excluded = append(excluded, types.StringValue(id))
		}
		excludedUsersIds = types.ListValueMust(types.StringType, excluded)
	} else {
		excludedUsersIds = types.ListNull(types.StringType)
	}

	var userIdList types.List
	if len(campaign.PrincipalScopeSettings.GetUserIds()) > 0 {
		userIds := make([]attr.Value, 0, len(campaign.PrincipalScopeSettings.GetUserIds()))
		for _, id := range campaign.PrincipalScopeSettings.GetUserIds() {
			userIds = append(userIds, types.StringValue(id))
		}
		userIdList = types.ListValueMust(types.StringType, userIds)
	} else {
		userIdList = types.ListNull(types.StringType)
	}

	var groupIdList types.List
	if len(campaign.PrincipalScopeSettings.GetGroupIds()) > 0 {
		groupIds := make([]attr.Value, 0, len(campaign.PrincipalScopeSettings.GetGroupIds()))
		for _, id := range campaign.PrincipalScopeSettings.GetGroupIds() {
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
				InactiveDays: types.Int32Value(campaign.PrincipalScopeSettings.PredefinedInactiveUsersScope.GetInactiveDays()),
			},
		}
	}

	var includeOnlyActiveUsers types.Bool
	if campaign.PrincipalScopeSettings.IncludeOnlyActiveUsers != nil {
		includeOnlyActiveUsers = types.BoolValue(campaign.PrincipalScopeSettings.GetIncludeOnlyActiveUsers())
	}

	var onlyIncludeUsersWithSODConflicts types.Bool
	if campaign.PrincipalScopeSettings.OnlyIncludeUsersWithSODConflicts != nil {
		onlyIncludeUsersWithSODConflicts = types.BoolValue(campaign.PrincipalScopeSettings.GetOnlyIncludeUsersWithSODConflicts())
	}

	var userScopeExpression types.String
	if campaign.PrincipalScopeSettings.UserScopeExpression != nil {
		userScopeExpression = types.StringValue(campaign.PrincipalScopeSettings.GetUserScopeExpression())
	}

	data.PrincipalScope = &principalScopeSettingsModel{
		Type:                             types.StringValue(string(campaign.PrincipalScopeSettings.GetType())),
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

func (d *campaignDataSource) createSettingModelDS(campaign *governance.CampaignFull) *reviewerSettingsModel {
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
		bulkDecisionDisabled = types.BoolValue(campaign.ReviewerSettings.GetBulkDecisionDisabled())
	}

	if campaign.ReviewerSettings.FallBackReviewerId != nil {
		fallbackReviewerId = types.StringValue(campaign.ReviewerSettings.GetFallBackReviewerId())
	}

	if campaign.ReviewerSettings.JustificationRequired != nil {
		justificationRequired = types.BoolValue(campaign.ReviewerSettings.GetJustificationRequired())
	}

	if campaign.ReviewerSettings.ReassignmentDisabled != nil {
		reassignmentDisabled = types.BoolValue(campaign.ReviewerSettings.GetReassignmentDisabled())
	}

	if campaign.ReviewerSettings.ReviewerGroupId != nil {
		reviewerGroupId = types.StringValue(campaign.ReviewerSettings.GetReviewerGroupId())
	}

	if campaign.ReviewerSettings.ReviewerId != nil {
		reviewerID = types.StringValue(campaign.ReviewerSettings.GetReviewerId())
	}

	if campaign.ReviewerSettings.ReviewerScopeExpression != nil {
		reviewerScopeExpression = types.StringValue(campaign.ReviewerSettings.GetReviewerScopeExpression())
	}

	if campaign.ReviewerSettings.SelfReviewDisabled != nil {
		selfReviewDisabled = types.BoolValue(campaign.ReviewerSettings.GetSelfReviewDisabled())
	}

	reviewerLevels := make([]reviewerLevelModel, len(campaign.ReviewerSettings.GetReviewerLevels()))
	if campaign.ReviewerSettings.ReviewerLevels != nil {
		for _, level := range campaign.ReviewerSettings.GetReviewerLevels() {
			reviewerLevel := reviewerLevelModel{
				Type:                    types.StringValue(string(level.GetType())),
				FallBackReviewerId:      types.StringValue(level.GetFallBackReviewerId()),
				ReviewerGroupId:         types.StringValue(level.GetReviewerGroupId()),
				ReviewerId:              types.StringValue(level.GetReviewerId()),
				ReviewerScopeExpression: types.StringValue(level.GetReviewerScopeExpression()),
				SelfReviewDisabled:      types.BoolValue(level.GetSelfReviewDisabled()),
			}
			startReviews := make([]startReviewModel, 1)
			startReviews[0].OnDay = types.Int32Value(level.StartReview.GetOnDay())
			startReviews[0].When = types.StringValue(string(level.StartReview.GetWhen()))

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
