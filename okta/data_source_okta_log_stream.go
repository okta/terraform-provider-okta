package okta

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceLogStream() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogStreamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "ID of the log stream to retrieve, conflicts with `name`.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Unique name for the Log Stream object, conflicts with `id`.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Streaming provider used - aws_eventbridge or splunk_cloud_logstreaming",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Log Stream Status - can either be ACTIVE or INACTIVE only",
			},
			"settings": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		Description: "Gets Okta Log Stream.",
	}
}

func dataSourceLogStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the log stream")
	}
	var (
		err       error
		logStream *sdk.LogStream
	)
	if id != "" {
		logStream, _, err = getOktaClientFromMetadata(m).LogStream.GetLogStream(ctx, id)
	} else {
		logStream, err = findLogStreamByName(ctx, m, name)
	}
	if err != nil {
		return diag.Errorf("failed to find log stream: %v", err)
	}
	d.SetId(logStream.Id)
	_ = d.Set("name", logStream.Name)
	_ = d.Set("type", logStream.Type)
	_ = d.Set("status", logStream.Status)

	settings := make(map[string]interface{})
	assignIfNotEmpty(&settings, "account_id", logStream.Settings.AccountId)
	assignIfNotEmpty(&settings, "event_source_name", logStream.Settings.EventSourceName)
	assignIfNotEmpty(&settings, "region", logStream.Settings.Region)
	assignIfNotEmpty(&settings, "edition", logStream.Settings.Edition)
	assignIfNotEmpty(&settings, "host", logStream.Settings.Host)
	_ = d.Set("settings", settings)

	return nil
}

func findLogStreamByName(ctx context.Context, m interface{}, name string) (*sdk.LogStream, error) {
	client := getOktaClientFromMetadata(m)
	logStreams, resp, err := client.LogStream.ListLogStreams(ctx, nil)
	if err != nil {
		return nil, err
	}
	for i := range logStreams {
		if logStreams[i].Name == name {
			return logStreams[i], nil
		}
	}
	for {
		var moreLogStreams []*sdk.LogStream
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &moreLogStreams)
			if err != nil {
				return nil, err
			}
			for i := range moreLogStreams {
				if moreLogStreams[i].Name == name {
					return moreLogStreams[i], nil
				}
			}
		} else {
			break
		}
	}
	return nil, fmt.Errorf("log stream with name '%s' does not exist", name)
}

func assignIfNotEmpty(m *map[string]interface{}, key string, value string) {
	if value != "" {
		(*m)[key] = value
	}
}
