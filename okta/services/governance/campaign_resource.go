package provider

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"github.com/okta/terraform-provider-okta/okta/config"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &campaignResource{}

func newCampaignResource() resource.Resource {
	return &campaignResource{}
}

type campaignResource struct {
	*config.Config
}

type campaignResourceModel struct {
	Id                  types.String                     `tfsdk:"id"`
	Name                types.String                     `tfsdk:"name"`
	RemediationSettings campaignRemediationSettingsModel `tfsdk:"remediation_settings"`
	ResourceSettings    campaignResourceSettingsModel    `tfsdk:"resource_settings"`
	ReviewerSettings    campaignReviewerSettingsModel    `tfsdk:"reviewer_settings"`
	ScheduleSettings    campaignScheduleSettingsModel    `tfsdk:"schedule_settings"`
}

type campaignRemediationSettingsModel struct {
	AccessApproved types.String `tfsdk:"access_approved"`
	AccessRevoked  types.String `tfsdk:"access_revoked"`
	NoResponse     types.String `tfsdk:"no_response"`
}

type campaignResourceSettingsModel struct {
	Type            types.String                  `tfsdk:"type"`
	TargetResources []campaignTargetResourceModel `tfsdk:"target_resources"`
}

type campaignTargetResourceModel struct {
	ResourceId   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
}

type campaignReviewerSettingsModel struct {
	Type                    types.String `tfsdk:"type"`
	ReviewerGroupId         types.String `tfsdk:"reviewer_group_id"`
	ReviewerId              types.String `tfsdk:"reviewer_id"`
	ReviewerScopeExpression types.String `tfsdk:"reviewer_scope_expression"`
	FallbackReviewerId      types.String `tfsdk:"fallback_reviewer_id"`
}

type campaignScheduleSettingsModel struct {
	StartDate      types.String `tfsdk:"start_date"`
	DurationInDays types.Int32  `tfsdk:"duration_in_days"`
	TimeZone       types.String `tfsdk:"time_zone"`
	Type           types.String `tfsdk:"type"`
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
		},
		Blocks: map[string]schema.Block{
			"remediation_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
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
				},
				Description: "Specify the action to be taken after a reviewer makes a decision to APPROVE or REVOKE the access, or if the campaign was CLOSED and there was no response from the reviewer.",
			},

			"resource_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
					},
					Blocks: map[string]schema.Block{
						"target_resources": schema.ListNestedBlock{
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

			"reviewer_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
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
						"fallback_reviewer_id": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},

			"schedule_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
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

	// Create API call logic
	r.OktaIDaaSClient.OktaIGSDKClientV5().CampaignsAPI.CreateCampaign(ctx).CampaignMutable(buildCampaign(data)).Execute()
	// Example data value setting
	data.Id = types.StringValue("example-id")

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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *campaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data campaignResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *campaignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data campaignResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}

func buildCampaign(d campaignResourceModel) oktaInternalGovernance.CampaignMutable {

	rt := oktaInternalGovernance.RESOURCETYPE_GROUP
	//reviewerGroupID := d.Get("reviewer_settings.0.reviewer_group_id").(string)
	reviewerGroupID := d.ReviewerSettings.ReviewerId.String()
	//reviewerId := d.Get("reviewer_settings.0.reviewer_id").(string)
	reviewerId := d.ReviewerSettings.ReviewerId.String()

	//reviewerScopeExpression := d.Get("reviewer_settings.0.reviewer_scope_expression").(string)
	reviewerScopeExpression := d.ReviewerSettings.ReviewerScopeExpression.String()

	//str := d.Get("schedule_settings.0.start_date").(string)
	startDate := d.ScheduleSettings.StartDate.String()
	parsedStartDate, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return oktaInternalGovernance.CampaignMutable{}
	}

	return oktaInternalGovernance.CampaignMutable{
		Name: d.Name.String(),
		RemediationSettings: oktaInternalGovernance.RemediationSettings{
			AccessApproved: oktaInternalGovernance.ApprovedRemediationAction(d.RemediationSettings.AccessApproved.String()),
			AccessRevoked:  oktaInternalGovernance.RevokedRemediationAction(d.RemediationSettings.AccessRevoked.String()),
			NoResponse:     oktaInternalGovernance.NoResponseRemediationAction(d.RemediationSettings.NoResponse.String()),
		},
		ResourceSettings: oktaInternalGovernance.ResourceSettingsMutable{
			Type: oktaInternalGovernance.CampaignResourceType(d.ResourceSettings.Type.String()),
			TargetResources: []oktaInternalGovernance.TargetResourcesRequestInner{
				{
					ResourceId:   d.ResourceSettings.TargetResources[0].ResourceId.String(),
					ResourceType: &rt,
				},
			},
		},
		ReviewerSettings: oktaInternalGovernance.ReviewerSettingsMutable{
			Type:                    oktaInternalGovernance.CampaignReviewerType(d.ReviewerSettings.Type.String()),
			ReviewerGroupId:         &reviewerGroupID,
			ReviewerId:              &reviewerId,
			ReviewerScopeExpression: &reviewerScopeExpression,
		},
		ScheduleSettings: oktaInternalGovernance.ScheduleSettingsMutable{
			StartDate:      parsedStartDate,
			DurationInDays: float32(d.ScheduleSettings.DurationInDays.ValueInt32()),
			TimeZone:       d.ScheduleSettings.TimeZone.String(),
			Type:           oktaInternalGovernance.ScheduleType(d.ScheduleSettings.Type.String()),
		},
	}
}
