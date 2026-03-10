package resources

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// WrapSDKDataSource wraps a terraform-plugin-sdk/v2 data source with panic recovery.
func WrapSDKDataSource(d *schema.Resource) *schema.Resource {
	if d == nil {
		return nil
	}

	// Data sources only have Read operations
	if original := d.ReadContext; original != nil {
		d.ReadContext = wrapSDKReadContextFunc(original)
	}
	if original := d.Read; original != nil {
		d.Read = wrapSDKReadFunc(original)
	}

	return d
}

// WrapSDKDataSources wraps all SDK data sources in a map with panic recovery.
func WrapSDKDataSources(dataSources map[string]*schema.Resource) map[string]*schema.Resource {
	wrapped := make(map[string]*schema.Resource, len(dataSources))
	for name, d := range dataSources {
		wrapped[name] = WrapSDKDataSource(d)
	}
	return wrapped
}
