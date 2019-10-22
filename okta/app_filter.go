package okta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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
		ApiFilter         string
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

func (a *appFilters) String() string {
	return fmt.Sprintf(`id: "%s", label: "%s", label_prefix: "%s"`, a.ID, a.Label, a.LabelPrefix)
}

func listApps(m interface{}, filters *appFilters) ([]*appID, error) {
	result := &searchResults{Apps: []*appID{}}
	qp := &query.Params{Limit: 200, Filter: filters.ApiFilter}
	return result.Apps, collectApps(getSupplementFromMetadata(m).RequestExecutor, filters, result, qp)
}

// Recursively list apps until no next links are returned
func collectApps(reqExe *okta.RequestExecutor, filters *appFilters, results *searchResults, qp *query.Params) error {
	req, err := reqExe.NewRequest("GET", fmt.Sprintf("/api/v1/apps?%s", qp.String()), nil)
	if err != nil {
		return err
	}
	var appList []*appID
	res, err := reqExe.Do(req, &appList)
	if err != nil {
		return err
	}

	results.Apps = append(results.Apps, filterApp(appList, filters)...)

	if after := getAfterParam(res); after != "" && !filters.shouldShortCircuit(results.Apps) {
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
	qp := &query.Params{Limit: 200, Filter: filters.ApiFilter}
	return result.SamlApps, collectSamlApps(getSupplementFromMetadata(m).RequestExecutor, filters, result, qp)
}

// Recursively list apps until no next links are returned
func collectSamlApps(reqExe *okta.RequestExecutor, filters *appFilters, results *searchResults, qp *query.Params) error {
	req, err := reqExe.NewRequest("GET", fmt.Sprintf("/api/v1/apps?%s", qp.String()), nil)
	if err != nil {
		return err
	}
	var appList []*okta.SamlApplication
	res, err := reqExe.Do(req, &appList)
	if err != nil {
		return err
	}

	results.SamlApps = append(results.SamlApps, filterSamlApp(appList, filters)...)

	if after := getAfterParam(res); after != "" && !filters.shouldShortCircuit(results.Apps) {
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

func (f *appFilters) shouldSamlShortCircuit(appList []*appID) bool {
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

func getAppFilters(d *schema.ResourceData) (*appFilters, error) {
	id := d.Get("id").(string)
	label := d.Get("label").(string)
	labelPrefix := d.Get("label_prefix").(string)
	filters := &appFilters{ID: id, Label: label, LabelPrefix: labelPrefix}

	if d.Get("active_only").(bool) {
		filters.ApiFilter = `status eq "ACTIVE"`
	}

	if id == "" && label == "" && labelPrefix == "" {
		return nil, errors.New("you must provide either an label_prefix, id, or label to search with")
	}

	return filters, nil
}
