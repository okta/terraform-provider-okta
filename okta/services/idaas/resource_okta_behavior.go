package idaas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
)

const (
	behaviorAnomalousLocation = "ANOMALOUS_LOCATION"
	behaviorAnomalousDevice   = "ANOMALOUS_DEVICE"
	behaviorAnomalousIP       = "ANOMALOUS_IP"
	behaviorVelocity          = "VELOCITY"
)

func resourceBehavior() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBehaviorCreate,
		ReadContext:   resourceBehaviorRead,
		UpdateContext: resourceBehaviorUpdate,
		DeleteContext: resourceBehaviorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource allows you to create and configure a behavior.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the behavior",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the behavior. Can be set to `ANOMALOUS_LOCATION`, `ANOMALOUS_DEVICE`, `ANOMALOUS_IP` or `VELOCITY`. Resource will be recreated when the type changes.e",
				ForceNew:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     StatusActive,
				Description: "Behavior status: ACTIVE or INACTIVE. Default: `ACTIVE`",
			},
			"location_granularity_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines the method and level of detail used to evaluate the behavior. Required for `ANOMALOUS_LOCATION` behavior type. Can be set to `LAT_LONG`, `CITY`, `COUNTRY` or `SUBDIVISION`.",
			},
			"radius_from_location": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Radius from location (in kilometers). Should be at least 5. Required when `location_granularity_type` is set to `LAT_LONG`.",
			},
			"number_of_authentications": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of recent authentications used to evaluate the behavior. Required for `ANOMALOUS_LOCATION`, `ANOMALOUS_DEVICE` and `ANOMALOUS_IP` behavior types.",
			},
			"velocity": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Velocity (in kilometers per hour). Should be at least 1. Required for `VELOCITY` behavior",
				ConflictsWith: []string{"number_of_authentications", "radius_from_location", "location_granularity_type"},
			},
		},
	}
}

func validateBehavior(d *schema.ResourceData) error {
	typ := d.Get("type").(string)
	if typ == behaviorAnomalousLocation || typ == behaviorAnomalousDevice || typ == behaviorAnomalousIP {
		_, ok := d.GetOk("number_of_authentications")
		if !ok {
			return fmt.Errorf("'number_of_authentications' should be set for '%s', '%s' and '%s' behavior types", behaviorAnomalousLocation, behaviorAnomalousDevice, behaviorAnomalousDevice)
		}
	}
	if typ == behaviorAnomalousLocation {
		lgt, ok := d.GetOk("location_granularity_type")
		if !ok {
			return fmt.Errorf("'location_granularity_type' should be provided for '%s' behavior type", behaviorAnomalousLocation)
		}
		_, ok = d.GetOk("radius_from_location")
		if lgt.(string) == "LAT_LONG" && !ok {
			return errors.New("'radius_from_location' should be set if location_granularity_type='LAT_LONG'")
		}
	}
	if typ == behaviorVelocity {
		_, ok := d.GetOk("velocity")
		if !ok {
			return fmt.Errorf("'velocity' should be set for '%s' behavior type", behaviorVelocity)
		}
	}
	return nil
}

func resourceBehaviorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("creating behavior", "name", d.Get("name").(string))
	err := validateBehavior(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, rawResp, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.CreateBehaviorDetectionRuleExecute(buildBehavior(ctx, d))
	if err != nil {
		if 200 <= rawResp.StatusCode && rawResp.StatusCode <= 299 {
			if strings.HasPrefix(err.Error(), "parsing time") {
				logger(meta).Info("error when parsing time, will process raw HTTP response")
			}
			if strings.Contains(err.Error(), "cannot unmarshal number") {
				logger(meta).Info("error when parsing number, will process raw HTTP response")
			}
		} else {
			return diag.Errorf("failed to create behavior: %v", err)
		}
	}
	rawRespBody, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(rawRespBody, &respMap)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(respMap["id"].(string))
	return resourceBehaviorRead(ctx, d, meta)
}

func resourceBehaviorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("getting behavior", "id", d.Id())
	_, rawResp, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.GetBehaviorDetectionRule(ctx, d.Id()).Execute()
	if err != nil {
		if 200 <= rawResp.StatusCode && rawResp.StatusCode <= 299 {
			if strings.HasPrefix(err.Error(), "parsing time") {
				logger(meta).Info("error when parsing time, will process raw HTTP response")
			}
			if strings.Contains(err.Error(), "cannot unmarshal number") {
				logger(meta).Info("error when parsing number, will process raw HTTP response")
			}
		} else {
			return diag.Errorf("failed to get behavior: %v", err)
		}
	}
	rawRespBody, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(rawRespBody, &respMap)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", respMap["name"])
	d.Set("type", respMap["type"])
	d.Set("status", respMap["status"])
	typ := respMap["type"].(string)
	settings := respMap["settings"].(map[string]any)
	if typ == behaviorAnomalousLocation || typ == behaviorAnomalousDevice || typ == behaviorAnomalousIP {
		_ = d.Set("number_of_authentications", settings["maxEventsUsedForEvaluation"])
	}
	if typ == behaviorAnomalousLocation {
		_ = d.Set("location_granularity_type", settings["granularity"])
		if settings["granularity"] == "LAT_LONG" {
			_ = d.Set("radius_from_location", settings["radiusKilometers"])
		}
	}
	if typ == behaviorVelocity {
		_ = d.Set("velocity", settings["velocityKph"])
	}
	d.SetId(respMap["id"].(string))
	return nil
}

func resourceBehaviorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("updating behavior", "name", d.Get("name").(string))
	err := validateBehavior(d)
	if err != nil {
		return diag.FromErr(err)
	}
	rule := buildRule(d)
	replaceBehaviorDetectionRule := getOktaV5ClientFromMetadata(meta).BehaviorAPI.ReplaceBehaviorDetectionRule(ctx, d.Id())
	replaceBehaviorDetectionRule = replaceBehaviorDetectionRule.Rule(rule)
	_, rawResp, err := replaceBehaviorDetectionRule.Execute()
	if err != nil {
		if 200 <= rawResp.StatusCode && rawResp.StatusCode <= 299 {
			if strings.HasPrefix(err.Error(), "parsing time") {
				logger(meta).Info("error when parsing time, will process raw HTTP response")
			}
			if strings.Contains(err.Error(), "cannot unmarshal number") {
				logger(meta).Info("error when parsing number, will process raw HTTP response")
			}
		} else {
			return diag.Errorf("failed to update behavior: %v", err)
		}
	}

	if d.HasChange("status") {
		err := handleBehaviorLifecycle(ctx, d, meta)
		if err != nil {
			return err
		}
	}
	return resourceBehaviorRead(ctx, d, meta)
}

func resourceBehaviorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("deleting behavior", "name", d.Get("name").(string))
	deleteBehaviorDetectionRule := getOktaV5ClientFromMetadata(meta).BehaviorAPI.DeleteBehaviorDetectionRule(ctx, d.Id())
	_, err := deleteBehaviorDetectionRule.Execute()
	if err != nil {
		return diag.Errorf("failed to delete behavior: %v", err)
	}
	return nil
}

