package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"fmt"
	"github.com/okta/terraform-provider-okta/okta/config"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &reviewDataSource{}

func newReviewDataSource() datasource.DataSource {
	return &reviewDataSource{}
}

type reviewDataSource struct {
	*config.Config
}

type ReviewerEntitlementValue struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type reviewDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	CampaignId           types.String `tfsdk:"campaign_id"`
	ResourceId           types.String `tfsdk:"resource_id"`
	Decision             types.String `tfsdk:"decision"`
	RemediationStatus    types.String `tfsdk:"remediation_status"`
	ReviewerType         types.String `tfsdk:"reviewer_type"`
	CurrentReviewerLevel types.String `tfsdk:"current_reviewer_level"`
	Created              types.String `tfsdk:"created"`
	CreatedBy            types.String `tfsdk:"created_by"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	LastUpdatedBy        types.String `tfsdk:"last_updated_by"`
	Decided              types.String `tfsdk:"decided"`

	PrincipalProfile  *PrincipalProfileModel    `tfsdk:"principal_profile"`
	ReviewerProfile   *PrincipalProfileModel    `tfsdk:"reviewer_profile"`
	EntitlementValue  *ReviewerEntitlementValue `tfsdk:"entitlement_value"`
	Note              *noteModel                `tfsdk:"note"`
	AllReviewerLevels []reviewLevelModel        `tfsdk:"all_reviewer_levels"`
	Links             *linksModel               `tfsdk:"links"`
}

type PrincipalProfileModel struct {
	Id        string                                        `tfsdk:"id"`
	Email     string                                        `tfsdk:"email"`
	FirstName *string                                       `tfsdk:"first_name"`
	LastName  *string                                       `tfsdk:"last_name"`
	Login     *string                                       `tfsdk:"login"`
	Status    oktaInternalGovernance.PrincipalProfileStatus `tfsdk:"status"`
}

type userProfileModel struct {
	Id        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Status    types.String `tfsdk:"status"`
}

type linksModel struct {
	SelfHref           types.String `tfsdk:"self_href"`
	ReassignReviewHref types.String `tfsdk:"reassign_review_href"`
}

type noteModel struct {
	Id   types.String `tfsdk:"id"`
	Note types.String `tfsdk:"note"`
}

type reviewLevelModel struct {
	Id                      types.String `tfsdk:"id"`
	CreatedBy               types.String `tfsdk:"created_by"`
	Created                 types.String `tfsdk:"created"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	LastUpdatedBy           types.String `tfsdk:"last_updated_by"`
	ReviewerLevel           types.String `tfsdk:"reviewer_level"`
	Decision                types.String `tfsdk:"decision"`
	ReviewerType            types.String `tfsdk:"reviewer_type"`
	ReviewerGroupResourceId types.String `tfsdk:"reviewer_group_resource_id"`

	ReviewerProfile      userProfileModel  `tfsdk:"reviewer_profile"`
	ReviewerGroupProfile groupProfileModel `tfsdk:"reviewer_group_profile"`
}

type groupProfileModel struct {
	GroupId    types.String `tfsdk:"group_id"`
	GlobalName types.String `tfsdk:"global_name"`
	Name       types.String `tfsdk:"name"`
	GroupType  types.String `tfsdk:"group_type"`
}

type linkModel struct {
	Href types.String `tfsdk:"href"`
}

func (d *reviewDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_review"
}

