package governance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

const defaultPrincipalType = "OKTA_USER"

var (
	_ resource.Resource                = &delegateAppointmentsResource{}
	_ resource.ResourceWithConfigure   = &delegateAppointmentsResource{}
	_ resource.ResourceWithImportState = &delegateAppointmentsResource{}
)

func newDelegateAppointmentsResource() resource.Resource {
	return &delegateAppointmentsResource{}
}

type delegateAppointmentsResource struct {
	*config.Config
}

type delegateAppointmentsResourceModel struct {
	Id            types.String                    `tfsdk:"id"`
	PrincipalId   types.String                    `tfsdk:"principal_id"`
	PrincipalType types.String                    `tfsdk:"principal_type"`
	Appointments  []delegateAppointmentBlockModel `tfsdk:"appointments"`
}

type delegateAppointmentBlockModel struct {
	DelegateId types.String `tfsdk:"delegate_id"`
	Note       types.String `tfsdk:"note"`
	StartTime  types.String `tfsdk:"start_time"`
	EndTime    types.String `tfsdk:"end_time"`
}

func (r *delegateAppointmentsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegate_appointments"
}

func (r *delegateAppointmentsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages delegate appointments for a principal in Okta Identity Governance. " +
			"This resource represents settings that always exist in Okta for a given principal. " +
			"Creating this resource adopts management of the principal's delegate appointments, " +
			"and destroying it resets the appointments to an empty state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this resource, required by Terraform to track state.",
			},
			"principal_id": schema.StringAttribute{
				Required:    true,
				Description: "The Okta ID of the principal whose delegate appointments are being managed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(defaultPrincipalType),
				Description: fmt.Sprintf("The type of principal. Defaults to %q.", defaultPrincipalType),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"appointments": schema.ListNestedBlock{
				Description: "The list of delegate appointments for this principal. The API currently supports a maximum of one appointment per principal.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"delegate_id": schema.StringAttribute{
							Required:    true,
							Description: "The Okta user ID of the delegate.",
						},
						"note": schema.StringAttribute{
							Optional:    true,
							Description: "A note that describes the delegate appointment.",
						},
						"start_time": schema.StringAttribute{
							Optional:    true,
							Description: "The start time of the delegate appointment in RFC3339 format.",
						},
						"end_time": schema.StringAttribute{
							Optional:    true,
							Description: "The end time of the delegate appointment in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func (r *delegateAppointmentsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *delegateAppointmentsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("principal_id"), req, resp)
}

func (r *delegateAppointmentsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data delegateAppointmentsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	patchable, diags := buildDelegateAppointmentsPatchable(data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	principalId := data.PrincipalId.ValueString()

	tflog.Debug(ctx, "creating delegate appointments", map[string]interface{}{
		"principal_id":      principalId,
		"appointment_count": len(patchable.Delegates.Appointments),
	})

	settingsResp, apiResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().PrincipalSettingsAPI.UpdatePrincipalSettings(ctx, principalId).PrincipalSettingsPatchable(patchable).Execute()
	if err != nil {
		errMsg := err.Error()
		if apiResp != nil && apiResp.Body != nil {
			defer apiResp.Body.Close()
			if body, readErr := io.ReadAll(apiResp.Body); readErr == nil && len(body) > 0 {
				errMsg = fmt.Sprintf("%s — response body: %s", errMsg, string(body))
			}
		}
		resp.Diagnostics.AddError(
			"Error creating Delegate Appointments",
			fmt.Sprintf("Could not create delegate appointments for principal %s: %s", principalId, errMsg),
		)
		return
	}

	applyDelegateAppointmentsResponseToState(&data, settingsResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *delegateAppointmentsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data delegateAppointmentsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	principalId := data.PrincipalId.ValueString()

	filter := fmt.Sprintf(`delegatorId eq "%s"`, principalId)
	listResp, apiResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().DelegatesAPI.ListDelegateAppointments(ctx).Filter(filter).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Delegate Appointments",
			fmt.Sprintf("Could not read delegate appointments for principal %s: %s", principalId, err.Error()),
		)
		return
	}

	allItems := listResp.Data
	for apiResp.HasNextPage() {
		var nextPage governance.DelegateAppointmentList
		apiResp, err = apiResp.Next(&nextPage)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading Delegate Appointments",
				fmt.Sprintf("Could not read next page of delegate appointments for principal %s: %s", principalId, err.Error()),
			)
			return
		}
		allItems = append(allItems, nextPage.Data...)
	}

	applyDelegateAppointmentListToState(&data, allItems)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *delegateAppointmentsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data delegateAppointmentsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	patchable, diags := buildDelegateAppointmentsPatchable(data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	principalId := data.PrincipalId.ValueString()

	tflog.Debug(ctx, "updating delegate appointments", map[string]interface{}{
		"principal_id":      principalId,
		"appointment_count": len(patchable.Delegates.Appointments),
	})

	settingsResp, apiResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().PrincipalSettingsAPI.UpdatePrincipalSettings(ctx, principalId).PrincipalSettingsPatchable(patchable).Execute()
	if err != nil {
		errMsg := err.Error()
		if apiResp != nil && apiResp.Body != nil {
			defer apiResp.Body.Close()
			if body, readErr := io.ReadAll(apiResp.Body); readErr == nil && len(body) > 0 {
				errMsg = fmt.Sprintf("%s — response body: %s", errMsg, string(body))
			}
		}
		resp.Diagnostics.AddError(
			"Error updating Delegate Appointments",
			fmt.Sprintf("Could not update delegate appointments for principal %s: %s", principalId, errMsg),
		)
		return
	}

	applyDelegateAppointmentsResponseToState(&data, settingsResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *delegateAppointmentsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data delegateAppointmentsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	principalId := data.PrincipalId.ValueString()
	patchable := governance.PrincipalSettingsPatchable{
		Delegates: &governance.DelegatesPatchable{
			Appointments: []governance.DelegatePatchable{},
		},
	}

	_, apiResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().PrincipalSettingsAPI.UpdatePrincipalSettings(ctx, principalId).PrincipalSettingsPatchable(patchable).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting Delegate Appointments",
			fmt.Sprintf("Could not remove delegate appointments for principal %s: %s", principalId, err.Error()),
		)
		return
	}
}

