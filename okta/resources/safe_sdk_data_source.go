package resources

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WrapSDKDataSource wraps a terraform-plugin-sdk/v2 data source with panic recovery.
func WrapSDKDataSource(d *schema.Resource) *schema.Resource {
	return wrapSDKDataSourceWithName(d, "unknown")
}

// wrapSDKDataSourceWithName wraps a terraform-plugin-sdk/v2 data source with panic recovery,
// including the data source name in error messages.
func wrapSDKDataSourceWithName(d *schema.Resource, dataSourceName string) *schema.Resource {
	if d == nil {
		return nil
	}

	// Data sources only have Read operations
	if original := d.ReadContext; original != nil {
		d.ReadContext = wrapSDKDataSourceReadContextFunc(original, dataSourceName)
	}
	if original := d.Read; original != nil {
		d.Read = wrapSDKDataSourceReadFunc(original, dataSourceName)
	}

	return d
}

// WrapSDKDataSources wraps all SDK data sources in a map with panic recovery.
func WrapSDKDataSources(dataSources map[string]*schema.Resource) map[string]*schema.Resource {
	wrapped := make(map[string]*schema.Resource, len(dataSources))
	for name, d := range dataSources {
		wrapped[name] = wrapSDKDataSourceWithName(d, name)
	}
	return wrapped
}

func wrapSDKDataSourceReadContextFunc(fn schema.ReadContextFunc, dataSourceName string) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = dataSourcePanicRecoveryDiagnostic("Read", dataSourceName, r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

// dataSourcePanicRecoveryDiagnostic creates a diagnostic error from a recovered data source panic.
func dataSourcePanicRecoveryDiagnostic(operation, dataSourceName string, r interface{}, stackTrace string) diag.Diagnostics {
	if dataSourceName == "" {
		dataSourceName = "unknown"
	}
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Provider Crash in %s operation of data source %s", operation, dataSourceName),
			Detail: fmt.Sprintf(
				"The Terraform Provider Okta crashed during the %s operation of data source %s.\n\n"+
					"Please check if this issue has already been reported on\n"+
					"https://github.com/okta/terraform-provider-okta/issues\n"+
					"or create a new issue with this stack trace.\n"+
					"Error: %v\n\nStack trace:\n%s\n\n",
				operation, dataSourceName, r, stackTrace,
			),
		},
	}
}

func dataSourcePanicRecoveryError(operation, dataSourceName string, r interface{}, stackTrace string) error {
	if dataSourceName == "" {
		dataSourceName = "unknown"
	}
	return fmt.Errorf(
		"The Terraform Provider Okta crashed during the %s operation of data source %s.\n\n"+
			"Please check if this issue has already been reported on\n"+
			"https://github.com/okta/terraform-provider-okta/issues\n"+
			"or create a new issue with this stack trace.\n"+
			"Error: %v\n\nStack trace:\n%s\n\n",
		operation, dataSourceName, r, stackTrace,
	)
}

func wrapSDKDataSourceReadFunc(fn schema.ReadFunc, dataSourceName string) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = dataSourcePanicRecoveryError("Read", dataSourceName, r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}
