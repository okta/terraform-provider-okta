package governance

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
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

type grantResourceModel struct {
	Id                  types.String                    `tfsdk:"id"`
	GrantType           types.String                    `tfsdk:"grant_type"`
	TargetPrincipalId   types.String                    `tfsdk:"target_principal_id"`
	TargetPrincipalType types.String                    `tfsdk:"target_principal_type"`
	TargetResourceOrn   types.String                    `tfsdk:"target_resource_orn"`
	EntitlementBundleId types.String                    `tfsdk:"entitlement_bundle_id"`
	Entitlements        []grantResourceEntitlementModel `tfsdk:"entitlements"`
	Action              types.String                    `tfsdk:"action"`
	Actor               types.String                    `tfsdk:"actor"`
	ExpirationDate      types.String                    `tfsdk:"expiration_date"`
	TimeZone            types.String                    `tfsdk:"time_zone"`
	Status              types.String                    `tfsdk:"status"`
	Created             types.String                    `tfsdk:"created"`
	CreatedBy           types.String                    `tfsdk:"created_by"`
	LastUpdated         types.String                    `tfsdk:"last_updated"`
	LastUpdatedBy       types.String                    `tfsdk:"last_updated_by"`
}

type grantResourceEntitlementModel struct {
	Id     types.String                         `tfsdk:"id"`
	Values []grantResourceEntitlementValueModel `tfsdk:"values"`
}

type grantResourceEntitlementValueModel struct {
	Id types.String `tfsdk:"id"`
}

func (r *grantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grant"
}

