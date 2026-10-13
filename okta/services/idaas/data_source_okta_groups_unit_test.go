package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oktav4 "github.com/okta/okta-sdk-golang/v4/okta"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
)

type testIDaaSClient struct {
	v5Client *v5okta.APIClient
}

func (m *testIDaaSClient) OktaSDKClientV6() *v6okta.APIClient        { return nil }
func (m *testIDaaSClient) OktaSDKClientV5() *v5okta.APIClient         { return m.v5Client }
func (m *testIDaaSClient) OktaSDKClientV3() *oktav4.APIClient         { return nil }
func (m *testIDaaSClient) OktaSDKClientV2() *sdk.Client               { return nil }
func (m *testIDaaSClient) OktaSDKSupplementClient() *sdk.APISupplement { return nil }

// newTestV5Client creates a V5 API client pointed at the given test server.
// We patch cfg.Host after creation because the V5 SDK's NewConfiguration
// calls url.Hostname() which strips the port.
func newTestV5Client(t *testing.T, serverURL string) *v5okta.APIClient {
	v5Cfg, err := v5okta.NewConfiguration(
		v5okta.WithOrgUrl(serverURL),
		v5okta.WithToken("test-token"),
		v5okta.WithAuthorizationMode("SSWS"),
		v5okta.WithTestingDisableHttpsCheck(true),
		v5okta.WithCache(false),
	)
	if err != nil {
		t.Fatalf("creating V5 config: %v", err)
	}
	purl, err := url.Parse(serverURL)
	if err != nil {
		t.Fatalf("parsing server URL: %v", err)
	}
	v5Cfg.Host = purl.Host
	return v5okta.NewAPIClient(v5Cfg)
}

func newTestMeta(t *testing.T, serverURL string) *config.Config {
	meta := &config.Config{}
	meta.SetIdaasAPIClient(&testIDaaSClient{v5Client: newTestV5Client(t, serverURL)})
	return meta
}

func makeGroupJSON(id, name string) map[string]interface{} {
	return map[string]interface{}{
		"id":   id,
		"type": "OKTA_GROUP",
		"profile": map[string]interface{}{
			"name":        name,
			"description": "test group",
		},
	}
}

// https://github.com/okta/terraform-provider-okta/issues/2673
func TestDataSourceGroupsRead_QParameterLimit(t *testing.T) {
	var mu sync.Mutex
	var receivedLimit string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedLimit = r.URL.Query().Get("limit")
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceGroups().Schema, map[string]interface{}{
		"q": "aws-",
	})
	diags := dataSourceGroupsRead(context.Background(), d, newTestMeta(t, server.URL))
	if diags.HasError() {
		t.Fatalf("dataSourceGroupsRead returned errors: %v", diags)
	}

	mu.Lock()
	defer mu.Unlock()
	if receivedLimit != "10000" {
		t.Fatalf("q-based search should send limit=10000 to Okta API, but sent limit=%s", receivedLimit)
	}
}

func TestDataSourceGroupsRead_QParameterReturnsAllGroups(t *testing.T) {
	const totalGroups = 300
	allGroups := make([]map[string]interface{}, totalGroups)
	for i := range allGroups {
		allGroups[i] = makeGroupJSON(fmt.Sprintf("group-%03d", i), fmt.Sprintf("aws-group-%03d", i))
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allGroups)
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceGroups().Schema, map[string]interface{}{
		"q": "aws-",
	})
	diags := dataSourceGroupsRead(context.Background(), d, newTestMeta(t, server.URL))
	if diags.HasError() {
		t.Fatalf("dataSourceGroupsRead returned errors: %v", diags)
	}

	groups := d.Get("groups").([]interface{})
	if len(groups) != totalGroups {
		t.Fatalf("expected %d groups, got %d", totalGroups, len(groups))
	}
}
