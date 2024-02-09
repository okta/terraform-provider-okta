package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

func NewLogStreamDataSource() datasource.DataSource {
	return &logStreamDataSource{}
}

type logStreamDataSource struct {
	config *Config
}

func (d *logStreamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_stream"
}

func (d *logStreamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Log Streams",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the log stream to retrieve, conflicts with `name`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("name"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Description: "Unique name for the Log Stream object, conflicts with `id`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("id"),
					}...),
				},
			},
			"type": schema.StringAttribute{
				Description: "Streaming provider used - aws_eventbridge or splunk_cloud_logstreaming",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Log Stream Status - can either be ACTIVE or INACTIVE only",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "AWS account ID. Required only for 'aws_eventbridge' type",
						Computed:    true,
					},
					"event_source_name": schema.StringAttribute{
						Description: "An alphanumeric name (no spaces) to identify this event source in AWS EventBridge. Required only for 'aws_eventbridge' type",
						Computed:    true,
					},
					"region": schema.StringAttribute{
						Description: "The destination AWS region where event source is located. Required only for 'aws_eventbridge' type",
						Computed:    true,
					},
					"edition": schema.StringAttribute{
						Description: "Edition of the Splunk Cloud instance. Could be one of: 'aws', 'aws_govcloud', 'gcp'. Required only for 'splunk_cloud_logstreaming' type",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "The domain name for Splunk Cloud instance. Don't include http or https in the string. For example: 'acme.splunkcloud.com'. Required only for 'splunk_cloud_logstreaming' type",
						Computed:    true,
					},
					"token": schema.StringAttribute{
						Description: "The HEC token for your Splunk Cloud HTTP Event Collector. Required only for 'splunk_cloud_logstreaming' type",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *logStreamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *logStreamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data logStreamModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var logStreamResp *okta.ListLogStreams200ResponseInner
	if data.ID.ValueString() != "" {
		logStreamResp, _, err = d.config.oktaSDKClientV3.LogStreamAPI.GetLogStream(ctx, data.ID.ValueString()).Execute()
	} else {
		logStreamResp, err = findLogStreamByName(ctx, d.config.oktaSDKClientV3, data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get log stream",
			err.Error(),
		)
		return
	}

	settings := &logStreamSettingsModel{}
	if logStreamResp.LogStreamAws != nil {
		data.ID = types.StringValue(logStreamResp.LogStreamAws.Id)
		data.Name = types.StringValue(logStreamResp.LogStreamAws.Name)
		data.Status = types.StringValue(logStreamResp.LogStreamAws.Status)
		data.Type = types.StringValue(string(logStreamResp.LogStreamAws.Type))

		lsSettings, ok := logStreamResp.LogStreamAws.GetSettingsOk()
		if ok {
			settings.AccountID = types.StringPointerValue(&lsSettings.AccountId)
			settings.EventSourceName = types.StringPointerValue(&lsSettings.EventSourceName)
			if region, ok := lsSettings.GetRegionOk(); ok {
				settings.Region = types.StringValue(string(*region))
			}
		}
	} else if logStreamResp.LogStreamSplunk != nil {
		data.ID = types.StringValue(logStreamResp.LogStreamSplunk.Id)
		data.Name = types.StringValue(logStreamResp.LogStreamSplunk.Name)
		data.Status = types.StringValue(logStreamResp.LogStreamSplunk.Status)
		data.Type = types.StringValue(string(logStreamResp.LogStreamSplunk.Type))

		lsSettings, ok := logStreamResp.LogStreamSplunk.GetSettingsOk()
		if ok {
			if edition, ok := lsSettings.GetEditionOk(); ok {
				settings.Edition = types.StringValue(string(*edition))
			}
			if host, ok := lsSettings.GetHostOk(); ok {
				settings.Host = types.StringValue(*host)
			}
		}
	}

	settingsValue, diags := types.ObjectValueFrom(ctx, data.Settings.AttributeTypes(ctx), settings)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	data.Settings = settingsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findLogStreamByName(ctx context.Context, client *okta.APIClient, name string) (*okta.ListLogStreams200ResponseInner, error) {
	var logStreamListResp []okta.ListLogStreams200ResponseInner
	logStreamListResp, resp, err := client.LogStreamAPI.ListLogStreams(ctx).Execute()
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
				return &logStreamResp, nil
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
