package governance

import (
	"context"
	"time"

	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &grantResource{}
	_ resource.ResourceWithConfigure   = &grantResource{}
	_ resource.ResourceWithImportState = &grantResource{}
)

func newGrantResource() resource.Resource {
	return &grantResource{}
}

type grantResource struct {
	*config.Config
}

type scheduleSettingModel struct {
	ExpirationDate types.String `tfsdk:"expiration_date"`
	Timezone       types.String `tfsdk:"timezone"`
}

type entitlementGrantModel struct {
	Id     types.String            `tfsdk:"id"`
	Values []entitlementValueModel `tfsdk:"values"`
}

type grantResourceModel struct {
	Id                  types.String            `tfsdk:"id"`
	EntitlementBundleId types.String            `tfsdk:"entitlement_bundle_id"`
	GrantType           types.String            `tfsdk:"grant_type"`
	Action              types.String            `tfsdk:"action"`
	TargetPrincipal     *principalModel         `tfsdk:"target_principal"`
	Actor               types.String            `tfsdk:"actor"`
	ScheduleSettings    *scheduleSettingModel   `tfsdk:"schedule_settings"`
	Target              *principalModel         `tfsdk:"target"`
	Entitlements        []entitlementGrantModel `tfsdk:"entitlements"`
}

func (r *grantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grant"
}

func (r *grantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *grantResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func (r *grantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"entitlement_bundle_id": schema.StringAttribute{
				Optional: true,
			},
			"grant_type": schema.StringAttribute{
				Required: true,
			},
			"action": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ALLOW"),
				Validators: []validator.String{
					stringvalidator.OneOf("ALLOW", "DENY"),
				},
			},
			"actor": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("API"),
				Validators: []validator.String{
					stringvalidator.OneOf("ACCESS_REQUEST", "ADMIN", "API", "NONE"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"target_principal": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"schedule_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"expiration_date": schema.StringAttribute{
						Optional: true,
					},
					"timezone": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("America/Toronto"),
					},
				},
			},
			"target": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			"entitlements": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Optional: true,
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

func (r *grantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data grantResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	createGrantResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.CreateGrant(ctx).GrantCreatable(*buildGrant(data)).Execute()
	if err != nil {
		return
	}

	// Example Data value setting
	r.applyGrantToState(ctx, &data, createGrantResp)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *grantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data grantResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *grantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan grantResourceModel
	var state grantResourceModel

	// Read Terraform plan Data into the model
	// Read plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id

	// Update API call logic
	updatedGrantResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.UpdateGrant(ctx, plan.Id.ValueString()).GrantPatch(buildGrantPatch(plan)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Grants",
			"Could not update Grants, unexpected error: "+err.Error(),
		)
		return
	}

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(r.applyGrantToState(ctx, &plan, updatedGrantResp)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *grantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Nothing to do – Terraform will forget the resource
}

func (r *grantResource) applyGrantToState(ctx context.Context, data *grantResourceModel, createGrantResp *governance.GrantFull) diag.Diagnostics {
	var diags diag.Diagnostics
	data.Id = types.StringValue(createGrantResp.Id)
	data.EntitlementBundleId = types.StringPointerValue(createGrantResp.EntitlementBundleId)
	data.GrantType = types.StringValue(string(createGrantResp.GrantType))
	data.Action = types.StringValue(string(createGrantResp.Action))
	data.TargetPrincipal = &principalModel{
		ExternalId: types.StringValue(createGrantResp.TargetPrincipal.ExternalId),
		Type:       types.StringValue(string(createGrantResp.TargetPrincipal.Type)),
	}
	data.Actor = types.StringValue(string(createGrantResp.Actor))
	if createGrantResp.ScheduleSettings != nil {
		data.ScheduleSettings = &scheduleSettingModel{}
		if createGrantResp.ScheduleSettings.ExpirationDate != nil {
			data.ScheduleSettings.ExpirationDate = types.StringValue(createGrantResp.ScheduleSettings.ExpirationDate.Format(time.RFC3339))
		}
		if createGrantResp.ScheduleSettings.TimeZone != nil {
			data.ScheduleSettings.Timezone = types.StringValue(*createGrantResp.ScheduleSettings.TimeZone)
		}
	}

	if createGrantResp.GrantType != governance.GRANTTYPE_ENTITLEMENT_BUNDLE {
		if target, ok := createGrantResp.GetTargetOk(); ok {
			data.Target = &principalModel{
				ExternalId: types.StringValue(target.ExternalId),
				Type:       types.StringValue(string(target.Type)),
			}
		} else {
			data.Target.ExternalId = types.StringNull()
			data.Target.Type = types.StringNull()
		}
	}
	if createGrantResp.Entitlements != nil {
		var entitlements []entitlementGrantModel
		for _, ent := range createGrantResp.Entitlements {
			var values []entitlementValueModel
			for _, val := range ent.Values {
				v := entitlementValueModel{
					Id: types.StringPointerValue(val.Id),
				}
				values = append(values, v)
			}
			entitlement := entitlementGrantModel{
				Id:     types.StringPointerValue(ent.Id),
				Values: values,
			}
			entitlements = append(entitlements, entitlement)
		}
		data.Entitlements = entitlements
	}
	return diags
}

