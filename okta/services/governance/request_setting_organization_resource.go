package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*requestSettingOrganizationResource)(nil)

func NewRequestSettingOrganizationResource() resource.Resource {
	return &requestSettingOrganizationResource{}
}

type requestSettingOrganizationResource struct {
	*config.Config
}

type requestSettingOrganizationResourceModel struct {
	LongTimePastProvisioned   types.Bool   `tfsdk:"long_time_past_provisioned"`
	ProvisioningStatus        types.String `tfsdk:"provisioning_status"`
	RequestExperiences        types.List   `tfsdk:"request_experiences"`
	SubprocessorsAcknowledged types.Bool   `tfsdk:"subprocessors_acknowledged"`
}

func (r *requestSettingOrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_organization"
}

func (r *requestSettingOrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"long_time_past_provisioned": schema.StringAttribute{
				Computed: true,
			},
			"provisioning_status": schema.StringAttribute{
				Computed: true,
			},
			"request_experiences": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"subprocessors_acknowledged": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *requestSettingOrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning(
		"Create Not Supported",
		"This resource cannot be created via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *requestSettingOrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestSettingOrganizationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readOrgRequestSettingResp, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().RequestSettingsAPI.GetOrgRequestSettingsV2(ctx).Execute()
	if err != nil {
		return
	}
	data.LongTimePastProvisioned = types.BoolValue(readOrgRequestSettingResp.GetLongTimePastProvisioned())
	data.ProvisioningStatus = types.StringValue(string(readOrgRequestSettingResp.GetProvisioningStatus()))
	data.SubprocessorsAcknowledged = types.BoolValue(readOrgRequestSettingResp.SubprocessorsAcknowledged)
	var reqExpVals []attr.Value
	for _, reqExp := range readOrgRequestSettingResp.GetRequestExperiences() {
		reqExpVals = append(reqExpVals, types.StringValue(string(reqExp)))
	}

	listVal, _ := types.ListValue(types.StringType, reqExpVals)
	data.RequestExperiences = listVal

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestSettingOrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data requestSettingOrganizationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	updateOrgSettingsResp, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().RequestSettingsAPI.UpdateOrgRequestSettingsV2(ctx).OrgRequestSettingsPatchable(createOrgRequestsSettings(data)).Execute()
	if err != nil {
		return
	}
	data.SubprocessorsAcknowledged = types.BoolValue(updateOrgSettingsResp.GetSubprocessorsAcknowledged())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func createOrgRequestsSettings(data requestSettingOrganizationResourceModel) governance.OrgRequestSettingsPatchable {
	return governance.OrgRequestSettingsPatchable{
		SubprocessorsAcknowledged: data.SubprocessorsAcknowledged.ValueBool(),
	}
}

func (r *requestSettingOrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}