func (d *reviewDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *reviewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     schema.StringAttribute{Required: true},
			"campaign_id":            schema.StringAttribute{Computed: true},
			"resource_id":            schema.StringAttribute{Computed: true},
			"decision":               schema.StringAttribute{Computed: true},
			"remediation_status":     schema.StringAttribute{Computed: true},
			"reviewer_type":          schema.StringAttribute{Computed: true},
			"current_reviewer_level": schema.StringAttribute{Computed: true},
			"created":                schema.StringAttribute{Computed: true},
			"created_by":             schema.StringAttribute{Computed: true},
			"last_updated":           schema.StringAttribute{Computed: true},
			"last_updated_by":        schema.StringAttribute{Computed: true},
			"decided":                schema.StringAttribute{Computed: true},
		},
		Blocks: map[string]schema.Block{
			"principal_profile": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true},
					"email":      schema.StringAttribute{Computed: true},
					"first_name": schema.StringAttribute{Computed: true},
					"last_name":  schema.StringAttribute{Computed: true},
					"status":     schema.StringAttribute{Computed: true},
					"login":      schema.StringAttribute{Computed: true},
				},
			},
			"reviewer_profile": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true},
					"email":      schema.StringAttribute{Computed: true},
					"first_name": schema.StringAttribute{Computed: true},
					"last_name":  schema.StringAttribute{Computed: true},
					"status":     schema.StringAttribute{Computed: true},
					"login":      schema.StringAttribute{Computed: true},
				},
			},
			"entitlement_value": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"note": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"id":   schema.StringAttribute{Computed: true},
					"note": schema.StringAttribute{Computed: true},
				},
			},
			"all_reviewer_levels": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                         schema.StringAttribute{Computed: true},
						"created_by":                 schema.StringAttribute{Computed: true},
						"created":                    schema.StringAttribute{Computed: true},
						"last_updated":               schema.StringAttribute{Computed: true},
						"last_updated_by":            schema.StringAttribute{Computed: true},
						"reviewer_level":             schema.StringAttribute{Computed: true},
						"decision":                   schema.StringAttribute{Computed: true},
						"reviewer_type":              schema.StringAttribute{Computed: true},
						"reviewer_group_resource_id": schema.StringAttribute{Computed: true},
					},
					Blocks: map[string]schema.Block{
						"reviewer_profile": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"id":         schema.StringAttribute{Computed: true},
								"email":      schema.StringAttribute{Computed: true},
								"first_name": schema.StringAttribute{Computed: true},
								"last_name":  schema.StringAttribute{Computed: true},
								"status":     schema.StringAttribute{Computed: true},
							},
						},
						"reviewer_group_profile": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"group_id":    schema.StringAttribute{Computed: true},
								"global_name": schema.StringAttribute{Computed: true},
								"name":        schema.StringAttribute{Computed: true},
								"group_type":  schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"links": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"self_href": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"reassign_review_href": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
				},
			},
		},
	}
}

