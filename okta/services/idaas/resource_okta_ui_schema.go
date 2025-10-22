package idaas

import (
	"context"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &uiSchemaResource{}
	_ resource.ResourceWithConfigure   = &uiSchemaResource{}
	_ resource.ResourceWithImportState = &uiSchemaResource{}
)

func newUISchemaResource() resource.Resource {
	return &uiSchemaResource{}
}

func (r *uiSchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ui_schema"
}

func (r *uiSchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *uiSchemaResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type uiSchemaResource struct {
	*config.Config
}

type options struct {
	Format types.String `tfsdk:"format"`
}

type elements struct {
	Label   types.String `tfsdk:"label"`
	Scope   types.String `tfsdk:"scope"`
	Type    types.String `tfsdk:"type"`
	Options *options     `tfsdk:"options"`
}

type uiSchema struct {
	ButtonLabel types.String `tfsdk:"button_label"`
	Label       types.String `tfsdk:"label"`
	Type        types.String `tfsdk:"type"`
	Elements    []elements   `tfsdk:"elements"`
}

type uiSchemaResourceModel struct {
	Id       types.String `tfsdk:"id"`
	UiSchema *uiSchema    `tfsdk:"ui_schema"`
}

func (r *uiSchemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id property of an entitlement.",
			},
		},
		Blocks: map[string]schema.Block{
			"ui_schema": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"button_label": schema.StringAttribute{
						Required:    true,
						Description: "The Okta app.id of the resource.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of resource.",
					},
					"label": schema.StringAttribute{
						Required:    true,
						Description: "The type of resource.",
					},
				},
				Blocks: map[string]schema.Block{
					"elements": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"label": schema.StringAttribute{
									Required:    true,
									Description: "The label of the element.",
								},
								"scope": schema.StringAttribute{
									Required:    true,
									Description: "The scope of the element.",
								},
								"type": schema.StringAttribute{
									Required:    true,
									Description: "The type of the element.",
								},
							},
							Blocks: map[string]schema.Block{
								"options": schema.SingleNestedBlock{
									Attributes: map[string]schema.Attribute{
										"format": schema.StringAttribute{
											Required:    true,
											Description: "The format of the option.",
										},
									},
								},
							},
						},
					},
				},
				Description: "Representation of a resource.",
			},
		},
	}
}

func (r *uiSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data uiSchemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	createUISchemaResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().UISchemaAPI.CreateUISchema(ctx).Uischemabody(createUISchemaBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating UISchema",
			"Could not create UISchema, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyUISchemaToState(createUISchemaResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyUISchemaToState(resp *v5okta.UISchemasResponseObject, data *uiSchemaResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	data.Id = types.StringValue(resp.GetId())
	data.UiSchema.ButtonLabel = types.StringValue(resp.UiSchema.GetButtonLabel())
	data.UiSchema.Label = types.StringValue(resp.UiSchema.GetLabel())
	data.UiSchema.Type = types.StringValue(resp.UiSchema.GetType())
	if resp.UiSchema.Elements != nil {
		data.UiSchema.Elements[0] = elements{}
		data.UiSchema.Elements[0].Label = types.StringValue(resp.UiSchema.Elements.GetLabel())
		data.UiSchema.Elements[0].Scope = types.StringValue(resp.UiSchema.Elements.GetScope())
		data.UiSchema.Elements[0].Type = types.StringValue(resp.UiSchema.Elements.GetType())
		if resp.UiSchema.Elements.Options != nil {
			data.UiSchema.Elements[0].Options = &options{}
			data.UiSchema.Elements[0].Options.Format = types.StringValue(resp.UiSchema.Elements.Options.GetFormat())
		}
	}
	return diags
}

func (r *uiSchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data uiSchemaResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readUISchemaResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().UISchemaAPI.GetUISchema(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Entitlement",
			"Could not create Entitlement, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyUISchemaToState(readUISchemaResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *uiSchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data uiSchemaResourceModel
	var state uiSchemaResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id
	// Update API call logic
	updateUISchemaResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().UISchemaAPI.ReplaceUISchemas(ctx, state.Id.ValueString()).UpdateUISchemaBody(updateUISchemaBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating UISchema",
			"An error occurred while updating the UISchema: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(applyUISchemaToState(updateUISchemaResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *uiSchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data uiSchemaResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaIDaaSClient.OktaSDKClientV5().UISchemaAPI.DeleteUISchemas(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting UISchema",
			"Could not delete UISchema, unexpected error: "+err.Error(),
		)
		return
	}
}

func updateUISchemaBody(data uiSchemaResourceModel) v5okta.UpdateUISchema {
	uiSchema := &v5okta.UISchemaObject{
		ButtonLabel: data.UiSchema.ButtonLabel.ValueStringPointer(),
		Label:       data.UiSchema.Label.ValueStringPointer(),
		Type:        data.UiSchema.Type.ValueStringPointer(),
	}
	options := v5okta.UIElementOptions{}
	options.Format = data.UiSchema.ButtonLabel.ValueStringPointer()
	elements := &v5okta.UIElement{}
	elements.Label = data.UiSchema.Label.ValueStringPointer()
	elements.Scope = data.UiSchema.Elements[0].Scope.ValueStringPointer()
	elements.Type = data.UiSchema.Elements[0].Type.ValueStringPointer()
	elements.Options = &options
	uiSchema.Elements = elements

	return v5okta.UpdateUISchema{UiSchema: uiSchema}
}

func createUISchemaBody(data uiSchemaResourceModel) v5okta.CreateUISchema {
	uiSchema := &v5okta.UISchemaObject{
		ButtonLabel: data.UiSchema.ButtonLabel.ValueStringPointer(),
		Label:       data.UiSchema.Label.ValueStringPointer(),
		Type:        data.UiSchema.Type.ValueStringPointer(),
	}
	options := v5okta.UIElementOptions{}
	options.Format = data.UiSchema.ButtonLabel.ValueStringPointer()
	elements := &v5okta.UIElement{}
	elements.Label = data.UiSchema.Label.ValueStringPointer()
	elements.Scope = data.UiSchema.Elements[0].Scope.ValueStringPointer()
	elements.Type = data.UiSchema.Elements[0].Type.ValueStringPointer()
	elements.Options = &options
	uiSchema.Elements = elements

	return v5okta.CreateUISchema{
		UiSchema: uiSchema,
	}
}