func buildDelegateAppointmentsPatchable(data delegateAppointmentsResourceModel) (governance.PrincipalSettingsPatchable, diag.Diagnostics) {
	var diags diag.Diagnostics

	appointments := make([]governance.DelegatePatchable, 0, len(data.Appointments))
	for _, appt := range data.Appointments {
		dp := governance.DelegatePatchable{
			Delegate: governance.DelegateAppointmentDelegate{
				ExternalId: appt.DelegateId.ValueString(),
				Type:       governance.PRINCIPALTYPE_OKTA_USER,
			},
		}
		if !appt.Note.IsNull() && !appt.Note.IsUnknown() {
			note := appt.Note.ValueString()
			dp.Note = &note
		}
		if !appt.StartTime.IsNull() && !appt.StartTime.IsUnknown() {
			t, err := time.Parse(time.RFC3339, appt.StartTime.ValueString())
			if err != nil {
				diags.AddAttributeError(
					path.Root("appointments"),
					"Invalid start_time",
					fmt.Sprintf("Could not parse start_time %q as RFC3339: %s", appt.StartTime.ValueString(), err.Error()),
				)
				continue
			}
			dp.StartTime = &t
		}
		if !appt.EndTime.IsNull() && !appt.EndTime.IsUnknown() {
			t, err := time.Parse(time.RFC3339, appt.EndTime.ValueString())
			if err != nil {
				diags.AddAttributeError(
					path.Root("appointments"),
					"Invalid end_time",
					fmt.Sprintf("Could not parse end_time %q as RFC3339: %s", appt.EndTime.ValueString(), err.Error()),
				)
				continue
			}
			dp.EndTime = &t
		}
		appointments = append(appointments, dp)
	}

	return governance.PrincipalSettingsPatchable{
		Delegates: &governance.DelegatesPatchable{
			Appointments: appointments,
		},
	}, diags
}

func applyDelegateAppointmentsResponseToState(data *delegateAppointmentsResourceModel, settingsResp *governance.PrincipalSettings) {
	data.Id = data.PrincipalId
	if data.PrincipalType.IsNull() || data.PrincipalType.IsUnknown() {
		data.PrincipalType = types.StringValue(defaultPrincipalType)
	}

	if settingsResp.Delegates == nil {
		data.Appointments = []delegateAppointmentBlockModel{}
		return
	}

	apiAppointments := settingsResp.Delegates.GetAppointments()
	reorderAppointmentsToMatchPlan(data, apiAppointments)
}

func applyDelegateAppointmentListToState(data *delegateAppointmentsResourceModel, items []governance.DelegateAppointment) {
	data.Id = data.PrincipalId
	if data.PrincipalType.IsNull() || data.PrincipalType.IsUnknown() {
		data.PrincipalType = types.StringValue(defaultPrincipalType)
	}

	if len(items) == 0 {
		data.Appointments = []delegateAppointmentBlockModel{}
		return
	}

	reorderAppointmentsToMatchPlan(data, items)
}

func reorderAppointmentsToMatchPlan(data *delegateAppointmentsResourceModel, apiAppointments []governance.DelegateAppointment) {
	apiByDelegateId := make(map[string]governance.DelegateAppointment, len(apiAppointments))
	for _, a := range apiAppointments {
		apiByDelegateId[a.Delegate.ExternalId] = a
	}

	matched := make(map[string]bool, len(data.Appointments))
	result := make([]delegateAppointmentBlockModel, 0, len(apiAppointments))
	for _, planned := range data.Appointments {
		delegateId := planned.DelegateId.ValueString()
		if a, ok := apiByDelegateId[delegateId]; ok {
			result = append(result, appointmentFromAPI(a))
			matched[delegateId] = true
		}
	}

	for _, a := range apiAppointments {
		if !matched[a.Delegate.ExternalId] {
			result = append(result, appointmentFromAPI(a))
		}
	}

	data.Appointments = result
}

func appointmentFromAPI(a governance.DelegateAppointment) delegateAppointmentBlockModel {
	m := delegateAppointmentBlockModel{
		DelegateId: types.StringValue(a.Delegate.ExternalId),
	}
	if a.Note != nil {
		m.Note = types.StringValue(*a.Note)
	} else {
		m.Note = types.StringNull()
	}
	if a.StartTime != nil {
		m.StartTime = types.StringValue(a.StartTime.Format(time.RFC3339))
	} else {
		m.StartTime = types.StringNull()
	}
	if a.EndTime != nil {
		m.EndTime = types.StringValue(a.EndTime.Format(time.RFC3339))
	} else {
		m.EndTime = types.StringNull()
	}
	return m
}
