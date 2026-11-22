package okta4aa

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	okta4aa "github.com/okta/okta-for-ai-agents-golang/okta4aa"

	"github.com/okta/terraform-provider-okta/okta/config"
)

type stubOkta4AIClient struct{ c *okta4aa.APIClient }

func (s stubOkta4AIClient) Okta4AISDKClient() *okta4aa.APIClient { return s.c }

// oauth_client is kept as configured because its request and response shapes differ, but its
// Optional+Computed "type" is still unknown in the plan when the user omits it. Create must resolve
// it from the response; leaving it unknown fails apply with "unknown value after apply".
func TestAiAgentCreate_ResolvesUnsetOauthClientType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"wlp1","status":"ACTIVE","source":"API",
			"profile":{"name":"agent"},
			"oauthClient":{"clientId":"https://agent.example.com/cimd.json","type":"CIMD"}}`))
	}))
	defer server.Close()

	sdkCfg, err := okta4aa.NewConfiguration(
		okta4aa.WithOrgUrl(server.URL),
		okta4aa.WithToken("test-token"),
		okta4aa.WithTestingDisableHttpsCheck(true),
		okta4aa.WithRateLimitMaxRetries(0),
	)
	if err != nil {
		t.Fatal(err)
	}
	sdkCfg.Host = strings.TrimPrefix(server.URL, "http://")
	sdkCfg.Scheme = "http"
	sdkCfg.Okta.Client.AuthorizationMode = "SSWS"
	cfg := &config.Config{}
	cfg.SetOkta4AI(stubOkta4AIClient{okta4aa.NewAPIClient(sdkCfg)})

	r := &aiAgentResource{Config: cfg}
	ctx := context.Background()
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	sch := schemaResp.Schema
	objType := sch.Type().TerraformType(ctx).(tftypes.Object)

	// Plan as Terraform builds it: oauth_client set, its type omitted (so unknown), other computed
	// attributes unknown, unset optional blocks null.
	attrs := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		attrs[name] = tftypes.NewValue(typ, nil)
	}
	for _, name := range []string{"id", "source", "status", "oauth_client_id"} {
		attrs[name] = tftypes.NewValue(objType.AttributeTypes[name], tftypes.UnknownValue)
	}
	oauthType := objType.AttributeTypes["oauth_client"].(tftypes.Object)
	attrs["oauth_client"] = tftypes.NewValue(oauthType, map[string]tftypes.Value{
		"client_id_match_pattern": tftypes.NewValue(tftypes.String, "https://agent.example.com/cimd.json"),
		"type":                    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})
	profileType := objType.AttributeTypes["profile"].(tftypes.Object)
	profile := map[string]tftypes.Value{}
	for name, typ := range profileType.AttributeTypes {
		profile[name] = tftypes.NewValue(typ, nil)
	}
	profile["name"] = tftypes.NewValue(tftypes.String, "agent")
	attrs["profile"] = tftypes.NewValue(profileType, profile)

	plan := tfsdk.Plan{Schema: sch, Raw: tftypes.NewValue(objType, attrs)}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("create failed: %v", resp.Diagnostics)
	}

	if !resp.State.Raw.IsFullyKnown() {
		t.Fatalf("state after Create still has unknown values: %s", resp.State.Raw)
	}
	var got aiAgentModel
	if diags := resp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("decode state: %v", diags)
	}
	if got.OauthClient == nil || got.OauthClient.Type.ValueString() != "CIMD" {
		t.Fatalf("oauth_client.type = %v, want CIMD from the response", got.OauthClient)
	}
	if got.OauthClient.ClientIdMatchPattern.ValueString() != "https://agent.example.com/cimd.json" {
		t.Fatalf("configured client_id_match_pattern changed: %v", got.OauthClient.ClientIdMatchPattern)
	}
}