func (r *grantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a grant. Supports CUSTOM and ENTITLEMENT-BUNDLE grant types. Note: POLICY grants are created automatically by collection assignments and cannot be created directly.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the grant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"grant_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of grant. Valid values: CUSTOM, ENTITLEMENT-BUNDLE. POLICY grants cannot be created directly.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_principal_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the principal (user or group) receiving the grant.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_principal_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of principal. Valid values: OKTA_USER, OKTA_GROUP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_resource_orn": schema.StringAttribute{
				Required:    true,
				Description: "The ORN of the target resource (e.g., app).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"entitlement_bundle_id": schema.StringAttribute{
				Optional:    true,
				Description: "The entitlement bundle ID. Required when grant_type is ENTITLEMENT-BUNDLE.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The grant action. Valid values: ALLOW, DENY. Default: ALLOW.",
			},
			"actor": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The actor creating the grant. Valid values: ACCESS_REQUEST, ADMIN, API, NONE. Default: API.",
			},
			"expiration_date": schema.StringAttribute{
				Optional:    true,
				Description: "The expiration date for the grant (ISO 8601 format).",
			},
			"time_zone": schema.StringAttribute{
				Optional:    true,
				Description: "The time zone in IANA format for the expiration date.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The grant status.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the grant was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User ID who created the grant.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the grant was last updated.",
			},
			"last_updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "User ID who last updated the grant.",
			},
		},
		Blocks: map[string]schema.Block{
			"entitlements": schema.ListNestedBlock{
				Description: "List of entitlements with their values. Required when grant_type is CUSTOM.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:    true,
							Description: "The entitlement ID.",
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.ListNestedBlock{
							Description: "List of entitlement value IDs.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:    true,
										Description: "The entitlement value ID.",
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
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...) // capture plan
	if resp.Diagnostics.HasError() {
		return
	}

	// Normalize defaults early
	if data.Action.IsNull() {
		data.Action = types.StringValue("ALLOW")
	}
	if data.Actor.IsNull() {
		data.Actor = types.StringValue("API")
	}

	grantType := data.GrantType.ValueString()
	if grantType != "CUSTOM" && grantType != "ENTITLEMENT-BUNDLE" {
		resp.Diagnostics.AddError(
			"Invalid grant type",
			fmt.Sprintf("grant_type must be CUSTOM or ENTITLEMENT-BUNDLE, got: %s. POLICY grants cannot be created directly.", grantType),
		)
		return
	}

	// Mutual exclusivity validation
	if grantType == "CUSTOM" {
		if len(data.Entitlements) == 0 {
			resp.Diagnostics.AddError("Missing entitlements", "entitlements are required for CUSTOM grant type")
			return
		}
		if !data.EntitlementBundleId.IsNull() {
			resp.Diagnostics.AddError("Invalid configuration", "entitlement_bundle_id must not be set when grant_type is CUSTOM")
			return
		}
	} else { // ENTITLEMENT-BUNDLE
		if data.EntitlementBundleId.IsNull() {
			resp.Diagnostics.AddError("Missing entitlement_bundle_id", "entitlement_bundle_id is required for ENTITLEMENT-BUNDLE grant type")
			return
		}
		if len(data.Entitlements) > 0 {
			resp.Diagnostics.AddError("Invalid configuration", "entitlements must not be set when grant_type is ENTITLEMENT-BUNDLE")
			return
		}
	}

	// target_resource_orn format validation (basic)
	orn := data.TargetResourceOrn.ValueString()
	if !strings.HasPrefix(orn, "orn:okta:idp:") || len(strings.Split(orn, ":")) < 6 {
		resp.Diagnostics.AddError("Invalid target_resource_orn", fmt.Sprintf("Unexpected format: %s", orn))
		return
	}

	// Cross-field validation for expiration/time_zone
	if !data.TimeZone.IsNull() && data.ExpirationDate.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "time_zone requires expiration_date to be set")
		return
	}
	if !data.ExpirationDate.IsNull() {
		parsed, err := time.Parse(time.RFC3339, data.ExpirationDate.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid expiration_date format", err.Error())
			return
		}
		if parsed.Before(time.Now().UTC()) {
			resp.Diagnostics.AddError("Invalid expiration_date", "expiration_date must be a future timestamp")
			return
		}
	}

	// Build the grant request based on grant type
	var grantCreatable governance.GrantCreatable
	if grantType == "CUSTOM" {
		customGrant := buildCustomGrant(data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		grantCreatable = governance.GrantTypeCustomWriteableAsGrantCreatable(&customGrant)
	} else { // ENTITLEMENT-BUNDLE
		bundleGrant := buildBundleGrant(data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		grantCreatable = governance.GrantTypeBundleWriteableAsGrantCreatable(&bundleGrant)
	}

	grant, httpResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.
		CreateGrant(ctx).
		GrantCreatable(grantCreatable).
		Execute()
	if err != nil {
		// Attempt to surface HTTP status for clarity
		if httpResp != nil {
			resp.Diagnostics.AddError("Error creating grant", fmt.Sprintf("%s (status %d)", err.Error(), httpResp.StatusCode))
		} else {
			resp.Diagnostics.AddError("Error creating grant", err.Error())
		}
		return
	}

	applyGrantToState(&data, grant)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *grantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data grantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	grantResp, httpResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.
		GetGrant(ctx, data.Id.ValueString()).
		Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource gone; clear state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading grant", err.Error())
		return
	}

	if grantResp.GrantFull != nil {
		applyGrantToState(&data, grantResp.GrantFull)
	} else if grantResp.GrantFullWithEntitlements != nil {
		applyGrantWithEntitlementsToState(&data, grantResp.GrantFullWithEntitlements)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *grantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data grantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only schedule settings can be updated
	scheduleSettings := governance.NewScheduleSettingsWriteable()

	if !data.TimeZone.IsNull() && data.ExpirationDate.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "time_zone requires expiration_date to be set")
		return
	}
	if !data.ExpirationDate.IsNull() {
		parsed, err := time.Parse(time.RFC3339, data.ExpirationDate.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid expiration_date format", err.Error())
			return
		}
		if parsed.Before(time.Now().UTC()) {
			resp.Diagnostics.AddError("Invalid expiration_date", "expiration_date must be a future timestamp")
			return
		}
		scheduleSettings.SetExpirationDate(parsed)
	}
	if !data.TimeZone.IsNull() {
		scheduleSettings.SetTimeZone(data.TimeZone.ValueString())
	}

	grantPatch := governance.NewGrantPatch(data.Id.ValueString(), *scheduleSettings)

	_, httpResp, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.
		UpdateGrant(ctx, data.Id.ValueString()).
		GrantPatch(*grantPatch).
		Execute()
	if err != nil {
		if httpResp != nil {
			resp.Diagnostics.AddError("Error updating grant", fmt.Sprintf("%s (status %d)", err.Error(), httpResp.StatusCode))
		} else {
			resp.Diagnostics.AddError("Error updating grant", err.Error())
		}
		return
	}

	// Re-read to get updated state
	var readReq resource.ReadRequest
	readReq.State = req.State
	var readResp resource.ReadResponse
	readResp.State = resp.State
	r.Read(ctx, readReq, &readResp)
	resp.Diagnostics = readResp.Diagnostics
	resp.State = readResp.State
}

func (r *grantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data grantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(
		"Grant deletion",
		"Grants cannot be directly deleted via the API. The grant has been removed from Terraform state, but may still exist in Okta. Remove collection assignments or revoke access through governance workflows to actually rescind access.",
	)
}

