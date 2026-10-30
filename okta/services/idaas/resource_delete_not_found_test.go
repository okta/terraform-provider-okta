package idaas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
)

// stubIDaaSClient only provides the v2 SDK client used by the SDKv2 resources under test.
type stubIDaaSClient struct {
	api.OktaIDaaSClient
	client *sdk.Client
}

func (s *stubIDaaSClient) OktaSDKClientV2() *sdk.Client { return s.client }

// newDeleteTestMeta returns provider meta whose v2 SDK client talks to a server
// that answers every request with the given status code.
func newDeleteTestMeta(t *testing.T, statusCode int) interface{} {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`{"errorCode":"E0000007","errorSummary":"Not found: Resource not found"}`))
	}))
	t.Cleanup(srv.Close)

	_, client, err := sdk.NewClient(
		context.Background(),
		sdk.WithOrgUrl(srv.URL),
		sdk.WithToken("test-token"),
		sdk.WithTestingDisableHttpsCheck(true),
		sdk.WithHttpClientPtr(srv.Client()),
		sdk.WithRateLimitMaxRetries(0),
		sdk.WithCache(false),
	)
	if err != nil {
		t.Fatalf("failed to create sdk client: %v", err)
	}
	return &config.Config{
		Logger:          hclog.NewNullLogger(),
		OktaIDaaSClient: &stubIDaaSClient{client: client},
	}
}

func TestResourceDeleteSuppressesNotFound(t *testing.T) {
	cases := []struct {
		name     string
		resource *schema.Resource
		raw      map[string]interface{}
	}{
		{
			name:     "okta_group",
			resource: resourceGroup(),
			raw:      map[string]interface{}{"name": "test"},
		},
		{
			name:     "okta_app_group_assignment",
			resource: resourceAppGroupAssignment(),
			raw:      map[string]interface{}{"app_id": "app-id", "group_id": "group-id"},
		},
		{
			name:     "okta_app_user",
			resource: resourceAppUser(),
			raw:      map[string]interface{}{"app_id": "app-id", "user_id": "user-id"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, sc := range []struct {
				statusCode int
				wantErr    bool
			}{
				{http.StatusNotFound, false},
				{http.StatusBadRequest, true},
			} {
				d := schema.TestResourceDataRaw(t, tc.resource.Schema, tc.raw)
				d.SetId("test-id")
				diags := tc.resource.DeleteContext(context.Background(), d, newDeleteTestMeta(t, sc.statusCode))
				if diags.HasError() != sc.wantErr {
					t.Errorf("status %d: got diagnostics %v, want error: %t", sc.statusCode, diags, sc.wantErr)
				}
			}
		})
	}
}
