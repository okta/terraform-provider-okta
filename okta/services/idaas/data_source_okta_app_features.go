package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*appFeatureDataSource)(nil)

func newAppFeaturesDataSource() datasource.DataSource {
	return &appFeatureDataSource{}
}

type appFeatureDataSource struct {
	*config.Config
}

func (d *appFeatureDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_features"
}

func (d *appFeatureDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *appFeatureDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "`id` used to specify the app feature ID",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "`app_id` used to specify the app ID",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the feature.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Key name of the feature.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Setting status.",
			},
		},
		Blocks: map[string]schema.Block{
			"capabilities": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"create": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"lifecycle_create": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether to update a user in the app when a user in Okta is updated.",
							},
						},
					},
					"update": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"lifecycle_deactivate": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether deprovisioning occurs when the app is unassigned.",
							},
							"password": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"change": schema.StringAttribute{
										Computed:    true,
										Description: "Determines whether a change in a user's password also updates the user's password in the app.",
										Validators: []validator.String{
											stringvalidator.OneOf("CHANGE", "KEEP_EXISTING"),
										},
									},
									"seed": schema.StringAttribute{
										Computed:    true,
										Description: "Determines whether the generated password is the user's Okta password or a randomly generated password.",
										Validators: []validator.String{
											stringvalidator.OneOf("OKTA", "RANDOM"),
										},
									},
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether Okta creates and pushes a password in the app for each assigned user.",
							},
							"profile": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
							},
						},
						Description: "Determines whether updates to a user's profile are pushed to the app.",
					},
					"import_rules": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"user_create_and_match": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"exact_match_criteria": schema.StringAttribute{
										Computed:    true,
										Description: "Determines the attribute to match users.",
									},
									"allow_partial_match": schema.BoolAttribute{
										Computed:    true,
										Description: "Allows user import upon partial matching. Partial matching occurs when the first and last names of an imported user match those of an existing Okta user, even if the username or email attributes don't match.",
									},
									"auto_activate_new_users": schema.BoolAttribute{
										Computed:    true,
										Description: "If set to true, imported new users are automatically activated.",
									},
									"autoconfirm_exact_match": schema.BoolAttribute{
										Computed:    true,
										Description: "If set to true, exact-matched users are automatically confirmed on activation. If set to false, exact-matched users need to be confirmed manually.",
									},
									"autoconfirm_new_users": schema.BoolAttribute{
										Computed:    true,
										Description: "If set to true, imported new users are automatically confirmed on activation. This doesn't apply to imported users that already exist in Okta.",
									},
									"autoconfirm_partial_match": schema.BoolAttribute{
										Computed:    true,
										Description: "If set to true, partially matched users are automatically confirmed on activation. If set to false, partially matched users need to be confirmed manually.",
									},
								},
								Description: "Rules for matching and creating users.",
							},
						},
					},
					"import_settings": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"username": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"username_format": schema.StringAttribute{
										Computed:    true,
										Description: "Determines the username format when users sign in to Okta.",
									},
									"username_expression": schema.StringAttribute{
										Computed:    true,
										Description: "For usernameFormat=CUSTOM, specifies the Okta Expression Language statement for a username format that imported users use to sign in to Okta.",
									},
								},
							},
							"schedule": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Setting status.",
									},
								},
								Blocks: map[string]schema.Block{
									"full_import": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Computed:    true,
												Description: "The import schedule in UNIX cron format.",
											},
											"timezone": schema.StringAttribute{
												Computed:    true,
												Description: "The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format.",
											},
										},
									},
									"incremental_import": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Computed:    true,
												Description: "The import schedule in UNIX cron format.",
											},
											"timezone": schema.StringAttribute{
												Computed:    true,
												Description: "The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format.",
											},
										},
										Description: "Determines the incremental import schedule.",
									},
								},
							},
						},
						Description: "Defines import settings.",
					},
				},
				Description: "Defines the configurations for the USER_PROVISIONING/INBOUND_PROVISIONING feature.",
			},
		},
	}
}

func (d *appFeatureDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data appFeaturesModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getAppFeatureResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.GetFeatureForApplication(ctx, data.AppID.ValueString(), data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app features",
			"Could not read app feature, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(updateAppFeatureState(&data, getAppFeatureResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
