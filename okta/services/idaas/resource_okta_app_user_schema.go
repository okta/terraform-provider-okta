package idaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	_ resource.Resource                = &appUserSchemaResource{}
	_ resource.ResourceWithConfigure   = &appUserSchemaResource{}
	_ resource.ResourceWithImportState = &appUserSchemaResource{}
)

type appUserSchemaResource struct {
	config *config.Config
}

type appUserSchemaResourceModel struct {
	ID             types.String                 `tfsdk:"id"`
	AppID          types.String                 `tfsdk:"app_id"`
	CustomProperty []appUserSchemaPropertyModel `tfsdk:"custom_property"`
}

type appUserSchemaPropertyModel struct {
	Index             types.String `tfsdk:"index"`
	Title             types.String `tfsdk:"title"`
	Type              types.String `tfsdk:"type"`
	ArrayType         types.String `tfsdk:"array_type"`
	Description       types.String `tfsdk:"description"`
	Required          types.Bool   `tfsdk:"required"`
	Scope             types.String `tfsdk:"scope"`
	MinLength         types.Int64  `tfsdk:"min_length"`
	MaxLength         types.Int64  `tfsdk:"max_length"`
	Enum              types.List   `tfsdk:"enum"`
	OneOf             []oneOfModel `tfsdk:"one_of"`
	ArrayEnum         types.List   `tfsdk:"array_enum"`
	ArrayOneOf        []oneOfModel `tfsdk:"array_one_of"`
	ExternalName      types.String `tfsdk:"external_name"`
	ExternalNamespace types.String `tfsdk:"external_namespace"`
	Master            types.String `tfsdk:"master"`
	Permissions       types.String `tfsdk:"permissions"`
	Union             types.Bool   `tfsdk:"union"`
	Unique            types.String `tfsdk:"unique"`
}

type oneOfModel struct {
	Const types.String `tfsdk:"const"`
	Title types.String `tfsdk:"title"`
}

func newAppUserSchemaResource() resource.Resource {
	return &appUserSchemaResource{}
}

func (r *appUserSchemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_user_schema"
}

func (r *appUserSchemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.config = resourceConfiguration(req, resp)
}

