package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestLogoStateFunc(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "../examples/okta_app_basic_auth/terraform_icon.png",
			expected: "188b6050b43d2fbc9be327e39bf5f7849b120bb4529bcd22cde78b02ccce6777", // compare to `shasum -a 256 filepath`
		},
		{
			input:    "invalid/file/path",
			expected: "",
		},
		{
			input:    "",
			expected: "",
		},
	}
	for _, c := range cases {
		result := localFileStateFunc(c.input)
		if result != c.expected {
			t.Errorf("Error matching logo, expected %q, got %q, for file %q", c.expected, result, c.input)
		}
	}
}

func deleteTestApps(client *testClient) error {
	appList, err := listApps(context.Background(), client.oktaClient, &appFilters{LabelPrefix: testResourcePrefix}, defaultPaginationLimit)
	if err != nil {
		return err
	}
	var warnings []string
	for _, app := range appList {
		warn := fmt.Sprintf("failed to sweep an application, there may be dangling resources. ID %s, label %s", app.Id, app.Label)
		_, err := client.oktaClient.Application.DeactivateApplication(context.Background(), app.Id)
		if err != nil {
			warnings = append(warnings, warn)
		}
		resp, err := client.oktaClient.Application.DeleteApplication(context.Background(), app.Id)
		if is404(resp) {
			warnings = append(warnings, warn)
		} else if err != nil {
			return err
		}
	}
	if len(warnings) > 0 {
		return fmt.Errorf("sweep failures: %s", strings.Join(warnings, ", "))
	}
	return nil
}
