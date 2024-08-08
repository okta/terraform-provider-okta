package okta

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/okta/okta-sdk-golang/v4/okta"
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

type logStreamResource struct {
	*Config
}

type logStreamModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Status   types.String `tfsdk:"status"`
	Settings types.Object `tfsdk:"settings"`
}
type logStreamSettingsModel struct {
	AccountID       types.String `tfsdk:"account_id"`
	EventSourceName types.String `tfsdk:"event_source_name"`
	Region          types.String `tfsdk:"region"`
	Edition         types.String `tfsdk:"edition"`
	Host            types.String `tfsdk:"host"`
	Token           types.String `tfsdk:"token"`
}

func NewLogStreamResource() resource.Resource {
	return &logStreamResource{}
}

func (r *logStreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_stream"
}

func (r *logStreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages log streams",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Log Stream ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Unique name for the Log Stream object",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Streaming provider used - 'aws_eventbridge' or 'splunk_cloud_logstreaming'",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					// force new
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						logStreamTypeEventBridge,
						logStreamTypeSplunk,
					}...),
				},
			},
			"status": schema.StringAttribute{
				Description: "Stream status",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						statusActive,
						statusInactive,
					}...),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "AWS account ID. Required only for 'aws_eventbridge' type",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthBetween(12, 12),
						},
					},
					"event_source_name": schema.StringAttribute{
						Description: "An alphanumeric name (no spaces) to identify this event source in AWS EventBridge. Required only for 'aws_eventbridge' type",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							// force new
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(awsEventBridgeEventSourceNameRegex, "Event Source must have an alphanumeric name (no spaces) shorter than 76 characters"),
						},
					},
					"region": schema.StringAttribute{
						Description: "The destination AWS region where event source is located. Required only for 'aws_eventbridge' type",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							// force new
							stringplanmodifier.RequiresReplace(),
						},
					},
					"edition": schema.StringAttribute{
						Description: "Edition of the Splunk Cloud instance. Could be one of: 'aws', 'aws_govcloud', 'gcp'. Required only for 'splunk_cloud_logstreaming' type",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								logStreamSplunkEditionAws,
								logStreamSplunkEditionAwsGovCloud,
								logStreamSplunkEditionGcp,
							}...),
						},
					},
					"host": schema.StringAttribute{
						Description: "The domain name for Splunk Cloud instance. Don't include http or https in the string. For example: 'acme.splunkcloud.com'. Required only for 'splunk_cloud_logstreaming' type",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								splunkHostRegex,
								"Splunk host must match the pattern: `^(?!(?:http-inputs-))([a-z0-9]+(-[a-z0-9]+)*){1,100}\\\\.splunkcloud(gc\\\\.com|fed\\\\.com|\\\\.com|\\\\.mil)$`",
							),
						},
					},
					"token": schema.StringAttribute{
						Description: "The HEC token for your Splunk Cloud HTTP Event Collector. Required only for 'splunk_cloud_logstreaming' type",
						Optional:    true,
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							// force new
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								splunkTokenRegex,
								"Splunk token must match the pattern: `(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`",
							),
						},
					},
				},
			},
		},
	}
}