func (d *reviewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data reviewDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reviewId := data.Id.ValueString()
	fmt.Println("reviewId:", reviewId)
	if reviewId == "" {
		resp.Diagnostics.AddError("Missing review ID", "The 'id' attribute must be set in the configuration.")
		return
	}

	// Call Okta API to fetch review details
	review, _, err := d.OktaGovernanceClient.OktaIGSDKClientV5().ReviewsAPI.GetReview(ctx, reviewId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read review",
			"Could not retrieve review from Okta Governance API: "+err.Error(),
		)
		return
	}

	// Map API response to Terraform state
	data = reviewDataSourceModel{
		Id:                types.StringValue(review.Id),
		CampaignId:        types.StringValue(review.CampaignId),
		ResourceId:        types.StringValue(review.ResourceId),
		Decision:          types.StringValue(string(review.Decision)),
		RemediationStatus: types.StringValue(string(review.RemediationStatus)),
		ReviewerType:      types.StringValue(string(review.ReviewerType)),
		CurrentReviewerLevel: func() types.String {
			if review.CurrentReviewerLevel != nil {
				return types.StringValue(string(*review.CurrentReviewerLevel))
			}
			return types.StringNull()
		}(),
		Created:       types.StringValue(review.Created.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(review.CreatedBy),
		LastUpdated:   types.StringValue(review.LastUpdated.Format(time.RFC3339)),
		LastUpdatedBy: types.StringValue(review.LastUpdatedBy),
		Decided: func() types.String {
			if review.Decided != nil {
				return types.StringValue(review.Decided.Format(time.RFC3339))
			}
			return types.StringNull()
		}(),
		PrincipalProfile:  convertPrincipalProfile(&review.PrincipalProfile),
		ReviewerProfile:   convertPrincipalProfile(review.ReviewerProfile),
		EntitlementValue:  convertEntitlementValue(review.EntitlementValue),
		Note:              convertNote(review.Note),
		AllReviewerLevels: convertReviewerLevels(review.AllReviewerLevels),
		Links:             convertLinks(&review.Links),
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertEntitlementValue(value *oktaInternalGovernance.ReviewerEntitlementValue) *ReviewerEntitlementValue {
	if value == nil || value.Id == "" || value.Name == "" {
		return nil
	}

	return &ReviewerEntitlementValue{
		Id:   types.StringValue(value.Id),
		Name: types.StringValue(value.Name),
	}
}

func convertNote(n *oktaInternalGovernance.Note) *noteModel {
	if n == nil || n.Note == nil || n.Id == nil {
		return nil
	}
	return &noteModel{
		Id:   types.StringPointerValue(n.Id),
		Note: types.StringPointerValue(n.Note),
	}
}

func convertReviewerLevels(levels []oktaInternalGovernance.ReviewerLevelInfoFull) []reviewLevelModel {
	var result []reviewLevelModel
	for _, l := range levels {

		level := reviewLevelModel{
			Id:            types.StringValue(l.Id),
			CreatedBy:     types.StringValue(l.CreatedBy),
			Created:       types.StringValue(l.Created.Format(time.RFC3339)),
			LastUpdated:   types.StringValue(l.LastUpdated.Format(time.RFC3339)),
			LastUpdatedBy: types.StringValue(l.LastUpdatedBy),
			ReviewerLevel: types.StringValue(string(l.ReviewerLevel)),
			Decision:      types.StringValue(string(l.Decision)),
			ReviewerType:  types.StringValue(string(l.ReviewerType)),
			//ReviewerGroupResourceId: types.StringValue(l.ReviewerGroupResourceId),
			ReviewerProfile:      buildUserProfileModel(l.ReviewerProfile),
			ReviewerGroupProfile: buildReviewerGroupProfile(l.ReviewerGroupProfile), // oktaInternalGovernance.NewReviewerGroupProfile(l.ReviewerGroupProfile.Name, l.ReviewerGroupProfile.GroupId, l.ReviewerGroupProfile.GroupType),
		}
		result = append(result, level)
	}
	return result
}

func buildReviewerGroupProfile(profile *oktaInternalGovernance.ReviewerGroupProfile) groupProfileModel {
	if profile == nil {
		return groupProfileModel{}
	}
	return groupProfileModel{
		GroupId:   types.StringValue(profile.GroupId),
		Name:      types.StringValue(profile.Name),
		GroupType: types.StringValue(string(profile.GroupType)),
	}
}

func buildUserProfileModel(profile *oktaInternalGovernance.PrincipalProfile) userProfileModel {
	userProfile := userProfileModel{}
	if profile == nil {
		return userProfileModel{}
	}
	userProfile.Id = types.StringValue(profile.Id)
	userProfile.Email = types.StringValue(profile.Email)
	if profile.FirstName != nil {
		userProfile.FirstName = types.StringValue(*profile.FirstName)
	}
	if profile.LastName != nil {
		userProfile.LastName = types.StringValue(*profile.LastName)
	}
	userProfile.Status = types.StringValue(string(profile.Status))
	return userProfile
}

//func convertLinks(links *oktaInternalGovernance.ReviewLinks) oktaInternalGovernance.Link {
//	if links == nil {
//		return oktaInternalGovernance.Link{}
//	}
//	return oktaInternalGovernance.Link{
//		Name: links.Self.Name,
//	}
//}

func convertLinks(links *oktaInternalGovernance.ReviewLinks) *linksModel {
	var selfHref, reassignHref string

	if links != nil {
		if links.Self.Href != "" {
			selfHref = links.Self.Href
		}
		if links.ReassignReview.Href != "" {
			reassignHref = links.ReassignReview.Href
		}
	}

	return &linksModel{
		SelfHref:           types.StringValue(selfHref),
		ReassignReviewHref: types.StringValue(reassignHref),
	}
}

func convertPrincipalProfile(p *oktaInternalGovernance.PrincipalProfile) *PrincipalProfileModel {
	if p == nil {
		return nil
	}

	return &PrincipalProfileModel{
		Id:        p.GetId(),
		Email:     p.GetEmail(),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Login:     p.Login,
		Status:    p.GetStatus(),
	}
}
