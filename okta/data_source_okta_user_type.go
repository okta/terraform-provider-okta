package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/sdk"
)

type userTypeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
}

func NewUserTypeDataSource() datasource.DataSource {
	return &userTypeDataSource{}
}

type userTypeDataSource struct {
	config *Config
}

func (d *userTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_type"
}

func (d *userTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a user type from Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the user type to retrieve, conflicts with `name`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("name"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of user type to retrieve, conflicts with `id`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("id"),
					}...),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "Display name of user type.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of user type.",
				Computed:    true,
			},
		},
	}
}

func (d *userTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *userTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data userTypeModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var userTypeResp *sdk.UserType
	if data.ID.ValueString() != "" {
		userTypeResp, _, err = d.config.oktaSDKClientV2.UserType.GetUserType(ctx, data.ID.ValueString())
	} else {
		userTypeResp, err = findUserTypeByName(ctx, d.config.oktaSDKClientV2, data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get user type",
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(userTypeResp.Id)
	data.Name = types.StringValue(userTypeResp.Name)
	data.DisplayName = types.StringValue(userTypeResp.DisplayName)
	data.Description = types.StringValue(userTypeResp.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findUserTypeByName(ctx context.Context, client *sdk.Client, name string) (*sdk.UserType, error) {
	var userType *sdk.UserType
	userTypeListResp, _, err := client.UserType.ListUserTypes(ctx)
	if err != nil {
		return nil, err
	}
	for _, ut := range userTypeListResp {
		if strings.EqualFold(name, ut.Name) {
			userType = ut
			return userType, nil
		}
	}
	return nil, fmt.Errorf("user type '%s' does not exist", name)
}