func (r *logStreamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *logStreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state logStreamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	logStreamCreateRequest := r.oktaSDKClientV3.LogStreamAPI.CreateLogStream(ctx)
	logStreamCreateRequestBody := buildLogStreamCreateBody(ctx, &state)
	if logStreamCreateRequestBody == nil {
		resp.Diagnostics.AddError(
			"failed to build log stream create request",
			fmt.Sprintf("unknown type %q", state.Type.ValueString()),
		)
		return
	}
	logStreamCreateRequest = logStreamCreateRequest.Instance(*logStreamCreateRequestBody)
	logStreamResp, _, err := logStreamCreateRequest.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create log stream",
			err.Error(),
		)
		return
	}
	logStream, err := normalizeLogSteamResponse(logStreamResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to normalize log stream",
			err.Error(),
		)
		return
	}

	// NOTE: Okta API doesn't allow directly setting the status on the log
	// stream when it is created. Therefore, we need to compare the operator's
	// intentions in the plan with the API result. See Create Log Stream:
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/LogStream/#tag/LogStream/operation/createLogStream
	planStatus := state.Status.ValueString()
	if planStatus != logStream.Status {
		if planStatus == statusActive {
			_, _, err = r.oktaSDKClientV3.LogStreamAPI.ActivateLogStream(ctx, logStream.Id).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to activate log stream",
					err.Error(),
				)
				return
			}
			logStream.Status = statusActive
		}
		if planStatus == statusInactive {
			_, _, err = r.oktaSDKClientV3.LogStreamAPI.DeactivateLogStream(ctx, logStream.Id).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to deactivate log stream",
					err.Error(),
				)
				return
			}
			logStream.Status = statusInactive
		}
	}

	applyLogStreamToState(ctx, logStream, &state)

	// need to set the "new" state of the log stream model, TF runtime does
	// change detection there
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	// don't need to check for error, we are returning already
}

func (r *logStreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data logStreamModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	logStreamResp, _, err := r.oktaSDKClientV3.LogStreamAPI.GetLogStream(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get log stream",
			err.Error(),
		)
		return
	}
	logStream, err := normalizeLogSteamResponse(logStreamResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to format log stream",
			err.Error(),
		)
		return
	}
	applyLogStreamToState(ctx, logStream, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *logStreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state logStreamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	logStreamReplaceRequest := r.oktaSDKClientV3.LogStreamAPI.ReplaceLogStream(ctx, state.ID.ValueString())
	logStreamReplaceBody := buildLogStreamReplaceBody(ctx, &state)
	if logStreamReplaceBody == nil {
		resp.Diagnostics.AddError(
			"failed to build log stream replace request",
			fmt.Sprintf("unknown type %q", state.Type.ValueString()),
		)
		return
	}
	logStreamReplaceRequest = logStreamReplaceRequest.Instance(*logStreamReplaceBody)
	logStreamResp, _, err := logStreamReplaceRequest.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to replace log stream",
			err.Error(),
		)
		return
	}
	logStream, err := normalizeLogSteamResponse(logStreamResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to normalize log stream",
			err.Error(),
		)
		return
	}

	// NOTE: Okta API doesn't allow directly setting the status on the log
	// stream when it is created. Therefore, we need to compare the operator's
	// intentions in the plan with the API result. See Create Log Stream:
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/LogStream/#tag/LogStream/operation/createLogStream
	planStatus := state.Status.ValueString()
	if planStatus != logStream.Status {
		if planStatus == statusActive {
			_, _, err = r.oktaSDKClientV3.LogStreamAPI.ActivateLogStream(ctx, logStream.Id).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to activate log stream",
					err.Error(),
				)
				return
			}
			logStream.Status = statusActive
		}
		if planStatus == statusInactive {
			_, _, err = r.oktaSDKClientV3.LogStreamAPI.DeactivateLogStream(ctx, logStream.Id).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to deactivate log stream",
					err.Error(),
				)
				return
			}
			logStream.Status = statusInactive
		}
	}

	applyLogStreamToState(ctx, logStream, &state)

	// need to set the "new" state of the log stream model, TF runtime does
	// change detection there
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	// don't need to check for error, we are returning already
}

func (r *logStreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data logStreamModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.oktaSDKClientV3.LogStreamAPI.DeactivateLogStream(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to deactivate log stream",
			err.Error(),
		)
		return
	}
	_, err = r.oktaSDKClientV3.LogStreamAPI.DeleteLogStream(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete log stream",
			err.Error(),
		)
		return
	}
}