func (r *appUserSchemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manages the entire app user schema for an application.

This resource manages all custom properties in an application's user schema as a single object. This approach aligns with how the Okta API actually works (single mutable schema object) and provides better visibility into auto-created properties when provisioning is enabled.

**Advantages over okta_app_user_schema_property:**
- Manages all properties in one place
- Auto-created properties (from provisioning) are visible in plan/state
- Single import operation captures entire schema
- Detects drift from auto-created properties

**IMPORTANT:** With 'enum', list its values as strings even though the 'type' may be something other than string. The provider handles type coercion when making Okta API calls. Same holds for the 'const' value of 'one_of' as well as the 'array_*' variation of 'enum' and 'one_of'.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID (same as app_id)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "The Application's ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			// TODO: Once the provider upgrades from protocol v5 to v6
			// (tf5to6server.UpgradeServer in main.go), this should be changed to
			// schema.MapNestedAttribute keyed by property name/index.
			// MapNestedAttribute is v6-only and avoids both the hash instability
			// of sets and the cascading index diffs of lists.
			"custom_property": schema.SetNestedBlock{
				Description: "Custom properties in the schema",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"index": schema.StringAttribute{
							Required:    true,
							Description: "The property name/index",
						},
						"title": schema.StringAttribute{
							Required:    true,
							Description: "Display name for the property",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the schema property. It can be `string`, `boolean`, `number`, `integer`, `array`, or `object`",
						},
						"array_type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of the array elements if `type` is set to `array`",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The description of the property",
						},
						"required": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the property is required",
						},
						"scope": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Determines whether an app user attribute can be set at the Personal `SELF` or Group `NONE` level.",
						},
						"min_length": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The minimum length of the property value. Only applies to type `string`",
						},
						"max_length": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The maximum length of the property value. Only applies to type `string`",
						},
						"enum": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Array of values a primitive property can be set to. See `array_enum` for arrays.",
						},
						"array_enum": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Array of values that an array property's items can be set to.",
						},
						"unique": schema.StringAttribute{
							Optional:    true,
							Description: "Whether the property should be unique. It can be set to `UNIQUE_VALIDATED` or `NOT_UNIQUE`.",
						},
						"external_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "External name of the property",
						},
						"external_namespace": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "External namespace of the property",
						},
						"master": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Master priority for the property. It can be set to `PROFILE_MASTER` or `OKTA`",
						},
						"permissions": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Access control permissions for the property. It can be set to `READ_WRITE`, `READ_ONLY`, or `HIDE`.",
						},
						"union": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "If `type` is set to `array`, used to set whether attribute value is determined by group priority `false`, or combine values across groups `true`. Can not be set to `true` if `scope` is set to `SELF`.",
						},
					},
					Blocks: map[string]schema.Block{
						"one_of": schema.ListNestedBlock{
							Description: "Array of maps containing a mapping for display name to enum value.\n  - `const` - (Required) value mapping to member of `enum`.\n  - `title` - (Required) display name for the enum value.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"const": schema.StringAttribute{
										Required:    true,
										Description: "Enum value",
									},
									"title": schema.StringAttribute{
										Required:    true,
										Description: "Enum title",
									},
								},
							},
						},
						"array_one_of": schema.ListNestedBlock{
							Description: "Display name and value an enum array can be set to.\n  - `const` - (Required) value mapping to member of `array_enum`.\n  - `title` - (Required) display name for the enum value.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"const": schema.StringAttribute{
										Required:    true,
										Description: "Value mapping to member of `array_enum`",
									},
									"title": schema.StringAttribute{
										Required:    true,
										Description: "Display name for the enum value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *appUserSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appUserSchemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := plan.AppID.ValueString()
	plan.ID = types.StringValue(appId)

	for _, prop := range plan.CustomProperty {
		if diags := validateUnionConstraint(prop); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	for _, prop := range plan.CustomProperty {
		index := prop.Index.ValueString()
		attr := expandPropertyModel(ctx, prop, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		custom := BuildCustomUserSchema(index, attr)
		retypeUserSchemaPropertyEnums(custom)
		apiResp, err := UpdateApplicationUserProfileWithRetry(ctx, r.config, appId, custom)
		if err != nil {
			if apiResp != nil && utils.SuppressErrorOn404(apiResp, err) == nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to create app user schema property %q: application not found", index),
					err.Error(),
				)
				return
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to create app user schema property %q", index),
				err.Error(),
			)
			return
		}
	}

	found := r.readIntoState(ctx, appId, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *appUserSchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appUserSchemaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := state.ID.ValueString()
	found := r.readIntoState(ctx, appId, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *appUserSchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan appUserSchemaResourceModel
	var state appUserSchemaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := plan.AppID.ValueString()

	for _, prop := range plan.CustomProperty {
		if diags := validateUnionConstraint(prop); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	oldProps := propsModelByIndex(state.CustomProperty)
	newProps := propsModelByIndex(plan.CustomProperty)

	// Delete removed properties first.
	for index := range oldProps {
		if _, ok := newProps[index]; ok {
			continue
		}
		custom := BuildCustomUserSchema(index, nil)
		apiResp, err := UpdateApplicationUserProfileWithRetry(ctx, r.config, appId, custom)
		if err != nil {
			if apiResp != nil && utils.SuppressErrorOn404(apiResp, err) == nil {
				// App deleted out-of-band; remove so Terraform recreates.
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to delete app user schema property %q", index),
				err.Error(),
			)
			return
		}
	}

	// Upsert desired properties.
	for index, prop := range newProps {
		attr := expandPropertyModel(ctx, prop, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		custom := BuildCustomUserSchema(index, attr)
		retypeUserSchemaPropertyEnums(custom)
		apiResp, err := UpdateApplicationUserProfileWithRetry(ctx, r.config, appId, custom)
		if err != nil {
			if apiResp != nil && utils.SuppressErrorOn404(apiResp, err) == nil {
				// App deleted out-of-band; remove so Terraform recreates.
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to update app user schema property %q", index),
				err.Error(),
			)
			return
		}
	}

	found := r.readIntoState(ctx, appId, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *appUserSchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appUserSchemaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := state.AppID.ValueString()

	for _, prop := range state.CustomProperty {
		index := prop.Index.ValueString()
		if index == "" {
			continue
		}
		custom := BuildCustomUserSchema(index, nil)
		apiResp, err := UpdateApplicationUserProfileWithRetry(ctx, r.config, appId, custom)
		if err != nil {
			if apiResp != nil && utils.SuppressErrorOn404(apiResp, err) == nil {
				return // app already gone
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to delete app user schema property %q", index),
				err.Error(),
			)
			return
		}
	}
}

func (r *appUserSchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), req.ID)...)
}

// readIntoState fetches the schema from the API and populates the model.
// Returns false only when the app is gone (404); callers should RemoveResource.
// Returns true on all other outcomes (success or non-404 error); callers must
// check diags.HasError() before using the state.
func (r *appUserSchemaResource) readIntoState(ctx context.Context, appId string, state *appUserSchemaResourceModel, diags *diag.Diagnostics) bool {
	client := r.config.OktaIDaaSClient.OktaSDKClientV2()
	us, resp, err := client.UserSchema.GetApplicationUserSchema(ctx, appId)
	if err != nil {
		if resp != nil && utils.SuppressErrorOn404(resp, err) == nil {
			return false
		}
		diags.AddError("failed to get application user schema", err.Error())
		return true // not a 404 — don't remove from state
	}

	state.ID = types.StringValue(appId)
	state.AppID = types.StringValue(appId)

	customProps := make([]appUserSchemaPropertyModel, 0)
	if us.Definitions != nil && us.Definitions.Custom != nil && us.Definitions.Custom.Properties != nil {
		for index, attr := range us.Definitions.Custom.Properties {
			customProps = append(customProps, flattenPropertyModel(ctx, index, attr, diags))
			if diags.HasError() {
				return true
			}
		}
	}
	state.CustomProperty = customProps

	return true
}

// flattenPropertyModel converts an SDK UserSchemaAttribute to a framework model.
// Only values returned by the API are set; omitted attributes use null types.
func flattenPropertyModel(ctx context.Context, index string, sdkAttr *sdk.UserSchemaAttribute, diags *diag.Diagnostics) appUserSchemaPropertyModel {
	m := appUserSchemaPropertyModel{
		Index: types.StringValue(index),
	}
	if sdkAttr == nil {
		return m
	}

	m.Title = types.StringValue(sdkAttr.Title)
	m.Type = types.StringValue(sdkAttr.Type)

	if sdkAttr.Items != nil && sdkAttr.Items.Type != "" {
		m.ArrayType = types.StringValue(sdkAttr.Items.Type)
	} else {
		m.ArrayType = types.StringNull()
	}

	if sdkAttr.Description != "" {
		m.Description = types.StringValue(sdkAttr.Description)
	} else {
		m.Description = types.StringNull()
	}

	if sdkAttr.ExternalName != "" {
		m.ExternalName = types.StringValue(sdkAttr.ExternalName)
	} else {
		m.ExternalName = types.StringNull()
	}

	if sdkAttr.ExternalNamespace != "" {
		m.ExternalNamespace = types.StringValue(sdkAttr.ExternalNamespace)
	} else {
		m.ExternalNamespace = types.StringNull()
	}

	if sdkAttr.Required != nil {
		m.Required = types.BoolValue(*sdkAttr.Required)
	} else {
		m.Required = types.BoolNull()
	}

	if sdkAttr.Scope != "" {
		m.Scope = types.StringValue(sdkAttr.Scope)
	} else {
		m.Scope = types.StringNull()
	}

	if sdkAttr.MinLengthPtr != nil {
		m.MinLength = types.Int64Value(*sdkAttr.MinLengthPtr)
	} else {
		m.MinLength = types.Int64Null()
	}

	if sdkAttr.MaxLengthPtr != nil {
		m.MaxLength = types.Int64Value(*sdkAttr.MaxLengthPtr)
	} else {
		m.MaxLength = types.Int64Null()
	}

	if len(sdkAttr.Enum) > 0 {
		stringifyEnumSlice(sdkAttr.Type, &sdkAttr.Enum)
		enumVals := make([]attr.Value, 0, len(sdkAttr.Enum))
		for _, e := range sdkAttr.Enum {
			if s, ok := e.(string); ok {
				enumVals = append(enumVals, types.StringValue(s))
			} else {
				enumVals = append(enumVals, types.StringValue(fmt.Sprintf("%v", e)))
			}
		}
		listVal, d := types.ListValue(types.StringType, enumVals)
		diags.Append(d...)
		m.Enum = listVal
	} else {
		m.Enum = types.ListNull(types.StringType)
	}

	if len(sdkAttr.OneOf) > 0 {
		stringifyOneOfSlice(sdkAttr.Type, &sdkAttr.OneOf)
		m.OneOf = flattenOneOfModels(sdkAttr.OneOf)
	}

	if sdkAttr.Items != nil && len(sdkAttr.Items.Enum) > 0 {
		stringifyEnumSlice(sdkAttr.Items.Type, &sdkAttr.Items.Enum)
		enumVals := make([]attr.Value, 0, len(sdkAttr.Items.Enum))
		for _, e := range sdkAttr.Items.Enum {
			switch v := e.(type) {
			case string:
				enumVals = append(enumVals, types.StringValue(v))
			case map[string]interface{}:
				b, _ := json.Marshal(v)
				enumVals = append(enumVals, types.StringValue(string(b)))
			default:
				enumVals = append(enumVals, types.StringValue(fmt.Sprintf("%v", v)))
			}
		}
		listVal, d := types.ListValue(types.StringType, enumVals)
		diags.Append(d...)
		m.ArrayEnum = listVal
	} else {
		m.ArrayEnum = types.ListNull(types.StringType)
	}

	if sdkAttr.Items != nil && len(sdkAttr.Items.OneOf) > 0 {
		stringifyOneOfSlice(sdkAttr.Items.Type, &sdkAttr.Items.OneOf)
		m.ArrayOneOf = flattenOneOfModels(sdkAttr.Items.OneOf)
	}

	if sdkAttr.Unique != "" {
		m.Unique = types.StringValue(sdkAttr.Unique)
	} else {
		m.Unique = types.StringNull()
	}

	if sdkAttr.Master != nil && sdkAttr.Master.Type != "" {
		m.Master = types.StringValue(sdkAttr.Master.Type)
	} else {
		m.Master = types.StringNull()
	}

	if len(sdkAttr.Permissions) > 0 {
		m.Permissions = types.StringNull()
		for _, perm := range sdkAttr.Permissions {
			if perm.Action != "" {
				m.Permissions = types.StringValue(perm.Action)
				break
			}
		}
	} else {
		m.Permissions = types.StringNull()
	}

	if sdkAttr.Union == "ENABLE" {
		m.Union = types.BoolValue(true)
	} else if sdkAttr.Union == "DISABLE" {
		m.Union = types.BoolValue(false)
	} else {
		m.Union = types.BoolNull()
	}

	return m
}

// expandPropertyModel converts a framework model to an SDK UserSchemaAttribute for API calls.
func expandPropertyModel(ctx context.Context, m appUserSchemaPropertyModel, diags *diag.Diagnostics) *sdk.UserSchemaAttribute {
	a := &sdk.UserSchemaAttribute{}

	if !m.Title.IsNull() && !m.Title.IsUnknown() {
		a.Title = m.Title.ValueString()
	}
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		a.Type = m.Type.ValueString()
	}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		a.Description = m.Description.ValueString()
	}
	if !m.Required.IsNull() && !m.Required.IsUnknown() {
		b := m.Required.ValueBool()
		a.Required = &b
	}
	if !m.Scope.IsNull() && !m.Scope.IsUnknown() {
		a.Scope = m.Scope.ValueString()
	}
	if !m.MinLength.IsNull() && !m.MinLength.IsUnknown() {
		v := m.MinLength.ValueInt64()
		a.MinLengthPtr = &v
	}
	if !m.MaxLength.IsNull() && !m.MaxLength.IsUnknown() {
		v := m.MaxLength.ValueInt64()
		a.MaxLengthPtr = &v
	}
	if !m.ArrayType.IsNull() && !m.ArrayType.IsUnknown() {
		a.Items = &sdk.UserSchemaAttributeItems{
			Type: m.ArrayType.ValueString(),
		}
	}
	if !m.Enum.IsNull() && !m.Enum.IsUnknown() {
		var enumStrings []string
		diags.Append(m.Enum.ElementsAs(ctx, &enumStrings, false)...)
		enumIface := make([]interface{}, len(enumStrings))
		for i, s := range enumStrings {
			enumIface[i] = s
		}
		a.Enum = enumIface
	}
	if len(m.OneOf) > 0 {
		a.OneOf = expandOneOfModels(m.OneOf, a.Type, diags)
		if diags.HasError() {
			return a
		}
	}
	if !m.ArrayEnum.IsNull() && !m.ArrayEnum.IsUnknown() {
		var enumStrings []string
		diags.Append(m.ArrayEnum.ElementsAs(ctx, &enumStrings, false)...)
		enumIface := make([]interface{}, len(enumStrings))
		for i, s := range enumStrings {
			enumIface[i] = s
		}
		if a.Items == nil {
			a.Items = &sdk.UserSchemaAttributeItems{}
		}
		a.Items.Enum = enumIface
	}
	if len(m.ArrayOneOf) > 0 {
		if a.Items == nil {
			a.Items = &sdk.UserSchemaAttributeItems{}
		}
		a.Items.OneOf = expandOneOfModels(m.ArrayOneOf, a.Items.Type, diags)
		if diags.HasError() {
			return a
		}
	}
	if !m.ExternalName.IsNull() && !m.ExternalName.IsUnknown() {
		a.ExternalName = m.ExternalName.ValueString()
	}
	if !m.ExternalNamespace.IsNull() && !m.ExternalNamespace.IsUnknown() {
		a.ExternalNamespace = m.ExternalNamespace.ValueString()
	}
	if !m.Master.IsNull() && !m.Master.IsUnknown() {
		a.Master = &sdk.UserSchemaAttributeMaster{Type: m.Master.ValueString()}
	}
	if !m.Permissions.IsNull() && !m.Permissions.IsUnknown() {
		a.Permissions = []*sdk.UserSchemaAttributePermission{
			{Action: m.Permissions.ValueString(), Principal: "SELF"},
		}
	}
	if !m.Unique.IsNull() && !m.Unique.IsUnknown() {
		a.Unique = m.Unique.ValueString()
	}
	// Only send union for array attributes — Okta ignores it for other types.
	if a.Type == "array" && !m.Union.IsNull() && !m.Union.IsUnknown() {
		if m.Union.ValueBool() {
			a.Union = "ENABLE"
		} else {
			a.Union = "DISABLE"
		}
	}

	return a
}

