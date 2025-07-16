package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/okta/terraform-provider-okta/okta/config"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &campaignResource{}
	_ resource.ResourceWithConfigure   = &campaignResource{}
	_ resource.ResourceWithImportState = &campaignResource{}
)

func newCampaignResource() resource.Resource {
	return &campaignResource{}
}

type campaignResource struct {
	*config.Config
}

func (r *campaignResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *campaignResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type campaignResourceModel struct {
	Id                   types.String                      `tfsdk:"id"`
	CampaignTier         types.String                      `tfsdk:"campaign_tier"` // Indicates the minimum required SKU to manage the campaign. Values can be `BASIC` and `PREMIUM`.
	Name                 types.String                      `tfsdk:"name"`
	LaunchCampaign       types.Bool                        `tfsdk:"launch_campaign"`
	CampaignType         types.String                      `tfsdk:"campaign_type"`
	Description          types.String                      `tfsdk:"description"`
	RemediationSettings  *campaignRemediationSettingsModel `tfsdk:"remediation_settings"`
	ResourceSettings     *resourceSettingsModel            `tfsdk:"resource_settings"`
	ReviewerSettings     *reviewerSettingsModel            `tfsdk:"reviewer_settings"`
	ScheduleSettings     *scheduleSettingsModel            `tfsdk:"schedule_settings"`
	NotificationSettings *notificationSettingsModel        `tfsdk:"notification_settings"`
	PrincipalScope       *principalScopeSettingsModel      `tfsdk:"principal_scope_settings"`
}

type campaignRemediationSettingsModel struct {
	AccessApproved          types.String                  `tfsdk:"access_approved"`
	AccessRevoked           types.String                  `tfsdk:"access_revoked"`
	NoResponse              types.String                  `tfsdk:"no_response"`
	AutoRemediationSettings *autoRemediationSettingsModel `tfsdk:"auto_remediation_settings"`
}

type autoRemediationSettingsModel struct {
	IncludeAllIndirectAssignments types.Bool            `tfsdk:"include_all_indirect_assignments"`
	IncludeOnly                   []targetResourceModel `tfsdk:"include_only"`
}

type principalScopeSettingsModel struct {
	Type                             types.String              `tfsdk:"type"`
	ExcludedUserIds                  types.List                `tfsdk:"excluded_user_ids"`
	GroupIds                         types.List                `tfsdk:"group_ids"`
	IncludeOnlyActiveUsers           types.Bool                `tfsdk:"include_only_active_users"`
	OnlyIncludeUsersWithSODConflicts types.Bool                `tfsdk:"only_include_users_with_sod_conflicts"`
	UserIds                          types.List                `tfsdk:"user_ids"`
	UserScopeExpression              types.String              `tfsdk:"user_scope_expression"`
	PredefinedInactiveUsersScope     []inactiveUsersScopeModel `tfsdk:"predefined_inactive_users_scope"`
}

type entitlementBundleModel struct {
	Id types.String `tfsdk:"id"`
}

type inactiveUsersScopeModel struct {
	InactiveDays types.Int32 `tfsdk:"inactive_days"`
}

type entitlementValueModel struct {
	Id types.String `tfsdk:"id"`
}

type entitlementModel struct {
	Id               types.String            `tfsdk:"id"`
	IncludeAllValues types.Bool              `tfsdk:"include_all_values"`
	Values           []entitlementValueModel `tfsdk:"values"`
}

type resourceSettingsModel struct {
	Type                               types.String          `tfsdk:"type"`
	IncludeAdminRoles                  types.Bool            `tfsdk:"include_admin_roles"`
	IncludeEntitlements                types.Bool            `tfsdk:"include_entitlements"`
	IndividuallyAssignedAppsOnly       types.Bool            `tfsdk:"individually_assigned_apps_only"`
	IndividuallyAssignedGroupsOnly     types.Bool            `tfsdk:"individually_assigned_groups_only"`
	OnlyIncludeOutOfPolicyEntitlements types.Bool            `tfsdk:"only_include_out_of_policy_entitlements"`
	ExcludedResources                  []targetResourceModel `tfsdk:"excluded_resources"`
	TargetResources                    []targetResourceModel `tfsdk:"target_resources"`
}

type targetResourceModel struct {
	ResourceId                       types.String             `tfsdk:"resource_id"`
	IncludeAllEntitlementsAndBundles types.Bool               `tfsdk:"include_all_entitlements_and_bundles"`
	ResourceType                     types.String             `tfsdk:"resource_type"`
	EntitlementBundles               []entitlementBundleModel `tfsdk:"entitlement_bundles"`
	Entitlements                     []entitlementModel       `tfsdk:"entitlements"`
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
	FallBackReviewerId      types.String       `tfsdk:"fallback_reviewer_id"`
	ReviewerGroupId         types.String       `tfsdk:"reviewer_group_id"`
	ReviewerId              types.String       `tfsdk:"reviewer_id"`
	ReviewerScopeExpression types.String       `tfsdk:"reviewer_scope_expression"`
	SelfReviewDisabled      types.Bool         `tfsdk:"self_review_disabled"`
	StartReview             []startReviewModel `tfsdk:"start_review"`
}

type startReviewModel struct {
	OnDay types.Int32  `tfsdk:"on_day"`
	When  types.String `tfsdk:"when"`
}

type scheduleSettingsModel struct {
	DurationInDays types.Int32       `tfsdk:"duration_in_days"`
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
	NotifyReviewerAtCampaignEnd                types.Bool `tfsdk:"notify_reviewer_at_campaign_end"`
	NotifyReviewerDuringMidpointOfReview       types.Bool `tfsdk:"notify_reviewer_during_midpoint_of_review"`
	NotifyReviewerWhenOverdue                  types.Bool `tfsdk:"notify_reviewer_when_overdue"`
	NotifyReviewerWhenReviewAssigned           types.Bool `tfsdk:"notify_reviewer_when_review_assigned"`
	NotifyReviewPeriodEnd                      types.Bool `tfsdk:"notify_review_period_end"`
	RemindersReviewerBeforeCampaignCloseInSecs types.List `tfsdk:"reminders_reviewer_before_campaign_close_in_secs"`
}

func (r *campaignResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaign"
}

func (r *campaignResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the campaign. Maintain some uniqueness when naming the campaign as it helps to identify and filter for campaigns when needed.",
			},
			"launch_campaign": schema.BoolAttribute{
				Optional:    true,
				Description: "Launch the campaign after creation. Defaults to false.",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"campaign_tier": schema.StringAttribute{
				Optional:    true,
				Description: "Indicates the minimum required SKU to manage the campaign. Values can be `BASIC` and `PREMIUM`.",
			},
			"campaign_type": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies if it is a resource campaign or a user campaign. By default it is RESOURCE.Values can be `RESOURCE` and `USER`.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description about the campaign.",
			},
		},
		Blocks: map[string]schema.Block{
			"remediation_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"access_approved": schema.StringAttribute{
						Required:    true,
						Description: "Specifies the action by default if the reviewer approves access. NO_ACTION indicates there is no remediation action and the user retains access.",
					},
					"access_revoked": schema.StringAttribute{
						Required:    true,
						Description: "Specifies the action if the reviewer revokes access. NO_ACTION indicates the user retains the same access. DENY indicates the user will have their access revoked as long as they are not assigned to a group through Group Rules.",
					},
					"no_response": schema.StringAttribute{
						Required: true,
					},
				},
				Blocks: map[string]schema.Block{
					"auto_remediation_settings": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"include_all_indirect_assignments": schema.BoolAttribute{
								Optional:    true,
								Description: "If true, all indirect assignments will be included in the campaign. If false, only direct assignments will be included.",
							},
						},
						Blocks: map[string]schema.Block{
							"include_only": schema.ListNestedBlock{
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"resource_id": schema.StringAttribute{
											Optional:    true,
											Description: "The ID of the resource to include in the campaign.",
										},
										"resource_type": schema.StringAttribute{
											Optional:    true,
											Description: "The type of the resource to include in the campaign. Valid values are 'APPLICATION', 'GROUP', 'ENTITLEMENT', 'ENTITLEMENT_BUNDLE'.",
										},
									},
								},
							},
						},
					},
				},
				Description: "Specify the action to be taken after a reviewer makes a decision to APPROVE or REVOKE the access, or if the campaign was CLOSED and there was no response from the reviewer.",
			},
			"resource_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"include_admin_roles": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"include_entitlements": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"individually_assigned_apps_only": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"individually_assigned_groups_only": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"only_include_out_of_policy_entitlements": schema.BoolAttribute{
						Computed: true,
						Optional: true,
						Default:  booldefault.StaticBool(false),
					},
				},
				Blocks: map[string]schema.Block{
					"excluded_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Optional: true,
								},
								"resource_type": schema.StringAttribute{
									Optional: true,
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
								"resource_type": schema.StringAttribute{
									Required: true,
								},
								"include_all_entitlements_and_bundles": schema.BoolAttribute{
									Optional: true,
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
												Optional: true,
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
						Required: true,
					},
					"bulk_decision_disabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"fallback_reviewer_id": schema.StringAttribute{
						Optional: true,
					},
					"justification_required": schema.BoolAttribute{
						Optional: true,
					},
					"reassignment_disabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"self_review_disabled": schema.BoolAttribute{
						Optional: true,
					},
					"reviewer_group_id": schema.StringAttribute{
						Optional: true,
					},
					"reviewer_id": schema.StringAttribute{
						Optional: true,
					},
					"reviewer_scope_expression": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"reviewer_levels": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},
								"fallback_reviewer_id": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"reviewer_group_id": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"reviewer_id": schema.StringAttribute{
									Required: true,
								},
								"reviewer_scope_expression": schema.StringAttribute{
									Computed: true,
								},
								"self_review_disabled": schema.BoolAttribute{
									Computed: true,
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"start_review": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"on_day": schema.Int64Attribute{
												Computed: true,
												Optional: true,
											},
											"when": schema.StringAttribute{
												Computed: true,
												Optional: true,
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
					"start_date": schema.StringAttribute{
						Required: true,
					},
					"duration_in_days": schema.Int32Attribute{
						Required: true,
					},
					"time_zone": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
					"end_date": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"recurrence": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"interval": schema.StringAttribute{
									Required: true,
								},
								"ends": schema.StringAttribute{
									Optional: true,
								},
								"repeat_on_type": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			"notification_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"notify_reviewer_at_campaign_end": schema.BoolAttribute{
						Optional: true,
					},
					"notify_reviewer_during_midpoint_of_review": schema.BoolAttribute{
						Optional: true,
					},
					"notify_reviewer_when_overdue": schema.BoolAttribute{
						Optional: true,
					},
					"notify_reviewer_when_review_assigned": schema.BoolAttribute{
						Optional: true,
					},
					"notify_review_period_end": schema.BoolAttribute{
						Optional: true,
					},
					"reminders_reviewer_before_campaign_close_in_secs": schema.ListAttribute{
						Optional:    true,
						ElementType: types.Int64Type,
						Description: "Specifies times (in seconds) to send reminders to reviewers before the campaign closes. Max 3 values. Example: [86400, 172800, 604800]",
					},
				},
			},
			"principal_scope_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"excluded_user_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"group_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"include_only_active_users": schema.BoolAttribute{
						Computed: true,
						Optional: true,
						Default:  booldefault.StaticBool(false),
					},
					"only_include_users_with_sod_conflicts": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"user_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"user_scope_expression": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"predefined_inactive_users_scope": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"inactive_days": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}

}

