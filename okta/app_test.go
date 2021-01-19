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
	appList, err := listApps(context.Background(), c, &appFilters{LabelPrefix: testResourcePrefix}, defaultPaginationLimit)
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
