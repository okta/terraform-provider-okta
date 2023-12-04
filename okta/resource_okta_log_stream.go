package okta

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

const (
	logStreamTypeEventBridge          = "aws_eventbridge"
	logStreamTypeSplunk               = "splunk_cloud_logstreaming"
	logStreamSplunkEditionAws         = "aws"
	logStreamSplunkEditionAwsGovCloud = "aws_govcloud"
	logStreamSplunkEditionGcp         = "gcp"
)

var (
	awsEventBridgeEventSourceNameRegex = regexp.MustCompile(`^[\\.\\-_A-Za-z0-9]{1,75}$`)
	splunkTokenRegex                   = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	splunkHostRegex                    = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*\.splunkcloud(\.gc\.com|\.fed\.com|\.com|\.mil)$`)
)

func resourceLogStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogStreamCreate,
		ReadContext:   resourceLogStreamRead,
		UpdateContext: resourceLogStreamUpdate,
		DeleteContext: resourceLogStreamDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name for the Log Stream object",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Streaming provider used - 'aws_eventbridge' or 'splunk_cloud_logstreaming'",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					logStreamTypeEventBridge, logStreamTypeSplunk,
				}, false)),
			},
			"status": statusSchema,
			"settings": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Description:      "AWS account ID. Required only for 'aws_eventbridge' type",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(12, 12)),
						},
						"event_source_name": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Description:      "An alphanumeric name (no spaces) to identify this event source in AWS EventBridge. Required only for 'aws_eventbridge' type",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(awsEventBridgeEventSourceNameRegex, "Event Source must have an alphanumeric name (no spaces) shorter than 76 characters")),
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The destination AWS region where event source is located. Required only for 'aws_eventbridge' type",
						},
						"edition": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Edition of the Splunk Cloud instance. Could be one of: 'aws', 'aws_govcloud', 'gcp'. Required only for 'splunk_cloud_logstreaming' type",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								logStreamSplunkEditionAws, logStreamSplunkEditionAwsGovCloud, logStreamSplunkEditionGcp,
							}, false)),
						},
						"host": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "The domain name for Splunk Cloud instance. Don't include http or https in the string. For example: 'acme.splunkcloud.com'. Required only for 'splunk_cloud_logstreaming' type",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(splunkHostRegex, "Splunk host must match the pattern: `^(?!(?:http-inputs-))([a-z0-9]+(-[a-z0-9]+)*){1,100}\\\\.splunkcloud(gc\\\\.com|fed\\\\.com|\\\\.com|\\\\.mil)$`")),
						},
						"token": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Sensitive:        true,
							Description:      "The HEC token for your Splunk Cloud HTTP Event Collector. Required only for 'splunk_cloud_logstreaming' type",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(splunkTokenRegex, "Splunk token must match the pattern: `(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`")),
						},
					},
				},
			},
		},
	}
}

func resourceLogStreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateLogStream(d)
	if err != nil {
		return diag.FromErr(err)
	}
	logStreamReq := getOktaV3ClientFromMetadata(m).LogStreamAPI.CreateLogStream(ctx)
	logStreamBody, err := buildLogStreamCreate(d)
	if err != nil {
		return diag.Errorf("invalid log stream format: %v", err)
	}
	logStreamReq = logStreamReq.Instance(*logStreamBody)
	resp, _, err := logStreamReq.Execute()
	if err != nil {
		return diag.Errorf("failed to create log stream: %v", err)
	}
	logStream, err := normalizeLogSteamResponse(resp)
	if err != nil {
		return diag.Errorf("failed to normalize log stream: %v", err)
	}

	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != "" && oldStatus != newStatus {
		var resp *okta.ListLogStreams200ResponseInner
		if newStatus == statusActive {
			resp, _, err = getOktaV3ClientFromMetadata(m).LogStreamAPI.ActivateLogStream(ctx, logStream.Id).Execute()
			if err != nil {
				return diag.Errorf("failed to activate log stream: %v", err)
			}
		} else {
			resp, _, err = getOktaV3ClientFromMetadata(m).LogStreamAPI.DeactivateLogStream(ctx, logStream.Id).Execute()
			if err != nil {
				return diag.Errorf("failed to deactivate log stream: %v", err)
			}
		}

		if resp.LogStreamAws != nil {
			logStream.Status = resp.LogStreamAws.Status
		}
		if resp.LogStreamSplunk != nil {
			logStream.Status = resp.LogStreamSplunk.Status
		}
	}

	diags := fillResourceDataFromLogStream(d, logStream)
	d.SetId(logStream.Id)

	return diags
}

func resourceLogStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logStreamResp, _, err := getOktaV3ClientFromMetadata(m).LogStreamAPI.GetLogStream(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to get log stream: %v", err)
	}
	logStream, err := normalizeLogSteamResponse(logStreamResp)
	if err != nil {
		return diag.Errorf("failed to read log stream properties: %v", err)
	}
	return fillResourceDataFromLogStream(d, logStream)
}

func resourceLogStreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateLogStream(d)
	if err != nil {
		return diag.FromErr(err)
	}
	logStreamReq := getOktaV3ClientFromMetadata(m).LogStreamAPI.ReplaceLogStream(ctx, d.Id())
	logStreamBody, err := buildLogStreamReplace(d)
	if err != nil {
		return diag.Errorf("invalid log stream format: %v", err)
	}
	logStreamReq = logStreamReq.Instance(*logStreamBody)
	resp, _, err := logStreamReq.Execute()
	if err != nil {
		return diag.Errorf("failed to update log stream: %v", err)
	}
	logStream, err := normalizeLogSteamResponse(resp)
	if err != nil {
		return diag.Errorf("failed to normalize log stream: %v", err)
	}

	diags := fillResourceDataFromLogStream(d, logStream)
	d.SetId(logStream.Id)
	return diags
}

func resourceLogStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaV3ClientFromMetadata(m).LogStreamAPI.DeactivateLogStream(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to deactivate log stream: %v", err)
	}
	_, err = getOktaV3ClientFromMetadata(m).LogStreamAPI.DeleteLogStream(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to delete log stream: %v", err)
	}
	return nil
}

func fillResourceDataFromLogStream(d *schema.ResourceData, logStream *providerLogStream) diag.Diagnostics {
	if logStream == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", logStream.Name)
	_ = d.Set("type", logStream.Type)
	_ = d.Set("status", logStream.Status)

	settingsList := d.Get("settings").(*schema.Set).List()
	if len(settingsList) == 1 {
		settings := settingsList[0].(map[string]interface{})
		settings["account_id"] = logStream.Settings.AccountId
		settings["event_source_name"] = logStream.Settings.EventSourceName
		settings["region"] = logStream.Settings.Region
		settings["edition"] = logStream.Settings.Edition
		settings["host"] = logStream.Settings.Host
		settings["token"] = logStream.Settings.Token
	}

	return nil
}

type providerLogStream struct {
	Id       string
	Name     string
	Status   string
	Type     string
	Settings struct {
		AccountId       string
		EventSourceName string
		Region          string
		Edition         string
		Host            string
		Token           string
	}
}

func normalizeLogSteamResponse(resp *okta.ListLogStreams200ResponseInner) (*providerLogStream, error) {
	ls := providerLogStream{}
	if resp.LogStreamAws != nil {
		ls.Id = resp.LogStreamAws.Id
		ls.Name = resp.LogStreamAws.Name
		ls.Status = resp.LogStreamAws.Status
		ls.Type = string(resp.LogStreamAws.Type)
		ls.Settings.AccountId = resp.LogStreamAws.Settings.AccountId
		ls.Settings.EventSourceName = resp.LogStreamAws.Settings.EventSourceName
		ls.Settings.Region = string(resp.LogStreamAws.Settings.Region)
	} else if resp.LogStreamSplunk != nil {
		ls.Id = resp.LogStreamSplunk.Id
		ls.Name = resp.LogStreamSplunk.Name
		ls.Status = resp.LogStreamSplunk.Status
		ls.Type = string(resp.LogStreamSplunk.Type)
		ls.Settings.Edition = string(resp.LogStreamSplunk.Settings.Edition)
		ls.Settings.Host = resp.LogStreamSplunk.Settings.Host
		ls.Settings.Token = resp.LogStreamSplunk.Settings.Token
	} else {
		return nil, fmt.Errorf("log stream is type other than aws or splunk, this is a provider bug")
	}

	return &ls, nil
}

func buildLogStreamReplace(d *schema.ResourceData) (*okta.ReplaceLogStreamRequest, error) {
	var _type interface{}
	var ok bool
	if _type, ok = d.GetOk("type"); !ok {
		return nil, fmt.Errorf("required argument %q not set", "type")
	}
	var lsps okta.LogStreamPutSchema
	if _type.(string) != "" {
		lsps = okta.LogStreamPutSchema{
			Name: d.Get("name").(string),
			Type: okta.LogStreamType(_type.(string)),
		}
	}

	settings := d.Get("settings").(*schema.Set).List()[0].(map[string]interface{})
	awsAccountId := settings["account_id"].(string)
	awsEventSourceName := settings["event_source_name"].(string)
	awsRegion := settings["region"].(string)
	splunkHost := settings["host"].(string)
	splunkEdition := settings["edition"].(string)

	switch _type.(string) {
	case logStreamTypeEventBridge:
		return &okta.ReplaceLogStreamRequest{
			LogStreamAwsPutSchema: &okta.LogStreamAwsPutSchema{
				LogStreamPutSchema: lsps,
				Settings: okta.LogStreamSettingsAws{
					AccountId:       awsAccountId,
					EventSourceName: awsEventSourceName,
					Region:          okta.AwsRegion(awsRegion),
				},
			},
		}, nil
	case logStreamTypeSplunk:
		return &okta.ReplaceLogStreamRequest{
			LogStreamSplunkPutSchema: &okta.LogStreamSplunkPutSchema{
				LogStreamPutSchema: lsps,
				Settings: okta.LogStreamSettingsSplunkPut{
					Edition: okta.SplunkEdition(splunkEdition),
					Host:    splunkHost,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown type %q argument", _type)
	}
}

func buildLogStreamCreate(d *schema.ResourceData) (*okta.ListLogStreams200ResponseInner, error) {
	var _type interface{}
	var ok bool
	if _type, ok = d.GetOk("type"); !ok {
		return nil, fmt.Errorf("required argument %q not set", "type")
	}
	var ls okta.LogStream
	if _type.(string) != "" {
		ls = okta.LogStream{
			Id:     d.Id(),
			Name:   d.Get("name").(string),
			Status: d.Get("status").(string),
			Type:   okta.LogStreamType(_type.(string)),
		}
	}

	settings := d.Get("settings").(*schema.Set).List()[0].(map[string]interface{})
	awsAccountId := settings["account_id"].(string)
	awsEventSourceName := settings["event_source_name"].(string)
	awsRegion := settings["region"].(string)
	splunkHost := settings["host"].(string)
	splunkToken := settings["token"].(string)
	splunkEdition := settings["edition"].(string)

	switch _type.(string) {
	case logStreamTypeEventBridge:
		return &okta.ListLogStreams200ResponseInner{
			LogStreamAws: &okta.LogStreamAws{
				LogStream: ls,
				Settings: okta.LogStreamSettingsAws{
					AccountId:       awsAccountId,
					EventSourceName: awsEventSourceName,
					Region:          okta.AwsRegion(awsRegion),
				},
			},
		}, nil
	case logStreamTypeSplunk:
		return &okta.ListLogStreams200ResponseInner{
			LogStreamSplunk: &okta.LogStreamSplunk{
				LogStream: ls,
				Settings: okta.LogStreamSettingsSplunk{
					Edition: okta.SplunkEdition(splunkEdition),
					Host:    splunkHost,
					Token:   splunkToken,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown type %q argument", _type)
	}
}

func validateLogStream(d *schema.ResourceData) error {
	logStreamType := d.Get("type").(string)
	settings := d.Get("settings").(*schema.Set).List()[0].(map[string]interface{})
	awsAccountId := settings["account_id"].(string)
	awsEventSourceName := settings["event_source_name"].(string)
	awsRegion := settings["region"].(string)
	splunkHost := settings["host"].(string)
	splunkToken := settings["token"].(string)
	splunkEdition := settings["edition"].(string)

	if logStreamType == logStreamTypeEventBridge {
		if awsAccountId == "" {
			return fmt.Errorf(`"account_id" is required for log stream type "aws_eventbridge"`)
		}
		if awsEventSourceName == "" {
			return fmt.Errorf(`"event_source_name" is required for log stream type "aws_eventbridge"`)
		}
		if awsRegion == "" {
			return fmt.Errorf(`"region" is required for log stream type "aws_eventbridge"`)
		}
		if splunkEdition != "" || splunkHost != "" || splunkToken != "" {
			return fmt.Errorf(`"edition", "host" and "token" are not required for log stream type "aws_eventbridge"`)
		}
	} else if logStreamType == logStreamTypeSplunk {
		if splunkHost == "" {
			return fmt.Errorf(`"host" is required for log stream type "splunk_cloud_logstreaming"`)
		}
		if splunkToken == "" {
			return fmt.Errorf(`"token" is required for log stream type "splunk_cloud_logstreaming"`)
		}
		if splunkEdition == "" {
			return fmt.Errorf(`"edition" is required for log stream type "splunk_cloud_logstreaming"`)
		}
		if awsAccountId != "" || awsEventSourceName != "" || awsRegion != "" {
			return fmt.Errorf(`"account_id", "event_source_name" and "region" are not required for log stream type "splunk_cloud_logstreaming"`)
		}

	}

	return nil
}
