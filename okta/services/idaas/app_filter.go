package idaas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

type AppFilters struct {
	Status      string
	ID          string
	Label       string
	LabelPrefix string
}

// Grabs application q query param
func (f *AppFilters) GetQ() string {
	if f.Label != "" {
		return f.Label
	}
	return f.LabelPrefix
}

func (f *AppFilters) String() string {
	return fmt.Sprintf(`id: "%s", label: "%s", label_prefix: "%s"`, f.ID, f.Label, f.LabelPrefix)
}

func ListApps(ctx context.Context, client *sdk.Client, filters *AppFilters, limit int64) ([]*sdk.Application, error) {
	params := &query.Params{Limit: limit}
	if filters != nil {
		params.Filter = filters.Status
		params.Q = filters.GetQ()
	}
	apps, resp, err := client.Application.ListApplications(ctx, params)
	if err != nil {
		return nil, err
	}
	resultingApps := make([]*sdk.Application, len(apps))
	for i := range apps {
		resultingApps[i] = apps[i].(*sdk.Application)
	}
	for resp.HasNextPage() {
		var nextApps []*sdk.Application
		resp, err = resp.Next(ctx, &nextApps)
		if err != nil {
			return nil, err
		}
		resultingApps = append(resultingApps, nextApps...)
	}
	return resultingApps, nil
}

func getAppFilters(d *schema.ResourceData) (*AppFilters, error) {
	id := d.Get("id").(string)
	label := d.Get("label").(string)
	labelPrefix := d.Get("label_prefix").(string)
	filters := &AppFilters{ID: id, Label: label, LabelPrefix: labelPrefix}
	if d.Get("active_only").(bool) {
		filters.Status = fmt.Sprintf(`status eq "%s"`, StatusActive)
	}
	if id == "" && label == "" && labelPrefix == "" {
		return nil, errors.New("you must provide either a 'label_prefix', 'id', or 'label' for application search")
	}
	return filters, nil
}
