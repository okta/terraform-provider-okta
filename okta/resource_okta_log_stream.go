package okta

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/okta/terraform-provider-okta/sdk"
	"regexp"
)

const (
	logStreamTypeEventBridge          = "aws_eventbridge"
	logStreamTypeSplunk               = "splunk_cloud_logstreaming"
	logStreamSplunkEditionAws         = "aws"
	logStreamSplunkEditionAwsGovCloud = "aws_govcloud"
	logStreamSplunkEditionGcp         = "gcp"
)

var awsEventBridgeEventSourceNameRegex = regexp.MustCompile(`^[\\.\\-_A-Za-z0-9]{1,75}$`)
var splunkTokenRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
var splunkHostRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*\.splunkcloud(\.gc\.com|\.fed\.com|\.com|\.mil)$`)

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
	requestLogStream := buildLogStream(d)
	logStream, _, err := getOktaClientFromMetadata(m).LogStream.CreateLogStream(ctx, requestLogStream)
	if err != nil {
		return diag.Errorf("failed to create log stream: %v", err)
	}

	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != "" && oldStatus != newStatus {
		if newStatus == statusActive {
			logStream, _, err = getOktaClientFromMetadata(m).LogStream.ActivateLogStream(ctx, logStream.Id)
		} else {
			logStream, _, err = getOktaClientFromMetadata(m).LogStream.DeactivateLogStream(ctx, logStream.Id)
		}
		if err != nil {
			return diag.Errorf("failed to change log stream status: %v", err)
		}
	}

	diags := fillResourceDataFromLogStream(d, logStream)
	d.SetId(logStream.Id)

	return diags
}

func resourceLogStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logStream, resp, err := getOktaClientFromMetadata(m).LogStream.GetLogStream(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get log stream: %v", err)
	}
	return fillResourceDataFromLogStream(d, logStream)
}

func resourceLogStreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateLogStream(d)
	if err != nil {
		return diag.FromErr(err)
	}
	logStream, _, err := getOktaClientFromMetadata(m).LogStream.UpdateLogStream(ctx, d.Id(), buildLogStream(d))
	if err != nil {
		return diag.Errorf("failed to update log stream: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == statusActive {
			_, _, err = getOktaClientFromMetadata(m).LogStream.ActivateLogStream(ctx, d.Id())
		} else {
			_, _, err = getOktaClientFromMetadata(m).LogStream.DeactivateLogStream(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change log stream status: %v", err)
		}
	}
	//return resourceLogStreamRead(ctx, d, m)
	return fillResourceDataFromLogStream(d, logStream)
}

func resourceLogStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, _ = getOktaClientFromMetadata(m).LogStream.DeactivateLogStream(ctx, d.Id())
	resp, err := getOktaClientFromMetadata(m).LogStream.DeleteLogStream(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete log stream: %v", err)
	}
	return nil
}

func fillResourceDataFromLogStream(d *schema.ResourceData, logStream *sdk.LogStream) diag.Diagnostics {
	if logStream == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", logStream.Name)
	_ = d.Set("type", logStream.Type)
	_ = d.Set("status", logStream.Status)

	settings := make(map[string]interface{})
	settingsList := d.Get("settings").(*schema.Set).List()
	if len(settingsList) == 1 {
		settings = settingsList[0].(map[string]interface{})

	}
	settings["account_id"] = logStream.Settings.AccountId
	settings["event_source_name"] = logStream.Settings.EventSourceName
	settings["region"] = logStream.Settings.Region
	settings["edition"] = logStream.Settings.Edition
	settings["host"] = logStream.Settings.Host

	return nil
}

func buildLogStream(d *schema.ResourceData) sdk.LogStream {
	settings := d.Get("settings").(*schema.Set).List()[0].(map[string]interface{})
	return sdk.LogStream{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		Settings: &sdk.LogStreamSettings{
			AccountId:       settings["account_id"].(string),
			EventSourceName: settings["event_source_name"].(string),
			Region:          settings["region"].(string),
			Host:            settings["host"].(string),
			Token:           settings["token"].(string),
			Edition:         settings["edition"].(string),
		},
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
