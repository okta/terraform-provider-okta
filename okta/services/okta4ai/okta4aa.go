package okta4ai

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func FWProviderResources() []func() resource.Resource {
	rawResources := []func() resource.Resource{
		newDelegationLinkResource,
		newAiAgentResource,
		newAiAgentProviderResource,
		newAgentGatewayResource,
		newClientAuthSettingResource,
		newMcpServerResource,
		newApiServerResource,
		newA2aServersAuthorizationServerResource,
		newMcpServersAuthorizationServerResource,
		newAiAgentsSchemaDefaultResource,
		newAiAgentsCredentialsJwksUriResource,
		newAiAgentsCredentialsPublicKeyModeResource,
		newAiAgentsCredentialsJwkEcResource,
		newAiAgentsCredentialsSecretResource,
		newAiAgentsCredentialsJwkRsaResource,
		newConnectionResource,
		newApiServersAuthorizationServerResource,
	}
	// Wrap all resources with SafeResource for panic recovery
	return resources.WrapResources(rawResources)
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newDelegationLinkDataSource,
		newAiAgentDataSource,
		newAgentGatewayDataSource,
		newAgentGatewaysAuthorizationServerDataSource,
		newAgentGatewaysToolDataSource,
		newMcpServersToolDataSource,
		newMcpServersMetadataDiscoveryDataSource,
		newAiAgentProvidersProfileMappingDataSource,
		newAiAgentProvidersStagingDataSource,
		newAiAgentProvidersStagingStatDataSource,
		newClientAuthSettingDataSource,
		newMcpServerDataSource,
		newA2aServersAuthorizationServerDataSource,
		newMcpServersAuthorizationServerDataSource,
		newAiAgentsSchemaDefaultDataSource,
		newAiAgentsCredentialsJwksUriDataSource,
		newAiAgentsCredentialsPublicKeyModeDataSource,
		newAiAgentsCredentialsSecretDataSource,
		newAiAgentsCredentialsJwkRsaDataSource,
		newAiAgentsCredentialsJwkEcDataSource,
		newApiServerDataSource,
		newConnectionDataSource,
		newAiAgentProviderDataSource,
	}
}
