package config

// OktaSKU represents an Okta product SKU that an org may or may not have provisioned.
//
// SKU detection approaches:
//
//   - Per-SKU endpoint probes: hit a known endpoint and check for 200 vs 401.
//     Current implementation uses this for governance.
//
//   - GET /api/internal/v1/admin/capabilities: internal admin API that returns
//     the full set of purchased capabilities/SKUs for the org. Requires admin
//     credentials. Could replace individual probes in the future.
//
//   - For governance specifically, the end-user catalog endpoint
//     GET /governance/api/v2/my/catalogs/default/entries?filter=not(parent%20pr)&limit=20
//     returns 200/401 regardless of admin permissions, making it a more robust
//     probe than /governance/api/v1/settings which may require admin scope.
type OktaSKU string

const (
	// SKUGovernance represents the Okta Identity Governance product.
	// Detection: GET /governance/api/v1/settings (200 = has SKU, 401 = does not)
	SKUGovernance OktaSKU = "governance"
)
