package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/stretchr/testify/require"
)

func TestBatchConfiguration(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		values := map[string]interface{}{}
		if enabled {
			values["app_read_batching"] = true
		}
		d := schema.TestResourceDataRaw(t, provider.Provider().Schema, values)
		require.Equal(t, enabled, config.NewConfig(d).AppReadBatching)
	}
}
