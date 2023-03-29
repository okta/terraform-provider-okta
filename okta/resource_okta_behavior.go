package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the behavior",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Behavior type",
				ForceNew:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     statusActive,
				Description: "Behavior status: ACTIVE or INACTIVE.",
			},
			"location_granularity_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines the method and level of detail used to evaluate the behavior.",
			},
			"radius_from_location": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Radius from location (in kilometers)",
			},
			"number_of_authentications": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of recent authentications used to evaluate the behavior.",
			},
			"velocity": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Velocity (in kilometers per hour).",
				ConflictsWith: []string{"number_of_authentications", "radius_from_location", "location_granularity_type"},
			},
		},
	}
}

func resourceBehaviorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating location behavior", "name", d.Get("name").(string))
	err := validateBehavior(d)
	if err != nil {
		return diag.FromErr(err)
	}
	behavior, _, err := getSupplementFromMetadata(m).CreateBehavior(ctx, buildBehavior(d))
	if err != nil {
		return diag.Errorf("failed to create location behavior: %v", err)
	}
	d.SetId(behavior.ID)
	return resourceBehaviorRead(ctx, d, m)
}

func resourceBehaviorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("getting behavior", "id", d.Id())
	behavior, resp, err := getSupplementFromMetadata(m).GetBehavior(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
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

func resourceBehaviorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating location behavior", "name", d.Get("name").(string))
	err := validateBehavior(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getSupplementFromMetadata(m).UpdateBehavior(ctx, d.Id(), buildBehavior(d))
	if err != nil {
		return diag.Errorf("failed to update location behavior: %v", err)
	}
	if d.HasChange("status") {
		err := handleBehaviorLifecycle(ctx, d, m)
		if err != nil {
			return err
		}
	}
	return resourceBehaviorRead(ctx, d, m)
}

func resourceBehaviorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("deleting location behavior", "name", d.Get("name").(string))
	_, err := getSupplementFromMetadata(m).DeleteBehavior(ctx, d.Id())
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

func handleBehaviorLifecycle(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getSupplementFromMetadata(m)
	if d.Get("status").(string) == statusActive {
		logger(m).Info("activating behavior", "name", d.Get("name").(string))
		_, err := client.ActivateBehavior(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to activate behavior: %v", err)
		}
		return nil
	}
	logger(m).Info("deactivating behavior", "name", d.Get("name").(string))
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
