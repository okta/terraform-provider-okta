package idaas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/okta/terraform-provider-okta/sdk"
)

const testSAMLMetadata = `<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://www.okta.com/exk123">
  <md:IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://example.okta.com/app/exk123/sso/saml"/>
  </md:IDPSSODescriptor>
</md:EntityDescriptor>`

// TestGetAppSAMLMetadata is a regression test for
// https://github.com/okta/terraform-provider-okta/issues/2937.
func TestGetAppSAMLMetadata(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		contentType string
		body        string
		wantFound   bool
		wantErr     bool
	}{
		{"metadata returned", http.StatusOK, "application/xml", testSAMLMetadata, true, false},
		{"404 with okta error body", http.StatusNotFound, "application/json", `{"errorCode":"E0000007","errorSummary":"Not found: Resource not found: 0oa123 (AppInstance)"}`, false, false},
		{"404 with empty body", http.StatusNotFound, "", "", false, false},
		{"500 is still an error", http.StatusInternalServerError, "application/json", `{"errorCode":"E0000009","errorSummary":"Internal Server Error"}`, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotKid, gotAccept string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotKid, gotAccept = r.URL.Path, r.URL.Query().Get("kid"), r.Header.Get("Accept")
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			_, client, err := sdk.NewClient(
				context.Background(),
				sdk.WithOrgUrl(server.URL),
				sdk.WithTestingDisableHttpsCheck(true),
				sdk.WithCache(false),
				sdk.WithToken("test-token"),
				sdk.WithAuthorizationMode("SSWS"),
				sdk.WithRateLimitMaxRetries(0),
			)
			if err != nil {
				t.Fatalf("building sdk client: %v", err)
			}
			supplement := &sdk.APISupplement{RequestExecutor: client.GetRequestExecutor()}

			raw, root, found, err := getAppSAMLMetadata(context.Background(), supplement, "0oa123", "kid456")

			if gotPath != "/api/v1/apps/0oa123/sso/saml/metadata" || gotKid != "kid456" || gotAccept != "application/xml" {
				t.Errorf("request = path %q kid %q accept %q", gotPath, gotKid, gotAccept)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if found != tt.wantFound {
				t.Fatalf("found = %v, want %v", found, tt.wantFound)
			}
			if !found {
				if raw != nil || root != nil {
					t.Errorf("expected nil metadata when not found, got raw=%q root=%v", raw, root)
				}
				return
			}
			if string(raw) != testSAMLMetadata {
				t.Errorf("raw metadata not returned unchanged")
			}
			if root.EntityID != "http://www.okta.com/exk123" || len(root.IDPSSODescriptors) != 1 {
				t.Errorf("unexpected parsed metadata: entityID %q, %d IDP descriptors", root.EntityID, len(root.IDPSSODescriptors))
			}
		})
	}
}
