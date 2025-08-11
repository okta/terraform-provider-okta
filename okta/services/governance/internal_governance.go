package governance

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/internal/mutexkv"
)

// oktaMutexKV is a global MutexKV for use within this plugin
var oktaMutexKV = mutexkv.NewMutexKV()

func logger(meta interface{}) hclog.Logger {
	return meta.(*config.Config).Logger
}

func FWProviderResources() []func() resource.Resource {
	return []func() resource.Resource{
		newCampaignResource,
		newEntitlementResource,
	}
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newCampaignDataSource,
		newEntitlementDataSource,
	}
}

func dataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

func resourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	p, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return p
}

func frameworkOIEOnlyFeatureError(kind, name string) fwdiag.Diagnostics {
	url := fmt.Sprintf("https://registry.terraform.io/providers/okta/okta/latest/docs/%s/%s", kind, string(name[5:]))
	if kind == "resources" {
		kind = "resource"
	}
	if kind == "data-sources" {
		kind = "datasource"
	}
	var diags fwdiag.Diagnostics
	diags.AddError(fmt.Sprintf("%q is a %s for OIE Orgs only", name, kind), fmt.Sprintf(", see %s", url))
	return diags
}
