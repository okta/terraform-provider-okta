package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

type (
	appID struct {
		ID          string `json:"id"`
		Label       string `json:"label"`
		Name        string `json:"name"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}

	appFilters struct {
		APIFilter         string
		ID                string
		Label             string
		LabelPrefix       string
		ShortCircuitCount int
	}

	searchResults struct {
		Apps     []*appID
		SamlApps []*okta.SamlApplication
		Users    []*okta.User
	}
)

func (f *appFilters) String() string {
	return fmt.Sprintf(`id: "%s", label: "%s", label_prefix: "%s"`, f.ID, f.Label, f.LabelPrefix)
}

func listApps(m interface{}, filters *appFilters) ([]*appID, error) {
	result := &searchResults{Apps: []*appID{}}
	qp := &query.Params{Limit: 200, Filter: filters.APIFilter, Q: filters.getQ()}

	return result.Apps, collectApps(getSupplementFromMetadata(m).RequestExecutor, filters, result, qp)
}

// Recursively list apps until no next links are returned
func collectApps(reqExe *okta.RequestExecutor, filters *appFilters, results *searchResults, qp *query.Params) error {
	req, err := reqExe.NewRequest("GET", fmt.Sprintf("/api/v1/apps%s", qp.String()), nil)
	if err != nil {
		return err
	}
	var appList []*appID
	res, err := reqExe.Do(context.Background(), req, &appList)
	if err != nil {
		return err
	}

	results.Apps = append(results.Apps, filterApp(appList, filters)...)

	// Never attempt to request more if the same "after" link is returned
	if after := sdk.GetAfterParam(res); after != "" && !filters.shouldShortCircuit(results.Apps) && after != qp.After {
		qp.After = after
		return collectApps(reqExe, filters, results, qp)
	}

	return nil
}

func filterApp(appList []*appID, filter *appFilters) []*appID {
	// No filters, return it all!
	if filter.Label == "" && filter.ID == "" && filter.LabelPrefix == "" {
		return appList
	}

	filteredList := []*appID{}
	for _, app := range appList {
		if (filter.ID != "" && filter.ID == app.ID) || (filter.Label != "" && filter.Label == app.Label) {
			filteredList = append(filteredList, app)
		}
		if filter.LabelPrefix != "" && strings.HasPrefix(app.Label, filter.LabelPrefix) {
			filteredList = append(filteredList, app)
		}
	}
	return filteredList
}

// Grabs application q query param
func (f *appFilters) getQ() string {
	if f.Label != "" {
		return f.Label
	}

	return ""
}

func (f *appFilters) shouldShortCircuit(appList []*appID) bool {
	if f.LabelPrefix != "" {
		return false
	}

	if f.ID != "" && f.Label != "" {
		return len(appList) > 1
	}

	if f.ID != "" || f.Label != "" {
		return len(appList) > 0
	}

	return false
}

// Basically a copy paste of listApps, considering adding some code generation but at this point, the juice is
// not worth the squeeze.
func listSamlApps(m interface{}, filters *appFilters) ([]*okta.SamlApplication, error) {
	result := &searchResults{SamlApps: []*okta.SamlApplication{}}
	qp := &query.Params{Limit: 200, Filter: filters.APIFilter}
	return result.SamlApps, collectSamlApps(getSupplementFromMetadata(m).RequestExecutor, filters, result, qp)
}

// Recursively list apps until no next links are returned
func collectSamlApps(reqExe *okta.RequestExecutor, filters *appFilters, results *searchResults, qp *query.Params) error {
	req, err := reqExe.NewRequest("GET", fmt.Sprintf("/api/v1/apps?%s", qp.String()), nil)
	if err != nil {
		return err
	}
	var appList []*okta.SamlApplication
	res, err := reqExe.Do(context.Background(), req, &appList)
	if err != nil {
		return err
	}

	results.SamlApps = append(results.SamlApps, filterSamlApp(appList, filters)...)

	if after := sdk.GetAfterParam(res); after != "" && !filters.shouldShortCircuit(results.Apps) {
		qp.After = after
		return collectApps(reqExe, filters, results, qp)
	}

	return nil
}

func filterSamlApp(appList []*okta.SamlApplication, filter *appFilters) []*okta.SamlApplication {
	// No filters, return it all!
	if filter.Label == "" && filter.ID == "" && filter.LabelPrefix == "" {
		return appList
	}

	filteredList := []*okta.SamlApplication{}
	for _, app := range appList {
		if (filter.ID != "" && filter.ID == app.Id) || (filter.Label != "" && filter.Label == app.Label) {
			filteredList = append(filteredList, app)
		}

		if filter.LabelPrefix != "" && strings.HasPrefix(app.Label, filter.LabelPrefix) {
			filteredList = append(filteredList, app)
		}
	}
	return filteredList
}

func getAppFilters(d *schema.ResourceData) (*appFilters, error) {
	id := d.Get("id").(string)
	label := d.Get("label").(string)
	labelPrefix := d.Get("label_prefix").(string)
	filters := &appFilters{ID: id, Label: label, LabelPrefix: labelPrefix}

	if d.Get("active_only").(bool) {
		filters.APIFilter = fmt.Sprintf(`status eq "%s"`, statusActive)
	}

	if id == "" && label == "" && labelPrefix == "" {
		return nil, errors.New("you must provide either an label_prefix, id, or label to search with")
	}

	return filters, nil
}
