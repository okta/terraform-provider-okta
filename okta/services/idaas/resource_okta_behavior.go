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
	"github.com/okta/terraform-provider-okta/okta/utils"
)

const (
	behaviorAnomalousLocation = "ANOMALOUS_LOCATION"
	behaviorAnomalousDevice   = "ANOMALOUS_DEVICE"
	behaviorAnomalousIP       = "ANOMALOUS_IP"
	behaviorVelocity          = "VELOCITY"
)

func resourceBehavior() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBehaviorCreateUsingSDK,
		ReadContext:   resourceBehaviorReadUsingSDK,
		UpdateContext: resourceBehaviorUpdateUsingSDK,
		DeleteContext: resourceBehaviorDeleteUsingSDK,
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

/*
	func resourceBehaviorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		logger(meta).Info("creating location behavior", "name", d.Get("name").(string))
		err := validateBehavior(d)
		if err != nil {
			return diag.FromErr(err)
		}
		behavior, _, err := getAPISupplementFromMetadata(meta).CreateBehavior(ctx, buildBehavior(d))
		if err != nil {
			return diag.Errorf("failed to create location behavior: %v", err)
		}
		d.SetId(behavior.ID)
		return resourceBehaviorRead(ctx, d, meta)
	}

	func resourceBehaviorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		logger(meta).Info("getting behavior", "id", d.Id())
		behavior, resp, err := getAPISupplementFromMetadata(meta).GetBehavior(ctx, d.Id())
		if err := utils.SuppressErrorOn404(resp, err); err != nil {
			return diag.Errorf("failed to find behavior: %v", err)
		}
		if behavior == nil {
			d.SetId("")
			return nil
		}
		_ = d.Set("name", behavior.Name)
		_ = d.Set("type", behavior.Type)
		_ = d.Set("status", behavior.Status)
		setSettings(d, behavior.Type, behavior.Settings)
		return nil
	}

	func resourceBehaviorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		logger(meta).Info("updating location behavior", "name", d.Get("name").(string))
		err := validateBehavior(d)
		if err != nil {
			return diag.FromErr(err)
		}
		_, _, err = getAPISupplementFromMetadata(meta).UpdateBehavior(ctx, d.Id(), buildBehavior(d))
		if err != nil {
			return diag.Errorf("failed to update location behavior: %v", err)
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
		logger(meta).Info("deleting location behavior", "name", d.Get("name").(string))
		_, err := getAPISupplementFromMetadata(meta).DeleteBehavior(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to delete location behavior: %v", err)
		}
		return nil
	}

	func buildBehavior(d *schema.ResourceData) sdk.Behavior {
		b := sdk.Behavior{
			Name:     d.Get("name").(string),
			Status:   d.Get("status").(string),
			Settings: make(map[string]interface{}),
			Type:     d.Get("type").(string),
		}
		if b.Type == behaviorAnomalousLocation || b.Type == behaviorAnomalousDevice || b.Type == behaviorAnomalousIP {
			b.Settings["maxEventsUsedForEvaluation"] = d.Get("number_of_authentications")
		}
		if b.Type == behaviorAnomalousLocation {
			b.Settings["granularity"] = d.Get("location_granularity_type")
			if d.Get("location_granularity_type").(string) == "LAT_LONG" {
				b.Settings["radiusKilometers"] = d.Get("radius_from_location")
			}
		}
		if b.Type == behaviorVelocity {
			b.Settings["velocityKph"] = d.Get("velocity")
		}
		return b
	}

	func handleBehaviorLifecycle(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := getAPISupplementFromMetadata(meta)
		if d.Get("status").(string) == StatusActive {
			logger(meta).Info("activating behavior", "name", d.Get("name").(string))
			_, err := client.ActivateBehavior(ctx, d.Id())
			if err != nil {
				return diag.Errorf("failed to activate behavior: %v", err)
			}
			return nil
		}
		logger(meta).Info("deactivating behavior", "name", d.Get("name").(string))
		_, err := client.DeactivateBehavior(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to deactivate behavior: %v", err)
		}
		return nil
	}

	func setSettings(d *schema.ResourceData, typ string, settings map[string]interface{}) {
		if typ == behaviorAnomalousLocation || typ == behaviorAnomalousDevice || typ == behaviorAnomalousIP {
			_ = d.Set("number_of_authentications", settings["maxEventsUsedForEvaluation"])
		}
		if typ == behaviorAnomalousLocation {
			_ = d.Set("location_granularity_type", settings["granularity"])
			if settings["granularity"].(string) == "LAT_LONG" {
				_ = d.Set("radius_from_location", settings["radiusKilometers"])
			}
		}
		if typ == behaviorVelocity {
			_ = d.Set("velocity", settings["velocityKph"])
		}
	}
*/
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

func resourceBehaviorCreateUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("creating behavior", "name", d.Get("name").(string))
	err := validateBehavior(d)
	if err != nil {
		return diag.FromErr(err)
	}
	behavior, rawResp, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.CreateBehaviorDetectionRuleExecute(buildBehaviorUsingSDK(ctx, d))
	if err != nil {
		if strings.HasPrefix(err.Error(), "parsing time") && (200 <= rawResp.StatusCode && rawResp.StatusCode <= 299) { // error is during parsing of time && http status code is 2xx
			logger(meta).Info("expected error due to unsupported timestamp format in API response, handling it by processing raw response.")
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
	fmt.Printf("DHIWAKAR rawRespMap = %+v\n", respMap)
	d.Set("name", respMap["name"])
	d.Set("type", respMap["type"])
	d.Set("status", respMap["status"])
	d.Set("settings", respMap["settings"])
	d.SetId(*behavior.Id)
	return resourceBehaviorReadUsingSDK(ctx, d, meta)
}

func resourceBehaviorReadUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("getting behavior", "id", d.Id())
	behavior, resp, err := getOktaV5ClientFromMetadata(meta).BehaviorAPI.GetBehaviorDetectionRule(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return diag.Errorf("failed to find behavior: %v", err)
	}
	// created primarily to reduce duplicate lines of code
	type rule interface {
		GetName() string
		GetType() string
		GetStatus() string
	}
	var behaviorRule rule
	switch {
	// Anomalous Device
	case behavior.BehaviorRuleAnomalousDevice != nil:
		behaviorRule = behavior.BehaviorRuleAnomalousDevice
		settings := behavior.BehaviorRuleAnomalousDevice.GetSettings()
		d.Set("number_of_authentications", settings.GetMaxEventsUsedForEvaluation())
	// Anomalous IP
	case behavior.BehaviorRuleAnomalousIP != nil:
		behaviorRule = behavior.BehaviorRuleAnomalousIP
		settings := behavior.BehaviorRuleAnomalousIP.GetSettings()
		d.Set("number_of_authentications", settings.GetMaxEventsUsedForEvaluation())
	// Anomalous Location
	case behavior.BehaviorRuleAnomalousLocation != nil:
		behaviorRule = behavior.BehaviorRuleAnomalousLocation
		settings := behavior.BehaviorRuleAnomalousLocation.GetSettings()
		d.Set("number_of_authentications", settings.GetMaxEventsUsedForEvaluation())
		d.Set("location_granularity_type", settings.GetGranularity())
		if settings.GetGranularity() == "LAT_LONG" {
			d.Set("radius_from_location", settings.GetRadiusKilometers())
		}
	// Anomalous Velocity
	case behavior.BehaviorRuleVelocity != nil:
		behaviorRule = behavior.BehaviorRuleVelocity
		settings := behavior.BehaviorRuleVelocity.GetSettings()
		d.Set("velocity", settings.GetVelocityKph())
	}
	d.Set("name", behaviorRule.GetName())
	d.Set("type", behaviorRule.GetType())
	d.Set("status", behaviorRule.GetStatus())
	return nil
}

func resourceBehaviorUpdateUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("updating location behavior", "name", d.Get("name").(string))
	err := validateBehavior(d)
	if err != nil {
		return diag.FromErr(err)
	}
	rule := buildRuleUsingSDK(d)
	replaceBehaviorDetectionRule := getOktaV5ClientFromMetadata(meta).BehaviorAPI.ReplaceBehaviorDetectionRule(ctx, d.Id())
	replaceBehaviorDetectionRule = replaceBehaviorDetectionRule.Rule(rule)
	_, _, err = replaceBehaviorDetectionRule.Execute()
	if err != nil {
		return diag.Errorf("failed to update location behavior: %v", err)
	}

	if d.HasChange("status") {
		err := handleBehaviorLifecycleUsingSDK(ctx, d, meta)
		if err != nil {
			return err
		}
	}
	return resourceBehaviorReadUsingSDK(ctx, d, meta)
}

func resourceBehaviorDeleteUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("deleting location behavior", "name", d.Get("name").(string))
	deleteBehaviorDetectionRule := getOktaV5ClientFromMetadata(meta).BehaviorAPI.DeleteBehaviorDetectionRule(ctx, d.Id())
	_, err := deleteBehaviorDetectionRule.Execute()
	if err != nil {
		return diag.Errorf("failed to delete location behavior: %v", err)
	}
	return nil
}

