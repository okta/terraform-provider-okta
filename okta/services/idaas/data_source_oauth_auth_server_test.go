package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccDataSourceOktaOAuthAuthorizationServer_read(t *testing.T) {
	mgr := newFixtureManager("data-sources", "okta_oauth_authorization_server", t.Name())
	config := mgr.GetFixtures("datasource.tf", t)
	configNoBaseURL := mgr.GetFixtures("datasource_no_base_url.tf", t)
	resourceName := "data.okta_oauth_authorization_server.test"
	resourceNameNoBaseURL := "data.okta_oauth_authorization_server.test_no_base_url"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "issuer"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "token_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "registration_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "response_types_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "response_modes_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "grant_types_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "subject_types_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "scopes_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "token_endpoint_auth_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "claims_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "code_challenge_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "introspection_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "introspection_endpoint_auth_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "revocation_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "revocation_endpoint_auth_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "end_session_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "request_parameter_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "request_object_signing_alg_values_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "device_authorization_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "pushed_authorization_request_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "backchannel_token_delivery_modes_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "backchannel_authentication_request_signing_alg_values_supported"),
					resource.TestCheckResourceAttrSet(resourceName, "dpop_signing_alg_values_supported"),
				),
			},
			{
				Config: configNoBaseURL,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "id"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "issuer"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "authorization_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "token_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "registration_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "response_types_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "response_modes_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "grant_types_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "subject_types_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "scopes_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "token_endpoint_auth_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "claims_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "code_challenge_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "introspection_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "introspection_endpoint_auth_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "revocation_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "revocation_endpoint_auth_methods_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "end_session_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "request_parameter_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "request_object_signing_alg_values_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "device_authorization_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "pushed_authorization_request_endpoint"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "backchannel_token_delivery_modes_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "backchannel_authentication_request_signing_alg_values_supported"),
					resource.TestCheckResourceAttrSet(resourceNameNoBaseURL, "dpop_signing_alg_values_supported"),
				),
			},
		},
	})
}
