package okta4ai

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func FWProviderResources() []func() resource.Resource {
	rawResources := []func() resource.Resource{
		NewDelegationLinkResource,
		NewAiAgentsCredentialsJwkRsaResource,
		NewAiAgentResource,
		NewAiAgentsCredentialsJwkEcResource,
		NewApiServerResource,
		NewClientAuthSettingResource,
		NewAiAgentProviderResource,
		NewA2aServersAuthorizationServerResource,
		NewApiServersAuthorizationServerResource,
		NewMcpServerResource,
	}
	// Wrap all resources with SafeResource for panic recovery
	return resources.WrapResources(rawResources)
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDelegationLinkDataSource,
		NewAiAgentsCredentialsJwkEcDataSource,
		NewAiAgentsCredentialsJwkRsaDataSource,
		NewOperationDataSource,
		NewAiAgentDataSource,
		NewApiServersAuthorizationServerDataSource,
		NewPotentialConnectionDataSource,
		NewApiServerDataSource,
		NewA2aServerDataSource,
		NewA2aServersAuthorizationServerDataSource,
		NewMcpServerDataSource,
		NewClientAuthSettingDataSource,
	}
}