func (r *campaignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data campaignResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	fmt.Println("Schedule type = ", data.ScheduleSettings.Type)

	if resp.Diagnostics.HasError() {
		return
	}

	campaign, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.CreateCampaign(ctx).CampaignMutable(buildCampaign(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Campaign",
			"Could not create Campaign, unexpected error: "+err.Error(),
		)
		return
	}
	// Example data value setting
	data.Id = types.StringValue(campaign.Id)
	if data.LaunchCampaign.ValueBool() {
		_, err = r.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.LaunchCampaign(ctx, campaign.Id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error launching Campaign",
				"Could not launch Campaign after creation, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *campaignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data campaignResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getCampaignResponse, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.GetCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading campaign",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyCampaignsToState(ctx, getCampaignResponse, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func applyCampaignsToState(ctx context.Context, resp *oktaInternalGovernance.CampaignFull, c *campaignResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	c.Id = types.StringValue(resp.Id)
	c.Name = types.StringValue(resp.Name)
	c.CampaignType = types.StringValue(string(*resp.CampaignType))
	c.Description = types.StringValue(*resp.Description)

	c.RemediationSettings = &campaignRemediationSettingsModel{}
	if resp.RemediationSettings.AccessRevoked != "" {
		c.RemediationSettings.AccessRevoked = types.StringValue(string(resp.RemediationSettings.AccessRevoked))
	}
	if resp.RemediationSettings.AccessApproved != "" {
		c.RemediationSettings.AccessApproved = types.StringValue(string(resp.RemediationSettings.AccessApproved))
	}
	if resp.RemediationSettings.NoResponse != "" {
		c.RemediationSettings.NoResponse = types.StringValue(string(resp.RemediationSettings.NoResponse))
	}
	if resp.RemediationSettings.AutoRemediationSettings != nil {
		c.RemediationSettings.AutoRemediationSettings.IncludeAllIndirectAssignments = types.BoolValue(*resp.RemediationSettings.AutoRemediationSettings.IncludeAllIndirectAssignments)
		for _, includeOnly := range resp.RemediationSettings.AutoRemediationSettings.IncludeOnly {
			targetResource := targetResourceModel{
				ResourceId:   types.StringValue(*includeOnly.ResourceId),
				ResourceType: types.StringValue(string(*includeOnly.ResourceType)),
			}
			c.RemediationSettings.AutoRemediationSettings.IncludeOnly = append(c.RemediationSettings.AutoRemediationSettings.IncludeOnly, targetResource)
		}
	}

	c.ResourceSettings = &resourceSettingsModel{}
	if resp.ResourceSettings.Type != "" {
		c.ResourceSettings.Type = types.StringValue(string(resp.ResourceSettings.Type))
	}
	if len(resp.ResourceSettings.TargetResources) > 0 {
		c.ResourceSettings.TargetResources = []targetResourceModel{
			{
				ResourceId:   types.StringValue(resp.ResourceSettings.TargetResources[0].ResourceId),
				ResourceType: types.StringValue(string(*resp.ResourceSettings.TargetResources[0].ResourceType)),
			},
		}
	}
	if len(resp.ResourceSettings.ExcludedResources) > 0 {
		c.ResourceSettings.ExcludedResources = []targetResourceModel{
			{
				ResourceId:   types.StringValue(*resp.ResourceSettings.ExcludedResources[0].ResourceId),
				ResourceType: types.StringValue(string(*resp.ResourceSettings.ExcludedResources[0].ResourceType)),
			},
		}
	}
	if resp.ResourceSettings.IncludeAdminRoles != nil {
		c.ResourceSettings.IncludeAdminRoles = types.BoolValue(*resp.ResourceSettings.IncludeAdminRoles)
	}
	if resp.ResourceSettings.IncludeEntitlements != nil {
		c.ResourceSettings.IncludeEntitlements = types.BoolValue(*resp.ResourceSettings.IncludeEntitlements)
	}
	if resp.ResourceSettings.IndividuallyAssignedAppsOnly != nil {
		c.ResourceSettings.IndividuallyAssignedAppsOnly = types.BoolValue(*resp.ResourceSettings.IndividuallyAssignedAppsOnly)
	}
	if resp.ResourceSettings.IndividuallyAssignedGroupsOnly != nil {
		c.ResourceSettings.IndividuallyAssignedGroupsOnly = types.BoolValue(*resp.ResourceSettings.IndividuallyAssignedGroupsOnly)
	}
	if resp.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements != nil {
		c.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements = types.BoolValue(*resp.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements)
	}

	c.ReviewerSettings = &reviewerSettingsModel{}
	c.ReviewerSettings.Type = types.StringValue(string(resp.ReviewerSettings.Type))
	if resp.ReviewerSettings.BulkDecisionDisabled != nil {
		c.ReviewerSettings.BulkDecisionDisabled = types.BoolValue(*resp.ReviewerSettings.BulkDecisionDisabled)
	}
	if resp.ReviewerSettings.FallBackReviewerId != nil {
		c.ReviewerSettings.FallbackReviewerId = types.StringValue(*resp.ReviewerSettings.FallBackReviewerId)
	}
	if resp.ReviewerSettings.JustificationRequired != nil {
		c.ReviewerSettings.JustificationRequired = types.BoolValue(*resp.ReviewerSettings.JustificationRequired)
	}
	if resp.ReviewerSettings.ReassignmentDisabled != nil {
		c.ReviewerSettings.ReassignmentDisabled = types.BoolValue(*resp.ReviewerSettings.ReassignmentDisabled)
	}
	if resp.ReviewerSettings.ReviewerGroupId != nil {
		c.ReviewerSettings.ReviewerGroupId = types.StringValue(*resp.ReviewerSettings.ReviewerGroupId)
	}
	if resp.ReviewerSettings.ReviewerId != nil {
		c.ReviewerSettings.ReviewerId = types.StringValue(*resp.ReviewerSettings.ReviewerId)
	}
	if resp.ReviewerSettings.ReviewerScopeExpression != nil {
		c.ReviewerSettings.ReviewerScopeExpression = types.StringValue(*resp.ReviewerSettings.ReviewerScopeExpression)
	}
	if resp.ReviewerSettings.SelfReviewDisabled != nil {
		c.ReviewerSettings.SelfReviewDisabled = types.BoolValue(*resp.ReviewerSettings.SelfReviewDisabled)
	}
	if resp.ReviewerSettings.ReviewerLevels != nil {
		c.ReviewerSettings.ReviewerLevels = make([]reviewerLevelModel, 0, len(resp.ReviewerSettings.ReviewerLevels))
		for _, level := range resp.ReviewerSettings.ReviewerLevels {
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

			c.ReviewerSettings.ReviewerLevels = append(c.ReviewerSettings.ReviewerLevels, reviewerLevel)
		}
	}

	fmt.Println("ReviewerSettings:", c.ReviewerSettings.FallbackReviewerId.ValueString(), c.ReviewerSettings.ReviewerGroupId.ValueString())

	c.ScheduleSettings = &scheduleSettingsModel{}
	c.ScheduleSettings.StartDate = types.StringValue(resp.ScheduleSettings.StartDate.UTC().Format("2006-01-02T15:04:05.000Z"))
	c.ScheduleSettings.DurationInDays = types.Int32Value(int32(resp.ScheduleSettings.DurationInDays))
	c.ScheduleSettings.TimeZone = types.StringValue(resp.ScheduleSettings.TimeZone)
	c.ScheduleSettings.Type = types.StringValue(string(resp.ScheduleSettings.Type))
	c.ScheduleSettings.DurationInDays = types.Int32Value(int32(resp.ScheduleSettings.DurationInDays))
	if resp.ScheduleSettings.Recurrence != nil {
		c.ScheduleSettings.Recurrence = make([]recurrenceModel, 1)
		rec := recurrenceModel{
			Interval:     types.StringValue(resp.ScheduleSettings.Recurrence.Interval),
			Ends:         types.StringValue(resp.ScheduleSettings.Recurrence.Ends.UTC().Format("2006-01-02T15:04:05.000Z")),
			RepeatOnType: types.StringValue(string(*resp.ScheduleSettings.Recurrence.RepeatOnType)),
		}
		c.ScheduleSettings.Recurrence = append(c.ScheduleSettings.Recurrence, rec)
	}

	c.NotificationSettings = &notificationSettingsModel{}
	if resp.NotificationSettings.NotifyReviewerAtCampaignEnd != nil {
		c.NotificationSettings.NotifyReviewerAtCampaignEnd = types.BoolValue(*resp.NotificationSettings.NotifyReviewerAtCampaignEnd)
	}
	if resp.NotificationSettings.NotifyReviewerDuringMidpointOfReview.Get() != nil {
		c.NotificationSettings.NotifyReviewerDuringMidpointOfReview = types.BoolValue(*resp.NotificationSettings.NotifyReviewerDuringMidpointOfReview.Get())
	}
	if resp.NotificationSettings.NotifyReviewerWhenOverdue.Get() != nil {
		c.NotificationSettings.NotifyReviewerWhenOverdue = types.BoolValue(*resp.NotificationSettings.NotifyReviewerWhenOverdue.Get())
	}
	if resp.NotificationSettings.NotifyReviewerWhenReviewAssigned != nil {
		c.NotificationSettings.NotifyReviewerWhenReviewAssigned = types.BoolValue(*resp.NotificationSettings.NotifyReviewerWhenReviewAssigned)
	}
	if resp.NotificationSettings.NotifyReviewPeriodEnd.Get() != nil {
		c.NotificationSettings.NotifyReviewPeriodEnd = types.BoolValue(*resp.NotificationSettings.NotifyReviewPeriodEnd.Get())
	}
	if resp.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs != nil {
		reminders := make([]int64, 0, len(resp.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs))
		for _, v := range resp.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs {
			reminders = append(reminders, int64(v))
		}

		listVal, _ := types.ListValueFrom(ctx, types.Int64Type, reminders)

		c.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs = listVal
	} else {
		// explicitly set an empty list with correct type
		c.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs = types.ListNull(types.Int64Type)
	}

	c.PrincipalScope = &principalScopeSettingsModel{}
	if resp.PrincipalScopeSettings.Type != "" {
		c.PrincipalScope.Type = types.StringValue(string(resp.PrincipalScopeSettings.Type))
	}
	if resp.PrincipalScopeSettings.ExcludedUserIds != nil && len(resp.PrincipalScopeSettings.ExcludedUserIds) > 0 {
		excluded := make([]attr.Value, 0, len(resp.PrincipalScopeSettings.ExcludedUserIds))
		for _, id := range resp.PrincipalScopeSettings.ExcludedUserIds {
			excluded = append(excluded, types.StringValue(id))
		}
		var diags diag.Diagnostics
		c.PrincipalScope.ExcludedUserIds, diags = types.ListValue(types.StringType, excluded)
		if diags.HasError() {
			println(diags.Errors())
		}
	} else {
		c.PrincipalScope.ExcludedUserIds = types.ListNull(types.StringType)
	}
	if resp.PrincipalScopeSettings.GroupIds != nil && len(resp.PrincipalScopeSettings.GroupIds) > 0 {
		groupIds := make([]attr.Value, 0, len(resp.PrincipalScopeSettings.GroupIds))
		for _, id := range resp.PrincipalScopeSettings.GroupIds {
			groupIds = append(groupIds, types.StringValue(id))
		}
		c.PrincipalScope.GroupIds = types.ListValueMust(types.StringType, groupIds)
	} else {
		c.PrincipalScope.GroupIds = types.ListNull(types.StringType)
	}
	if resp.PrincipalScopeSettings.IncludeOnlyActiveUsers != nil {
		c.PrincipalScope.IncludeOnlyActiveUsers = types.BoolValue(*resp.PrincipalScopeSettings.IncludeOnlyActiveUsers)
	}
	if resp.PrincipalScopeSettings.OnlyIncludeUsersWithSODConflicts != nil {
		c.PrincipalScope.OnlyIncludeUsersWithSODConflicts = types.BoolValue(*resp.PrincipalScopeSettings.OnlyIncludeUsersWithSODConflicts)
	}
	if resp.PrincipalScopeSettings.UserIds != nil && len(resp.PrincipalScopeSettings.UserIds) > 0 {
		listVal, _ := types.ListValueFrom(ctx, types.StringType, resp.PrincipalScopeSettings.UserIds)
		c.PrincipalScope.UserIds = listVal
	} else {
		c.PrincipalScope.UserIds = types.ListNull(types.StringType)
	}
	if resp.PrincipalScopeSettings.UserScopeExpression != nil {
		c.PrincipalScope.UserScopeExpression = types.StringValue(*resp.PrincipalScopeSettings.UserScopeExpression)
	}
	if resp.PrincipalScopeSettings.PredefinedInactiveUsersScope != nil {
		c.PrincipalScope.PredefinedInactiveUsersScope = []inactiveUsersScopeModel{
			{
				InactiveDays: types.Int32Value(*resp.PrincipalScopeSettings.PredefinedInactiveUsersScope.InactiveDays),
			},
		}
	}

	return diags
}

func (r *campaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data campaignResourceModel
	var state campaignResourceModel

	// Load both planned and current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if only 'launch_campaign' changed
	if !data.LaunchCampaign.Equal(state.LaunchCampaign) && isOnlyLaunchChanged(data, state) {
		// Call the create API again with updated launch flag
		log.Println("launch_campaign updated — calling CreateCampaign API again")
		data.Id = types.StringValue(state.Id.ValueString())
		if data.LaunchCampaign.ValueBool() {
			_, err := r.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.LaunchCampaign(ctx, state.Id.ValueString()).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"Error launching Campaign",
					"Could not launch Campaign after creation, unexpected error: "+err.Error(),
				)
				return
			}

			state.LaunchCampaign = data.LaunchCampaign
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		}
		return
	}

	resp.Diagnostics.AddError(
		"Update Not Supported",
		"No other fields other than launch_campaign and end_campaign are updatable for this resource. Terraform will retain the existing state.",
	)
}

func (r *campaignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data campaignResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Println("Deleting Campaign with ID:", data.Id.ValueString())
	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.DeleteCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	return
}

func buildCampaign(d campaignResourceModel) oktaInternalGovernance.CampaignMutable {
	startDate := d.ScheduleSettings.StartDate.ValueString()
	parsedStartDate, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		log.Printf("invalid start_date format: %v", err)
		return oktaInternalGovernance.CampaignMutable{}
	}

	// Convert target resources
	var targetResources []oktaInternalGovernance.TargetResourcesRequestInner
	for _, tr := range d.ResourceSettings.TargetResources {
		rt := oktaInternalGovernance.ResourceType(tr.ResourceType.ValueString())

		targetResources = append(targetResources, oktaInternalGovernance.TargetResourcesRequestInner{
			ResourceId:   tr.ResourceId.ValueString(),
			ResourceType: &rt,
		})
	}

	var ExcludedResources = make([]oktaInternalGovernance.ResourceSettingsMutableExcludedResourcesInner, 0, len(d.ResourceSettings.ExcludedResources))
	for _, ex := range d.ResourceSettings.ExcludedResources {
		rt := oktaInternalGovernance.ResourceType(ex.ResourceType.ValueString())
		ExcludedResources = append(ExcludedResources, oktaInternalGovernance.ResourceSettingsMutableExcludedResourcesInner{
			ResourceId:   ex.ResourceId.ValueStringPointer(),
			ResourceType: &rt,
		})
	}

	//todo check if this is correct, the r[0] part
	var recur oktaInternalGovernance.RecurrenceDefinitionMutable
	//for _, r := range d.ScheduleSettings.Recurrence {
	r := d.ScheduleSettings.Recurrence // Assuming only one recurrence for simplicity
	if len(r) != 0 {
		endStr := r[0].Ends.ValueString()
		parsedTime, _ := time.Parse(time.RFC3339, endStr)
		repeatStr := oktaInternalGovernance.RecurrenceRepeatOnType(r[0].RepeatOnType.ValueString())
		recur = oktaInternalGovernance.RecurrenceDefinitionMutable{
			Interval:     r[0].Interval.ValueString(),
			Ends:         &parsedTime,
			RepeatOnType: &repeatStr,
		}
	}
	//}
	var remindersReviewerBeforeCampaignCloseInSecs []int32
	d.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs.ElementsAs(context.Background(), &remindersReviewerBeforeCampaignCloseInSecs, false)

	var autoRemediationSettings *oktaInternalGovernance.AutoRemediationSettings
	if d.RemediationSettings.AutoRemediationSettings != nil {
		var includeOnlyConverted []oktaInternalGovernance.AutoRemediationSettingsIncludeOnlyInner

		for _, tr := range d.RemediationSettings.AutoRemediationSettings.IncludeOnly {
			rt := oktaInternalGovernance.AutoRemediationResourceType(tr.ResourceType.ValueString())
			includeOnlyConverted = append(includeOnlyConverted, oktaInternalGovernance.AutoRemediationSettingsIncludeOnlyInner{
				ResourceId:   tr.ResourceId.ValueStringPointer(),
				ResourceType: &rt,
			})
		}
		autoRemediationSettings = &oktaInternalGovernance.AutoRemediationSettings{
			IncludeAllIndirectAssignments: d.RemediationSettings.AutoRemediationSettings.IncludeAllIndirectAssignments.ValueBoolPointer(),
			IncludeOnly:                   includeOnlyConverted,
		}

		for _, includeOnly := range d.RemediationSettings.AutoRemediationSettings.IncludeOnly {
			rt := oktaInternalGovernance.AutoRemediationResourceType(includeOnly.ResourceType.ValueString())
			targetResource := oktaInternalGovernance.AutoRemediationSettingsIncludeOnlyInner{
				ResourceId:   includeOnly.ResourceId.ValueStringPointer(),
				ResourceType: &rt,
			}
			autoRemediationSettings.IncludeOnly = append(autoRemediationSettings.IncludeOnly, targetResource)
		}
	}

	var reviewerLevels []oktaInternalGovernance.ReviewerLevelSettingsMutable
	for i, level := range d.ReviewerSettings.ReviewerLevels {
		reviewerLevel := oktaInternalGovernance.ReviewerLevelSettingsMutable{
			Type:                    oktaInternalGovernance.ReviewerType(level.Type.ValueString()),
			FallBackReviewerId:      level.FallBackReviewerId.ValueStringPointer(),
			ReviewerGroupId:         level.ReviewerGroupId.ValueStringPointer(),
			ReviewerId:              level.ReviewerId.ValueStringPointer(),
			ReviewerScopeExpression: level.ReviewerScopeExpression.ValueStringPointer(),
			SelfReviewDisabled:      level.SelfReviewDisabled.ValueBoolPointer(),
			StartReview:             oktaInternalGovernance.ReviewerLevelStartReview{OnDay: level.StartReview[i].OnDay.ValueInt32(), When: (*oktaInternalGovernance.ReviewerLowerLevelCondition)(level.StartReview[i].When.ValueStringPointer())},
		}
		reviewerLevels = append(reviewerLevels, reviewerLevel)
	}

	//if d.PrincipalScope.ExcludedUserIds != nil && len(d.PrincipalScope.ExcludedUserIds) > 0 {
	//	var excluded make([]string,0)

	// IDs related to principal scope settings
	var (
		excludedUserIDs       []string // Users to be excluded from the campaign
		groupIDs              []string // Groups included in the campaign scope
		principalScopeUserIDs []string // Users explicitly included in the campaign scope
	)
	_ = d.PrincipalScope.ExcludedUserIds.ElementsAs(context.Background(), &excludedUserIDs, false)
	_ = d.PrincipalScope.GroupIds.ElementsAs(context.Background(), &groupIDs, false)
	_ = d.PrincipalScope.GroupIds.ElementsAs(context.Background(), &principalScopeUserIDs, false)

	// Build and return CampaignMutable
	return oktaInternalGovernance.CampaignMutable{
		Name:         d.Name.ValueString(),
		CampaignTier: (*oktaInternalGovernance.CampaignTier)(d.CampaignTier.ValueStringPointer()),
		CampaignType: (*oktaInternalGovernance.CampaignType)(d.CampaignType.ValueStringPointer()),
		Description:  d.Description.ValueStringPointer(),
		RemediationSettings: oktaInternalGovernance.RemediationSettings{
			AccessApproved:          oktaInternalGovernance.ApprovedRemediationAction(d.RemediationSettings.AccessApproved.ValueString()),
			AccessRevoked:           oktaInternalGovernance.RevokedRemediationAction(d.RemediationSettings.AccessRevoked.ValueString()),
			NoResponse:              oktaInternalGovernance.NoResponseRemediationAction(d.RemediationSettings.NoResponse.ValueString()),
			AutoRemediationSettings: autoRemediationSettings,
		},
		ResourceSettings: oktaInternalGovernance.ResourceSettingsMutable{
			Type:                               oktaInternalGovernance.CampaignResourceType(d.ResourceSettings.Type.ValueString()),
			TargetResources:                    targetResources,
			IncludeAdminRoles:                  d.ResourceSettings.IncludeAdminRoles.ValueBoolPointer(),
			IncludeEntitlements:                d.ResourceSettings.IncludeEntitlements.ValueBoolPointer(),
			IndividuallyAssignedAppsOnly:       d.ResourceSettings.IndividuallyAssignedAppsOnly.ValueBoolPointer(),
			IndividuallyAssignedGroupsOnly:     d.ResourceSettings.IndividuallyAssignedGroupsOnly.ValueBoolPointer(),
			OnlyIncludeOutOfPolicyEntitlements: d.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements.ValueBoolPointer(),
			ExcludedResources:                  ExcludedResources,
		},
		ReviewerSettings: oktaInternalGovernance.ReviewerSettingsMutable{
			Type:                    oktaInternalGovernance.CampaignReviewerType(d.ReviewerSettings.Type.ValueString()),
			ReviewerGroupId:         d.ReviewerSettings.ReviewerGroupId.ValueStringPointer(),
			ReviewerId:              d.ReviewerSettings.ReviewerId.ValueStringPointer(),
			ReviewerScopeExpression: d.ReviewerSettings.ReviewerScopeExpression.ValueStringPointer(),
			FallBackReviewerId:      d.ReviewerSettings.FallbackReviewerId.ValueStringPointer(),
			JustificationRequired:   d.ReviewerSettings.JustificationRequired.ValueBoolPointer(),
			SelfReviewDisabled:      d.ReviewerSettings.SelfReviewDisabled.ValueBoolPointer(),
			ReassignmentDisabled:    d.ReviewerSettings.ReassignmentDisabled.ValueBoolPointer(),
			BulkDecisionDisabled:    d.ReviewerSettings.BulkDecisionDisabled.ValueBoolPointer(),
			ReviewerLevels:          reviewerLevels,
		},
		ScheduleSettings: oktaInternalGovernance.ScheduleSettingsMutable{
			StartDate:      parsedStartDate,
			DurationInDays: float32(d.ScheduleSettings.DurationInDays.ValueInt32()),
			TimeZone:       d.ScheduleSettings.TimeZone.ValueString(),
			Type:           oktaInternalGovernance.ScheduleType(d.ScheduleSettings.Type.ValueString()),
			Recurrence:     &recur,
		},
		NotificationSettings: &oktaInternalGovernance.NotificationSettings{
			NotifyReviewerAtCampaignEnd:                d.NotificationSettings.NotifyReviewerAtCampaignEnd.ValueBoolPointer(),
			NotifyReviewerDuringMidpointOfReview:       *toNullableBool(d.NotificationSettings.NotifyReviewerDuringMidpointOfReview.ValueBoolPointer()),
			NotifyReviewerWhenOverdue:                  *toNullableBool(d.NotificationSettings.NotifyReviewerWhenOverdue.ValueBoolPointer()),
			NotifyReviewerWhenReviewAssigned:           d.NotificationSettings.NotifyReviewerWhenReviewAssigned.ValueBoolPointer(),
			NotifyReviewPeriodEnd:                      *toNullableBool(d.NotificationSettings.NotifyReviewPeriodEnd.ValueBoolPointer()),
			RemindersReviewerBeforeCampaignCloseInSecs: remindersReviewerBeforeCampaignCloseInSecs,
		},
		PrincipalScopeSettings: &oktaInternalGovernance.PrincipalScopeSettingsMutable{
			Type:                             oktaInternalGovernance.PrincipalScopeType(d.PrincipalScope.Type.ValueString()),
			ExcludedUserIds:                  excludedUserIDs,
			GroupIds:                         groupIDs,
			IncludeOnlyActiveUsers:           d.PrincipalScope.IncludeOnlyActiveUsers.ValueBoolPointer(),
			OnlyIncludeUsersWithSODConflicts: d.PrincipalScope.OnlyIncludeUsersWithSODConflicts.ValueBoolPointer(),
			UserIds:                          principalScopeUserIDs,
			UserScopeExpression:              d.PrincipalScope.UserScopeExpression.ValueStringPointer(),
		},
	}
}

func isOnlyLaunchChanged(plan, state campaignResourceModel) bool {
	return plan.Name.Equal(state.Name) &&
		equalRemediation(plan.RemediationSettings, state.RemediationSettings) &&
		equalResourceSettings(plan.ResourceSettings, state.ResourceSettings) &&
		equalReviewerSettings(plan.ReviewerSettings, state.ReviewerSettings) &&
		equalScheduleSettings(plan.ScheduleSettings, state.ScheduleSettings)
}

func equalScheduleSettings(planSchedule, stateSchedule *scheduleSettingsModel) bool {
	return planSchedule.Type.Equal(stateSchedule.Type) &&
		planSchedule.TimeZone.Equal(stateSchedule.TimeZone) &&
		planSchedule.DurationInDays.Equal(stateSchedule.DurationInDays) &&
		planSchedule.StartDate.Equal(stateSchedule.StartDate)
}

func equalReviewerSettings(planReviewer, stateReviewer *reviewerSettingsModel) bool {
	return planReviewer.Type.Equal(stateReviewer.Type) &&
		planReviewer.ReviewerGroupId.Equal(stateReviewer.ReviewerGroupId) &&
		planReviewer.ReviewerId.Equal(stateReviewer.ReviewerId) &&
		planReviewer.ReviewerScopeExpression.Equal(stateReviewer.ReviewerScopeExpression) &&
		planReviewer.FallbackReviewerId.Equal(stateReviewer.FallbackReviewerId)
}

func equalResourceSettings(planResource, stateResource *resourceSettingsModel) bool {
	if len(planResource.TargetResources) == len(stateResource.TargetResources) {
		for i := 0; i < len(planResource.TargetResources); i++ {
			if !(planResource.TargetResources[i].ResourceId.Equal(stateResource.TargetResources[i].ResourceId) &&
				planResource.TargetResources[i].ResourceType.Equal(stateResource.TargetResources[i].ResourceType)) {
				return false
			}
		}
		return true
	}
	return false
}

func equalRemediation(a, b *campaignRemediationSettingsModel) bool {
	return a.AccessApproved.Equal(b.AccessApproved) &&
		a.AccessRevoked.Equal(b.AccessRevoked) &&
		a.NoResponse.Equal(b.NoResponse)
}

func toNullableBool(v *bool) *oktaInternalGovernance.NullableBool {
	if v == nil {
		return nil
	}
	return oktaInternalGovernance.NewNullableBool(v)
}
