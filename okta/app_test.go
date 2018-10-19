package okta

import (
	"fmt"
	"strings"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/cache"
)

type AppID struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func deleteTestApps(client *testClient) error {
	// Due to https://github.com/okta/okta-sdk-golang/issues/41, have to manually make request to Okta. What a pain!
	// I did not open a PR with a fix to them mostly due to the fact that it would require a design decision.
	c, err := oktaConfig()

	if err != nil {
		return err
	}

	config := okta.NewConfig().
		WithOrgUrl(fmt.Sprintf("https://%v.%v", c.orgName, c.domain)).
		WithToken(c.apiToken).
		WithCache(false)
	requestExecutor := okta.NewRequestExecutor(nil, cache.NewNoOpCache(), config)
	req, err := requestExecutor.NewRequest("GET", "/api/v1/apps", nil)
	if err != nil {
		return err
	}

	var appList []AppID
	_, err = requestExecutor.Do(req, &appList)
	if err != nil {
		return err
	}

	var warnings []string
	for _, app := range appList {
		if strings.HasPrefix(app.Label, testResourcePrefix) {
			warn := fmt.Sprintf("failed to sweep an application, there may be dangling resources. ID %s, label %s", app.ID, app.Label)
			_, err := client.oktaClient.Application.DeactivateApplication(app.ID)
			if err != nil {
				warnings = append(warnings, warn)
			}

			_, err = client.oktaClient.Application.DeleteApplication(app.ID)

			if err != nil {
				warnings = append(warnings, warn)
			}
		}
	}

	if len(warnings) > 0 {
		return fmt.Errorf("sweep failures: %s", strings.Join(warnings, ", "))
	}

	return nil
}
