package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*UISchemaDataSource)(nil)

func newUISchemaDataSource() datasource.DataSource {
	return &UISchemaDataSource{}
}

type UISchemaDataSource struct {
	*config.Config
}

type uiSchemaDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Created     types.String `tfsdk:"created"`
	LastUpdated types.String `tfsdk:"last_updated"`
	UISchema    *uiSchema    `tfsdk:"ui_schema"`
}

func (d *UISchemaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ui_schema"
}

func (d *UISchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *UISchemaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The id property of an UI Schema.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the UI Schema was created (ISO 86001).",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the UI Schema was last modified (ISO 86001).",
			},
		},
		Blocks: map[string]schema.Block{
			"ui_schema": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"button_label": schema.StringAttribute{
						Computed:    true,
						Description: "The Okta app.id of the resource.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of resource.",
					},
					"label": schema.StringAttribute{
						Computed:    true,
						Description: "The type of resource.",
					},
				},
				Blocks: map[string]schema.Block{
					"elements": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"label": schema.StringAttribute{
									Computed:    true,
									Description: "The label of the element.",
								},
								"scope": schema.StringAttribute{
									Computed:    true,
									Description: "The scope of the element.",
								},
								"type": schema.StringAttribute{
									Computed:    true,
									Description: "The type of the element.",
								},
							},
							Blocks: map[string]schema.Block{
								"options": schema.SingleNestedBlock{
									Attributes: map[string]schema.Attribute{
										"format": schema.StringAttribute{
											Computed:    true,
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

func (d *UISchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data uiSchemaDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readUISchemaResp, _, err := d.OktaIDaaSClient.OktaSDKClientV6().UISchemaAPI.GetUISchema(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading UISchema",
			"Could not read UISchema, unexpected error: "+err.Error(),
		)
		return
	}

	data = uiSchemaDataSourceModel{
		ID:          types.StringValue(readUISchemaResp.GetId()),
		Created:     types.StringValue(readUISchemaResp.GetCreated().Format(time.RFC3339)),
		LastUpdated: types.StringValue(readUISchemaResp.GetLastUpdated().Format(time.RFC3339)),
		UISchema:    &uiSchema{},
	}

	if readUISchemaResp.UiSchema.ButtonLabel != nil {
		data.UISchema.ButtonLabel = types.StringValue(readUISchemaResp.UiSchema.GetButtonLabel())
	}
	if readUISchemaResp.UiSchema.Label != nil {
		data.UISchema.Label = types.StringValue(readUISchemaResp.UiSchema.GetLabel())
	}

	if readUISchemaResp.UiSchema.Type != nil {
		data.UISchema.Type = types.StringValue(readUISchemaResp.UiSchema.GetType())
	}
	if readUISchemaResp.UiSchema.Elements != nil {
		var elems []elements
		for _, elem := range readUISchemaResp.UiSchema.Elements {
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