func buildRuleUsingSDK(d *schema.ResourceData) v5okta.ListBehaviorDetectionRules200ResponseInner {
	behaviorType := d.Get("type").(string)
	behaviorName := d.Get("name").(string)
	behaviorStatus := d.Get("status").(string)
	listBehaviorDetectionRules := v5okta.ListBehaviorDetectionRules200ResponseInner{}
	switch behaviorType {

	case behaviorAnomalousDevice:
		behaviorRuleAnomalousDevice := v5okta.NewBehaviorRuleAnomalousDevice(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleAnomalousDeviceAsListBehaviorDetectionRules200ResponseInner(behaviorRuleAnomalousDevice)
		settings := v5okta.NewBehaviorRuleSettingsAnomalousDevice()
		settings.SetMaxEventsUsedForEvaluation(int32(d.Get("number_of_authentications").(int))) // this needs to be validated, potential for panic
		listBehaviorDetectionRules.BehaviorRuleAnomalousDevice.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleAnomalousDevice.SetStatus(behaviorStatus)

	case behaviorAnomalousIP:
		behaviorRuleAnomalousIP := v5okta.NewBehaviorRuleAnomalousIP(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleAnomalousIPAsListBehaviorDetectionRules200ResponseInner(behaviorRuleAnomalousIP)
		settings := v5okta.NewBehaviorRuleSettingsAnomalousIP()
		settings.SetMaxEventsUsedForEvaluation(int32(d.Get("number_of_authentications").(int))) // this needs to be validated, potential for panic
		listBehaviorDetectionRules.BehaviorRuleAnomalousIP.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleAnomalousIP.SetStatus(behaviorStatus)

	case behaviorAnomalousLocation:
		behaviorRuleAnomalousLocation := v5okta.NewBehaviorRuleAnomalousLocation(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleAnomalousLocationAsListBehaviorDetectionRules200ResponseInner(behaviorRuleAnomalousLocation)
		granularity := d.Get("location_granularity_type").(string)             // this needs to be validated, potential for panic
		locationGranularityType := d.Get("location_granularity_type").(string) // this needs to be validated, potential for panic
		settings := v5okta.NewBehaviorRuleSettingsAnomalousLocation(granularity)
		if locationGranularityType == "LAT_LONG" {
			settings.SetRadiusKilometers(int32(d.Get("radius_from_location").(int))) // this needs to be validated, potential for panic
		}
		settings.SetMaxEventsUsedForEvaluation(int32(d.Get("number_of_authentications").(int))) // this needs to be validated, potential for panic
		listBehaviorDetectionRules.BehaviorRuleAnomalousLocation.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleAnomalousLocation.SetStatus(behaviorStatus)

	case behaviorVelocity:
		behaviorRuleVelocity := v5okta.NewBehaviorRuleVelocity(behaviorName, behaviorType)
		listBehaviorDetectionRules = v5okta.BehaviorRuleVelocityAsListBehaviorDetectionRules200ResponseInner(behaviorRuleVelocity)
		velocity := d.Get("velocity").(int) // this needs to be validated, potential for panic
		settings := v5okta.NewBehaviorRuleSettingsVelocity(int32(velocity))
		listBehaviorDetectionRules.BehaviorRuleVelocity.SetSettings(*settings)
		listBehaviorDetectionRules.BehaviorRuleVelocity.SetStatus(behaviorStatus)
	}
	return listBehaviorDetectionRules
}

func buildBehaviorUsingSDK(ctx context.Context, d *schema.ResourceData) v5okta.ApiCreateBehaviorDetectionRuleRequest {
	listBehaviorDetectionRules := buildRuleUsingSDK(d)
	behaviorAPIService := v5okta.BehaviorAPIService{}
	behaviorDetectionRuleRequest := behaviorAPIService.CreateBehaviorDetectionRule(ctx)
	behaviorDetectionRuleRequest = behaviorDetectionRuleRequest.Rule(listBehaviorDetectionRules)
	return behaviorDetectionRuleRequest
}

func handleBehaviorLifecycleUsingSDK(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	if d.Get("status").(string) == StatusActive {
		logger(meta).Info("activating behavior", "name", d.Get("name").(string))
		activateBehaviorDetectionRuleRequest := client.BehaviorAPI.ActivateBehaviorDetectionRule(ctx, d.Id())
		_, _, err := activateBehaviorDetectionRuleRequest.Execute()
		if err != nil {
			return diag.Errorf("failed to activate behavior: %v", err)
		}
		return nil
	}
	logger(meta).Info("deactivating behavior", "name", d.Get("name").(string))
	deactivateBehaviorDetectionRuleRequest := client.BehaviorAPI.DeactivateBehaviorDetectionRule(ctx, d.Id())
	_, _, err := deactivateBehaviorDetectionRuleRequest.Execute()
	if err != nil {
		return diag.Errorf("failed to deactivate behavior: %v", err)
	}
	return nil
}