// propsModelByIndex builds an index-keyed map from the CustomProperty slice.
func propsModelByIndex(props []appUserSchemaPropertyModel) map[string]appUserSchemaPropertyModel {
	out := make(map[string]appUserSchemaPropertyModel, len(props))
	for _, p := range props {
		idx := p.Index.ValueString()
		if idx != "" {
			out[idx] = p
		}
	}
	return out
}

// validateUnionConstraint checks that union is only enabled for array-type properties
// and not when scope is SELF. Error messages are kept verbatim for test compatibility.
func validateUnionConstraint(prop appUserSchemaPropertyModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if prop.Union.IsNull() || prop.Union.IsUnknown() {
		return diags
	}

	typ := prop.Type.ValueString()
	if typ != "array" {
		diags.AddError(
			"Invalid union configuration",
			fmt.Sprintf("custom_property %q: union can only be set when type is \"array\"", prop.Index.ValueString()),
		)
		return diags
	}

	if prop.Union.ValueBool() {
		if !prop.Scope.IsNull() && prop.Scope.ValueString() == "SELF" {
			diags.AddError(
				"Invalid union configuration",
				fmt.Sprintf("custom_property %q: union cannot be enabled when scope is \"SELF\"", prop.Index.ValueString()),
			)
		}
	}

	return diags
}