func buildGrant(data grantResourceModel) *governance.GrantCreatable {
	if data.GrantType.ValueString() == string(governance.GRANTTYPE_ENTITLEMENT_BUNDLE) {
		var expirationDate *time.Time
		var timezone *string
		if data.ScheduleSettings != nil {
			ptr := data.ScheduleSettings.ExpirationDate.ValueStringPointer()
			if ptr != nil && *ptr != "" {
				t, _ := time.Parse(time.RFC3339, *ptr)
				expirationDate = &t
			}
			timezone = data.ScheduleSettings.Timezone.ValueStringPointer()
		}
		grant := &governance.GrantCreatable{
			GrantTypeBundleWriteable: &governance.GrantTypeBundleWriteable{
				GrantType:           data.GrantType.ValueString(),
				EntitlementBundleId: data.EntitlementBundleId.ValueString(),
				Action:              (*governance.GrantAction)(data.Action.ValueStringPointer()),
				TargetPrincipal: governance.TargetPrincipal{
					ExternalId: data.TargetPrincipal.ExternalId.ValueString(),
					Type:       governance.PrincipalType(data.TargetPrincipal.Type.ValueString()),
				},
				Actor: (*governance.GrantActor)(data.Actor.ValueStringPointer()),
			},
		}
		if expirationDate != nil && timezone != nil {
			grant.GrantTypeBundleWriteable.ScheduleSettings = governance.NewScheduleSettingsWriteableWithDefaults()
			grant.GrantTypeBundleWriteable.ScheduleSettings.ExpirationDate = expirationDate
			grant.GrantTypeBundleWriteable.ScheduleSettings.TimeZone = timezone
		} else if expirationDate != nil {
			grant.GrantTypeBundleWriteable.ScheduleSettings = governance.NewScheduleSettingsWriteableWithDefaults()
			grant.GrantTypeBundleWriteable.ScheduleSettings.ExpirationDate = expirationDate
		}
		return grant

	} else if data.GrantType.ValueString() == string(governance.GRANTTYPE_POLICY) {
		// Handle other grant types if necessary
		var expirationDate *time.Time
		var timezone *string
		if data.ScheduleSettings != nil {
			ptr := data.ScheduleSettings.ExpirationDate.ValueStringPointer()
			if ptr != nil && *ptr != "" {
				t, _ := time.Parse(time.RFC3339, *ptr)
				expirationDate = &t
			}
			timezone = data.ScheduleSettings.Timezone.ValueStringPointer()
		}
		grant := &governance.GrantCreatable{
			GrantTypePolicyWriteable: &governance.GrantTypePolicyWriteable{
				// Populate fields as needed
				GrantType: data.GrantType.ValueString(),
				Target: governance.TargetResource{
					ExternalId: data.Target.ExternalId.ValueString(),
					Type:       governance.ResourceType2(data.Target.Type.ValueString()),
				},
				Action: (*governance.GrantAction)(data.Action.ValueStringPointer()),
				TargetPrincipal: governance.TargetPrincipal{
					ExternalId: data.TargetPrincipal.ExternalId.ValueString(),
					Type:       governance.PrincipalType(data.TargetPrincipal.Type.ValueString()),
				},
				Actor: (*governance.GrantActor)(data.Actor.ValueStringPointer()),
			},
		}
		if expirationDate != nil && timezone != nil {
			grant.GrantTypeCustomWriteable.ScheduleSettings = governance.NewScheduleSettingsWriteableWithDefaults()
			grant.GrantTypeCustomWriteable.ScheduleSettings.ExpirationDate = expirationDate
			grant.GrantTypeCustomWriteable.ScheduleSettings.TimeZone = timezone
		} else if expirationDate != nil {
			grant.GrantTypePolicyWriteable.ScheduleSettings = governance.NewScheduleSettingsWriteableWithDefaults()
			grant.GrantTypePolicyWriteable.ScheduleSettings.ExpirationDate = expirationDate
		}
		return grant
	} else if data.GrantType.ValueString() == string(governance.GRANTTYPE_CUSTOM) {
		var expirationDate *time.Time
		var timezone *string
		if data.ScheduleSettings != nil {
			ptr := data.ScheduleSettings.ExpirationDate.ValueStringPointer()
			if ptr != nil && *ptr != "" {
				t, _ := time.Parse(time.RFC3339, *ptr)
				expirationDate = &t
			}
			timezone = data.ScheduleSettings.Timezone.ValueStringPointer()
		}
		// Handle custom grant type if necessary
		grant := &governance.GrantCreatable{
			GrantTypeCustomWriteable: &governance.GrantTypeCustomWriteable{
				GrantType: data.GrantType.ValueString(),
				Target: governance.TargetResource{
					ExternalId: data.Target.ExternalId.ValueString(),
					Type:       governance.ResourceType2(data.Target.Type.ValueString()),
				},
				Entitlements: buildEntitlements(data.Entitlements),
				TargetPrincipal: governance.TargetPrincipal{
					ExternalId: data.TargetPrincipal.ExternalId.ValueString(),
					Type:       governance.PrincipalType(data.TargetPrincipal.Type.ValueString()),
				},
				Action: (*governance.GrantAction)(data.Action.ValueStringPointer()),
				Actor:  (*governance.GrantActor)(data.Actor.ValueStringPointer()),
			},
		}
		if expirationDate != nil && timezone != nil {
			grant.GrantTypeCustomWriteable.ScheduleSettings = governance.NewScheduleSettingsWriteableWithDefaults()
			grant.GrantTypeCustomWriteable.ScheduleSettings.ExpirationDate = expirationDate
			grant.GrantTypeCustomWriteable.ScheduleSettings.TimeZone = timezone
		} else if expirationDate != nil {
			grant.GrantTypeCustomWriteable.ScheduleSettings = governance.NewScheduleSettingsWriteableWithDefaults()
			grant.GrantTypeCustomWriteable.ScheduleSettings.ExpirationDate = expirationDate
		}
		return grant
	}
	return nil
}

