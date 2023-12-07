package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
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
		logStream *providerLogStream
	)
	if id != "" {
		logStreamResp, _, err := getOktaV3ClientFromMetadata(m).LogStreamAPI.GetLogStream(ctx, id).Execute()
		if err != nil {
			return diag.Errorf("failed to get log stream %s: %v", id, err)
		}
		logStream, err = normalizeLogSteamResponse(logStreamResp)
		if err != nil {
			return diag.Errorf("failed to read log stream properties: %v", err)
		}
	} else {
		logStream, err = findLogStreamByName(ctx, m, name)
		if err != nil {
			return diag.Errorf("failed to find log stream %q: %v", name, err)
		}
	}

	d.SetId(logStream.Id)
	_ = d.Set("name", logStream.Name)
	_ = d.Set("type", logStream.Type)
	_ = d.Set("status", logStream.Status)

	settings := make(map[string]interface{})
	// aws
	assignIfNotEmpty(&settings, "account_id", logStream.Settings.AccountID)
	assignIfNotEmpty(&settings, "event_source_name", logStream.Settings.EventSourceName)
	assignIfNotEmpty(&settings, "region", logStream.Settings.Region)

	// splunk
	assignIfNotEmpty(&settings, "edition", logStream.Settings.Edition)
	assignIfNotEmpty(&settings, "host", logStream.Settings.Host)

	_ = d.Set("settings", settings)

	return nil
}

func findLogStreamByName(ctx context.Context, m interface{}, name string) (*providerLogStream, error) {
	var logStreamListResp []okta.ListLogStreams200ResponseInner
	logStreamListResp, resp, err := getOktaV3ClientFromMetadata(m).LogStreamAPI.ListLogStreams(ctx).Execute()
	if err != nil {
		return nil, err
	}
	moreStreams := true
	for moreStreams {
		moreStreams = false
		for _, logStreamResp := range logStreamListResp {
			var streamName string
			if logStreamResp.LogStreamAws != nil {
				streamName = logStreamResp.LogStreamAws.Name
			}
			if logStreamResp.LogStreamSplunk != nil {
				streamName = logStreamResp.LogStreamSplunk.Name
			}

			if streamName == name {
				return normalizeLogSteamResponse(&logStreamResp)
			}
		}
		if resp.HasNextPage() {
			moreStreams = true
			var nextLogStreamListResp []okta.ListLogStreams200ResponseInner
			resp, err = resp.Next(&nextLogStreamListResp)
			if err != nil {
				return nil, err
			}
			logStreamListResp = nextLogStreamListResp
		}
	}
	return nil, fmt.Errorf("log stream with name '%s' does not exist", name)
}

func assignIfNotEmpty(m *map[string]interface{}, key string, value string) {
	if value != "" {
		(*m)[key] = value
	}
}