// flattenOneOfModels converts SDK OneOf enums to framework model slices.
// Const values are stringified before this is called; object-type maps are marshaled to JSON.
func flattenOneOfModels(oneOf []*sdk.UserSchemaAttributeEnum) []oneOfModel {
	result := make([]oneOfModel, 0, len(oneOf))
	for _, v := range oneOf {
		if v == nil {
			continue
		}
		var constStr string
		if obj, ok := v.Const.(map[string]interface{}); ok {
			b, _ := json.Marshal(obj)
			constStr = string(b)
		} else if s, ok := v.Const.(string); ok {
			constStr = s
		} else {
			constStr = fmt.Sprintf("%v", v.Const)
		}
		result = append(result, oneOfModel{
			Const: types.StringValue(constStr),
			Title: types.StringValue(v.Title),
		})
	}
	return result
}

// expandOneOfModels converts framework model slices to SDK OneOf enums.
// For object types, JSON strings are unmarshaled to maps.
func expandOneOfModels(models []oneOfModel, elemType string, diags *diag.Diagnostics) []*sdk.UserSchemaAttributeEnum {
	result := make([]*sdk.UserSchemaAttributeEnum, len(models))
	for i, m := range models {
		constVal := m.Const.ValueString()
		var constIface interface{} = constVal
		if elemType == "object" {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(constVal), &obj); err != nil {
				diags.AddError("Invalid JSON in one_of const",
					fmt.Sprintf("Failed to parse JSON for one_of[%d].const %q: %s", i, constVal, err))
				return nil
			}
			constIface = obj
		}
		result[i] = &sdk.UserSchemaAttributeEnum{
			Const: constIface,
			Title: m.Title.ValueString(),
		}
	}
	return result
}
