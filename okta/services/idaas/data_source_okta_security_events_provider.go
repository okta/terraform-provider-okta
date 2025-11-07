package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &securityEventsProviderDataSource{}

func newSecurityEventsProviderDataSource() datasource.DataSource {
	return &securityEventsProviderDataSource{}
}

type securityEventsProviderDataSource struct {
	*config.Config
}

func (d *securityEventsProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_events_provider"
}

func (d *securityEventsProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *securityEventsProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Security Events Provider instance for signal ingestion.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of this instance.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Security Events Provider instance.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The application type of the Security Events Provider.",
			},
			"is_enabled": schema.StringAttribute{
				Computed:    true,
				Description: "Whether or not the Security Events Provider is enabled.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Indicates whether the Security Events Provider is active or not.",
			},
		},
		Blocks: map[string]schema.Block{
			// The required 'settings' block that houses the "one-of" logic
			"settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					// --- Well-Known URL Setting ---
					"well_known_url": schema.StringAttribute{
						Computed:    true, // Optional because it's part of a 'one-of'
						Description: "The published well-known URL of the Security Events Provider (the SSF transmitter).",
						Validators: []validator.String{
							stringvalidator.LengthAtMost(1000),
						},
					},
					// --- Issuer and JWKS Settings ---
					"issuer": schema.StringAttribute{
						Computed:    true, // Optional because it's part of a 'one-of'
						Description: "Issuer URL.",
						Validators: []validator.String{
							stringvalidator.LengthAtMost(700),
						},
					},
					"jwks_url": schema.StringAttribute{
						Computed:    true, // Optional because it's part of a 'one-of'
						Description: "The public URL where the JWKS public key is uploaded.",
						Validators: []validator.String{
							stringvalidator.LengthAtMost(1000),
						},
					},
				},
				Description: "Information about the Security Events Provider for signal ingestion.",
			},
		},
	}
}

func (d *securityEventsProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data securityEventsProviderModel

	// Read the state from the Terraform configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getSecurityEventsProviderResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.GetSecurityEventsProviderInstance(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Security Events Provider",
			"Could not read Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applySecurityEventsProviderToState(getSecurityEventsProviderResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
