package idaas

//
//import (
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/okta/terraform-provider-okta/okta/config"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//func TestValidateMemberURLs_WithStandardOrgDomain(t *testing.T) {
//	d := schema.TestResourceDataRaw(t, resourceAdminRoleCustomAssignments().Schema, map[string]interface{}{
//		"members": schema.NewSet(schema.HashString, []interface{}{
//			"https://mycompany.okta.com/api/v1/users/user123",
//			"https://mycompany.okta.com/api/v1/groups/group456",
//		}),
//	})
//
//	cfg := &config.Config{}
//	// Mock the HTTP client URL to return standard org domain
//	cfg.Domain = "mycompany.okta.com"
//
//	diags := validateMemberURLs(d, cfg)
//	assert.False(t, diags.HasError(), "Should not have errors with standard org domain")
//}
//
//func TestValidateMemberURLs_WithCustomDomain(t *testing.T) {
//	d := schema.TestResourceDataRaw(t, resourceAdminRoleCustomAssignments().Schema, map[string]interface{}{
//		"members": schema.NewSet(schema.HashString, []interface{}{
//			"https://custom.domain.com/api/v1/users/user123",
//		}),
//	})
//
//	cfg := &config.Config{}
//	cfg.Domain = "mycompany.okta.com"
//
//	diags := validateMemberURLs(d, cfg)
//	require.True(t, diags.HasError(), "Should have errors with custom domain")
//	assert.Contains(t, diags[0].Detail, "custom domain not supported")
//	assert.Contains(t, diags[0].Detail, "mycompany.okta.com")
//}
//
//func TestValidateMemberURLs_WithMixedDomains(t *testing.T) {
//	d := schema.TestResourceDataRaw(t, resourceAdminRoleCustomAssignments().Schema, map[string]interface{}{
//		"members": schema.NewSet(schema.HashString, []interface{}{
//			"https://mycompany.okta.com/api/v1/users/user123",
//			"https://custom.domain.com/api/v1/groups/group456",
//		}),
//	})
//
//	cfg := &config.Config{}
//	cfg.Domain = "mycompany.okta.com"
//
//	diags := validateMemberURLs(d, cfg)
//	require.True(t, diags.HasError(), "Should have errors with mixed domains")
//	assert.Equal(t, 1, len(diags), "Should have exactly one error for the custom domain URL")
//}
//
//func TestValidateMemberURLs_WithEmptyMembers(t *testing.T) {
//	d := schema.TestResourceDataRaw(t, resourceAdminRoleCustomAssignments().Schema, map[string]interface{}{
//		"members": schema.NewSet(schema.HashString, []interface{}{}),
//	})
//
//	cfg := &config.Config{}
//	cfg.Domain = "mycompany.okta.com"
//
//	diags := validateMemberURLs(d, cfg)
//	assert.False(t, diags.HasError(), "Should not have errors with empty members")
//}