func buildEntitlements(entitlements []entitlementGrantModel) []governance.EntitlementCreatable {
	var entitlementCreatables []governance.EntitlementCreatable
	for _, ent := range entitlements {
		var values []governance.EntitlementValueCreatable
		for _, val := range ent.Values {
			creatable := governance.EntitlementValueCreatable{Id: val.Id.ValueStringPointer()}
			values = append(values, creatable)
		}
		entitlementCreatable := governance.EntitlementCreatable{
			Id:     ent.Id.ValueStringPointer(),
			Values: values,
		}
		entitlementCreatables = append(entitlementCreatables, entitlementCreatable)
	}
	return entitlementCreatables
}

func buildGrantPatch(plan grantResourceModel) governance.GrantPatch {
	var expirationDate *time.Time
	if ptr := plan.ScheduleSettings.ExpirationDate.ValueStringPointer(); ptr != nil && *ptr != "" {
		t, _ := time.Parse(time.RFC3339, *ptr)
		expirationDate = &t
	}
	timezone := plan.ScheduleSettings.Timezone.ValueStringPointer()
	scheduleSettings := governance.ScheduleSettingsWriteable{}
	if expirationDate != nil {
		scheduleSettings.ExpirationDate = expirationDate
	}
	if timezone != nil {
		scheduleSettings.TimeZone = timezone
	}
	return governance.GrantPatch{
		Id:               plan.Id.ValueString(),
		ScheduleSettings: scheduleSettings,
	}
}
