package governance

import (
	"context"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &campaignsDataSource{}

func newCampaignsDataSource() datasource.DataSource {
	return &campaignsDataSource{}
}

type campaignsDataSource struct {
	*config.Config
}

type campaignsDataSourceModel struct {
	Campaigns types.List   `tfsdk:"campaigns"`
	Filter    types.String `tfsdk:"filter"`
	Limit     types.Int64  `tfsdk:"limit"`
	OrderBy   types.List   `tfsdk:"order_by"`
}

func (d *campaignsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_campaigns"
}

func (d *campaignsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *campaignsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all or a subset of campaigns in your organization. Use the `filter` parameter to narrow results by name, status, scheduleType, reviewerType, and recurringCampaignId. Use `order_by` to sort by name, created, startDate, endDate, or status.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "Filter expression to narrow results. Supported campaign properties: name, status, scheduleType, reviewerType, recurringCampaignId.",
			},
			"limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of campaigns to return per page. Pagination is handled automatically to return all campaigns when not set.",
			},
			"order_by": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Order results by campaign properties: name, created, startDate, endDate, status. By default results are sorted by created.",
			},
			"campaigns": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the campaign.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the campaign.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Human readable description.",
						},
						"start_date": schema.StringAttribute{
							Computed:    true,
							Description: "The date on which the campaign is scheduled to start (ISO 8601).",
						},
						"end_date": schema.StringAttribute{
							Computed:    true,
							Description: "The date on which the campaign is supposed to end (ISO 8601).",
						},
						"schedule_type": schema.StringAttribute{
							Computed:    true,
							Description: "Schedule type of the campaign (e.g. ONE_OFF, RECURRING).",
						},
						"recurring_campaign_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the recurring campaign if this campaign was created as part of a recurring schedule.",
						},
						"reviewer_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of reviewer for the campaign.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the campaign.",
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
					},
				},
				Description: "List of campaigns.",
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

	apiReq := d.OktaGovernanceClient.OktaGovernanceSDKClient().CampaignsAPI.ListCampaigns(ctx)

	if !data.Filter.IsNull() && !data.Filter.IsUnknown() {
		apiReq = apiReq.Filter(data.Filter.ValueString())
	}
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() && data.Limit.ValueInt64() > 0 {
		apiReq = apiReq.Limit(int32(data.Limit.ValueInt64()))
	}
	if !data.OrderBy.IsNull() && !data.OrderBy.IsUnknown() {
		var orderByElems []types.String
		resp.Diagnostics.Append(data.OrderBy.ElementsAs(ctx, &orderByElems, false)...)
		if !resp.Diagnostics.HasError() && len(orderByElems) > 0 {
			orderBy := make([]string, 0, len(orderByElems))
			for _, s := range orderByElems {
				orderBy = append(orderBy, s.ValueString())
			}
			apiReq = apiReq.OrderBy(orderBy)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var allCampaigns []governance.CampaignSparse
	after := ""

	for {
		pageReq := apiReq
		if after != "" {
			pageReq = pageReq.After(after)
		}

		list, _, err := pageReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing campaigns",
				"Could not list campaigns: "+err.Error(),
			)
			return
		}

		page := list.GetData()
		allCampaigns = append(allCampaigns, page...)

		if len(page) == 0 {
			break
		}
		links := list.GetLinks()
		if !links.HasNext() {
			break
		}
		nextLink := links.GetNext()
		nextHref := nextLink.GetHref()
		after = parseAfterFromURL(nextHref)
		if after == "" {
			break
		}
	}

	// Build state
	attrTypes := campaignListItemAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	campaignElems := make([]attr.Value, 0, len(allCampaigns))
	for _, c := range allCampaigns {
		attrs := map[string]attr.Value{
			"id":              types.StringValue(c.GetId()),
			"name":            types.StringValue(c.GetName()),
			"schedule_type":   types.StringValue(string(c.GetScheduleType())),
			"reviewer_type":   types.StringValue(string(c.GetReviewerType())),
			"status":          types.StringValue(string(c.GetStatus())),
			"created_by":      types.StringValue(c.GetCreatedBy()),
			"last_updated_by": types.StringValue(c.GetLastUpdatedBy()),
		}
		attrs["start_date"] = timeToAttrValue(c.GetStartDateOk())
		attrs["end_date"] = timeToAttrValue(c.GetEndDateOk())
		attrs["created"] = timeToAttrValue(c.GetCreatedOk())
		attrs["last_updated"] = timeToAttrValue(c.GetLastUpdatedOk())
		attrs["description"] = optionalStringToAttrValue(c.GetDescription(), c.HasDescription())
		attrs["recurring_campaign_id"] = optionalStringToAttrValue(c.GetRecurringCampaignId(), c.HasRecurringCampaignId())
		obj, diags := types.ObjectValue(objType.AttrTypes, attrs)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		campaignElems = append(campaignElems, obj)
	}

	campaignsList, diags := types.ListValue(objType, campaignElems)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Campaigns = campaignsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// timeToAttrValue converts an optional API time (e.g. from GetStartDateOk()) to a
// Terraform attribute value: RFC3339 string when present, null when absent.
func timeToAttrValue(t *time.Time, ok bool) attr.Value {
	if !ok || t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(time.RFC3339))
}

// optionalStringToAttrValue converts an optional API string (e.g. from GetDescription() + HasDescription()) to a
// Terraform attribute value: string when present, null when absent.
func optionalStringToAttrValue(value string, ok bool) attr.Value {
	if !ok {
		return types.StringNull()
	}
	return types.StringValue(value)
}

// parseAfterFromURL extracts the "after" query parameter from a next-page href.
func parseAfterFromURL(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return u.Query().Get("after")
}

func campaignListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                    types.StringType,
		"name":                  types.StringType,
		"description":           types.StringType,
		"start_date":            types.StringType,
		"end_date":              types.StringType,
		"schedule_type":         types.StringType,
		"recurring_campaign_id": types.StringType,
		"reviewer_type":         types.StringType,
		"status":                types.StringType,
		"created":               types.StringType,
		"created_by":            types.StringType,
		"last_updated":          types.StringType,
		"last_updated_by":       types.StringType,
	}
}
