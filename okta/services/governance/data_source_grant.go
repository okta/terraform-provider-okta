package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*grantDataSource)(nil)

func newGrantDataSource() datasource.DataSource {
	return &grantDataSource{}
}

type grantDataSource struct {
	*config.Config
}

type grantDataSourceModel struct {
	Id                  types.String                      `tfsdk:"id"`
	GrantType           types.String                      `tfsdk:"grant_type"`
	TargetPrincipalId   types.String                      `tfsdk:"target_principal_id"`
	TargetPrincipalType types.String                      `tfsdk:"target_principal_type"`
	TargetResourceOrn   types.String                      `tfsdk:"target_resource_orn"`
	EntitlementBundleId types.String                      `tfsdk:"entitlement_bundle_id"`
	Entitlements        []grantDataSourceEntitlementModel `tfsdk:"entitlements"`
	Action              types.String                      `tfsdk:"action"`
	Actor               types.String                      `tfsdk:"actor"`
	ExpirationDate      types.String                      `tfsdk:"expiration_date"`
	TimeZone            types.String                      `tfsdk:"time_zone"`
	Status              types.String                      `tfsdk:"status"`
	Created             types.String                      `tfsdk:"created"`
	CreatedBy           types.String                      `tfsdk:"created_by"`
	LastUpdated         types.String                      `tfsdk:"last_updated"`
	LastUpdatedBy       types.String                      `tfsdk:"last_updated_by"`
}

type grantDataSourceEntitlementModel struct {
	Id     types.String                           `tfsdk:"id"`
	Values []grantDataSourceEntitlementValueModel `tfsdk:"values"`
}

type grantDataSourceEntitlementValueModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *grantDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_grant"
}

func (d *grantDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *grantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a grant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the grant.",
			},
			"grant_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of grant (CUSTOM, ENTITLEMENT-BUNDLE, or POLICY).",
			},
			"target_principal_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the principal (user or group) receiving the grant.",
			},
			"target_principal_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of principal (OKTA_USER or OKTA_GROUP).",
			},
			"target_resource_orn": schema.StringAttribute{
				Computed:    true,
				Description: "The ORN of the target resource (e.g., app).",
			},
			"entitlement_bundle_id": schema.StringAttribute{
				Computed:    true,
				Description: "The entitlement bundle ID (for ENTITLEMENT-BUNDLE grants).",
			},
			"action": schema.StringAttribute{
				Computed:    true,
				Description: "The grant action (ALLOW or DENY).",
			},
			"actor": schema.StringAttribute{
				Computed:    true,
				Description: "The actor who created the grant.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "The expiration date for the grant (ISO 8601 format).",
			},
			"time_zone": schema.StringAttribute{
				Computed:    true,
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
				Description: "List of entitlements with their values (for CUSTOM grants).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The entitlement ID.",
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.ListNestedBlock{
							Description: "List of entitlement value IDs.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
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

func (d *grantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data grantDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	grantId := data.Id.ValueString()
	if grantId == "" {
		resp.Diagnostics.AddError("Missing grant ID", "The 'id' attribute must be set in the configuration.")
		return
	}

	// Read API call
	grantResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().GrantsAPI.
		GetGrant(ctx, grantId).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read grant",
			"Could not retrieve grant, unexpected error: "+err.Error(),
		)
		return
	}

	// The response is a union type, need to extract the correct grant type
	if grantResp.GrantFull != nil {
		applyGrantDataSourceToState(&data, grantResp.GrantFull)
	} else if grantResp.GrantFullWithEntitlements != nil {
		applyGrantWithEntitlementsDataSourceToState(&data, grantResp.GrantFullWithEntitlements)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyGrantDataSourceToState(data *grantDataSourceModel, grant *governance.GrantFull) {
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
			data.ExpirationDate = types.StringValue(settings.GetExpirationDate().Format("2006-01-02T15:04:05Z07:00"))
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
		entitlements := grant.GetEntitlements()
		data.Entitlements = make([]grantDataSourceEntitlementModel, len(entitlements))
		for i, ent := range entitlements {
			entModel := grantDataSourceEntitlementModel{
				Id: types.StringValue(ent.GetId()),
			}
			if ent.HasValues() {
				values := ent.GetValues()
				entModel.Values = make([]grantDataSourceEntitlementValueModel, len(values))
				for j, val := range values {
					entModel.Values[j] = grantDataSourceEntitlementValueModel{
						Id: types.StringValue(val.GetId()),
					}
				}
			}
			data.Entitlements[i] = entModel
		}
	}

	data.Created = types.StringValue(grant.Created.Format("2006-01-02T15:04:05Z07:00"))
	data.CreatedBy = types.StringValue(grant.GetCreatedBy())
	data.LastUpdated = types.StringValue(grant.LastUpdated.Format("2006-01-02T15:04:05Z07:00"))
	data.LastUpdatedBy = types.StringValue(grant.GetLastUpdatedBy())
}

func applyGrantWithEntitlementsDataSourceToState(data *grantDataSourceModel, grant *governance.GrantFullWithEntitlements) {
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
			data.ExpirationDate = types.StringValue(settings.GetExpirationDate().Format("2006-01-02T15:04:05Z07:00"))
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

	// GrantFullWithEntitlements has Entitlements field directly
	entitlements := grant.GetEntitlements()
	data.Entitlements = make([]grantDataSourceEntitlementModel, len(entitlements))
	for i, ent := range entitlements {
		entModel := grantDataSourceEntitlementModel{
			Id: types.StringValue(ent.GetId()),
		}
		if ent.HasValues() {
			values := ent.GetValues()
			entModel.Values = make([]grantDataSourceEntitlementValueModel, len(values))
			for j, val := range values {
				entModel.Values[j] = grantDataSourceEntitlementValueModel{
					Id: types.StringValue(val.GetId()),
				}
			}
		}
		data.Entitlements[i] = entModel
	}

	data.Created = types.StringValue(grant.Created.Format("2006-01-02T15:04:05Z07:00"))
	data.CreatedBy = types.StringValue(grant.GetCreatedBy())
	data.LastUpdated = types.StringValue(grant.LastUpdated.Format("2006-01-02T15:04:05Z07:00"))
	data.LastUpdatedBy = types.StringValue(grant.GetLastUpdatedBy())
}
