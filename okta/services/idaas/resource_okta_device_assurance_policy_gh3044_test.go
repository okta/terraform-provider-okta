package idaas_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

// GH-3044: Okta sets lastUpdatedBy to the caller on every update, so last_updated_by must
// not be planned as "unchanged" (UseStateForUnknown); otherwise any update of a policy last
// modified by someone else fails with "Provider produced inconsistent result after apply".
func TestDeviceAssuranceLastUpdatedByNotPinnedToState_GH3044(t *testing.T) {
	found := 0
	for _, newResource := range idaas.FWProviderResources() {
		r := newResource()
		meta := &resource.MetadataResponse{}
		r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "okta"}, meta)
		if !strings.HasPrefix(meta.TypeName, "okta_policy_device_assurance_") {
			continue
		}
		found++
		sch := &resource.SchemaResponse{}
		r.Schema(context.Background(), resource.SchemaRequest{}, sch)
		attr, ok := sch.Schema.Attributes["last_updated_by"].(schema.StringAttribute)
		require.Truef(t, ok, "%s: last_updated_by is not a string attribute", meta.TypeName)
		assert.Emptyf(t, attr.PlanModifiers, "%s: last_updated_by must not keep its prior state on update", meta.TypeName)
	}
	assert.Equal(t, 5, found, "expected the 5 okta_policy_device_assurance_* resources")
}
