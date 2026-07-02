package resources

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WrapSDKResource wraps a terraform-plugin-sdk/v2 resource with panic recovery.
// This handles SDK-based resources like app_bookmark that use schema.Resource.
func WrapSDKResource(r *schema.Resource) *schema.Resource {
	return wrapSDKResourceWithName(r, "unknown")
}

// wrapSDKResourceWithName wraps a terraform-plugin-sdk/v2 resource with panic recovery,
// including the resource name in error messages.
func wrapSDKResourceWithName(r *schema.Resource, resourceName string) *schema.Resource {
	if r == nil {
		return nil
	}

	// Wrap context-aware functions (preferred)
	if original := r.CreateContext; original != nil {
		r.CreateContext = wrapSDKCreateContextFunc(original, resourceName)
	}
	if original := r.ReadContext; original != nil {
		r.ReadContext = wrapSDKReadContextFunc(original, resourceName)
	}
	if original := r.UpdateContext; original != nil {
		r.UpdateContext = wrapSDKUpdateContextFunc(original, resourceName)
	}
	if original := r.DeleteContext; original != nil {
		r.DeleteContext = wrapSDKDeleteContextFunc(original, resourceName)
	}

	// Wrap legacy functions (deprecated but still used in some resources)
	if original := r.Create; original != nil {
		r.Create = wrapSDKCreateFunc(original, resourceName)
	}
	if original := r.Read; original != nil {
		r.Read = wrapSDKReadFunc(original, resourceName)
	}
	if original := r.Update; original != nil {
		r.Update = wrapSDKUpdateFunc(original, resourceName)
	}
	if original := r.Delete; original != nil {
		r.Delete = wrapSDKDeleteFunc(original, resourceName)
	}

	return r
}

// WrapSDKResources wraps all SDK resources in a map with panic recovery.
func WrapSDKResources(resources map[string]*schema.Resource) map[string]*schema.Resource {
	wrapped := make(map[string]*schema.Resource, len(resources))
	for name, r := range resources {
		wrapped[name] = wrapSDKResourceWithName(r, name)
	}
	return wrapped
}

// panicRecoveryDiagnostic creates a diagnostic error from a recovered panic.
func resourcePanicRecoveryDiagnostic(operation, resourceName string, r interface{}, stackTrace string) diag.Diagnostics {
	if resourceName == "" {
		resourceName = "unknown"
	}
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Provider Crash in %s operation of resource %s", operation, resourceName),
			Detail: fmt.Sprintf(
				"The Terraform Provider Okta crashed during the %s operation of resource %s.\n\n"+
					"Please check if this issue has already been reported on\n"+
					"https://github.com/okta/terraform-provider-okta/issues\n"+
					"or create a new issue with this stack trace.\n"+
					"Error: %v\n\nStack trace:\n%s\n\n",
				operation, resourceName, r, stackTrace,
			),
		},
	}
}

// panicRecoveryError creates an error from a recovered panic for legacy functions.
func panicRecoveryError(operation, resName string, r interface{}, stackTrace string) error {
	if resName == "" {
		resName = "unknown"
	}
	return fmt.Errorf(
		"The Terraform Provider Okta crashed during the %s operation of resource %s.\n\n"+
			"Please check if this issue has already been reported on\n"+
			"https://github.com/okta/terraform-provider-okta/issues\n"+
			"or create a new issue with this stack trace.\n"+
			"Error: %v\n\nStack trace:\n%s\n\n",
		operation, resName, r, stackTrace,
	)
}

// Context-aware function wrappers (each is a distinct type in the SDK)

func wrapSDKCreateContextFunc(fn schema.CreateContextFunc, resourceName string) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = resourcePanicRecoveryDiagnostic("Create", resourceName, r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

func wrapSDKReadContextFunc(fn schema.ReadContextFunc, resourceName string) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = resourcePanicRecoveryDiagnostic("Read", resourceName, r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

func wrapSDKUpdateContextFunc(fn schema.UpdateContextFunc, resourceName string) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = resourcePanicRecoveryDiagnostic("Update", resourceName, r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

func wrapSDKDeleteContextFunc(fn schema.DeleteContextFunc, resourceName string) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = resourcePanicRecoveryDiagnostic("Delete", resourceName, r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

// Legacy function wrappers (each is a distinct type in the SDK)

func wrapSDKCreateFunc(fn schema.CreateFunc, resourceName string) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Create", resourceName, r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}

func wrapSDKReadFunc(fn schema.ReadFunc, resourceName string) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Read", resourceName, r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}

func wrapSDKUpdateFunc(fn schema.UpdateFunc, resourceName string) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Update", resourceName, r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}

func wrapSDKDeleteFunc(fn schema.DeleteFunc, resourceName string) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Delete", resourceName, r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}
