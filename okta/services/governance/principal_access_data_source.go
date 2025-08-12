package governance

import (
	"context"
	"fmt"
	"time"

	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &principalAccessDataSource{}

func newPrincipalAccessDataSource() datasource.DataSource {
	return &principalAccessDataSource{}
}

func (d *principalAccessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type principalAccessDataSource struct {
	*config.Config
}

type principalAccessDataSourceModel struct {
	TargetPrincipalOrn types.String    `tfsdk:"target_principal_orn"`
	ParentResourceOrn  types.String    `tfsdk:"parent_resource_orn"`
	ExpirationTime     types.String    `tfsdk:"expiration_time"`
	TimeZone           types.String    `tfsdk:"time_zone"`
	TargetPrincipal    *principalModel `tfsdk:"target_principal"`
	Parent             *parentModel    `tfsdk:"parent"`
	Base               *grantModel     `tfsdk:"base"`
	Additional         []grantModel    `tfsdk:"additional"`
}

type grantModel struct {
	GrantType      types.String                      `tfsdk:"grant_type"`
	GrantMethod    types.String                      `tfsdk:"grant_method"`
	ExpirationTime types.String                      `tfsdk:"expiration_time"`
	StartTime      types.String                      `tfsdk:"start_time"` // Only used in `additional`
	TimeZone       types.String                      `tfsdk:"time_zone"`
	Grant          *grantDetailModel                 `tfsdk:"grant"`
	Bundle         *bundleModel                      `tfsdk:"bundle"` // Only used in `additional`
	Entitlements   []entitlementPrincipalAccessModel `tfsdk:"entitlements"`
}

type grantDetailModel struct {
	Id       types.String        `tfsdk:"id"`
	Metadata *grantMetadataModel `tfsdk:"metadata"`
	Self     *selfLinkModel      `tfsdk:"self"`
}

type bundleModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type grantMetadataModel struct {
	Collection *collectionModel `tfsdk:"collection"`
}

type selfLinkModel struct {
	Href types.String `tfsdk:"href"`
}

type collectionModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type entitlementPrincipalAccessModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ExternalValue types.String `tfsdk:"external_value"`
	Description   types.String `tfsdk:"description"`
	MultiValue    types.Bool   `tfsdk:"multi_value"`
	Required      types.Bool   `tfsdk:"required"`
	DataType      types.String `tfsdk:"data_type"`
	Values        []valueModel `tfsdk:"values"`
}

type valueModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ExternalValue types.String `tfsdk:"external_value"`
}

func (d *principalAccessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_principal_access"
}

func (d *principalAccessDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target_principal_orn": schema.StringAttribute{Computed: true},
			"parent_resource_orn":  schema.StringAttribute{Computed: true},
			"expiration_time":      schema.StringAttribute{Computed: true, Optional: true},
			"time_zone":            schema.StringAttribute{Computed: true, Optional: true},
		},
		Blocks: map[string]schema.Block{
			"target_principal": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{Optional: true, Computed: true},
					"type":        schema.StringAttribute{Optional: true, Computed: true},
				},
			},
			"parent": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{Optional: true, Computed: true},
					"type":        schema.StringAttribute{Optional: true, Computed: true},
				},
			},
			"base": grantBlock(),
			"additional": schema.ListNestedBlock{
				NestedObject: grantBlockObject(),
			},
		},
	}
}

func grantBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"grant_type":      schema.StringAttribute{Computed: true},
			"grant_method":    schema.StringAttribute{Computed: true},
			"expiration_time": schema.StringAttribute{Computed: true, Optional: true},
			"start_time":      schema.StringAttribute{Computed: true, Optional: true},
			"time_zone":       schema.StringAttribute{Computed: true, Optional: true},
		},
		Blocks: grantNestedBlocks(),
	}
}

func grantBlockObject() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"grant_type":      schema.StringAttribute{Computed: true},
			"grant_method":    schema.StringAttribute{Computed: true},
			"expiration_time": schema.StringAttribute{Computed: true, Optional: true},
			"start_time":      schema.StringAttribute{Computed: true, Optional: true},
			"time_zone":       schema.StringAttribute{Computed: true, Optional: true},
		},
		Blocks: grantNestedBlocks(),
	}
}

func grantNestedBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"grant": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{Computed: true},
			},
			Blocks: map[string]schema.Block{
				"metadata": schema.SingleNestedBlock{
					Attributes: map[string]schema.Attribute{},
					Blocks: map[string]schema.Block{
						"collection": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"id":   schema.StringAttribute{Computed: true},
								"name": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
				"self": schema.SingleNestedBlock{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
		"bundle": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"id":   schema.StringAttribute{Computed: true},
				"name": schema.StringAttribute{Computed: true},
			},
		},
		"entitlements": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id":             schema.StringAttribute{Computed: true},
					"name":           schema.StringAttribute{Computed: true},
					"external_value": schema.StringAttribute{Computed: true},
					"description":    schema.StringAttribute{Computed: true},
					"multi_value":    schema.BoolAttribute{Computed: true},
					"required":       schema.BoolAttribute{Computed: true},
					"data_type":      schema.StringAttribute{Computed: true},
				},
				Blocks: map[string]schema.Block{
					"values": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id":             schema.StringAttribute{Computed: true},
								"name":           schema.StringAttribute{Computed: true},
								"external_value": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
		},
	}
}

