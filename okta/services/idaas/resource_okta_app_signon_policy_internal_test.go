package idaas

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/api"
	"github.com/okta/terraform-provider-okta/okta/config"
)

type v5OnlyClient struct {
	api.OktaIDaaSClient
	v5 *v5okta.APIClient
}

func (c v5OnlyClient) OktaSDKClientV5() *v5okta.APIClient { return c.v5 }

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func readAppSignOnPolicy(t *testing.T, status int) (*resource.ReadResponse, tftypes.Value) {
	t.Helper()
	stub := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"errorCode":"E0000007","errorSummary":"Not found: Resource not found: 00p1 (Policy)"}`)),
		}, nil
	})}

	cfg, err := v5okta.NewConfiguration(
		v5okta.WithOrgUrl("https://example.okta.com"),
		v5okta.WithToken("token"),
		v5okta.WithCache(false),
		v5okta.WithHttpClientPtr(stub),
	)
	if err != nil {
		t.Fatal(err)
	}
	r := &appSignOnPolicyResource{Config: &config.Config{
		OktaIDaaSClient:  v5OnlyClient{v5: v5okta.NewAPIClient(cfg)},
		QueriedWellKnown: true,
	}}

	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	objType := schemaResp.Schema.Type().TerraformType(context.Background())
	attrTypes := objType.(tftypes.Object).AttributeTypes
	values := map[string]tftypes.Value{}
	for name, typ := range attrTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	values["id"] = tftypes.NewValue(tftypes.String, "00p1")
	values["name"] = tftypes.NewValue(tftypes.String, "policy")
	values["description"] = tftypes.NewValue(tftypes.String, "desc")
	raw := tftypes.NewValue(objType, values)

	req := resource.ReadRequest{State: tfsdk.State{Schema: schemaResp.Schema, Raw: raw}}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: raw}}
	r.Read(context.Background(), req, resp)
	return resp, raw
}

func TestAppSignOnPolicyReadRemovesPolicyDeletedOutsideTerraform(t *testing.T) {
	resp, _ := readAppSignOnPolicy(t, http.StatusNotFound)
	if resp.Diagnostics.HasError() {
		t.Fatalf("a 404 should not be an error, got: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatalf("a 404 should remove the resource from state, got: %v", resp.State.Raw)
	}
}

func TestAppSignOnPolicyReadStillFailsOnOtherErrors(t *testing.T) {
	resp, raw := readAppSignOnPolicy(t, http.StatusForbidden)
	if !resp.Diagnostics.HasError() {
		t.Fatal("a 403 should still be reported")
	}
	if got := resp.Diagnostics[0].Summary(); got != "failed to read access policy" {
		t.Fatalf("unexpected diagnostic %q", got)
	}
	if resp.State.Raw.IsNull() || !resp.State.Raw.Equal(raw) {
		t.Fatalf("state should be untouched on a 403, got: %v", resp.State.Raw)
	}
}
