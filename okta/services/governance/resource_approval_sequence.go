package governance

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &approvalSequenceResource{}
	_ resource.ResourceWithConfigure   = &approvalSequenceResource{}
	_ resource.ResourceWithImportState = &approvalSequenceResource{}
)

func newApprovalSequenceResource() resource.Resource {
	return &approvalSequenceResource{}
}

func (r *approvalSequenceResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	parts := strings.Split(request.ID, "/")
	if len(parts) != 2 {
		response.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: resource_id/sequence_id",
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("resource_id"), parts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *approvalSequenceResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type approvalSequenceResource struct {
	*config.Config
}

type approvalSequenceResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	ResourceId              types.String `tfsdk:"resource_id"`
	Description             types.String `tfsdk:"description"`
	Link                    types.String `tfsdk:"link"`
	Name                    types.String `tfsdk:"name"`
	CompatibleResourceTypes types.List   `tfsdk:"compatible_resource_types"`
}

func (r *approvalSequenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_approval_sequence"
}

func approvalSequenceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages an Okta Approval Sequence (also known as Request Sequence in the Okta API).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the approval sequence.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the approval sequence.",
			},
			"link": schema.StringAttribute{
				Computed:    true,
				Description: "Link to edit the approval sequence.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the approval sequence.",
			},
			"compatible_resource_types": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("APP", "GROUP"),
					),
				},
			},
		},
	}
}

func (r *approvalSequenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = approvalSequenceResourceSchema()
}

func (r *approvalSequenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning(
		"Create Not Supported",
		"This resource cannot be created via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *approvalSequenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data approvalSequenceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.GetResourceRequestSequenceV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Approval Sequence",
			"Could not read Approval Sequence, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(readResp.Id)
	data.Name = types.StringValue(readResp.Name)
	data.ResourceId = types.StringValue(data.ResourceId.ValueString())
	data.Description = types.StringValue(readResp.Description)
	data.Link = types.StringValue(readResp.Link)
	data.CompatibleResourceTypes = setCompatibleResourceType(readResp.CompatibleResourceTypes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func setCompatibleResourceType(resourceTypes []governance.CompatibleResourceTypes) types.List {
	values := make([]attr.Value, len(resourceTypes))
	for i, resourceType := range resourceTypes {
		values[i] = types.StringValue(string(resourceType))
	}
	listVal, _ := types.ListValue(types.StringType, values)
	return listVal
}

func (r *approvalSequenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource cannot be updated via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *approvalSequenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data approvalSequenceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.DeleteRequestSequenceV2(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Approval Sequence",
			"Could not delete Approval Sequence, unexpected error: "+err.Error(),
		)
		return
	}
}
