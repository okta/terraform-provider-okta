package resources

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WrapSDKResource wraps a terraform-plugin-sdk/v2 resource with panic recovery.
// This handles SDK-based resources like app_bookmark that use schema.Resource.
func WrapSDKResource(r *schema.Resource) *schema.Resource {
	if r == nil {
		return nil
	}

	// Wrap context-aware functions (preferred)
	if original := r.CreateContext; original != nil {
		r.CreateContext = wrapSDKCreateContextFunc(original)
	}
	if original := r.ReadContext; original != nil {
		r.ReadContext = wrapSDKReadContextFunc(original)
	}
	if original := r.UpdateContext; original != nil {
		r.UpdateContext = wrapSDKUpdateContextFunc(original)
	}
	if original := r.DeleteContext; original != nil {
		r.DeleteContext = wrapSDKDeleteContextFunc(original)
	}

	// Wrap legacy functions (deprecated but still used in some resources)
	if original := r.Create; original != nil {
		r.Create = wrapSDKCreateFunc(original)
	}
	if original := r.Read; original != nil {
		r.Read = wrapSDKReadFunc(original)
	}
	if original := r.Update; original != nil {
		r.Update = wrapSDKUpdateFunc(original)
	}
	if original := r.Delete; original != nil {
		r.Delete = wrapSDKDeleteFunc(original)
	}

	return r
}

// WrapSDKResources wraps all SDK resources in a map with panic recovery.
func WrapSDKResources(resources map[string]*schema.Resource) map[string]*schema.Resource {
	wrapped := make(map[string]*schema.Resource, len(resources))
	for name, r := range resources {
		wrapped[name] = WrapSDKResource(r)
	}
	return wrapped
}

// panicRecoveryDiagnostic creates a diagnostic error from a recovered panic.
func panicRecoveryDiagnostic(operation string, r interface{}, stackTrace string) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Provider Crash in %s", operation),
			Detail: fmt.Sprintf(
				"The provider crashed during the %s operation.\n\n"+
					"This is always a bug in the provider and should be reported to the provider developers.\n\n"+
					"Error: %v\n\nStack trace:\n%s",
				operation, r, stackTrace),
		},
	}
}

// panicRecoveryError creates an error from a recovered panic for legacy functions.
func panicRecoveryError(operation string, r interface{}, stackTrace string) error {
	return fmt.Errorf(
		"provider crashed during the %s operation.\n\n"+
			"This is always a bug in the provider and should be reported to the provider developers.\n\n"+
			"Error: %v\n\nStack trace:\n%s",
		operation, r, stackTrace)
}

// Context-aware function wrappers (each is a distinct type in the SDK)

func wrapSDKCreateContextFunc(fn schema.CreateContextFunc) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = panicRecoveryDiagnostic("Create", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Create operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

func wrapSDKReadContextFunc(fn schema.ReadContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = panicRecoveryDiagnostic("Read", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Read operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

func wrapSDKUpdateContextFunc(fn schema.UpdateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = panicRecoveryDiagnostic("Update", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Update operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

func wrapSDKDeleteContextFunc(fn schema.DeleteContextFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagResult diag.Diagnostics) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				diagResult = panicRecoveryDiagnostic("Delete", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Delete operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(ctx, d, meta)
	}
}

// Legacy function wrappers (each is a distinct type in the SDK)

func wrapSDKCreateFunc(fn schema.CreateFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Create", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Create operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}

func wrapSDKReadFunc(fn schema.ReadFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Read", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Read operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}

func wrapSDKUpdateFunc(fn schema.UpdateFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Update", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Update operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}

func wrapSDKDeleteFunc(fn schema.DeleteFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				err = panicRecoveryError("Delete", r, stackTrace)
				log.Printf("[CRITICAL] Provider panic in Delete operation: %v\nStack trace:\n%s", r, stackTrace)
			}
		}()
		return fn(d, meta)
	}
}
