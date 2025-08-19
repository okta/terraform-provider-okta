package governance

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"strings"
)

var _ datasource.DataSource = &campaignsDataSource{}

func newCampaignsDataSource() datasource.DataSource {
	return &campaignsDataSource{}
}

type campaignsDataSource struct {
	*config.Config
}

func (d *campaignsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaigns"
}

func (d *campaignsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type campaignsModel struct {
	Id                  types.String `tfsdk:"id"`
	Created             types.String `tfsdk:"created"`
	CreatedBy           types.String `tfsdk:"created_by"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	LastUpdatedBy       types.String `tfsdk:"last_updated_by"`
	Name                types.String `tfsdk:"name"`
	Status              types.String `tfsdk:"status"`
	CampaignType        types.String `tfsdk:"campaign_type"`
	Description         types.String `tfsdk:"description"`
	RecurringCampaignId types.String `tfsdk:"recurring_campaign_id"`
}
type campaignsDataSourceModel struct {
	CampaignName           types.String     `tfsdk:"campaign_name"`
	CampaignStatus         types.String     `tfsdk:"campaign_status"`
	MultipleCampaignStatus types.List       `tfsdk:"multiple_campaign_status"`
	StartDate              types.String     `tfsdk:"start_date"`
	EndDate                types.String     `tfsdk:"end_date"`
	ScheduleTypeOneOff     types.String     `tfsdk:"schedule_type_one_off"`
	ScheduleTypeRecurring  types.String     `tfsdk:"schedule_type_recurring"`
	ReviewerTypeUser       types.String     `tfsdk:"reviewer_type_user"`
	ReviewerTypeGroup      types.String     `tfsdk:"reviewer_type_group"`
	ReviewerTypeGroupOwner types.String     `tfsdk:"reviewer_type_group_owner"`
	ReviewerTypeMultiLevel types.String     `tfsdk:"reviewer_type_multi_level"`
	RecurringCampaignId    types.String     `tfsdk:"recurring_campaign_id"`
	Data                   []campaignsModel `tfsdk:"data"`
}

func (d *campaignsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"campaign_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the campaign.",
			},
			"campaign_status": schema.StringAttribute{
				Optional:    true,
				Description: "Status of the campaign.",
			},
			"multiple_campaign_status": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of status of the campaign.",
			},
			"start_date": schema.StringAttribute{
				Optional:    true,
				Description: "Start date of the campaign.",
			},
			"end_date": schema.StringAttribute{
				Optional:    true,
				Description: "End date of the campaign.",
			},
			"schedule_type_one_off": schema.StringAttribute{
				Optional:    true,
				Description: "Type of a campaign.",
			},
			"schedule_type_recurring": schema.StringAttribute{
				Optional:    true,
				Description: "Type of a campaign.",
			},
			"reviewer_type_user": schema.StringAttribute{
				Optional:    true,
				Description: "Type of a campaign reviewer.",
			},
			"reviewer_type_group": schema.StringAttribute{
				Optional:    true,
				Description: "Type of a campaign reviewer.",
			},
			"reviewer_type_group_owner": schema.StringAttribute{
				Optional:    true,
				Description: "Type of a campaign reviewer.",
			},
			"reviewer_type_multi_level": schema.StringAttribute{
				Optional:    true,
				Description: "Type of a campaign reviewer.",
			},
			"recurring_campaign_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the recurring campaign if this campaign was created as part of a recurring schedule.",
			},
		},
		Blocks: map[string]schema.Block{
			"data": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
							Computed: true,
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
							Description: "The last updated date and time when the resource was last updated.",
						},
						"last_updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the last updated user who created the resource.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the campaign.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the campaign.",
						},
						"campaign_type": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Type of the campaign.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Description of the campaign.",
						},
						"recurring_campaign_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the recurring campaign if this campaign was created.",
						},
					},
				},
			},
		},
	}
}

func (d *campaignsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data campaignsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	campaigns, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().CampaignsAPI.ListCampaigns(ctx).Filter(buildFilterForCampaigns(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Campaigns",
			"Could not read Campaigns, unexpected error: "+err.Error(),
		)
		return
	}

	var campaignList []campaignsModel
	for _, campaign := range campaigns.GetData() {
		campaignModel := campaignsModel{
			Id:                  types.StringValue(campaign.GetId()),
			Name:                types.StringValue(campaign.GetName()),
			Description:         types.StringValue(campaign.GetDescription()),
			Status:              types.StringValue(string(campaign.GetStatus())),
			Created:             types.StringValue(campaign.GetCreated().String()),
			CreatedBy:           types.StringValue(campaign.GetCreatedBy()),
			LastUpdated:         types.StringValue(campaign.GetLastUpdated().String()),
			LastUpdatedBy:       types.StringValue(campaign.GetLastUpdatedBy()),
			RecurringCampaignId: types.StringValue(campaign.GetRecurringCampaignId()),
		}
		fmt.Println("CAMPAIGN id", campaignModel.Id.ValueString())
		campaignList = append(campaignList, campaignModel)
	}
	data.Data = campaignList
	fmt.Println("campaign data size", len(data.Data))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildFilterForCampaigns(d campaignsDataSourceModel) string {
	if d.CampaignName.ValueString() != "" {
		return fmt.Sprintf("name eq \"%s\"", d.CampaignName.ValueString())
	} else if d.CampaignStatus.ValueString() != "" {
		return fmt.Sprintf("stauts eq \"%s\"", d.CampaignName.ValueString())
	} else if d.StartDate.ValueString() != "" {
		return fmt.Sprintf("startDate eq \"%s\"", d.CampaignName.ValueString())
	} else if d.EndDate.ValueString() != "" {
		return fmt.Sprintf("startDate eq \"%s\"", d.CampaignName.ValueString())
	} else if d.ScheduleTypeOneOff.ValueString() != "" {
		return fmt.Sprintf("scheduleType eq \"%s\"", d.ScheduleTypeOneOff.ValueString())
	} else if d.ScheduleTypeRecurring.ValueString() != "" {
		return fmt.Sprintf("scheduleType eq \"%s\"", d.ScheduleTypeRecurring.ValueString())
	} else if d.RecurringCampaignId.ValueString() != "" {
		return fmt.Sprintf("recurringCampaignId eq \"%s\"", d.RecurringCampaignId.ValueString())
	} else if d.ReviewerTypeUser.ValueString() != "" {
		return fmt.Sprintf("reviewerType eq \"%s\"", d.ReviewerTypeUser.ValueString())
	} else if d.ReviewerTypeGroup.ValueString() != "" {
		return fmt.Sprintf("reviewerType eq \"%s\"", d.ReviewerTypeGroup.ValueString())
	} else if d.ReviewerTypeGroupOwner.ValueString() != "" {
		return fmt.Sprintf("reviewerType eq \"%s\"", d.ReviewerTypeGroupOwner.ValueString())
	} else if d.ReviewerTypeMultiLevel.ValueString() != "" {
		return fmt.Sprintf("reviewerType eq \"%s\"", d.ReviewerTypeMultiLevel.ValueString())
	} else if d.MultipleCampaignStatus.Elements() != nil {
		var statuses []string
		// Convert types.List → []string
		if diags := d.MultipleCampaignStatus.ElementsAs(context.Background(), &statuses, false); diags.HasError() {
			// handle error, maybe return diag
			return ""
		}

		// Build the filter query
		var filterParts []string
		for _, status := range statuses {
			filterParts = append(filterParts, fmt.Sprintf(`status eq "%s"`, status))
		}

		return strings.Join(filterParts, " OR ")
	}
	return ""
}
