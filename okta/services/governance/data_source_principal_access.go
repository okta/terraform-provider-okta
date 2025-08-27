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
	Id                 types.String    `tfsdk:"id"`
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
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this data source, required by Terraform to track state. This field does not exist in the Okta API response.",
			},
			"target_principal_orn": schema.StringAttribute{Computed: true, Description: "The target principal orn for this data source."},
			"parent_resource_orn":  schema.StringAttribute{Computed: true, Description: "The parent resource orn for this data source."},
			"expiration_time":      schema.StringAttribute{Computed: true, Optional: true, Description: "The date on which the user access expires. Date in ISO 8601 format."},
			"time_zone":            schema.StringAttribute{Computed: true, Optional: true, Description: "The time zone, in IANA format, for the end date of the user access."},
		},
		Blocks: map[string]schema.Block{
			"target_principal": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{Optional: true, Computed: true, Description: "The Okta user id."},
					"type":        schema.StringAttribute{Optional: true, Computed: true, Description: "The type of principal."},
				},
			},
			"parent": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{Optional: true, Computed: true, Description: "The Okta app id of the resource."},
					"type":        schema.StringAttribute{Optional: true, Computed: true, Description: "The type of resource."},
				},
				Description: "Representation of a resource.",
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
			"grant_type":      schema.StringAttribute{Computed: true, Description: "The grant type."},
			"grant_method":    schema.StringAttribute{Computed: true, Description: "Type of grant assignment method."},
			"expiration_time": schema.StringAttribute{Computed: true, Optional: true, Description: "The date on which the user access expires. Date in ISO 8601 format."},
			"start_time":      schema.StringAttribute{Computed: true, Optional: true, Description: "The date on which the user received an access. Date in ISO 8601 format."},
			"time_zone":       schema.StringAttribute{Computed: true, Optional: true, Description: "The time zone, in IANA format."},
		},
		Blocks: grantNestedBlocks(),
	}
}

func grantBlockObject() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"grant_type":      schema.StringAttribute{Computed: true, Description: "The grant type."},
			"grant_method":    schema.StringAttribute{Computed: true, Description: "Type of grant assignment method."},
			"expiration_time": schema.StringAttribute{Computed: true, Optional: true, Description: "The date on which the user access expires. Date in ISO 8601 format."},
			"start_time":      schema.StringAttribute{Computed: true, Optional: true, Description: "The date on which the user received an access. Date in ISO 8601 format."},
			"time_zone":       schema.StringAttribute{Computed: true, Optional: true, Description: "The time zone, in IANA format."},
		},
		Blocks: grantNestedBlocks(),
	}
}

func grantNestedBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"grant": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{Computed: true, Description: "The Grant id."},
			},
			Blocks: map[string]schema.Block{
				"metadata": schema.SingleNestedBlock{
					Blocks: map[string]schema.Block{
						"collection": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"id":   schema.StringAttribute{Computed: true, Description: "The resource collection id."},
								"name": schema.StringAttribute{Computed: true, Description: "The name of a resource collection."},
							},
							Description: "Collection metadata properties.",
						},
					},
				},
				"self": schema.SingleNestedBlock{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{Computed: true, Description: "Link URI"},
					},
				},
			},
		},
		"bundle": schema.SingleNestedBlock{
			Attributes: map[string]schema.Attribute{
				"id":   schema.StringAttribute{Computed: true, Description: "The entitlement bundle id."},
				"name": schema.StringAttribute{Computed: true, Description: "The unique name of the entitlement bundle."},
			},
		},
		"entitlements": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id":             schema.StringAttribute{Computed: true, Description: "The id property of an entitlement."},
					"name":           schema.StringAttribute{Computed: true, Description: "The display name for an entitlement property."},
					"external_value": schema.StringAttribute{Computed: true, Description: "The value of an entitlement property."},
					"description":    schema.StringAttribute{Computed: true, Description: "The description of an entitlement property."},
					"multi_value":    schema.BoolAttribute{Computed: true, Description: "The property that determines if the entitlement property can hold multiple values."},
					"required":       schema.BoolAttribute{Computed: true, Description: "The property that determines if the entitlement property is a required attribute."},
					"data_type":      schema.StringAttribute{Computed: true, Description: "The data type of the entitlement property."},
				},
				Blocks: map[string]schema.Block{
					"values": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id":             schema.StringAttribute{Computed: true, Description: "The id property of an entitlement value."},
								"name":           schema.StringAttribute{Computed: true, Description: "The display name for an entitlement value."},
								"external_value": schema.StringAttribute{Computed: true, Description: "The value of an entitlement value."},
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

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	principalAccessResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().PrincipalAccessAPI.GetPrincipalAccess(ctx).Filter(buildFilterForPrincipalAccess(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Principal Access",
			"Could not read Principal Access, unexpected error: "+err.Error(),
		)
		return
	}
	// Set top-level fields
	data.Id = types.StringValue("principal_access")
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
