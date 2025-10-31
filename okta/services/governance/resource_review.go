package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &reviewResource{}
	_ resource.ResourceWithConfigure   = &reviewResource{}
	_ resource.ResourceWithImportState = &reviewResource{}
)

func newReviewResource() resource.Resource {
	return &reviewResource{}
}

type reviewResource struct {
	*config.Config
}

func (r *reviewResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *reviewResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type reviewResourceModel struct {
	Id            types.String `tfsdk:"id"`
	CampaignId    types.String `tfsdk:"campaign_id"`
	ReviewerId    types.String `tfsdk:"reviewer_id"`
	Note          types.String `tfsdk:"note"`
	ReviewerLevel types.String `tfsdk:"reviewer_level"`
	ReviewIds     types.List   `tfsdk:"review_ids"`
	ResourceId    types.String `tfsdk:"resource_id"`
	Decision      types.String `tfsdk:"decision"`
	ReviewerType  types.String `tfsdk:"reviewer_type"`
	Created       types.String `tfsdk:"created"`
	CreatedBy     types.String `tfsdk:"created_by"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	LastUpdatedBy types.String `tfsdk:"last_updated_by"`
}

func (r *reviewResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_review"
}

func (r *reviewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the Review.",
			},
			"note": schema.StringAttribute{
				Required:    true,
				Description: "A note to justify the reassignment decision for the specified review.",
			},
			"campaign_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the campaign.",
			},
			"reviewer_id": schema.StringAttribute{
				Required:    true,
				Description: "The Okta user id of the new reviewer.",
			},
			"reviewer_level": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the reviewer level of each reviews during access certification. Applicable for multi level campaigns only.",
			},
			"review_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "A list of reviews (review id values) that are reassigned to the new reviewer.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"resource_id":     schema.StringAttribute{Computed: true},
			"decision":        schema.StringAttribute{Computed: true, Description: "The decision of the reviewer."},
			"reviewer_type":   schema.StringAttribute{Computed: true, Description: "The type of reviewer to which the review is assigned."},
			"created":         schema.StringAttribute{Computed: true, Description: "The ISO 8601 formatted date and time when the resource was created"},
			"created_by":      schema.StringAttribute{Computed: true, Description: "The id of user who created the resource."},
			"last_updated":    schema.StringAttribute{Computed: true, Description: "The ISO 8601 formatted date and time when the object was last updated."},
			"last_updated_by": schema.StringAttribute{Computed: true, Description: "The id of user who last updated the object."},
		},
	}
}

func (r *reviewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data reviewResourceModel

	// Read the plan into our model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := buildReassignReviewRequest(data)
	// Call the Okta API
	reassignedReview, _, err := r.OktaGovernanceClient.
		OktaGovernanceSDKClient().
		ReviewsAPI.
		ReassignReviews(ctx, data.CampaignId.ValueString()).
		ReviewsReassign(request).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reassigning reviews",
			"Could not reassign reviews: "+err.Error(),
		)
		return
	}

	applyReviewToState(ctx, &data, reassignedReview)
	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyReviewToState(ctx context.Context, data *reviewResourceModel, reassignedReview *governance.ReviewReassignList) {
	review := reassignedReview.Data[0]
	data.Id = types.StringValue(review.Id)
	data.ReviewerId = types.StringValue(review.ReviewerProfile.GetId())
	data.CampaignId = types.StringValue(review.GetCampaignId())
	data.ResourceId = types.StringValue(review.GetResourceId())
	data.Decision = types.StringValue(string(review.GetDecision()))
	data.ReviewerType = types.StringValue(string(review.GetReviewerType()))
	data.Created = types.StringValue(review.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(review.GetCreatedBy())
	data.LastUpdated = types.StringValue(review.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(review.GetLastUpdatedBy())
}

func (r *reviewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data reviewResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getReview, _, err := r.OktaGovernanceClient.
		OktaGovernanceSDKClient().
		ReviewsAPI.GetReview(ctx, data.CampaignId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Review",
			"Could not reading Review, unexpected error: "+err.Error(),
		)
		return
	}
	data.Id = types.StringValue(getReview.Id)
	data.ReviewerId = types.StringValue(getReview.ReviewerProfile.Id)
	data.CampaignId = types.StringValue(getReview.CampaignId)
	data.ResourceId = types.StringValue(getReview.ResourceId)
	data.Decision = types.StringValue(string(getReview.Decision))
	data.ReviewerType = types.StringValue(string(getReview.ReviewerType))
	data.Created = types.StringValue(getReview.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(getReview.CreatedBy)
	data.LastUpdated = types.StringValue(getReview.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(getReview.LastUpdatedBy)

	// No-op: the reassignment is not persisted on Okta's side
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reviewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data reviewResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := buildReassignReviewRequest(data)
	// Call the Okta API
	reassignedReview, _, err := r.OktaGovernanceClient.
		OktaGovernanceSDKClient().
		ReviewsAPI.
		ReassignReviews(ctx, data.CampaignId.ValueString()).
		ReviewsReassign(request).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reassigning reviews",
			"Could not reassign reviews: "+err.Error(),
		)
		return
	}

	review := reassignedReview.Data[0]
	data.Id = types.StringValue(review.Id)
	data.ReviewerId = types.StringValue(review.ReviewerProfile.Id)
	data.CampaignId = types.StringValue(review.CampaignId)
	data.ResourceId = types.StringValue(review.ResourceId)
	data.Decision = types.StringValue(string(review.Decision))
	data.ReviewerType = types.StringValue(string(review.ReviewerType))
	data.Created = types.StringValue(review.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(review.CreatedBy)
	data.LastUpdated = types.StringValue(review.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(review.LastUpdatedBy)
	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reviewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No delete call exists in Okta — review reassignments are immutable
}

func buildReassignReviewRequest(data reviewResourceModel) governance.ReviewsReassign {
	reassignReviewBody := governance.ReviewsReassign{}
	reviewerLevel, _ := governance.NewReviewerLevelTypeFromValue(data.ReviewerLevel.ValueString())
	reviewIds := buildReviewIds(data.ReviewIds)
	if len(reviewIds) > 0 {
		reassignReviewBody.ReviewIds = reviewIds
	}
	reassignReviewBody.ReviewerLevel = reviewerLevel
	reassignReviewBody.ReviewerId = data.ReviewerId.ValueString()
	reassignReviewBody.Note = data.Note.ValueString()
	return reassignReviewBody
}

func buildReviewIds(list types.List) []string {
	var reviewIds []string

	// Make sure list has elements
	if list.IsNull() || list.IsUnknown() {
		return reviewIds
	}

	// Extract each element as a string
	for _, elem := range list.Elements() {
		if strElem, ok := elem.(types.String); ok {
			reviewIds = append(reviewIds, strElem.ValueString())
		}
	}

	return reviewIds
}
