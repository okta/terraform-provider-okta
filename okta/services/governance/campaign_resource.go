package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"fmt"
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
	Id                  types.String                     `tfsdk:"id"`
	Name                types.String                     `tfsdk:"name"`
	LaunchCampaign      types.Bool                       `tfsdk:"launch_campaign"`
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
			"launch_campaign": schema.BoolAttribute{
				Optional:    true,
				Description: "Launch the campaign after creation. Defaults to false.",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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
				Description: "Specify the action to be taken after a reviewer makes a decision to APPROVE or REVOKE the access, or if the campaign was CLOSED and there was no response from the reviewer.",
			},

			"resource_settings": schema.SingleNestedBlock{
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

			"reviewer_settings": schema.SingleNestedBlock{
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

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getCampaignResponse, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().CampaignsAPI.GetCampaign(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	applyCampaignsToState(ctx, getCampaignResponse, &data, resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyCampaignsToState(ctx context.Context, resp *oktaInternalGovernance.CampaignFull, c *campaignResourceModel, diagnostics diag.Diagnostics) {

	c.Id = types.StringValue(resp.Id)
	c.Name = types.StringValue(resp.Name)

	c.RemediationSettings.AccessApproved = types.StringValue(string(resp.RemediationSettings.AccessApproved))
	c.RemediationSettings.AccessRevoked = types.StringValue(string(resp.RemediationSettings.AccessRevoked))
	c.RemediationSettings.NoResponse = types.StringValue(string(resp.RemediationSettings.NoResponse))

	c.ResourceSettings.Type = types.StringValue(string(resp.ResourceSettings.Type))
	if len(resp.ResourceSettings.TargetResources) > 0 {
		c.ResourceSettings.TargetResources = []campaignTargetResourceModel{
			{
				ResourceId:   types.StringValue(resp.ResourceSettings.TargetResources[0].ResourceId),
				ResourceType: types.StringValue(string(*resp.ResourceSettings.TargetResources[0].ResourceType)),
			},
		}
	}

	c.ReviewerSettings.Type = types.StringValue(string(resp.ReviewerSettings.Type))
	if resp.ReviewerSettings.ReviewerGroupId != nil {
		c.ReviewerSettings.ReviewerGroupId = types.StringValue(*resp.ReviewerSettings.ReviewerGroupId)
	}
	if resp.ReviewerSettings.ReviewerId != nil {
		c.ReviewerSettings.ReviewerId = types.StringValue(*resp.ReviewerSettings.ReviewerId)
	}
	if resp.ReviewerSettings.ReviewerScopeExpression != nil {
		c.ReviewerSettings.ReviewerScopeExpression = types.StringValue(*resp.ReviewerSettings.ReviewerScopeExpression)
	}
	if resp.ReviewerSettings.FallBackReviewerId != nil {
		c.ReviewerSettings.FallbackReviewerId = types.StringValue(*resp.ReviewerSettings.FallBackReviewerId)
	}

	fmt.Println("Time read from API:", resp.ScheduleSettings.StartDate)
	t, err := time.Parse(time.RFC3339, resp.ScheduleSettings.StartDate.String())
	if err != nil {
		return
	}
	c.ScheduleSettings.StartDate = types.StringValue(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	c.ScheduleSettings.DurationInDays = types.Int32Value(int32(resp.ScheduleSettings.DurationInDays))
	c.ScheduleSettings.TimeZone = types.StringValue(resp.ScheduleSettings.TimeZone)
	c.ScheduleSettings.Type = types.StringValue(string(resp.ScheduleSettings.Type))

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

	// Build and return CampaignMutable
	return oktaInternalGovernance.CampaignMutable{
		Name: d.Name.ValueString(),

		RemediationSettings: oktaInternalGovernance.RemediationSettings{
			AccessApproved: oktaInternalGovernance.ApprovedRemediationAction(d.RemediationSettings.AccessApproved.ValueString()),
			AccessRevoked:  oktaInternalGovernance.RevokedRemediationAction(d.RemediationSettings.AccessRevoked.ValueString()),
			NoResponse:     oktaInternalGovernance.NoResponseRemediationAction(d.RemediationSettings.NoResponse.ValueString()),
		},

		ResourceSettings: oktaInternalGovernance.ResourceSettingsMutable{
			Type:            oktaInternalGovernance.CampaignResourceType(d.ResourceSettings.Type.ValueString()),
			TargetResources: targetResources,
		},

		ReviewerSettings: oktaInternalGovernance.ReviewerSettingsMutable{
			Type:                    oktaInternalGovernance.CampaignReviewerType(d.ReviewerSettings.Type.ValueString()),
			ReviewerGroupId:         ptrString(d.ReviewerSettings.ReviewerGroupId.ValueString()),
			ReviewerId:              ptrString(d.ReviewerSettings.ReviewerId.ValueString()),
			ReviewerScopeExpression: ptrString(d.ReviewerSettings.ReviewerScopeExpression.ValueString()),
		},

		ScheduleSettings: oktaInternalGovernance.ScheduleSettingsMutable{
			StartDate:      parsedStartDate,
			DurationInDays: float32(d.ScheduleSettings.DurationInDays.ValueInt32()),
			TimeZone:       d.ScheduleSettings.TimeZone.ValueString(),
			Type:           oktaInternalGovernance.ScheduleType(d.ScheduleSettings.Type.ValueString()),
		},
	}
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func isOnlyLaunchChanged(plan, state campaignResourceModel) bool {
	return plan.Name.Equal(state.Name) &&
		equalRemediation(plan.RemediationSettings, state.RemediationSettings) &&
		equalResourceSettings(plan.ResourceSettings, state.ResourceSettings) &&
		equalReviewerSettings(plan.ReviewerSettings, state.ReviewerSettings) &&
		equalScheduleSettings(plan.ScheduleSettings, state.ScheduleSettings)
}

func equalScheduleSettings(planSchedule, stateSchedule campaignScheduleSettingsModel) bool {
	return planSchedule.Type.Equal(stateSchedule.Type) &&
		planSchedule.TimeZone.Equal(stateSchedule.TimeZone) &&
		planSchedule.DurationInDays.Equal(stateSchedule.DurationInDays) &&
		planSchedule.StartDate.Equal(stateSchedule.StartDate)
}

func equalReviewerSettings(planReviewer, stateReviewer campaignReviewerSettingsModel) bool {
	return planReviewer.Type.Equal(stateReviewer.Type) &&
		planReviewer.ReviewerGroupId.Equal(stateReviewer.ReviewerGroupId) &&
		planReviewer.ReviewerId.Equal(stateReviewer.ReviewerId) &&
		planReviewer.ReviewerScopeExpression.Equal(stateReviewer.ReviewerScopeExpression) &&
		planReviewer.FallbackReviewerId.Equal(stateReviewer.FallbackReviewerId)
}

func equalResourceSettings(planResource, stateResource campaignResourceSettingsModel) bool {
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

func equalRemediation(a, b campaignRemediationSettingsModel) bool {
	return a.AccessApproved.Equal(b.AccessApproved) &&
		a.AccessRevoked.Equal(b.AccessRevoked) &&
		a.NoResponse.Equal(b.NoResponse)
}
