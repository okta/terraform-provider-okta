package governance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &delegateAppointmentsDataSource{}

func newDelegateAppointmentsDataSource() datasource.DataSource {
	return &delegateAppointmentsDataSource{}
}

type delegateAppointmentsDataSource struct {
	*config.Config
}

type delegateAppointmentsDataSourceModel struct {
	Id          types.String                         `tfsdk:"id"`
	PrincipalId types.String                         `tfsdk:"principal_id"`
	Data        []delegateAppointmentDataSourceModel `tfsdk:"data"`
}

type delegateAppointmentDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	DelegatorId   types.String `tfsdk:"delegator_id"`
	DelegatorType types.String `tfsdk:"delegator_type"`
	DelegateId    types.String `tfsdk:"delegate_id"`
	DelegateType  types.String `tfsdk:"delegate_type"`
	Note          types.String `tfsdk:"note"`
	StartTime     types.String `tfsdk:"start_time"`
	EndTime       types.String `tfsdk:"end_time"`
	CreatedBy     types.String `tfsdk:"created_by"`
	Created       types.String `tfsdk:"created"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	LastUpdatedBy types.String `tfsdk:"last_updated_by"`
}

func (d *delegateAppointmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *delegateAppointmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegate_appointments"
}

func (d *delegateAppointmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists delegate appointments in Okta Identity Governance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this data source, required by Terraform to track state.",
			},
			"principal_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Okta principal ID to filter delegate appointments by delegator. If not specified, all delegate appointments in the org are returned.",
			},
		},
		Blocks: map[string]schema.Block{
			"data": schema.ListNestedBlock{
				Description: "The list of delegate appointments.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the delegate appointment.",
						},
						"delegator_id": schema.StringAttribute{
							Computed:    true,
							Description: "The Okta ID of the delegator.",
						},
						"delegator_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the delegator principal.",
						},
						"delegate_id": schema.StringAttribute{
							Computed:    true,
							Description: "The Okta ID of the delegate.",
						},
						"delegate_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the delegate principal.",
						},
						"note": schema.StringAttribute{
							Computed:    true,
							Description: "A note that describes the delegate appointment.",
						},
						"start_time": schema.StringAttribute{
							Computed:    true,
							Description: "The start time of the delegate appointment in RFC3339 format.",
						},
						"end_time": schema.StringAttribute{
							Computed:    true,
							Description: "The end time of the delegate appointment in RFC3339 format.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The Okta user ID of the user who created the appointment.",
						},
						"created": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 formatted date and time when the appointment was created.",
						},
						"last_updated": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 formatted date and time when the appointment was last updated.",
						},
						"last_updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "The Okta user ID of the user who last updated the appointment.",
						},
					},
				},
			},
		},
	}
}

func (d *delegateAppointmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data delegateAppointmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := d.OktaGovernanceClient.OktaGovernanceSDKClient().DelegatesAPI.ListDelegateAppointments(ctx)
	if !data.PrincipalId.IsNull() && !data.PrincipalId.IsUnknown() {
		filter := fmt.Sprintf(`delegatorId eq "%s"`, data.PrincipalId.ValueString())
		apiReq = apiReq.Filter(filter)
	}

	listResp, _, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Delegate Appointments",
			"Could not read delegate appointments, unexpected error: "+err.Error(),
		)
		return
	}

	appointments := make([]delegateAppointmentDataSourceModel, 0, len(listResp.Data))
	for _, item := range listResp.Data {
		appt := delegateAppointmentDataSourceModel{
			Id:            types.StringValue(item.Id),
			DelegatorId:   types.StringValue(item.Delegator.ExternalId),
			DelegatorType: types.StringValue(string(item.Delegator.Type)),
			DelegateId:    types.StringValue(item.Delegate.ExternalId),
			DelegateType:  types.StringValue(string(item.Delegate.Type)),
			CreatedBy:     types.StringValue(item.CreatedBy),
			Created:       types.StringValue(item.Created.Format(time.RFC3339)),
			LastUpdated:   types.StringValue(item.LastUpdated.Format(time.RFC3339)),
			LastUpdatedBy: types.StringValue(item.LastUpdatedBy),
		}
		if item.Note != nil {
			appt.Note = types.StringValue(*item.Note)
		} else {
			appt.Note = types.StringNull()
		}
		if item.StartTime != nil {
			appt.StartTime = types.StringValue(item.StartTime.Format(time.RFC3339))
		} else {
			appt.StartTime = types.StringNull()
		}
		if item.EndTime != nil {
			appt.EndTime = types.StringValue(item.EndTime.Format(time.RFC3339))
		} else {
			appt.EndTime = types.StringNull()
		}
		appointments = append(appointments, appt)
	}

	data.Data = appointments
	data.Id = types.StringValue("delegate-appointments")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