func buildRule(d *schema.ResourceData) v5okta.ListBehaviorDetectionRules200ResponseInner {
	behaviorType := d.Get("type").(string)
	behaviorName := d.Get("name").(string)
	behaviorStatus := d.Get("status").(string)
	listBehaviorDetectionRules := v5okta.ListBehaviorDetectionRules200ResponseInner{}
	switch behaviorType {

	case behaviorAnomalousDevice:
		behaviorRuleAnomalousDevice := v5okta.NewBehaviorRuleAnomalousDevice(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleAnomalousDeviceAsListBehaviorDetectionRules200ResponseInner(behaviorRuleAnomalousDevice)
		settings := v5okta.NewBehaviorRuleSettingsAnomalousDevice()
		settings.SetMaxEventsUsedForEvaluation(int32(d.Get("number_of_authentications").(int)))
		listBehaviorDetectionRules.BehaviorRuleAnomalousDevice.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleAnomalousDevice.SetStatus(behaviorStatus)

	case behaviorAnomalousIP:
		behaviorRuleAnomalousIP := v5okta.NewBehaviorRuleAnomalousIP(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleAnomalousIPAsListBehaviorDetectionRules200ResponseInner(behaviorRuleAnomalousIP)
		settings := v5okta.NewBehaviorRuleSettingsAnomalousIP()
		settings.SetMaxEventsUsedForEvaluation(int32(d.Get("number_of_authentications").(int)))
		listBehaviorDetectionRules.BehaviorRuleAnomalousIP.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleAnomalousIP.SetStatus(behaviorStatus)

	case behaviorAnomalousLocation:
		behaviorRuleAnomalousLocation := v5okta.NewBehaviorRuleAnomalousLocation(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleAnomalousLocationAsListBehaviorDetectionRules200ResponseInner(behaviorRuleAnomalousLocation)
		granularity := d.Get("location_granularity_type").(string)
		locationGranularityType := d.Get("location_granularity_type").(string)
		settings := v5okta.NewBehaviorRuleSettingsAnomalousLocation(granularity)
		if locationGranularityType == "LAT_LONG" {
			settings.SetRadiusKilometers(int32(d.Get("radius_from_location").(int)))
		}
		settings.SetMaxEventsUsedForEvaluation(int32(d.Get("number_of_authentications").(int)))
		listBehaviorDetectionRules.BehaviorRuleAnomalousLocation.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleAnomalousLocation.SetStatus(behaviorStatus)

	case behaviorVelocity:
		behaviorRuleVelocity := v5okta.NewBehaviorRuleVelocity(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleVelocityAsListBehaviorDetectionRules200ResponseInner(behaviorRuleVelocity)
		velocity := d.Get("velocity").(int)
		settings := v5okta.NewBehaviorRuleSettingsVelocity(int32(velocity))
		listBehaviorDetectionRules.BehaviorRuleVelocity.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleVelocity.SetStatus(behaviorStatus)
	}
	return listBehaviorDetectionRules
}

func buildBehavior(ctx context.Context, d *schema.ResourceData) v5okta.ApiCreateBehaviorDetectionRuleRequest {
	listBehaviorDetectionRules := buildRule(d)
	behaviorAPIService := v5okta.BehaviorAPIService{}
	behaviorDetectionRuleRequest := behaviorAPIService.CreateBehaviorDetectionRule(ctx)
	behaviorDetectionRuleRequest = behaviorDetectionRuleRequest.Rule(listBehaviorDetectionRules)
	return behaviorDetectionRuleRequest
}

func handleBehaviorLifecycle(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	if d.Get("status").(string) == StatusActive {
		logger(meta).Info("activating behavior", "name", d.Get("name").(string))
		activateBehaviorDetectionRuleRequest := client.BehaviorAPI.ActivateBehaviorDetectionRule(ctx, d.Id())
		_, rawResp, err := activateBehaviorDetectionRuleRequest.Execute()
		if err != nil {
			if 200 <= rawResp.StatusCode && rawResp.StatusCode <= 299 {
				if strings.HasPrefix(err.Error(), "parsing time") {
					logger(meta).Info("error when parsing time, will process raw HTTP response")
				}
				if strings.Contains(err.Error(), "cannot unmarshal number") {
					logger(meta).Info("error when parsing number, will process raw HTTP response")
				}
			} else {
				return diag.Errorf("failed to activate behavior: %v", err)
			}
		}
		return nil
	}
	logger(meta).Info("deactivating behavior", "name", d.Get("name").(string))
	deactivateBehaviorDetectionRuleRequest := client.BehaviorAPI.DeactivateBehaviorDetectionRule(ctx, d.Id())
	_, rawResp, err := deactivateBehaviorDetectionRuleRequest.Execute()
	if err != nil {
		if 200 <= rawResp.StatusCode && rawResp.StatusCode <= 299 {
			if strings.HasPrefix(err.Error(), "parsing time") {
				logger(meta).Info("error when parsing time, will process raw HTTP response")
			}
			if strings.Contains(err.Error(), "cannot unmarshal number") {
				logger(meta).Info("error when parsing number, will process raw HTTP response")
			}
		} else {
			return diag.Errorf("failed to activate behavior: %v", err)
		}
	}
	return nil
}
