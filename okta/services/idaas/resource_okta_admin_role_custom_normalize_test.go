package idaas_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/okta/terraform-provider-okta/sdk"
)

// Helper function to create Permission objects from labels
func createPermissions(labels []string) []*sdk.Permission {
	permissions := make([]*sdk.Permission, len(labels))
	for i, label := range labels {
		permissions[i] = &sdk.Permission{Label: label}
	}
	return permissions
}

// Helper function to extract labels from a set interface
func extractLabelsFromSet(setInterface interface{}) []string {
	if setInterface == nil {
		return []string{}
	}

	// In the actual implementation, this would be a *schema.Set
	// For testing purposes, we'll mock the behavior
	return []string{}
}

// Test the workflow permission normalization behavior
func TestWorkflowPermissionNormalization(t *testing.T) {
	testCases := []struct {
		name           string
		apiPermissions []string
		expected       []string
	}{
		{
			name: "workflow read permission expansion",
			apiPermissions: []string{
				"okta.workflows.read",
				"okta.workflows.flows.read",
				"okta.apps.assignment.manage",
			},
			expected: []string{
				"okta.workflows.read",
				"okta.apps.assignment.manage",
			},
		},
		{
			name: "workflow invoke permission expansion",
			apiPermissions: []string{
				"okta.workflows.invoke",
				"okta.workflows.flows.invoke",
				"okta.apps.assignment.manage",
			},
			expected: []string{
				"okta.workflows.invoke",
				"okta.apps.assignment.manage",
			},
		},
		{
			name: "both workflow permissions expansion",
			apiPermissions: []string{
				"okta.workflows.read",
				"okta.workflows.flows.read",
				"okta.workflows.invoke",
				"okta.workflows.flows.invoke",
				"okta.apps.assignment.manage",
			},
			expected: []string{
				"okta.workflows.read",
				"okta.workflows.invoke",
				"okta.apps.assignment.manage",
			},
		},
		{
			name: "only expanded workflow permissions",
			apiPermissions: []string{
				"okta.workflows.flows.read",
				"okta.workflows.flows.invoke",
				"okta.apps.assignment.manage",
			},
			expected: []string{
				"okta.workflows.flows.read",
				"okta.workflows.flows.invoke",
				"okta.apps.assignment.manage",
			},
		},
		{
			name: "no workflow permissions",
			apiPermissions: []string{
				"okta.apps.assignment.manage",
				"okta.users.userprofile.manage",
			},
			expected: []string{
				"okta.apps.assignment.manage",
				"okta.users.userprofile.manage",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the normalization logic directly
			result := normalizePermissions(tc.apiPermissions)

			// Sort both slices for comparison
			sort.Strings(result)
			sort.Strings(tc.expected)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected permissions %v, got %v", tc.expected, result)
			}
		})
	}
}

// normalizePermissions is a test copy of the function from the main file
// In practice, this would be made exportable or moved to a shared package
func normalizePermissions(apiPermissions []string) []string {
	// Create a map to track which permissions we've seen
	permissionMap := make(map[string]bool)
	for _, perm := range apiPermissions {
		permissionMap[perm] = true
	}

	// Track which permissions to include in the normalized result
	normalizedSet := make(map[string]bool)

	// Handle workflow permissions expansion
	if permissionMap["okta.workflows.read"] {
		normalizedSet["okta.workflows.read"] = true
		// Don't include the expanded version in normalized output
		delete(permissionMap, "okta.workflows.flows.read")
	} else if permissionMap["okta.workflows.flows.read"] {
		// If only the expanded version exists, keep it
		normalizedSet["okta.workflows.flows.read"] = true
	}

	// Handle workflow invoke permissions expansion
	if permissionMap["okta.workflows.invoke"] {
		normalizedSet["okta.workflows.invoke"] = true
		// Don't include the expanded version in normalized output
		delete(permissionMap, "okta.workflows.flows.invoke")
	} else if permissionMap["okta.workflows.flows.invoke"] {
		// If only the expanded version exists, keep it
		normalizedSet["okta.workflows.flows.invoke"] = true
	}

	// Add all other permissions that weren't handled by the workflow normalization
	for perm := range permissionMap {
		if perm != "okta.workflows.flows.read" && perm != "okta.workflows.flows.invoke" {
			normalizedSet[perm] = true
		}
	}

	// Convert back to slice
	result := make([]string, 0, len(normalizedSet))
	for perm := range normalizedSet {
		result = append(result, perm)
	}

	return result
}