func (r *grantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *grantResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func buildCustomGrant(data grantResourceModel, diags *diag.Diagnostics) governance.GrantTypeCustomWriteable {
	targetPrincipal := *governance.NewTargetPrincipal(
		data.TargetPrincipalId.ValueString(),
		governance.PrincipalType(data.TargetPrincipalType.ValueString()),
	)
	targetResource := *governance.NewTargetResource(data.TargetResourceOrn.ValueString(), governance.RESOURCETYPE2_APPLICATION)
	grant := *governance.NewGrantTypeCustomWriteable(
		"CUSTOM",
		targetResource,
		targetPrincipal,
	)

	// Sort entitlements and values for deterministic ordering
	if len(data.Entitlements) > 0 {
		entitlements := make([]governance.EntitlementCreatable, len(data.Entitlements))
		// Sort entitlements by id
		sort.SliceStable(data.Entitlements, func(i, j int) bool {
			return data.Entitlements[i].Id.ValueString() < data.Entitlements[j].Id.ValueString()
		})
		for i, ent := range data.Entitlements {
			entitlement := governance.NewEntitlementCreatable()
			entitlement.SetId(ent.Id.ValueString())
			if len(ent.Values) > 0 {
				sort.SliceStable(ent.Values, func(a, b int) bool { return ent.Values[a].Id.ValueString() < ent.Values[b].Id.ValueString() })
				values := make([]governance.EntitlementValueCreatable, len(ent.Values))
				for j, val := range ent.Values {
					value := governance.NewEntitlementValueCreatable()
					value.SetId(val.Id.ValueString())
					values[j] = *value
				}
				entitlement.SetValues(values)
			}
			entitlements[i] = *entitlement
		}
		grant.SetEntitlements(entitlements)
	}

	if !data.Action.IsNull() {
		action := governance.GrantAction(data.Action.ValueString())
		grant.SetAction(action)
	}
	if !data.Actor.IsNull() {
		actor := governance.GrantActor(data.Actor.ValueString())
		grant.SetActor(actor)
	}
	if !data.ExpirationDate.IsNull() || !data.TimeZone.IsNull() {
		scheduleSettings := governance.NewScheduleSettingsWriteable()
		if !data.ExpirationDate.IsNull() {
			parsed, err := time.Parse(time.RFC3339, data.ExpirationDate.ValueString())
			if err != nil {
				diags.AddError("Invalid expiration_date format", err.Error())
				return grant
			}
			scheduleSettings.SetExpirationDate(parsed)
		}
		if !data.TimeZone.IsNull() {
			scheduleSettings.SetTimeZone(data.TimeZone.ValueString())
		}
		grant.SetScheduleSettings(*scheduleSettings)
	}
	return grant
}

func buildBundleGrant(data grantResourceModel, diags *diag.Diagnostics) governance.GrantTypeBundleWriteable {
	targetPrincipal := *governance.NewTargetPrincipal(
		data.TargetPrincipalId.ValueString(),
		governance.PrincipalType(data.TargetPrincipalType.ValueString()),
	)
	grant := *governance.NewGrantTypeBundleWriteable(
		"ENTITLEMENT-BUNDLE",
		data.EntitlementBundleId.ValueString(),
		targetPrincipal,
	)
	if !data.Action.IsNull() {
		action := governance.GrantAction(data.Action.ValueString())
		grant.SetAction(action)
	}
	if !data.Actor.IsNull() {
		actor := governance.GrantActor(data.Actor.ValueString())
		grant.SetActor(actor)
	}
	if !data.ExpirationDate.IsNull() || !data.TimeZone.IsNull() {
		scheduleSettings := governance.NewScheduleSettingsWriteable()
		if !data.ExpirationDate.IsNull() {
			parsed, err := time.Parse(time.RFC3339, data.ExpirationDate.ValueString())
			if err != nil {
				diags.AddError("Invalid expiration_date format", err.Error())
				return grant
			}
			scheduleSettings.SetExpirationDate(parsed)
		}
		if !data.TimeZone.IsNull() {
			scheduleSettings.SetTimeZone(data.TimeZone.ValueString())
		}
		grant.SetScheduleSettings(*scheduleSettings)
	}
	return grant
}

func applyGrantToState(data *grantResourceModel, grant *governance.GrantFull) {
	if grant == nil {
		return
	}
	data.Id = types.StringValue(grant.GetId())
	data.GrantType = types.StringValue(string(grant.GetGrantType()))
	data.Status = types.StringValue(string(grant.GetStatus()))
	if grant.HasEntitlementBundleId() {
		data.EntitlementBundleId = types.StringValue(grant.GetEntitlementBundleId())
	} else {
		data.EntitlementBundleId = types.StringNull()
	}
	data.TargetPrincipalId = types.StringValue(grant.TargetPrincipal.GetExternalId())
	data.TargetPrincipalType = types.StringValue(string(grant.TargetPrincipal.GetType()))
	data.TargetResourceOrn = types.StringValue(grant.GetTargetResourceOrn())
	data.Action = types.StringValue(string(grant.GetAction()))
	data.Actor = types.StringValue(string(grant.GetActor()))
	if grant.HasScheduleSettings() {
		settings := grant.GetScheduleSettings()
		if settings.HasExpirationDate() {
			data.ExpirationDate = types.StringValue(settings.GetExpirationDate().Format(time.RFC3339))
		} else {
			data.ExpirationDate = types.StringNull()
		}
		if settings.HasTimeZone() {
			data.TimeZone = types.StringValue(settings.GetTimeZone())
		} else {
			data.TimeZone = types.StringNull()
		}
	} else {
		data.ExpirationDate = types.StringNull()
		data.TimeZone = types.StringNull()
	}
	if grant.HasEntitlements() {
		ents := grant.GetEntitlements()
		// Sort for deterministic state
		sort.SliceStable(ents, func(i, j int) bool { return ents[i].GetId() < ents[j].GetId() })
		data.Entitlements = make([]grantResourceEntitlementModel, len(ents))
		for i, ent := range ents {
			entModel := grantResourceEntitlementModel{Id: types.StringValue(ent.GetId())}
			if ent.HasValues() {
				vals := ent.GetValues()
				sort.SliceStable(vals, func(a, b int) bool { return vals[a].GetId() < vals[b].GetId() })
				entModel.Values = make([]grantResourceEntitlementValueModel, len(vals))
				for j, val := range vals {
					entModel.Values[j] = grantResourceEntitlementValueModel{Id: types.StringValue(val.GetId())}
				}
			}
			data.Entitlements[i] = entModel
		}
	} else {
		data.Entitlements = []grantResourceEntitlementModel{}
	}
	data.Created = types.StringValue(grant.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(grant.GetCreatedBy())
	data.LastUpdated = types.StringValue(grant.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(grant.GetLastUpdatedBy())
}

func applyGrantWithEntitlementsToState(data *grantResourceModel, grant *governance.GrantFullWithEntitlements) {
	if grant == nil {
		return
	}
	data.Id = types.StringValue(grant.GetId())
	data.GrantType = types.StringValue(string(grant.GetGrantType()))
	data.Status = types.StringValue(string(grant.GetStatus()))
	if grant.HasEntitlementBundleId() {
		data.EntitlementBundleId = types.StringValue(grant.GetEntitlementBundleId())
	} else {
		data.EntitlementBundleId = types.StringNull()
	}
	data.TargetPrincipalId = types.StringValue(grant.TargetPrincipal.GetExternalId())
	data.TargetPrincipalType = types.StringValue(string(grant.TargetPrincipal.GetType()))
	data.TargetResourceOrn = types.StringValue(grant.GetTargetResourceOrn())
	data.Action = types.StringValue(string(grant.GetAction()))
	data.Actor = types.StringValue(string(grant.GetActor()))
	if grant.HasScheduleSettings() {
		settings := grant.GetScheduleSettings()
		if settings.HasExpirationDate() {
			data.ExpirationDate = types.StringValue(settings.GetExpirationDate().Format(time.RFC3339))
		} else {
			data.ExpirationDate = types.StringNull()
		}
		if settings.HasTimeZone() {
			data.TimeZone = types.StringValue(settings.GetTimeZone())
		} else {
			data.TimeZone = types.StringNull()
		}
	} else {
		data.ExpirationDate = types.StringNull()
		data.TimeZone = types.StringNull()
	}
	ents := grant.GetEntitlements()
	sort.SliceStable(ents, func(i, j int) bool { return ents[i].GetId() < ents[j].GetId() })
	data.Entitlements = make([]grantResourceEntitlementModel, len(ents))
	for i, ent := range ents {
		entModel := grantResourceEntitlementModel{Id: types.StringValue(ent.GetId())}
		if ent.HasValues() {
			vals := ent.GetValues()
			sort.SliceStable(vals, func(a, b int) bool { return vals[a].GetId() < vals[b].GetId() })
			entModel.Values = make([]grantResourceEntitlementValueModel, len(vals))
			for j, val := range vals {
				entModel.Values[j] = grantResourceEntitlementValueModel{Id: types.StringValue(val.GetId())}
			}
		}
		data.Entitlements[i] = entModel
	}
	data.Created = types.StringValue(grant.Created.Format(time.RFC3339))
	data.CreatedBy = types.StringValue(grant.GetCreatedBy())
	data.LastUpdated = types.StringValue(grant.LastUpdated.Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(grant.GetLastUpdatedBy())
}
