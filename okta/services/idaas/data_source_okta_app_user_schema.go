package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

var (
	_ datasource.DataSource              = &appUserSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &appUserSchemaDataSource{}
)

type appUserSchemaDataSource struct {
	config *config.Config
}

type appUserSchemaDataSourceModel struct {
	ID             types.String                 `tfsdk:"id"`
	AppID          types.String                 `tfsdk:"app_id"`
	CustomProperty []appUserSchemaPropertyModel `tfsdk:"custom_property"`
}

func newAppUserSchemaDataSource() datasource.DataSource {
	return &appUserSchemaDataSource{}
}

func (d *appUserSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_user_schema"
}

func (d *appUserSchemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *appUserSchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Gets the entire application user schema for an application.

This data source allows you to retrieve all custom properties in an application's user schema without managing them in Terraform.

~> **Note:** This is useful for referencing auto-created properties (from provisioning features like PUSH_NEW_USERS) or reading the complete schema configuration for use in other resources.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The Application's ID (same as app_id)",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "The Application's ID",
			},
		},
		Blocks: map[string]schema.Block{
			"custom_property": schema.SetNestedBlock{
				Description: "Custom properties in the schema",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"index": schema.StringAttribute{
							Computed:    true,
							Description: "The property name/index",
						},
						"title": schema.StringAttribute{
							Computed:    true,
							Description: "Display name for the property",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the schema property. It can be `string`, `boolean`, `number`, `integer`, `array`, or `object`",
						},
						"array_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the array elements if `type` is set to `array`",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the property",
						},
						"required": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the property is required",
						},
						"scope": schema.StringAttribute{
							Computed:    true,
							Description: "Determines whether an app user attribute can be set at the Personal `SELF` or Group `NONE` level.",
						},
						"min_length": schema.Int64Attribute{
							Computed:    true,
							Description: "The minimum length of the property value. Only applies to type `string`",
						},
						"max_length": schema.Int64Attribute{
							Computed:    true,
							Description: "The maximum length of the property value. Only applies to type `string`",
						},
						"enum": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Array of values a primitive property can be set to. See `array_enum` for arrays.",
						},
						"array_enum": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Array of values that an array property's items can be set to.",
						},
						"unique": schema.StringAttribute{
							Computed:    true,
							Description: "Whether the property should be unique. It can be set to `UNIQUE_VALIDATED` or `NOT_UNIQUE`.",
						},
						"external_name": schema.StringAttribute{
							Computed:    true,
							Description: "External name of the property",
						},
						"external_namespace": schema.StringAttribute{
							Computed:    true,
							Description: "External namespace of the property",
						},
						"master": schema.StringAttribute{
							Computed:    true,
							Description: "Master priority for the property. It can be set to `PROFILE_MASTER` or `OKTA`",
						},
						"permissions": schema.StringAttribute{
							Computed:    true,
							Description: "Access control permissions for the property. It can be set to `READ_WRITE`, `READ_ONLY`, or `HIDE`.",
						},
						"union": schema.BoolAttribute{
							Computed:    true,
							Description: "If `type` is set to `array`, used to set whether attribute value is determined by group priority `false`, or combine values across groups `true`.",
						},
					},
					Blocks: map[string]schema.Block{
						"one_of": schema.ListNestedBlock{
							Description: "Array of maps containing a mapping for display name to enum value.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"const": schema.StringAttribute{
										Computed:    true,
										Description: "Enum value",
									},
									"title": schema.StringAttribute{
										Computed:    true,
										Description: "Enum title",
									},
								},
							},
						},
						"array_one_of": schema.ListNestedBlock{
							Description: "Display name and value an enum array can be set to.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"const": schema.StringAttribute{
										Computed:    true,
										Description: "Value mapping to member of `array_enum`",
									},
									"title": schema.StringAttribute{
										Computed:    true,
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

func (d *appUserSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appUserSchemaDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := state.AppID.ValueString()
	client := d.config.OktaIDaaSClient.OktaSDKClientV2()
	us, apiResp, err := client.UserSchema.GetApplicationUserSchema(ctx, appId)
	if err != nil {
		if apiResp != nil && utils.SuppressErrorOn404(apiResp, err) == nil {
			resp.Diagnostics.AddError(
				"Application not found",
				fmt.Sprintf("Application with ID '%s' not found or has no user schema", appId),
			)
			return
		}
		resp.Diagnostics.AddError("Failed to get application user schema", err.Error())
		return
	}

	state.ID = types.StringValue(appId)

	customProps := make([]appUserSchemaPropertyModel, 0)
	if us.Definitions != nil && us.Definitions.Custom != nil && us.Definitions.Custom.Properties != nil {
		for index, attr := range us.Definitions.Custom.Properties {
			customProps = append(customProps, flattenPropertyModel(ctx, index, attr, &resp.Diagnostics))
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}
	state.CustomProperty = customProps

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
