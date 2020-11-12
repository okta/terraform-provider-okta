package okta

import (
	"context"
	"fmt"
	"strings"
)

func deleteTestApps(client *testClient) error {
	c, err := oktaConfig()
	if err != nil {
		return err
	}
	appList, err := listApps(c, &appFilters{LabelPrefix: testResourcePrefix})

	if err != nil {
		return err
	}

	var warnings []string
	for _, app := range appList {
		warn := fmt.Sprintf("failed to sweep an application, there may be dangling resources. ID %s, label %s", app.ID, app.Label)
		_, err := client.oktaClient.Application.DeactivateApplication(context.Background(), app.ID)
		if err != nil {
			warnings = append(warnings, warn)
		}

		resp, err := client.oktaClient.Application.DeleteApplication(context.Background(), app.ID)

		if err != nil && is404(resp.StatusCode) {
			warnings = append(warnings, warn)
		}
	}

	if len(warnings) > 0 {
		return fmt.Errorf("sweep failures: %s", strings.Join(warnings, ", "))
	}

	return nil
}
