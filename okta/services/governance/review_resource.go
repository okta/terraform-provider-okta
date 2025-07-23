package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/okta/terraform-provider-okta/okta/config"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	//CurrentReviewerLevel types.String `tfsdk:"current_reviewer_level"`
	Created       types.String `tfsdk:"created"`
	CreatedBy     types.String `tfsdk:"created_by"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	LastUpdatedBy types.String `tfsdk:"last_updated_by"`
	//Decided              types.String `tfsdk:"decided"`
	//PrincipalProfile     *PrincipalProfileModel `tfsdk:"principal_profile"`
	//ReviewerProfile      *PrincipalProfileModel `tfsdk:"reviewer_profile"`
	//InternalLinks        *linksModel            `tfsdk:"internal_links"`
	//ExternalLinks        *linksModel            `tfsdk:"external_links"`
}

func (r *reviewResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_review"
}

func (r *reviewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"note": schema.StringAttribute{
				Required: true,
			},
			"campaign_id": schema.StringAttribute{
				Required: true,
			},
			"reviewer_id": schema.StringAttribute{
				Required: true,
			},
			"reviewer_level": schema.StringAttribute{
				Optional: true,
			},
			"review_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"resource_id":     schema.StringAttribute{Computed: true},
			"decision":        schema.StringAttribute{Computed: true},
			"reviewer_type":   schema.StringAttribute{Computed: true},
			"created":         schema.StringAttribute{Computed: true},
			"created_by":      schema.StringAttribute{Computed: true},
			"last_updated":    schema.StringAttribute{Computed: true},
			"last_updated_by": schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			//"principal_profile": schema.SingleNestedBlock{
			//	Attributes: map[string]schema.Attribute{
			//		"id":         schema.StringAttribute{Computed: true},
			//		"email":      schema.StringAttribute{Computed: true},
			//		"first_name": schema.StringAttribute{Computed: true},
			//		"last_name":  schema.StringAttribute{Computed: true},
			//		"status":     schema.StringAttribute{Computed: true},
			//		"login":      schema.StringAttribute{Computed: true},
			//	},
			//},
			//"reviewer_profile": schema.SingleNestedBlock{
			//	Attributes: map[string]schema.Attribute{
			//		"id":         schema.StringAttribute{Computed: true},
			//		"email":      schema.StringAttribute{Computed: true},
			//		"first_name": schema.StringAttribute{Computed: true},
			//		"last_name":  schema.StringAttribute{Computed: true},
			//		"status":     schema.StringAttribute{Computed: true},
			//		"login":      schema.StringAttribute{Computed: true},
			//	},
			//},
			//"internal_links": schema.SingleNestedBlock{
			//	Attributes: map[string]schema.Attribute{
			//		"self_href": schema.StringAttribute{
			//			Computed: true,
			//			Optional: true,
			//		},
			//		"reassign_review_href": schema.StringAttribute{
			//			Computed: true,
			//			Optional: true,
			//		},
			//	},
			//},
			//"external_links": schema.SingleNestedBlock{
			//	Attributes: map[string]schema.Attribute{
			//		"self_href": schema.StringAttribute{
			//			Computed: true,
			//			Optional: true,
			//		},
			//		"reassign_review_href": schema.StringAttribute{
			//			Computed: true,
			//			Optional: true,
			//		},
			//	},
			//},
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
		OktaIGSDKClientV5().
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
	fmt.Println("Final reveiwerID", review.ReviewerProfile.Id)
	data.Id = types.StringValue(review.Id)
	data.ReviewerId = types.StringValue(review.ReviewerProfile.Id)
	data.CampaignId = types.StringValue(review.CampaignId)
	data.ResourceId = types.StringValue(review.ResourceId)
	data.Decision = types.StringValue(string(review.Decision))
	data.ReviewerType = types.StringValue(string(review.ReviewerType))
	//data.CurrentReviewerLevel = func() types.String {
	//	if review.CurrentReviewerLevel != nil {
	//		return types.StringValue(string(*review.CurrentReviewerLevel))
	//	}
	//	return types.StringNull()
	//}()
	data.Created = types.StringValue(review.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(review.CreatedBy)
	data.LastUpdated = types.StringValue(review.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(review.LastUpdatedBy)
	//data.Decided = func() types.String {
	//	if review.Decided != nil {
	//		return types.StringValue(review.Decided.Format(time.RFC3339))
	//	}
	//	return types.StringNull()
	//}()
	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

//	func applyReassignReviewToState(ctx context.Context, data *oktaInternalGovernance.ReviewFull, state *reviewResourceModel) diag.Diagnostics {
//		var diags diag.Diagnostics
//
//		state.Id = types.StringValue(data.Id)
//		state.CampaignId = types.StringValue(data.CampaignId)
//		state.ResourceId = types.StringValue(data.ResourceId)
//		state.Decision = types.StringValue(string(data.Decision))
//		state.ReviewerType = types.StringValue(string(data.ReviewerType))
//
//		state.CurrentReviewerLevel = func() types.String {
//			if data.CurrentReviewerLevel != nil {
//				return types.StringValue(string(*data.CurrentReviewerLevel))
//			}
//			return types.StringNull()
//		}()
//
//		state.Created = types.StringValue(data.Created.Format(time.RFC3339))
//		state.CreatedBy = types.StringValue(data.CreatedBy)
//		state.LastUpdated = types.StringValue(data.LastUpdated.Format(time.RFC3339))
//		state.LastUpdatedBy = types.StringValue(data.LastUpdatedBy)
//
//		state.Decided = func() types.String {
//			if data.Decided != nil {
//				return types.StringValue(data.Decided.Format(time.RFC3339))
//			}
//			return types.StringNull()
//		}()
//
//		state.PrincipalProfile = convertPrincipalProfile(&data.PrincipalProfile)
//		state.ReviewerProfile = convertPrincipalProfile(data.ReviewerProfile)
//		state.Links = convertLinks(&data.Links)
//
//		return diags
//	}
func (r *reviewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data reviewResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getReview, _, err := r.OktaGovernanceClient.
		OktaIGSDKClientV5().
		ReviewsAPI.GetReview(ctx, data.CampaignId.ValueString()).Execute()
	if err != nil {
		return
	}
	data.Id = types.StringValue(getReview.Id)
	data.ReviewerId = types.StringValue(getReview.ReviewerProfile.Id)
	data.CampaignId = types.StringValue(getReview.CampaignId)
	data.ResourceId = types.StringValue(getReview.ResourceId)
	data.Decision = types.StringValue(string(getReview.Decision))
	data.ReviewerType = types.StringValue(string(getReview.ReviewerType))
	//data.CurrentReviewerLevel = func() types.String {
	//	if getReview.CurrentReviewerLevel != nil {
	//		return types.StringValue(string(*getReview.CurrentReviewerLevel))
	//	}
	//	return types.StringNull()
	//}()
	data.Created = types.StringValue(getReview.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(getReview.CreatedBy)
	data.LastUpdated = types.StringValue(getReview.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(getReview.LastUpdatedBy)
	//data.Decided = func() types.String {
	//	if getReview.Decided != nil {
	//		return types.StringValue(getReview.Decided.Format(time.RFC3339))
	//	}
	//	return types.StringNull()
	//}()

	// No-op: the reassignment is not persisted on Okta's side
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reviewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data reviewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := buildReassignReviewRequest(data)
	// Call the Okta API
	reassignedReview, _, err := r.OktaGovernanceClient.
		OktaIGSDKClientV5().
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
	fmt.Println("Final reveiwerID", review.ReviewerProfile.Id)
	data.Id = types.StringValue(review.Id)
	data.ReviewerId = types.StringValue(review.ReviewerProfile.Id)
	data.CampaignId = types.StringValue(review.CampaignId)
	data.ResourceId = types.StringValue(review.ResourceId)
	data.Decision = types.StringValue(string(review.Decision))
	data.ReviewerType = types.StringValue(string(review.ReviewerType))
	//data.CurrentReviewerLevel = func() types.String {
	//	if review.CurrentReviewerLevel != nil {
	//		return types.StringValue(string(*review.CurrentReviewerLevel))
	//	}
	//	return types.StringNull()
	//}()
	data.Created = types.StringValue(review.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(review.CreatedBy)
	data.LastUpdated = types.StringValue(review.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(review.LastUpdatedBy)
	//data.Decided = func() types.String {
	//	if review.Decided != nil {
	//		return types.StringValue(review.Decided.Format(time.RFC3339))
	//	}
	//	return types.StringNull()
	//}()
	////// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reviewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No delete call exists in Okta — review reassignments are immutable
}

func buildReassignReviewRequest(data reviewResourceModel) oktaInternalGovernance.ReviewsReassign {
	reassignReviewBody := oktaInternalGovernance.ReviewsReassign{}
	reviewerLevel, _ := oktaInternalGovernance.NewReviewerLevelTypeFromValue(data.ReviewerLevel.ValueString())
	reviewIds := buildReviewIds(data.ReviewIds)
	if len(reviewIds) > 0 {
		reassignReviewBody.ReviewIds = reviewIds
	}
	reassignReviewBody.ReviewerLevel = reviewerLevel
	reassignReviewBody.ReviewerId = data.ReviewerId.ValueString()
	reassignReviewBody.Note = data.Note.ValueString()
	fmt.Println("reassignReviewBody ReviewerID:", reassignReviewBody.ReviewerId)
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
