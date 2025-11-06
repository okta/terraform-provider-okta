package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &hookKey{}
	_ resource.ResourceWithConfigure   = &hookKey{}
	_ resource.ResourceWithImportState = &hookKey{}
)

type hookKey struct {
	*config.Config
}

type hookKeyModel struct {
	Id          types.String `tfsdk:"id"`
	KeyId       types.String `tfsdk:"key_id"`
	Name        types.String `tfsdk:"name"`
	Created     types.String `tfsdk:"created"`
	LastUpdated types.String `tfsdk:"last_updated"`
	IsUsed      types.Bool   `tfsdk:"is_used"`
}

// HookKeyResponse represents the response from Okta Hook Key API
type HookKeyResponse struct {
	Id          string    `json:"id"`
	KeyId       string    `json:"keyId"`
	Name        string    `json:"name"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
	IsUsed      string    `json:"isUsed"`
}

// HookKeyRequest represents the request to create/update Hook Key
type HookKeyRequest struct {
	Name string `json:"name"`
}

func newHookKeyResource() resource.Resource {
	return &hookKey{}
}

func (r *hookKey) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hook_key"
}

func (r *hookKey) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *hookKey) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *hookKey) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a Hook Key for use with Okta Inline Hooks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique Okta ID of this key record",
			},
			"key_id": schema.StringAttribute{
				Computed:    true,
				Description: "The alias of the public key.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
				Description: "Display name for the key.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the key was created.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the key was updated.",
			},
			"is_used": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this key is currently in use by other applications.",
			},
		},
	}
}

func (r *hookKey) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data hookKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) // Read Terraform plan data into the model
	if resp.Diagnostics.HasError() {
		return
	}

	hookKeyRequest := r.OktaIDaaSClient.OktaSDKClientV5().HookKeyAPI.CreateHookKey(ctx)
	keyRequest := v5okta.KeyRequest{}
	keyRequest.SetName(data.Name.ValueString())
	hookKey, _, err := hookKeyRequest.KeyRequest(keyRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating hook key", "Could not create hook key, unexpected error: "+err.Error())
		return
	}

	applyHookKeyToState(&data, hookKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) // Save data into Terraform state
}

func (r *hookKey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data hookKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) // Read prior state data into the model
	if resp.Diagnostics.HasError() {
		return
	}

	hookKey, _, err := r.OktaIDaaSClient.OktaSDKClientV5().HookKeyAPI.GetHookKey(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading hook key", "Could not read hook key, unexpected error: "+err.Error())
		return
	}
	applyHookKeyToState(&data, hookKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) // Save updated data into Terraform state
}

func (r *hookKey) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data hookKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) // Read Terraform plan data into the model
	if resp.Diagnostics.HasError() {
		return
	}

	updateHookKeyRequest := r.OktaIDaaSClient.OktaSDKClientV5().HookKeyAPI.ReplaceHookKey(ctx, data.Id.ValueString())
	updateKeyRequest := v5okta.KeyRequest{}
	updateKeyRequest.SetName(data.Name.ValueString())
	hookKey, _, err := updateHookKeyRequest.KeyRequest(updateKeyRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error replacing hook key", "Could not replace hook key, unexpected error: "+err.Error())
		return
	}

	applyHookKeyToState(&data, hookKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) // Save updated data into Terraform state
}

func (r *hookKey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data hookKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...) // Read prior state data into the model
	if resp.Diagnostics.HasError() {
		return
	}

	deleteHookKeyRequest := r.OktaIDaaSClient.OktaSDKClientV5().HookKeyAPI.DeleteHookKey(ctx, data.Id.ValueString())
	_, err := deleteHookKeyRequest.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting hook key", "Could not delete  hook key, unexpected error: "+err.Error())
		return
	}
}

func applyHookKeyToState(data *hookKeyModel, hookKey *v5okta.HookKey) {
	data.Id = types.StringPointerValue(hookKey.Id)
	data.KeyId = types.StringPointerValue(hookKey.KeyId)
	data.Name = types.StringPointerValue(hookKey.Name)
	data.Created = types.StringValue(hookKey.Created.Format(time.RFC3339))
	data.LastUpdated = types.StringValue(hookKey.LastUpdated.Format(time.RFC3339))
	data.IsUsed = types.BoolPointerValue(hookKey.IsUsed)
}
