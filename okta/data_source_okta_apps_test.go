package okta

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-framework/diag"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// )

// func TestAccDataSourceOktaApps_read(t *testing.T) {
// 	cfg := helper.TestCaseConfig{
// 		InitProvider: func() *helper.TestProvider {
// 			return &helper.TestProvider{
// 				Provider: NewProvider(),
// 				CheckConfigure: func(context.Context, resource.ProviderConfigRequest, *resource.ProviderConfigResponse) diag.Diagnostics {
// 					return nil
// 				},
// 			}
// 		},
// 		Steps: []helper.TestStep{
// 			{
// 				Config: appsResources,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					helper.TestCheckResourceAttrSet("okta_app_oauth.test1", "id"),
// 					helper.TestCheckResourceAttrSet("okta_app_oauth.test2", "id"),
// 					helper.TestCheckResourceAttrSet("okta_app_oauth.test3", "id"),
// 				),
// 			},
// 			{
// 				Config: appsRead,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					helper.TestCheckResourceAttr("data.okta_apps.test_by_exact_match", "apps.#", "1"),
// 					helper.TestCheckResourceAttrSet("data.okta_apps.test_by_exact_match", "apps.#.id"),
// 					helper.TestCheckResourceAttr("data.okta_apps.test_by_exact_match", "apps.#.label", fmt.Sprintf("testApp_%s_one", buildResourceName(mgr.Seed))),
// 					helper.TestCheckResourceAttr("data.okta_apps.test_by_exact_match", "apps.#.status", "ACTIVE"),

// 					helper.TestCheckResourceAttr("data.okta_apps.test_by_prefix", "apps.#", "2"),

// 					helper.TestCheckResourceAttr("data.okta_apps.test_by_no_match", "apps.#", "0"),
// 				),
// 			},
// 		},
// 	}

// 	helper.Test(t, cfg)
// }

// const appsResources = `
// resource "okta_app_oauth" "test1" {
// 	label          = "testApp_testAcc_replace_with_uuid_one"
// 	type           = "web"
// 	grant_types    = ["implicit", "authorization_code"]
// 	redirect_uris  = ["http://a.com/"]
// 	response_types = ["code", "token", "id_token"]
// 	issuer_mode    = "ORG_URL"
// 	consent_method = "TRUSTED"
// }
// resource "okta_app_oauth" "test2" {
// 	label          = "testApp_testAcc_replace_with_uuid_two"
// 	type           = "web"
// 	grant_types    = ["implicit", "authorization_code"]
// 	redirect_uris  = ["http://b.com/"]
// 	response_types = ["code", "token", "id_token"]
// 	issuer_mode    = "ORG_URL"
// 	consent_method = "TRUSTED"
// }
// resource "okta_app_oauth" "test3" {
// 	label          = "testAppInvalid_testAcc_replace_with_uuid"
// 	type           = "web"
// 	grant_types    = ["implicit", "authorization_code"]
// 	redirect_uris  = ["http://c.com/"]
// 	response_types = ["code", "token", "id_token"]
// 	issuer_mode    = "ORG_URL"
// 	consent_method = "TRUSTED"
// }
// `

// const appsRead = `
// data "okta_apps" "test_by_exact_match" {
// 	label = "testApp_testAcc_replace_with_uuid_one"
// }

// data "okta_apps" "test_by_prefix" {
// 	label_prefix = "testApp_testAcc_replace_with_uuid_"
// }

// data "okta_apps" "test_by_no_match" {
// 	label = "invalidApp_replace_with_uuid"
// }
// `
