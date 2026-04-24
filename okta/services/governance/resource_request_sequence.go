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
	_ resource.Resource                = &requestSequenceResource{}
	_ resource.ResourceWithConfigure   = &requestSequenceResource{}
	_ resource.ResourceWithImportState = &requestSequenceResource{}
)

func newRequestSequenceResource() resource.Resource {
	return &requestSequenceResource{}
}

func (r *requestSequenceResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
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

func (r *requestSequenceResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type requestSequenceResource struct {
	*config.Config
}

type requestSequenceResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	ResourceId              types.String `tfsdk:"resource_id"`
	Description             types.String `tfsdk:"description"`
	Link                    types.String `tfsdk:"link"`
	Name                    types.String `tfsdk:"name"`
	CompatibleResourceTypes types.List   `tfsdk:"compatible_resource_types"`
}

func (r *requestSequenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_sequence"
}

func (r *requestSequenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the request sequence. This is typically the sequence ID in Okta.",
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the request sequence",
			},
			"link": schema.StringAttribute{
				Computed:    true,
				Description: "Link to edit the request sequence.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the request sequence.",
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

func (r *requestSequenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning(
		"Create Not Supported",
		"This resource cannot be created via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *requestSequenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestSequenceResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readResourceRequestSeqResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.GetResourceRequestSequenceV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Sequence",
			"Could not read Request Sequence, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(readResourceRequestSeqResp.Id)
	data.Name = types.StringValue(readResourceRequestSeqResp.Name)
	data.ResourceId = types.StringValue(data.ResourceId.ValueString())
	data.Description = types.StringValue(readResourceRequestSeqResp.Description)
	data.Link = types.StringValue(readResourceRequestSeqResp.Link)
	data.CompatibleResourceTypes = setCompatibleResourceType(readResourceRequestSeqResp.CompatibleResourceTypes)

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

func (r *requestSequenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource cannot be updated via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *requestSequenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestSequenceResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSequencesAPI.DeleteRequestSequenceV2(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Request Sequence",
			"Could not delete Request Sequence, unexpected error: "+err.Error(),
		)
		return

	}
}
