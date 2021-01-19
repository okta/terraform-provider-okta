package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type appFilters struct {
	Status      string
	ID          string
	Label       string
	LabelPrefix string
}

// Grabs application q query param
func (f *appFilters) getQ() string {
	if f.Label != "" {
		return f.Label
	}
	return f.LabelPrefix
}

func (f *appFilters) String() string {
	return fmt.Sprintf(`id: "%s", label: "%s", label_prefix: "%s"`, f.ID, f.Label, f.LabelPrefix)
}

func listApps(ctx context.Context, m interface{}, filters *appFilters, limit int64) ([]*okta.Application, error) {
	apps, resp, err := getOktaClientFromMetadata(m).Application.
		ListApplications(ctx, &query.Params{Limit: limit, Filter: filters.Status, Q: filters.getQ()})
	if err != nil {
		return nil, err
	}
	resultingApps := make([]*okta.Application, len(apps))
	for i := range apps {
		resultingApps[i] = apps[i].(*okta.Application)
	}
	for resp.HasNextPage() {
		var nextApps []*okta.Application
		resp, err = resp.Next(ctx, &nextApps)
		if err != nil {
			return nil, err
		}
		for i := range nextApps {
			resultingApps = append(resultingApps, nextApps[i])
		}
	}
	return resultingApps, nil
}

func getAppFilters(d *schema.ResourceData) (*appFilters, error) {
	id := d.Get("id").(string)
	label := d.Get("label").(string)
	labelPrefix := d.Get("label_prefix").(string)
	filters := &appFilters{ID: id, Label: label, LabelPrefix: labelPrefix}
	if d.Get("active_only").(bool) {
		filters.Status = fmt.Sprintf(`status eq "%s"`, statusActive)
	}
	if id == "" && label == "" && labelPrefix == "" {
		return nil, errors.New("you must provide either a 'label_prefix', 'id', or 'label' for application search")
	}
	return filters, nil
}
