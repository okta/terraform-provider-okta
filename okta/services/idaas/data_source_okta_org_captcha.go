// Copyright 2025 - Present Okta, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idaas

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &orgCaptchaDataSource{}
	_ datasource.DataSourceWithConfigure = &orgCaptchaDataSource{}
)

// OrgCaptchaDataSource defines the data source implementation.
type orgCaptchaDataSource struct {
	Config *config.Config
}

// OrgCaptchaDataSourceModel describes the data source data model.
type orgCaptchaDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	CaptchaId    types.String `tfsdk:"captcha_id"`
	EnabledPages types.List   `tfsdk:"enabled_pages"`
}

func newOrgCaptchaDataSource() datasource.DataSource {
	return &orgCaptchaDataSource{}
}

func (d *orgCaptchaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_captcha"
}

func (d *orgCaptchaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *orgCaptchaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the CAPTCHA settings object for your organization > **Note**: If the current organization hasn't configured CAPTCHA Settings, the request returns an empty object.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the org_captcha.",
				Optional:            true,
				Computed:            true,
			},
			"captcha_id": schema.StringAttribute{
				MarkdownDescription: "The unique key of the associated CAPTCHA instance",
				Computed:            true,
			},
			"enabled_pages": schema.ListAttribute{
				MarkdownDescription: "An array of pages that have CAPTCHA enabled",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *orgCaptchaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state orgCaptchaDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.CAPTCHAAPI.GetOrgCaptchaSettings(ctx).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Not Found", "org_captcha with the given ID was not found.")
			return
		}
		resp.Diagnostics.AddError("Error reading org_captcha", err.Error())
		return
	}
	state.ID = types.StringValue("org_captcha")
	_ = result // IDExpr may not reference result; prevent "declared and not used"
	state.CaptchaId = types.StringValue(string(result.GetCaptchaId()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
