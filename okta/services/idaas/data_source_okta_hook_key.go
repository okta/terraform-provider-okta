package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &hookKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &hookKeyDataSource{}
)

type hookKeyDataSource struct {
	*config.Config
}

func newHookKeyDataSource() datasource.DataSource {
	return &hookKeyDataSource{}
}

func (d *hookKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hook_key"
}

func (d *hookKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *hookKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a Hook Key by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name")),
				},
				Description: "The unique identifier of the Hook Key. Conflicts with name.",
			},
			"key_id": schema.StringAttribute{
				Computed:    true,
				Description: "The alias of the public key that can be used to retrieve the public key data.",
			},
			"name": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id")),
				},
				Description: "Display name for the key. Conflicts with id.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the Hook Key was created.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the Hook Key was last updated.",
			},
			"is_used": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Hook Key is currently being used.",
			},
		},
	}
}

func (d *hookKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateData hookKeyModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hookKey, _, err := d.OktaIDaaSClient.OktaSDKClientV5().HookKeyAPI.GetHookKey(ctx, stateData.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading hook key", "Could not read hook key, unexpected error: "+err.Error())
		return
	}
	applyHookKeyToState(&stateData, hookKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...) // Save updated data into Terraform state
}