func (r *logStreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func applyLogStreamToState(ctx context.Context, ls *providerLogStream, m *logStreamModel) {
	m.ID = types.StringValue(ls.Id)
	m.Name = types.StringValue(ls.Name)
	m.Status = types.StringValue(ls.Status)
	m.Type = types.StringValue(ls.Type)

	settings := &logStreamSettingsModel{}
	m.Settings.As(ctx, settings, basetypes.ObjectAsOptions{})

	if ls.Settings.AccountID != "" {
		settings.AccountID = types.StringValue(ls.Settings.AccountID)
	}
	if ls.Settings.EventSourceName != "" {
		settings.EventSourceName = types.StringValue(ls.Settings.EventSourceName)
	}
	if ls.Settings.Region != "" {
		settings.Region = types.StringValue(ls.Settings.Region)
	}
	if ls.Settings.Edition != "" {
		settings.Edition = types.StringValue(ls.Settings.Edition)
	}
	if ls.Settings.Host != "" {
		settings.Host = types.StringValue(ls.Settings.Host)
	}
	if ls.Settings.Token != "" {
		settings.Token = types.StringValue(ls.Settings.Token)
	}
}

func buildLogStreamCreateBody(ctx context.Context, m *logStreamModel) *okta.ListLogStreams200ResponseInner {
	_type := m.Type.ValueString()
	ls := okta.LogStream{
		Id:     m.ID.ValueString(),
		Name:   m.Name.ValueString(),
		Status: m.Status.ValueString(),
		Type:   m.Type.ValueString(),
	}

	settings := &logStreamSettingsModel{}
	m.Settings.As(ctx, settings, basetypes.ObjectAsOptions{})

	switch _type {
	case logStreamTypeEventBridge:
		var settingAws okta.LogStreamSettingsAws
		settingAws.AccountId = settings.AccountID.ValueString()
		settingAws.EventSourceName = settings.EventSourceName.ValueString()
		settingAws.Region = settings.Region.ValueString()
		return &okta.ListLogStreams200ResponseInner{
			LogStreamAws: &okta.LogStreamAws{
				LogStream: ls,
				Settings:  settingAws,
			},
		}
	case logStreamTypeSplunk:
		var settingSplunk okta.LogStreamSettingsSplunk
		settingSplunk.Edition = settings.Edition.ValueString()
		settingSplunk.Host = settings.Host.ValueString()
		settingSplunk.Token = settings.Token.ValueString()
		return &okta.ListLogStreams200ResponseInner{
			LogStreamSplunk: &okta.LogStreamSplunk{
				LogStream: ls,
				Settings:  settingSplunk,
			},
		}
	}
	return nil
}

func buildLogStreamReplaceBody(ctx context.Context, m *logStreamModel) *okta.ReplaceLogStreamRequest {
	_type := m.Type.ValueString()
	ls := okta.LogStreamPutSchema{
		Name: m.Name.ValueString(),
		Type: m.Type.ValueString(),
	}

	settings := &logStreamSettingsModel{}
	m.Settings.As(ctx, settings, basetypes.ObjectAsOptions{})

	switch _type {
	case logStreamTypeEventBridge:
		var settingAws okta.LogStreamSettingsAws
		settingAws.AccountId = settings.AccountID.ValueString()
		settingAws.EventSourceName = settings.EventSourceName.ValueString()
		settingAws.Region = settings.Region.ValueString()
		return &okta.ReplaceLogStreamRequest{
			LogStreamAwsPutSchema: &okta.LogStreamAwsPutSchema{
				LogStreamPutSchema: ls,
				Settings:           settingAws,
			},
		}
	case logStreamTypeSplunk:
		var settingSplunk okta.LogStreamSettingsSplunkPut
		settingSplunk.Edition = settings.Edition.ValueString()
		settingSplunk.Host = settings.Host.ValueString()
		return &okta.ReplaceLogStreamRequest{
			LogStreamSplunkPutSchema: &okta.LogStreamSplunkPutSchema{
				LogStreamPutSchema: ls,
				Settings:           settingSplunk,
			},
		}
	}
	return nil
}

type providerLogStream struct {
	Id       string
	Name     string
	Status   string
	Type     string
	Settings struct {
		AccountID       string
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
		ls.Settings.AccountID = resp.LogStreamAws.Settings.AccountId
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