func (d *principalAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data principalAccessDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	principalAccessResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().PrincipalAccessAPI.GetPrincipalAccess(ctx).Filter(buildFilterForPrincipalAccess(data)).Execute()
	if err != nil {
		return
	}
	// Set top-level fields
	data.TargetPrincipalOrn = types.StringValue(principalAccessResp.TargetPrincipalOrn)
	data.ParentResourceOrn = types.StringValue(principalAccessResp.ParentResourceOrn)
	if principalAccessResp.ExpirationTime != nil {
		data.ExpirationTime = types.StringValue(principalAccessResp.ExpirationTime.Format(time.RFC3339))
	}
	data.TimeZone = types.StringPointerValue(principalAccessResp.TimeZone)

	// Set nested principal
	if &principalAccessResp.TargetPrincipal != nil {
		data.TargetPrincipal = &principalModel{
			ExternalId: types.StringValue(principalAccessResp.TargetPrincipal.ExternalId),
			Type:       types.StringValue(string(principalAccessResp.TargetPrincipal.Type)),
		}
	}

	// Set nested parent
	if &principalAccessResp.Parent != nil {
		data.Parent = &parentModel{
			ExternalId: types.StringValue(principalAccessResp.Parent.ExternalId),
			Type:       types.StringValue(string(principalAccessResp.Parent.Type)),
		}
	}

	// Set base grant
	if principalAccessResp.Base != nil {
		data.Base = convertGrant(*principalAccessResp.Base)
	}

	// Set additional grants
	for _, g := range principalAccessResp.Additional {
		data.Additional = append(data.Additional, *convertGrant(g))
	}

	// Save to Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertGrant(g governance.Grant) *grantModel {
	if &g == nil {
		return nil
	}

	grant := &grantModel{
		GrantType: types.StringValue(string(g.GrantType)),
	}

	if g.GrantMethod != nil {
		grant.GrantMethod = types.StringValue(string(*g.GrantMethod))
	} else {
		grant.GrantMethod = types.StringNull()
	}

	if g.ExpirationTime != nil && !g.ExpirationTime.IsZero() {
		grant.ExpirationTime = types.StringValue(g.ExpirationTime.Format(time.RFC3339))
	} else {
		grant.ExpirationTime = types.StringNull()
	}

	if g.StartTime != nil && !g.StartTime.IsZero() {
		grant.StartTime = types.StringValue(g.StartTime.Format(time.RFC3339))
	} else {
		grant.StartTime = types.StringNull()
	}

	grant.TimeZone = types.StringPointerValue(g.TimeZone)

	// Convert grant details
	if &g.Grant != nil {
		detail := &grantDetailModel{
			Id: types.StringPointerValue(g.Grant.Id),
		}
		//if g.Grant.Self != nil {
		//	detail.Self = &selfLinkModel{
		//		Href: types.StringValue(g.Grant.Self.Href),
		//	}
		//}
		if g.Grant.Metadata != nil && g.Grant.Metadata.Collection != nil {
			detail.Metadata = &grantMetadataModel{
				Collection: &collectionModel{
					Id:   types.StringPointerValue(g.Grant.Metadata.Collection.Id),
					Name: types.StringPointerValue(g.Grant.Metadata.Collection.Name),
				},
			}
		}
		grant.Grant = detail
	}

	// Convert bundle
	if g.Bundle != nil {
		grant.Bundle = &bundleModel{
			Id:   types.StringValue(g.Bundle.Id),
			Name: types.StringValue(g.Bundle.Name),
		}
	}

	// Convert entitlements
	for _, e := range g.Entitlements {
		ent := entitlementPrincipalAccessModel{
			Id:            types.StringValue(e.Id),
			Name:          types.StringValue(e.Name),
			ExternalValue: types.StringValue(e.ExternalValue),
			Description:   types.StringPointerValue(e.Description),
			MultiValue:    types.BoolValue(e.MultiValue),
			Required:      types.BoolValue(e.Required),
			DataType:      types.StringValue(string(e.DataType)),
		}

		for _, v := range e.Values {
			ent.Values = append(ent.Values, valueModel{
				Id:            types.StringValue(v.Id),
				Name:          types.StringValue(v.Name),
				ExternalValue: types.StringValue(v.ExternalValue),
			})
		}

		grant.Entitlements = append(grant.Entitlements, ent)
	}

	return grant
}

func buildFilterForPrincipalAccess(d principalAccessDataSourceModel) string {
	parentExternalId := d.Parent.ExternalId.ValueString()
	targetPrincipalExternalId := d.TargetPrincipal.ExternalId.ValueString()
	parentType := d.Parent.Type.ValueString()
	targetPrincipalType := d.TargetPrincipal.Type.ValueString()
	return fmt.Sprintf(
		`parent.externalId eq "%s" AND parent.type eq "%s" AND targetPrincipal.externalId eq "%s" AND targetPrincipal.type eq "%s"`,
		parentExternalId, // or extract from parent_resource_orn
		parentType,
		targetPrincipalExternalId, // or extract externalId from ORN
		targetPrincipalType,
	)
}
