package idaas

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

// expectedTypeForIDPrefix maps known Okta object ID prefixes to their owner type.
// These prefixes are consistently used across all Okta tenants but are not
// formally documented as a stable API contract.
var expectedTypeForIDPrefix = map[string]string{
	"00u": "USER",
	"00g": "GROUP",
}

// isAlreadyAssignedOwnerError returns true if AssignGroupOwner returned the specific
// 400 error indicating the owner is already assigned to the group.
func isAlreadyAssignedOwnerError(apiResp *okta.APIResponse, err error) bool {
	if err == nil {
		return false
	}
	if apiResp == nil || apiResp.Response == nil || apiResp.StatusCode != 400 {
		return false
	}
	needle := "provided owner is already assigned to this group"
	var oae *okta.GenericOpenAPIError
	if errors.As(err, &oae) {
		if m := oae.Model(); m != nil {
			if oe, ok := m.(okta.Error); ok {
				for _, cause := range oe.GetErrorCauses() {
					if strings.Contains(strings.ToLower(cause.GetErrorSummary()), needle) {
						return true
					}
				}
			}
		}
		// Fallback: check raw body
		if strings.Contains(strings.ToLower(string(oae.Body())), needle) {
			return true
		}
	}
	// Fallback: substring search in error string
	return strings.Contains(strings.ToLower(err.Error()), needle)
}

// apiErrorBody extracts the response body from a GenericOpenAPIError for diagnostic output.
func apiErrorBody(err error) string {
	var oae *okta.GenericOpenAPIError
	if errors.As(err, &oae) {
		return string(oae.Body())
	}
	return ""
}

func safeString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// warnIfIDTypeMismatch emits a warning if the owner ID prefix doesn't match
// the declared type. This is a best-effort check — the prefixes are not
// formally guaranteed by Okta but have been stable for 10+ years.
func warnIfIDTypeMismatch(diags *diag.Diagnostics, id, typ string) {
	if len(id) < 3 {
		return
	}
	prefix := id[:3]
	expected, known := expectedTypeForIDPrefix[prefix]
	if !known {
		return
	}
	if expected != typ {
		diags.AddWarning(
			"owner type may not match ID",
			fmt.Sprintf("Owner %q has prefix %q which typically indicates type %q, but type %q was declared. Verify the owner ID and type are correct.", id, prefix, expected, typ),
		)
	}
}

// assignOwnerErrorDetail returns an error message that always includes the
// raw API error. When the failure is a 404 and the owner ID prefix suggests
// a type mismatch, a hint is prepended to help the user find the root cause.
func assignOwnerErrorDetail(apiResp *okta.APIResponse, err error, ownerID, declaredType string) string {
	body := apiErrorBody(err)
	raw := fmt.Sprintf("error=%s", err)
	if body != "" {
		raw = fmt.Sprintf("error=%s body=%s", err, body)
	}

	is404 := apiResp != nil && apiResp.Response != nil && apiResp.StatusCode == 404
	if is404 && len(ownerID) >= 3 {
		prefix := ownerID[:3]
		if expected, known := expectedTypeForIDPrefix[prefix]; known && expected != declaredType {
			return fmt.Sprintf(
				"Owner %q was not found as type %s. The ID prefix %q typically indicates a %s entity — check that the type is correct. %s",
				ownerID, declaredType, prefix, expected, raw,
			)
		}
	}
	if is404 {
		return fmt.Sprintf(
			"Owner %q of type %s was not found — the entity may not exist or may have been deleted. %s",
			ownerID, declaredType, raw,
		)
	}

	return raw
}
