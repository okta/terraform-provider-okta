package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	skuCacheMu sync.Mutex
	skuCache   = make(map[config.OktaSKU]bool)
	skuProbed  = make(map[config.OktaSKU]bool)
)

// RequireSKU skips the test if the current org does not have the given SKU.
// During VCR playback, SKU detection is skipped (the cassette was recorded
// against an org that had the SKU).
// Results are cached per process for the lifetime of the test run.
func RequireSKU(t *testing.T, sku config.OktaSKU) {
	t.Helper()

	if os.Getenv("OKTA_VCR_TF_ACC") == "play" {
		return
	}

	skuCacheMu.Lock()
	defer skuCacheMu.Unlock()

	if skuProbed[sku] {
		if !skuCache[sku] {
			t.Skipf("org does not have required SKU %q, skipping test", sku)
		}
		return
	}

	has := probeSKU(t, sku)
	skuCache[sku] = has
	skuProbed[sku] = true

	if !has {
		t.Skipf("org does not have required SKU %q, skipping test", sku)
	}
}

func probeSKU(t *testing.T, sku config.OktaSKU) bool {
	t.Helper()

	switch sku {
	case config.SKUGovernance:
		return probeGovernanceSKU(t)
	default:
		t.Fatalf("unknown SKU %q for test detection", sku)
		return false
	}
}

// probeGovernanceSKU checks for the Governance SKU by hitting the settings endpoint.
// NOTE: /governance/api/v1/settings may require admin credentials. A more robust
// alternative is the end-user catalog endpoint:
//
//	GET /governance/api/v2/my/catalogs/default/entries?filter=not(parent%20pr)&limit=20
//
// which returns 200/401 regardless of admin permissions.
func probeGovernanceSKU(t *testing.T) bool {
	t.Helper()

	orgName := os.Getenv("OKTA_ORG_NAME")
	baseURL := os.Getenv("OKTA_BASE_URL")
	apiToken := os.Getenv("OKTA_API_TOKEN")

	if orgName == "" || baseURL == "" {
		t.Log("OKTA_ORG_NAME or OKTA_BASE_URL not set, cannot probe governance SKU")
		return false
	}

	url := fmt.Sprintf("https://%s.%s/governance/api/v1/settings", orgName, baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Logf("error building governance SKU probe request: %v", err)
		return false
	}

	if apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("SSWS %s", apiToken))
	} else if accessToken := os.Getenv("OKTA_ACCESS_TOKEN"); accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	} else {
		t.Log("no API token or access token available for SKU detection")
		return false
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Logf("error probing governance SKU: %v", err)
		return false
	}
	defer resp.Body.Close()

	has := resp.StatusCode == http.StatusOK
	t.Logf("governance SKU probe: status=%d has_sku=%v", resp.StatusCode, has)
	return has
}
