package governance

import (
	"context"
	"fmt"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/internal/mutexkv"
	"github.com/okta/terraform-provider-okta/sdk"
	"log"

	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

// oktaMutexKV is a global MutexKV for use within this plugin
var oktaMutexKV = mutexkv.NewMutexKV()

//func getOktaClientFromMetadata(meta interface{}) *sdk.Client {
//	return meta.(*config.Config).OktaIDaaSClient.OktaSDKClientV2()
//}

//func getOktaV3ClientFromMetadata(meta interface{}) *oktaInternalGovernance.IGAPIClient {
//	return meta.(*config.Config).OktaIDaaSClient.OktaIGSDKClientV3()
//}

func getAPISupplementFromMetadata(meta interface{}) *sdk.APISupplement {
	return meta.(*config.Config).OktaIDaaSClient.OktaSDKSupplementClient()
}

func getOktaV5ClientFromMetadata(meta interface{}) *oktaInternalGovernance.IGAPIClient {
	c := meta.(*config.Config)
	log.Println("[INFO]Inside oktaV5ClientFromMetadata", c.ApiToken)
	log.Println("[INFO]Printing meta inside getOktaV5ClientFromMetadata", meta)
	return meta.(*config.Config).OktaGovernanceClient.OktaIGSDKClientV5()
}

func logger(meta interface{}) hclog.Logger {
	return meta.(*config.Config).Logger
}

//func getRequestExecutor(m interface{}) *sdk.RequestExecutor {
//	return getOktaClientFromMetadata(m).GetRequestExecutor()
//}

func fwproviderIsClassicOrg(ctx context.Context, config *config.Config) bool {
	return config.IsClassicOrg(ctx)
}

func providerIsClassicOrg(ctx context.Context, m interface{}) bool {
	if config, ok := m.(*config.Config); ok && config.IsClassicOrg(ctx) {
		return true
	}
	return false
}

func FWProviderResources() []func() resource.Resource {
	return []func() resource.Resource{
		newCampaignResource,
	}
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newCampaignDataSource,
	}
}

func ProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func ProviderDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{

		//resources.OktaInternalGovernanceCampaigns:    dataSourceCampaigns(),
	}
}

func stringIsJSON(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %q JSON to not be empty, got %v", k, i)
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		return diag.Errorf("%q contains an invalid JSON: %s", k, err)
	}
	return nil
}

// doNotRetry helper function to flag if provider should be using backoff.Retry
func doNotRetry(m interface{}, err error) bool {
	return m.(*config.Config).TimeOperations.DoNotRetry(err)
}

func datasourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("data-sources", name)
}

func oieOnlyFeatureError(kind, name string) diag.Diagnostics {
	url := fmt.Sprintf("https://registry.terraform.io/providers/okta/okta/latest/docs/%s/%s", kind, string(name[5:]))
	if kind == "resources" {
		kind = "resource"
	}
	if kind == "data-sources" {
		kind = "datasource"
	}
	return diag.Errorf("%q is a %s for OIE Orgs only, see %s", name, kind, url)
}

func resourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("resources", name)
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

func frameworkResourceOIEOnlyFeatureError(name string) fwdiag.Diagnostics {
	return frameworkOIEOnlyFeatureError("resources", name)
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
