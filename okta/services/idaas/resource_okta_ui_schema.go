package idaas

import (
	"context"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"

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
	ID       types.String `tfsdk:"id"`
	UISchema *uiSchema    `tfsdk:"ui_schema"`
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
	createUISchemaResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UISchemaAPI.CreateUISchema(ctx).Uischemabody(createUISchemaBody(data)).Execute()
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

func applyUISchemaToState(resp *v6okta.UISchemasResponseObject, data *uiSchemaResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	data.ID = types.StringValue(resp.GetId())
	data.UISchema = &uiSchema{}
	if resp.UiSchema.ButtonLabel != nil {
		data.UISchema.ButtonLabel = types.StringValue(resp.UiSchema.GetButtonLabel())
	}
	if resp.UiSchema.Label != nil {
		data.UISchema.Label = types.StringValue(resp.UiSchema.GetLabel())
	}

	if resp.UiSchema.Type != nil && resp.UiSchema.GetType() != "" {
		data.UISchema.Type = types.StringValue(resp.UiSchema.GetType())
	}

	if resp.UiSchema.Elements != nil {
		var elems []elements
		for _, elem := range resp.UiSchema.Elements {
			e := elements{}
			e.Label = types.StringValue(elem.GetLabel())
			e.Scope = types.StringValue(elem.GetScope())
			e.Type = types.StringValue(elem.GetType())
			if elem.Options != nil {
				e.Options = &options{}
				e.Options.Format = types.StringValue(elem.Options.GetFormat())
			}
			elems = append(elems, e)
		}
		data.UISchema.Elements = elems
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
	readUISchemaResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UISchemaAPI.GetUISchema(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading UISchema",
			"Could not read UISchema, unexpected error: "+err.Error(),
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
	data.ID = state.ID
	// Update API call logic
	updateUISchemaResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().UISchemaAPI.ReplaceUISchemas(ctx, state.ID.ValueString()).UpdateUISchemaBody(updateUISchemaBody(data)).Execute()
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
	_, err := r.OktaIDaaSClient.OktaSDKClientV6().UISchemaAPI.DeleteUISchemas(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting UISchema",
			"Could not delete UISchema, unexpected error: "+err.Error(),
		)
		return
	}
}

func updateUISchemaBody(data uiSchemaResourceModel) v6okta.UpdateUISchema {
	uiSchema := &v6okta.UISchemaObject{
		ButtonLabel: data.UISchema.ButtonLabel.ValueStringPointer(),
		Label:       data.UISchema.Label.ValueStringPointer(),
		Type:        data.UISchema.Type.ValueStringPointer(),
	}
	var elements []v6okta.UIElement
	for _, elem := range data.UISchema.Elements {
		options := v6okta.UIElementOptions{}
		options.Format = elem.Options.Format.ValueStringPointer()
		element := v6okta.UIElement{}
		element.Label = elem.Label.ValueStringPointer()
		element.Scope = elem.Scope.ValueStringPointer()
		element.Type = elem.Type.ValueStringPointer()
		element.Options = &options
		elements = append(elements, element)
	}

	uiSchema.Elements = elements

	return v6okta.UpdateUISchema{UiSchema: uiSchema}
}

func createUISchemaBody(data uiSchemaResourceModel) v6okta.CreateUISchema {
	uiSchema := &v6okta.UISchemaObject{
		ButtonLabel: data.UISchema.ButtonLabel.ValueStringPointer(),
		Label:       data.UISchema.Label.ValueStringPointer(),
		Type:        data.UISchema.Type.ValueStringPointer(),
	}
	var elements []v6okta.UIElement
	for _, elem := range data.UISchema.Elements {
		options := v6okta.UIElementOptions{}
		options.Format = elem.Options.Format.ValueStringPointer()
		element := v6okta.UIElement{}
		element.Label = elem.Label.ValueStringPointer()
		element.Scope = elem.Scope.ValueStringPointer()
		element.Type = elem.Type.ValueStringPointer()
		element.Options = &options
		elements = append(elements, element)
	}

	uiSchema.Elements = elements

	return v6okta.CreateUISchema{
		UiSchema: uiSchema,
	}
}
