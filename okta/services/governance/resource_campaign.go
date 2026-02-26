package governance

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
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
	SkipRemediation      types.Bool                        `tfsdk:"skip_remediation"`
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
	Type                               types.String            `tfsdk:"type"`
	IncludeAdminRoles                  types.Bool              `tfsdk:"include_admin_roles"`
	IncludeEntitlements                types.Bool              `tfsdk:"include_entitlements"`
	IndividuallyAssignedAppsOnly       types.Bool              `tfsdk:"individually_assigned_apps_only"`
	IndividuallyAssignedGroupsOnly     types.Bool              `tfsdk:"individually_assigned_groups_only"`
	OnlyIncludeOutOfPolicyEntitlements types.Bool              `tfsdk:"only_include_out_of_policy_entitlements"`
	ExcludedResources                  []excludedResourceModel `tfsdk:"excluded_resources"`
	TargetResources                    []targetResourceModel   `tfsdk:"target_resources"`
}

type excludedResourceModel struct {
	ResourceId   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the campaign. Maintain some uniqueness when naming the campaign as it helps to identify and filter for campaigns when needed.",
			},
			"campaign_tier": schema.StringAttribute{
				Optional:    true,
				Description: "Indicates the minimum required SKU to manage the campaign. Values can be `BASIC` and `PREMIUM`.",
			},
			"campaign_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Identifies if it is a resource campaign or a user campaign. By default it is RESOURCE.Values can be `RESOURCE` and `USER`.",
				Validators: []validator.String{
					stringvalidator.OneOf("RESOURCE", "USER"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description about the campaign.",
			},
			"skip_remediation": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, skip remediation when ending the campaign (only applicable if remediationSetting.noResponse=DENY).",
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
						Required:    true,
						Description: "Specifies the action if the reviewer doesn't respond to the request or if the campaign is closed before an action is taken.",
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
						Required:    true,
						Description: "The type of Okta resource.",
						Validators: []validator.String{
							stringvalidator.OneOf("APPLICATION", "APPLICATION_AND_GROUP", "GROUP"),
						},
					},
					"include_admin_roles": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Include admin roles.",
						Default:     booldefault.StaticBool(false),
					},
					"include_entitlements": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Include entitlements for this application. This property is only applicable if resource_type = APPLICATION and Entitlement Management is enabled.",
					},
					"individually_assigned_apps_only": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Only include individually assigned apps. This is only applicable if campaign type is USER.",
					},
					"individually_assigned_groups_only": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Only include individually assigned groups. This is only applicable if campaign type is USER.",
					},
					"only_include_out_of_policy_entitlements": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Only include out-of-policy entitlements. Only applicable if resource_type = APPLICATION and Entitlement Management is enabled.",
					},
				},
				Blocks: map[string]schema.Block{
					"excluded_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Optional:    true,
									Description: "The ID of the resource to exclude in the campaign.",
								},
								"resource_type": schema.StringAttribute{
									Optional:    true,
									Description: "The type of resource to exclude in the campaign.",
									Validators: []validator.String{
										stringvalidator.OneOf("APPLICATION", "GROUP"),
									},
								},
							},
						},
						Description: "An array of resources that are excluded from the review.",
					},
					"target_resources": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.StringAttribute{
									Required:    true,
									Description: "The resource ID that is being reviewed.",
								},
								"resource_type": schema.StringAttribute{
									Required:    true,
									Description: "The type of Okta resource.",
									Validators: []validator.String{
										stringvalidator.OneOf("APPLICATION", "GROUP"),
									},
								},
								"include_all_entitlements_and_bundles": schema.BoolAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Include all entitlements and entitlement bundles for this application. Only applicable if the resourcetype = APPLICATION and Entitlement Management is enabled.",
								},
							},
							Blocks: map[string]schema.Block{
								"entitlement_bundles": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required:    true,
												Description: "The ID of the entitlement bundle.",
											},
										},
									},
									Description: "An array of entitlement bundles for this application.",
								},
								"entitlements": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Required:    true,
												Description: "The entitlement id.",
											},
											"include_all_values": schema.BoolAttribute{
												Optional:    true,
												Description: "Whether to include all entitlement values. If false we must provide the values property.",
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
											},
										},
									},
									Description: "An array of entitlements associated with resourceId that should be chosen as target when creating reviews",
								},
							},
						},
						Description: "Represents a resource that will be part of Access certifications. If the app is enabled for Access Certifications, it's possible to review entitlements and entitlement bundles.",
					},
				},
				Description: "Resource specific properties.",
			},
			"reviewer_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Identifies the kind of reviewer for Access Certification.",
						Validators: []validator.String{
							stringvalidator.OneOf("GROUP", "MULTI_LEVEL", "RESOURCE_OWNER", "REVIEWER_EXPRESSION", "USER"),
						},
					},
					"bulk_decision_disabled": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "When approving or revoking review items, bulk actions are disabled if true.",
					},
					"fallback_reviewer_id": schema.StringAttribute{
						Optional:    true,
						Description: "The ID of the fallback reviewer. Required when the type=`REVIEWER_EXPRESSION` or type=`RESOURCE_OWNER`",
					},
					"justification_required": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "When approving or revoking review items, a justification is required if true.",
					},
					"reassignment_disabled": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Reassignment is disabled for reviewers if true.",
					},
					"self_review_disabled": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "This property is required to be true for resource-centric campaigns when the Okta Admin Console is one of the resources.",
					},
					"reviewer_group_id": schema.StringAttribute{
						Optional:    true,
						Description: "The ID of the reviewer group to which the reviewer is assigned.",
					},
					"reviewer_id": schema.StringAttribute{
						Optional: true,
					},
					"reviewer_scope_expression": schema.StringAttribute{
						Optional:    true,
						Description: "This property is required when type=`USER`",
					},
				},
				Blocks: map[string]schema.Block{
					"reviewer_levels": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "Identifies the kind of reviewer.",
									Validators: []validator.String{
										stringvalidator.OneOf("GROUP", "RESOURCE_OWNER", "REVIEWER_EXPRESSION", "USER"),
									},
								},
								"fallback_reviewer_id": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "Required when the type=`REVIEWER_EXPRESSION` or type=`RESOURCE_OWNER`",
								},
								"reviewer_group_id": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "The ID of the reviewer group to which the reviewer is assigned.This property is required when type=`GROUP`",
								},
								"reviewer_id": schema.StringAttribute{
									Optional:    true,
									Description: "The ID of the reviewer to which the reviewer is assigned.This property is required when type=`USER`.",
								},
								"reviewer_scope_expression": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "This property is required when type=`REVIEWER_EXPRESSION`",
								},
								"self_review_disabled": schema.BoolAttribute{
									Computed:    true,
									Optional:    true,
									Description: "This property is used to prevent self review.",
								},
							},
							Blocks: map[string]schema.Block{
								"start_review": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"on_day": schema.Int32Attribute{
												Computed:    true,
												Optional:    true,
												Default:     int32default.StaticInt32(0),
												Description: "The day of the campaign when the review starts. 0 means the first day of the campaign.",
											},
											"when": schema.StringAttribute{
												Computed:    true,
												Optional:    true,
												Description: "The condition for which, the lower level reviews will move to that level for further review.",
											},
										},
									},
									Validators: []validator.List{
										listvalidator.SizeAtMost(2),
									},
									Description: "The rules for which the reviews can move to that level.",
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtMost(2),
						},
						Description: "Definition of reviewer level for a given campaign. Each reviewer level defines the kind of reviewer who is going to review.",
					},
				},
				Description: "Identifies the kind of reviewer for Access Certification.",
			},
			"schedule_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"start_date": schema.StringAttribute{
						Required:    true,
						Description: "The date on which the campaign is supposed to start. Accepts date in ISO 8601 format.",
					},
					"duration_in_days": schema.Int32Attribute{
						Required:    true,
						Description: "The duration (in days) that the campaign is active.",
					},
					"time_zone": schema.StringAttribute{
						Required:    true,
						Description: "The time zone in which the campaign is active.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of campaign being scheduled.",
						Validators: []validator.String{
							stringvalidator.OneOf("ONE_OFF", "RECURRING"),
						},
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
									Required:    true,
									Description: "Recurrence interval specified according to ISO8061 notation for durations.",
								},
								"ends": schema.StringAttribute{
									Optional:    true,
									Description: "Specifies when the recurring schedule can have an end.",
								},
								"repeat_on_type": schema.StringAttribute{
									Optional:    true,
									Description: "Specifies when the recurring schedule can have an end.",
									Validators: []validator.String{
										stringvalidator.OneOf("LAST_WEEKDAY_AS_START_DATE", "SAME_DAY_AS_START_DATE", "SAME_WEEKDAY_AS_START_DATE"),
									},
								},
							},
						},
					},
				},
				Description: "Scheduler specific settings.",
			},
			"notification_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"notify_reviewer_at_campaign_end": schema.BoolAttribute{
						Required:    true,
						Description: "To indicate whether a notification should be sent to the reviewers when campaign has come to an end.",
					},
					"notify_reviewer_during_midpoint_of_review": schema.BoolAttribute{
						Required:    true,
						Description: "To indicate whether a notification should be sent to the reviewer during the midpoint of the review process.",
					},
					"notify_reviewer_when_overdue": schema.BoolAttribute{
						Required:    true,
						Description: "To indicate whether a notification should be sent to the reviewer when the review is overdue.",
					},
					"notify_reviewer_when_review_assigned": schema.BoolAttribute{
						Required:    true,
						Description: "To indicate whether a notification should be sent to the reviewer when actionable reviews are assigned.",
					},
					"notify_review_period_end": schema.BoolAttribute{
						Required:    true,
						Description: "To indicate whether a notification should be sent to the reviewer when a given reviewer level period is about to end.",
					},
					"reminders_reviewer_before_campaign_close_in_secs": schema.ListAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.Int64Type,
						Description: "Specifies times (in seconds) to send reminders to reviewers before the campaign closes. Max 3 values. Example: [86400, 172800, 604800]",
					},
				},
			},
			"principal_scope_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Specifies the type for principal_scope_settings.",
					},
					"excluded_user_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "An array of Okta user IDs excluded from access certification or the campaign. This field is optional. A maximum of 50 users can be specified in the array.",
					},
					"group_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "An array of Okta group IDs included from access certification or the campaign. userIds, groupIds or userScopeExpression is required if campaign type is USER. A maximum of 5 groups can be specified in the array.",
					},
					"include_only_active_users": schema.BoolAttribute{
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Description: "If set to true, only active Okta users are included in the campaign.",
					},
					"only_include_users_with_sod_conflicts": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "If set to true, only includes users that have at least one SOD conflict that was caused due to entitlement(s) within Campaign scope.",
					},
					"user_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "An array of Okta user IDs included from access certification or the campaign. userIds, groupIds or userScopeExpression is required if campaign type is USER. A maximum of 100 users can be specified in the array.",
					},
					"user_scope_expression": schema.StringAttribute{
						Optional:    true,
						Description: "The Okta expression language user expression on the resourceSettings to include users in the campaign.",
					},
				},
				Blocks: map[string]schema.Block{
					"predefined_inactive_users_scope": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"inactive_days": schema.Int32Attribute{
									Optional:    true,
									Description: "The duration the users have not used single sign on (SSO) to access their account within the specific time frame. Minimum 30 days and maximum 365 days are supported.",
									Validators: []validator.Int32{
										int32validator.AtLeast(30),
										int32validator.AtMost(365),
									},
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

	if resp.Diagnostics.HasError() {
		return
	}

	campaign, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CampaignsAPI.CreateCampaign(ctx).CampaignMutable(buildCampaign(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Campaign",
			"Could not create Campaign, unexpected error: "+err.Error(),
		)
		return
	}
	data.Id = types.StringValue(campaign.Id)

	resp.Diagnostics.Append(applyCampaignsToState(ctx, campaign, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	getCampaignResponse, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CampaignsAPI.GetCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading campaign",
			"Could not read Campaign, unexpected error: "+err.Error(),
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

func (r *campaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data campaignResourceModel
	var state campaignResourceModel

	// Load both planned and current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
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
	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().CampaignsAPI.DeleteCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Campaign",
			"Could not delete Campaign with ID"+data.Id.ValueString()+" unexpected error: "+err.Error(),
		)
		return
	}
}

func applyCampaignsToState(ctx context.Context, resp *governance.CampaignFull, c *campaignResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	c.Id = types.StringValue(resp.GetId())
	c.Name = types.StringValue(resp.GetName())
	c.CampaignType = types.StringValue(string(resp.GetCampaignType()))
	c.Description = types.StringValue(resp.GetDescription())

	c.RemediationSettings = &campaignRemediationSettingsModel{}
	if resp.RemediationSettings.GetAccessRevoked() != "" {
		c.RemediationSettings.AccessRevoked = types.StringValue(string(resp.RemediationSettings.GetAccessRevoked()))
	}
	if resp.RemediationSettings.GetAccessApproved() != "" {
		c.RemediationSettings.AccessApproved = types.StringValue(string(resp.RemediationSettings.GetAccessApproved()))
	}
	if resp.RemediationSettings.GetNoResponse() != "" {
		c.RemediationSettings.NoResponse = types.StringValue(string(resp.RemediationSettings.GetNoResponse()))
	}
	if autoRemediationSettings, ok := resp.RemediationSettings.GetAutoRemediationSettingsOk(); ok {
		c.RemediationSettings.AutoRemediationSettings = &autoRemediationSettingsModel{}
		c.RemediationSettings.AutoRemediationSettings.IncludeAllIndirectAssignments = types.BoolValue(autoRemediationSettings.GetIncludeAllIndirectAssignments())
		for _, includeOnly := range autoRemediationSettings.GetIncludeOnly() {
			targetResource := targetResourceModel{
				ResourceId:   types.StringValue(includeOnly.GetResourceId()),
				ResourceType: types.StringValue(string(includeOnly.GetResourceType())),
			}
			c.RemediationSettings.AutoRemediationSettings.IncludeOnly = append(c.RemediationSettings.AutoRemediationSettings.IncludeOnly, targetResource)
		}
	}

	c.ResourceSettings = &resourceSettingsModel{}
	if resp.ResourceSettings.GetType() != "" {
		c.ResourceSettings.Type = types.StringValue(string(resp.ResourceSettings.GetType()))
	}
	if len(resp.ResourceSettings.GetTargetResources()) > 0 {
		var targets []targetResourceModel
		for _, targetResource := range resp.ResourceSettings.GetTargetResources() {
			target := targetResourceModel{
				ResourceId: types.StringValue(targetResource.GetResourceId()),
			}
			if targetResource.ResourceType != nil {
				target.ResourceType = types.StringValue(string(targetResource.GetResourceType()))
			}
			if targetResource.IncludeAllEntitlementsAndBundles != nil {
				target.IncludeAllEntitlementsAndBundles = types.BoolValue(targetResource.GetIncludeAllEntitlementsAndBundles())
			}
			if targetResource.GetEntitlements() != nil {
				var entitlements []entitlementModel
				for _, entitlement := range targetResource.GetEntitlements() {
					var values []entitlementValueModel
					for _, val := range entitlement.GetValues() {
						v := entitlementValueModel{
							Id: types.StringValue(val.GetId()),
						}
						values = append(values, v)
					}
					e := entitlementModel{
						Id:               types.StringValue(entitlement.GetId()),
						IncludeAllValues: types.BoolValue(entitlement.GetIncludeAllValues()),
						Values:           values,
					}
					entitlements = append(entitlements, e)
				}
				target.Entitlements = entitlements
			}

			if target.EntitlementBundles != nil {
				var entitlementBundles []entitlementBundleModel
				for _, entitlementBundle := range target.EntitlementBundles {
					bundle := entitlementBundleModel{
						Id: entitlementBundle.Id,
					}
					entitlementBundles = append(entitlementBundles, bundle)
				}
				target.EntitlementBundles = entitlementBundles
			}

			targets = append(targets, target)
		}
		c.ResourceSettings.TargetResources = targets
	}
	if len(resp.ResourceSettings.ExcludedResources) > 0 {
		var excluded []excludedResourceModel
		for _, res := range resp.ResourceSettings.ExcludedResources {
			excludedRes := excludedResourceModel{
				ResourceId: types.StringValue(*res.ResourceId),
			}
			if res.ResourceType != nil {
				excludedRes.ResourceType = types.StringValue(string(*res.ResourceType))
			} else {
				excludedRes.ResourceType = types.StringNull()
			}
			excluded = append(excluded, excludedRes)
		}

		c.ResourceSettings.ExcludedResources = excluded
	}
	if resp.ResourceSettings.IncludeAdminRoles != nil {
		c.ResourceSettings.IncludeAdminRoles = types.BoolValue(resp.ResourceSettings.GetIncludeAdminRoles())
	}
	if resp.ResourceSettings.IncludeEntitlements != nil {
		c.ResourceSettings.IncludeEntitlements = types.BoolValue(resp.ResourceSettings.GetIncludeEntitlements())
	}
	if resp.ResourceSettings.IndividuallyAssignedAppsOnly != nil {
		c.ResourceSettings.IndividuallyAssignedAppsOnly = types.BoolValue(resp.ResourceSettings.GetIndividuallyAssignedAppsOnly())
	}
	if resp.ResourceSettings.IndividuallyAssignedGroupsOnly != nil {
		c.ResourceSettings.IndividuallyAssignedGroupsOnly = types.BoolValue(resp.ResourceSettings.GetIndividuallyAssignedGroupsOnly())
	}
	if resp.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements != nil {
		c.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements = types.BoolValue(resp.ResourceSettings.GetOnlyIncludeOutOfPolicyEntitlements())
	}

	c.ReviewerSettings = &reviewerSettingsModel{}
	c.ReviewerSettings.Type = types.StringValue(string(resp.ReviewerSettings.GetType()))
	if resp.ReviewerSettings.BulkDecisionDisabled != nil {
		c.ReviewerSettings.BulkDecisionDisabled = types.BoolValue(resp.ReviewerSettings.GetBulkDecisionDisabled())
	} else {
		c.ReviewerSettings.BulkDecisionDisabled = types.BoolValue(false)
	}
	if resp.ReviewerSettings.FallBackReviewerId != nil {
		c.ReviewerSettings.FallbackReviewerId = types.StringValue(resp.ReviewerSettings.GetFallBackReviewerId())
	}
	if resp.ReviewerSettings.JustificationRequired != nil {
		c.ReviewerSettings.JustificationRequired = types.BoolValue(resp.ReviewerSettings.GetJustificationRequired())
	}
	if resp.ReviewerSettings.ReassignmentDisabled != nil {
		c.ReviewerSettings.ReassignmentDisabled = types.BoolValue(resp.ReviewerSettings.GetReassignmentDisabled())
	} else {
		c.ReviewerSettings.ReassignmentDisabled = types.BoolValue(false)
	}
	if resp.ReviewerSettings.ReviewerGroupId != nil {
		c.ReviewerSettings.ReviewerGroupId = types.StringValue(resp.ReviewerSettings.GetReviewerGroupId())
	}
	if resp.ReviewerSettings.ReviewerId != nil {
		c.ReviewerSettings.ReviewerId = types.StringValue(resp.ReviewerSettings.GetReviewerId())
	}
	if resp.ReviewerSettings.ReviewerScopeExpression != nil {
		c.ReviewerSettings.ReviewerScopeExpression = types.StringValue(resp.ReviewerSettings.GetReviewerScopeExpression())
	}
	if resp.ReviewerSettings.SelfReviewDisabled != nil {
		c.ReviewerSettings.SelfReviewDisabled = types.BoolValue(resp.ReviewerSettings.GetSelfReviewDisabled())
	} else {
		c.ReviewerSettings.SelfReviewDisabled = types.BoolValue(false)
	}
	if resp.ReviewerSettings.ReviewerLevels != nil {
		c.ReviewerSettings.ReviewerLevels = make([]reviewerLevelModel, 0, len(resp.ReviewerSettings.GetReviewerLevels()))
		for _, level := range resp.ReviewerSettings.ReviewerLevels {
			reviewerLevel := reviewerLevelModel{}
			reviewerLevel.Type = types.StringValue(string(level.GetType()))
			fallbackReviewerId := level.GetFallBackReviewerId()
			if fallbackReviewerId != "" {
				reviewerLevel.FallBackReviewerId = types.StringValue(fallbackReviewerId)
			}
			reviewerGroupId := level.GetReviewerGroupId()
			if reviewerGroupId != "" {
				reviewerLevel.ReviewerGroupId = types.StringValue(reviewerGroupId)
			}
			reviewerId := level.GetReviewerId()
			if reviewerId != "" {
				reviewerLevel.ReviewerId = types.StringValue(reviewerId)
			}
			reviewerScopeExpression := level.GetReviewerScopeExpression()
			if reviewerScopeExpression != "" {
				reviewerLevel.ReviewerScopeExpression = types.StringValue(reviewerScopeExpression)
			}
			if level.SelfReviewDisabled != nil {
				reviewerLevel.SelfReviewDisabled = types.BoolValue(level.GetSelfReviewDisabled())
			}

			startReviews := make([]startReviewModel, 1)
			startReviews[0].OnDay = types.Int32Value(level.StartReview.GetOnDay())
			if level.StartReview.When != nil {
				startReviews[0].When = types.StringValue(string(level.StartReview.GetWhen()))
			}
			reviewerLevel.StartReview = startReviews

			c.ReviewerSettings.ReviewerLevels = append(c.ReviewerSettings.ReviewerLevels, reviewerLevel)

		}
	}

	c.ScheduleSettings = &scheduleSettingsModel{}
	c.ScheduleSettings.StartDate = types.StringValue(resp.ScheduleSettings.GetStartDate().UTC().Format("2006-01-02T15:04:05.000Z"))
	c.ScheduleSettings.DurationInDays = types.Int32Value(int32(resp.ScheduleSettings.GetDurationInDays()))
	c.ScheduleSettings.TimeZone = types.StringValue(resp.ScheduleSettings.GetTimeZone())
	c.ScheduleSettings.Type = types.StringValue(string(resp.ScheduleSettings.Type))
	c.ScheduleSettings.DurationInDays = types.Int32Value(int32(resp.ScheduleSettings.GetDurationInDays()))
	if resp.ScheduleSettings.Recurrence != nil {
		c.ScheduleSettings.Recurrence = make([]recurrenceModel, 0)
		rec := getRecurrence(resp)
		c.ScheduleSettings.Recurrence = append(c.ScheduleSettings.Recurrence, rec)
	}

	c.NotificationSettings = &notificationSettingsModel{}
	if resp.NotificationSettings != nil {
		if resp.NotificationSettings.NotifyReviewerAtCampaignEnd != nil {
			c.NotificationSettings.NotifyReviewerAtCampaignEnd = types.BoolValue(resp.NotificationSettings.GetNotifyReviewerAtCampaignEnd())
		}
		if resp.NotificationSettings.NotifyReviewerDuringMidpointOfReview.Get() != nil {
			c.NotificationSettings.NotifyReviewerDuringMidpointOfReview = types.BoolValue(resp.NotificationSettings.GetNotifyReviewerDuringMidpointOfReview())
		}
		if resp.NotificationSettings.NotifyReviewerWhenOverdue.Get() != nil {
			c.NotificationSettings.NotifyReviewerWhenOverdue = types.BoolValue(resp.NotificationSettings.GetNotifyReviewerWhenOverdue())
		}
		if resp.NotificationSettings.NotifyReviewerWhenReviewAssigned != nil {
			c.NotificationSettings.NotifyReviewerWhenReviewAssigned = types.BoolValue(resp.NotificationSettings.GetNotifyReviewerWhenReviewAssigned())
		}
		if resp.NotificationSettings.NotifyReviewPeriodEnd.Get() != nil {
			c.NotificationSettings.NotifyReviewPeriodEnd = types.BoolValue(resp.NotificationSettings.GetNotifyReviewPeriodEnd())
		}
		if len(resp.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs) > 0 {
			reminders := make([]int64, 0, len(resp.NotificationSettings.GetRemindersReviewerBeforeCampaignCloseInSecs()))
			for _, v := range resp.NotificationSettings.GetRemindersReviewerBeforeCampaignCloseInSecs() {
				reminders = append(reminders, int64(v))
			}

			listVal, _ := types.ListValueFrom(ctx, types.Int64Type, reminders)

			c.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs = listVal
		} else {
			// explicitly set an empty list with correct type
			c.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs = types.ListNull(types.Int64Type)
		}
	}

	c.PrincipalScope = &principalScopeSettingsModel{}
	if resp.PrincipalScopeSettings.GetType() != "" {
		c.PrincipalScope.Type = types.StringValue(string(resp.PrincipalScopeSettings.GetType()))
	}
	if len(resp.PrincipalScopeSettings.GetExcludedUserIds()) > 0 {
		excluded := make([]attr.Value, 0, len(resp.PrincipalScopeSettings.GetExcludedUserIds()))
		for _, id := range resp.PrincipalScopeSettings.GetExcludedUserIds() {
			excluded = append(excluded, types.StringValue(id))
		}
		c.PrincipalScope.ExcludedUserIds, _ = types.ListValue(types.StringType, excluded)
	} else {
		c.PrincipalScope.ExcludedUserIds = types.ListNull(types.StringType)
	}
	if len(resp.PrincipalScopeSettings.GetGroupIds()) > 0 {
		groupIds := make([]attr.Value, 0, len(resp.PrincipalScopeSettings.GetGroupIds()))
		for _, id := range resp.PrincipalScopeSettings.GetGroupIds() {
			groupIds = append(groupIds, types.StringValue(id))
		}
		c.PrincipalScope.GroupIds = types.ListValueMust(types.StringType, groupIds)
	} else {
		c.PrincipalScope.GroupIds = types.ListNull(types.StringType)
	}
	if resp.PrincipalScopeSettings.IncludeOnlyActiveUsers != nil {
		c.PrincipalScope.IncludeOnlyActiveUsers = types.BoolValue(resp.PrincipalScopeSettings.GetIncludeOnlyActiveUsers())
	}
	if resp.PrincipalScopeSettings.OnlyIncludeUsersWithSODConflicts != nil {
		c.PrincipalScope.OnlyIncludeUsersWithSODConflicts = types.BoolValue(resp.PrincipalScopeSettings.GetOnlyIncludeUsersWithSODConflicts())
	}
	if len(resp.PrincipalScopeSettings.GetUserIds()) > 0 {
		listVal, diags := types.ListValueFrom(ctx, types.StringType, resp.PrincipalScopeSettings.GetUserIds())
		if diags.HasError() {
			diags.Append(diags...)
			return diags
		}
		c.PrincipalScope.UserIds = listVal
	} else {
		c.PrincipalScope.UserIds = types.ListNull(types.StringType)
	}
	if resp.PrincipalScopeSettings.UserScopeExpression != nil {
		c.PrincipalScope.UserScopeExpression = types.StringValue(resp.PrincipalScopeSettings.GetUserScopeExpression())
	}
	if resp.PrincipalScopeSettings.PredefinedInactiveUsersScope != nil {
		c.PrincipalScope.PredefinedInactiveUsersScope = []inactiveUsersScopeModel{
			{
				InactiveDays: types.Int32Value(resp.PrincipalScopeSettings.PredefinedInactiveUsersScope.GetInactiveDays()),
			},
		}
	}

	return diags
}

func getRecurrence(resp *governance.CampaignFull) recurrenceModel {
	recurrence := recurrenceModel{}
	recurrence.Interval = types.StringValue(resp.ScheduleSettings.Recurrence.Interval)
	if resp.ScheduleSettings.Recurrence.Ends != nil && !resp.ScheduleSettings.Recurrence.Ends.IsZero() {
		recurrence.Ends = types.StringValue(resp.ScheduleSettings.Recurrence.Ends.UTC().Format("2006-01-02T15:04:05.000Z"))
	}

	if resp.ScheduleSettings.Recurrence.RepeatOnType != nil {
		recurrence.RepeatOnType = types.StringValue(string(*resp.ScheduleSettings.Recurrence.RepeatOnType))
	}
	return recurrence
}

func buildCampaign(d campaignResourceModel) governance.CampaignMutable {
	startDate := d.ScheduleSettings.StartDate.ValueString()
	parsedStartDate, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		log.Printf("invalid start_date format: %v", err)
		return governance.CampaignMutable{}
	}

	// Convert target resources
	targetResources := buildTargetResources(d)

	ExcludedResources := make([]governance.ResourceSettingsMutableExcludedResourcesInner, 0, len(d.ResourceSettings.ExcludedResources))
	for _, ex := range d.ResourceSettings.ExcludedResources {
		x := ex.ResourceId.ValueString()
		var resourceType *governance.ResourceType
		if !ex.ResourceType.IsNull() && ex.ResourceType.ValueString() != "" {
			rt := governance.ResourceType(ex.ResourceType.ValueString())
			resourceType = &rt
		}
		excludedRes := governance.ResourceSettingsMutableExcludedResourcesInner{
			ResourceId:   &x,
			ResourceType: resourceType,
		}
		ExcludedResources = append(ExcludedResources, excludedRes)
	}

	var recur governance.RecurrenceDefinitionMutable
	r := d.ScheduleSettings.Recurrence
	if len(r) != 0 {
		endStr := r[0].Ends.ValueString()
		parsedTime, _ := time.Parse(time.RFC3339, endStr)
		if !parsedTime.IsZero() {
			recur.Ends = &parsedTime
		}
		recur.Interval = r[0].Interval.ValueString()
		repeatStr := governance.RecurrenceRepeatOnType(r[0].RepeatOnType.ValueString())
		if repeatStr != "" {
			recur.RepeatOnType = &repeatStr
		}
	}
	var remindersReviewerBeforeCampaignCloseInSecs []int32

	if d.NotificationSettings != nil && !d.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs.IsNull() && !d.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs.IsUnknown() {
		err := d.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs.ElementsAs(
			context.Background(),
			&remindersReviewerBeforeCampaignCloseInSecs,
			false,
		)
		if err != nil {
			log.Printf("Error remindersReviewerBeforeCampaignCloseInSecs: %v", err)
		}
	}
	var autoRemediationSettings *governance.AutoRemediationSettings
	if d.RemediationSettings.AutoRemediationSettings != nil {
		var includeOnlyConverted []governance.AutoRemediationSettingsIncludeOnlyInner

		for _, tr := range d.RemediationSettings.AutoRemediationSettings.IncludeOnly {
			rt := governance.AutoRemediationResourceType(tr.ResourceType.ValueString())
			includeOnlyConverted = append(includeOnlyConverted, governance.AutoRemediationSettingsIncludeOnlyInner{
				ResourceId:   tr.ResourceId.ValueStringPointer(),
				ResourceType: &rt,
			})
		}
		autoRemediationSettings = &governance.AutoRemediationSettings{
			IncludeAllIndirectAssignments: d.RemediationSettings.AutoRemediationSettings.IncludeAllIndirectAssignments.ValueBoolPointer(),
			IncludeOnly:                   includeOnlyConverted,
		}

		for _, includeOnly := range d.RemediationSettings.AutoRemediationSettings.IncludeOnly {
			rt := governance.AutoRemediationResourceType(includeOnly.ResourceType.ValueString())
			targetResource := governance.AutoRemediationSettingsIncludeOnlyInner{
				ResourceId:   includeOnly.ResourceId.ValueStringPointer(),
				ResourceType: &rt,
			}
			autoRemediationSettings.IncludeOnly = append(autoRemediationSettings.IncludeOnly, targetResource)
		}
	}

	var reviewerLevels []governance.ReviewerLevelSettingsMutable
	for _, level := range d.ReviewerSettings.ReviewerLevels {
		var startReview governance.ReviewerLevelStartReview
		if len(level.StartReview) > 0 {
			start := level.StartReview[0]
			startReview = governance.ReviewerLevelStartReview{
				OnDay: start.OnDay.ValueInt32(),
			}
			if !start.When.IsNull() && start.When.ValueString() != "" {
				when := start.When.ValueString()
				startReview.When = (*governance.ReviewerLowerLevelCondition)(&when)
			}
		}

		var reviewerLevel governance.ReviewerLevelSettingsMutable
		reviewerGroupId := level.ReviewerGroupId.ValueStringPointer()
		reviewerId := level.ReviewerId.ValueString()
		reviewerScopeExpression := level.ReviewerScopeExpression.ValueStringPointer()
		fallBackReviewerId := level.FallBackReviewerId.ValueStringPointer()
		selfReviewDisabled := level.SelfReviewDisabled.ValueBoolPointer()

		reviewerLevel.SetType(governance.ReviewerType(level.Type.ValueString()))
		if fallBackReviewerId != nil && *fallBackReviewerId != "" {
			reviewerLevel.FallBackReviewerId = fallBackReviewerId
		}
		if reviewerGroupId != nil && *reviewerGroupId != "" {
			reviewerLevel.ReviewerGroupId = reviewerGroupId
		}
		if reviewerId != "" {
			reviewerLevel.SetReviewerId(reviewerId)
		}
		if reviewerScopeExpression != nil && *reviewerScopeExpression != "" {
			reviewerLevel.ReviewerScopeExpression = reviewerScopeExpression
		}
		if selfReviewDisabled != nil {
			reviewerLevel.SelfReviewDisabled = selfReviewDisabled
		}
		reviewerLevel.StartReview = startReview

		reviewerLevels = append(reviewerLevels, reviewerLevel)
	}

	var (
		excludedUserIDs []string
		groupIDs        []string
		userIds         []string
	)

	_ = d.PrincipalScope.ExcludedUserIds.ElementsAs(context.Background(), &excludedUserIDs, false)
	_ = d.PrincipalScope.GroupIds.ElementsAs(context.Background(), &groupIDs, false)
	_ = d.PrincipalScope.UserIds.ElementsAs(context.Background(), &userIds, false)
	var campaignType *governance.CampaignType
	if v := d.CampaignType.ValueString(); v != "" {
		cType := governance.CampaignType(v)
		campaignType = &cType
	}

	var campaignTier *governance.CampaignTier
	if v := d.CampaignTier.ValueString(); v != "" {
		tier := governance.CampaignTier(v)
		campaignTier = &tier
	}
	predefinedInactiveUsersScope := buildPredefinedUserScope(d.PrincipalScope.PredefinedInactiveUsersScope)
	return governance.CampaignMutable{
		Name:         d.Name.ValueString(),
		CampaignTier: campaignTier,
		CampaignType: campaignType,
		Description:  d.Description.ValueStringPointer(),
		RemediationSettings: governance.RemediationSettings{
			AccessApproved:          governance.ApprovedRemediationAction(d.RemediationSettings.AccessApproved.ValueString()),
			AccessRevoked:           governance.RevokedRemediationAction(d.RemediationSettings.AccessRevoked.ValueString()),
			NoResponse:              governance.NoResponseRemediationAction(d.RemediationSettings.NoResponse.ValueString()),
			AutoRemediationSettings: autoRemediationSettings,
		},
		ResourceSettings: governance.ResourceSettingsMutable{
			Type:                               governance.CampaignResourceType(d.ResourceSettings.Type.ValueString()),
			TargetResources:                    targetResources,
			IncludeAdminRoles:                  d.ResourceSettings.IncludeAdminRoles.ValueBoolPointer(),
			IncludeEntitlements:                d.ResourceSettings.IncludeEntitlements.ValueBoolPointer(),
			IndividuallyAssignedAppsOnly:       d.ResourceSettings.IndividuallyAssignedAppsOnly.ValueBoolPointer(),
			IndividuallyAssignedGroupsOnly:     d.ResourceSettings.IndividuallyAssignedGroupsOnly.ValueBoolPointer(),
			OnlyIncludeOutOfPolicyEntitlements: d.ResourceSettings.OnlyIncludeOutOfPolicyEntitlements.ValueBoolPointer(),
			ExcludedResources:                  ExcludedResources,
		},
		ReviewerSettings:       *getReviewerSettingForRequests(d, reviewerLevels),
		ScheduleSettings:       *getScheduleSettingForRequests(d, parsedStartDate, recur),
		NotificationSettings:   getNotificationSettingsForRequest(d, remindersReviewerBeforeCampaignCloseInSecs),
		PrincipalScopeSettings: buildPrincipalScopeSettings(d, excludedUserIDs, groupIDs, userIds, predefinedInactiveUsersScope),
	}
}

func buildPrincipalScopeSettings(d campaignResourceModel, excludedUserIDs []string, groupIDs []string, userIds []string, predefinedInactiveUsersScope governance.PredefinedInactiveUsersScopeSettings) *governance.PrincipalScopeSettingsMutable {
	var principalScopeSettings governance.PrincipalScopeSettingsMutable
	principalScopeSettings.SetType(governance.PrincipalScopeType(d.PrincipalScope.Type.ValueString()))
	principalScopeSettings.SetExcludedUserIds(excludedUserIDs)
	principalScopeSettings.SetGroupIds(groupIDs)
	principalScopeSettings.SetUserIds(userIds)
	principalScopeSettings.SetIncludeOnlyActiveUsers(d.PrincipalScope.IncludeOnlyActiveUsers.ValueBool())
	principalScopeSettings.SetOnlyIncludeUsersWithSODConflicts(d.PrincipalScope.OnlyIncludeUsersWithSODConflicts.ValueBool())
	if d.PrincipalScope.UserScopeExpression.ValueStringPointer() != nil {
		principalScopeSettings.SetUserScopeExpression(d.PrincipalScope.UserScopeExpression.ValueString())
	}
	if predefinedInactiveUsersScope.InactiveDays != nil {
		principalScopeSettings.SetPredefinedInactiveUsersScope(predefinedInactiveUsersScope)
	}
	return &principalScopeSettings
}

func buildPredefinedUserScope(principalScopeSettings []inactiveUsersScopeModel) governance.PredefinedInactiveUsersScopeSettings {
	var predefinedScopeSettings governance.PredefinedInactiveUsersScopeSettings
	for _, scope := range principalScopeSettings {
		predefinedScopeSettings.SetInactiveDays(scope.InactiveDays.ValueInt32())
	}
	return predefinedScopeSettings
}

func buildTargetResources(d campaignResourceModel) []governance.TargetResourcesRequestInner {
	var targetResources []governance.TargetResourcesRequestInner
	for _, tr := range d.ResourceSettings.TargetResources {
		rt := governance.ResourceType(tr.ResourceType.ValueString())
		x := governance.TargetResourcesRequestInner{
			ResourceId:   tr.ResourceId.ValueString(),
			ResourceType: &rt,
		}
		if tr.IncludeAllEntitlementsAndBundles.ValueBoolPointer() != nil {
			x.IncludeAllEntitlementsAndBundles = tr.IncludeAllEntitlementsAndBundles.ValueBoolPointer()
		}
		if len(tr.Entitlements) > 0 {
			var entitlements []governance.EntitlementsInner
			for _, entitlement := range tr.Entitlements {
				v := getEntitlementValue(entitlement.Values)
				entitlementInner := governance.EntitlementsInner{
					Id:               entitlement.Id.ValueString(),
					IncludeAllValues: entitlement.IncludeAllValues.ValueBoolPointer(),
					Values:           v,
				}
				entitlements = append(entitlements, entitlementInner)
			}
			x.Entitlements = entitlements
		}
		targetResources = append(targetResources, x)
	}
	return targetResources
}

func getEntitlementValue(entitlements []entitlementValueModel) []governance.EntitlementValue {
	var value []governance.EntitlementValue
	for _, e := range entitlements {
		v := governance.EntitlementValue{
			Id: e.Id.ValueString(),
		}
		value = append(value, v)
	}
	return value
}

func getScheduleSettingForRequests(d campaignResourceModel, parsedStartDate time.Time, recur governance.RecurrenceDefinitionMutable) *governance.ScheduleSettingsMutable {
	scheduleSettings := governance.NewScheduleSettingsMutableWithDefaults()
	scheduleSettings.SetRecurrence(recur)
	scheduleSettings.SetStartDate(parsedStartDate)
	scheduleSettings.SetDurationInDays(float32(d.ScheduleSettings.DurationInDays.ValueInt32()))
	scheduleSettings.SetTimeZone(d.ScheduleSettings.TimeZone.ValueString())
	scheduleSettings.SetType(governance.ScheduleType(d.ScheduleSettings.Type.ValueString()))
	return scheduleSettings
}

func getReviewerSettingForRequests(d campaignResourceModel, reviewerLevels []governance.ReviewerLevelSettingsMutable) *governance.ReviewerSettingsMutable {
	reviewerSettings := governance.NewReviewerSettingsMutableWithDefaults()
	if d.ReviewerSettings != nil {
		if !d.ReviewerSettings.Type.IsNull() {
			reviewerSettings.Type = governance.CampaignReviewerType(d.ReviewerSettings.Type.ValueString())
		}
		if !d.ReviewerSettings.BulkDecisionDisabled.IsNull() {
			reviewerSettings.BulkDecisionDisabled = d.ReviewerSettings.BulkDecisionDisabled.ValueBoolPointer()
		}
		if !d.ReviewerSettings.FallbackReviewerId.IsNull() {
			reviewerSettings.FallBackReviewerId = d.ReviewerSettings.FallbackReviewerId.ValueStringPointer()
		}
		if !d.ReviewerSettings.JustificationRequired.IsNull() {
			reviewerSettings.JustificationRequired = d.ReviewerSettings.JustificationRequired.ValueBoolPointer()
		}
		if !d.ReviewerSettings.ReassignmentDisabled.IsNull() {
			reviewerSettings.ReassignmentDisabled = d.ReviewerSettings.ReassignmentDisabled.ValueBoolPointer()
		}
		if !d.ReviewerSettings.ReviewerGroupId.IsNull() {
			reviewerSettings.ReviewerGroupId = d.ReviewerSettings.ReviewerGroupId.ValueStringPointer()
		}
		if !d.ReviewerSettings.ReviewerId.IsNull() {
			reviewerSettings.ReviewerId = d.ReviewerSettings.ReviewerId.ValueStringPointer()
		}
		if !d.ReviewerSettings.ReviewerScopeExpression.IsNull() {
			reviewerSettings.ReviewerScopeExpression = d.ReviewerSettings.ReviewerScopeExpression.ValueStringPointer()
		}
		if !d.ReviewerSettings.SelfReviewDisabled.IsNull() {
			reviewerSettings.SelfReviewDisabled = d.ReviewerSettings.SelfReviewDisabled.ValueBoolPointer()
		}
		reviewerSettings.ReviewerLevels = reviewerLevels
	}
	return reviewerSettings
}

func getNotificationSettingsForRequest(d campaignResourceModel, remindersReviewerBeforeCampaignCloseInSecs []int32) *governance.NotificationSettings {
	notificationSettings := governance.NewNotificationSettingsWithDefaults()
	if d.NotificationSettings != nil {
		if !d.NotificationSettings.NotifyReviewerAtCampaignEnd.IsNull() {
			notificationSettings.NotifyReviewerAtCampaignEnd = d.NotificationSettings.NotifyReviewerAtCampaignEnd.ValueBoolPointer()
		}
		if !d.NotificationSettings.NotifyReviewerDuringMidpointOfReview.IsNull() {
			notificationSettings.NotifyReviewerDuringMidpointOfReview = *toNullableBool(d.NotificationSettings.NotifyReviewerDuringMidpointOfReview.ValueBoolPointer())
		}
		if !d.NotificationSettings.NotifyReviewerWhenOverdue.IsNull() {
			notificationSettings.NotifyReviewerWhenOverdue = *toNullableBool(d.NotificationSettings.NotifyReviewerWhenOverdue.ValueBoolPointer())
		}
		if !d.NotificationSettings.NotifyReviewerWhenReviewAssigned.IsNull() {
			notificationSettings.NotifyReviewerWhenReviewAssigned = d.NotificationSettings.NotifyReviewerWhenReviewAssigned.ValueBoolPointer()
		}
		if !d.NotificationSettings.NotifyReviewPeriodEnd.IsNull() {
			notificationSettings.NotifyReviewPeriodEnd = *toNullableBool(d.NotificationSettings.NotifyReviewPeriodEnd.ValueBoolPointer())
		}
		if !d.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs.IsNull() && !d.NotificationSettings.RemindersReviewerBeforeCampaignCloseInSecs.IsUnknown() {
			notificationSettings.RemindersReviewerBeforeCampaignCloseInSecs = remindersReviewerBeforeCampaignCloseInSecs
		}
	}
	return notificationSettings
}

func toNullableBool(v *bool) *governance.NullableBool {
	if v == nil {
		return nil
	}
	return governance.NewNullableBool(v)
}
