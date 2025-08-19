package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &requestSettingOrganizationResource{}
	_ resource.ResourceWithConfigure   = &requestSettingOrganizationResource{}
	_ resource.ResourceWithImportState = &requestSettingOrganizationResource{}
)

func newRequestSettingOrganizationResource() resource.Resource {
	return &requestSettingOrganizationResource{}
}

func (r *requestSettingOrganizationResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *requestSettingOrganizationResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type requestSettingOrganizationResource struct {
	*config.Config
}

type experience struct {
	ExperienceType string `tfsdk:"experience_type"`
}

type requestSettingOrganizationResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	SubprocessorsAcknowledged types.Bool   `tfsdk:"subprocessors_acknowledged"`
}

func (r *requestSettingOrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_organization"
}

func (r *requestSettingOrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"subprocessors_acknowledged": schema.BoolAttribute{
				Optional: true,
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

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readOrgRequestSettingResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSettingsAPI.GetOrgRequestSettingsV2(ctx).Execute()
	if err != nil {
		return
	}
	data.SubprocessorsAcknowledged = types.BoolValue(readOrgRequestSettingResp.SubprocessorsAcknowledged)
	var experiences []experience
	for _, exp := range readOrgRequestSettingResp.GetRequestExperiences() {
		experiences = append(experiences, experience{
			ExperienceType: string(exp),
		})
	}

	data.Id = types.StringValue("default")

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestSettingOrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data requestSettingOrganizationResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	updateOrgSettingsResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSettingsAPI.UpdateOrgRequestSettingsV2(ctx).OrgRequestSettingsPatchable(createOrgRequestsSettings(data)).Execute()
	if err != nil {
		return
	}
	data.SubprocessorsAcknowledged = types.BoolValue(updateOrgSettingsResp.GetSubprocessorsAcknowledged())

	// Save updated Data into Terraform state
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
		"This resource cannot be deleted via Terraform.",
	)
}
