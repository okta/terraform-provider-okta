package resources

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// typeBaseName strips pointer prefixes and package qualifiers from a reflect type string.
// e.g. "*idaas.deviceDataSource" → "deviceDataSource"
func typeBaseName(t reflect.Type) string {
	s := t.String()
	// dereference pointer types
	for strings.HasPrefix(s, "*") {
		s = s[1:]
	}
	// strip package qualifier
	if idx := strings.LastIndex(s, "."); idx >= 0 {
		s = s[idx+1:]
	}
	return s
}

// Ensure SafeResource implements all required interfaces
var (
	_ resource.Resource                   = &SafeResource{}
	_ resource.ResourceWithConfigure      = &SafeResource{}
	_ resource.ResourceWithImportState    = &SafeResource{}
	_ resource.ResourceWithValidateConfig = &SafeResource{}
	_ resource.ResourceWithModifyPlan     = &SafeResource{}
	_ resource.ResourceWithUpgradeState   = &SafeResource{}
)

// SafeResource wraps a resource with panic recovery to prevent provider crashes
type SafeResource struct {
	underlying   resource.Resource
	nameOnce     sync.Once
	resourceName atomic.Value // string
}

// NewSafeResource creates a new SafeResource wrapper around the given resource
func NewSafeResource(r resource.Resource) resource.Resource {
	return &SafeResource{underlying: r}
}

// WrapResources wraps multiple resource constructors with SafeResource
func WrapResources(constructors []func() resource.Resource) []func() resource.Resource {
	wrapped := make([]func() resource.Resource, len(constructors))
	for i, constructor := range constructors {
		c := constructor // capture loop variable
		wrapped[i] = func() resource.Resource {
			return NewSafeResource(c())
		}
	}
	return wrapped
}

// recoverPanic handles panic recovery and adds appropriate diagnostics
func (s *SafeResource) recoverPanic(diags *diag.Diagnostics, operation string) {
	if r := recover(); r != nil {
		stackTrace := string(debug.Stack())
		resName, _ := s.resourceName.Load().(string)
		if resName == "" && s.underlying != nil {
			resName = typeBaseName(reflect.TypeOf(s.underlying))
		}
		if resName == "" {
			resName = "unknown"
		}

		diags.AddError(
			fmt.Sprintf("Provider Crash in %s operation of resource %s", operation, resName),
			fmt.Sprintf(
				"The Terraform Provider Okta crashed during the %s operation of resource %s.\n\n"+
					"Please check if this issue has already been reported on\n"+
					"https://github.com/okta/terraform-provider-okta/issues\n"+
					"or create a new issue with this stack trace.\n"+
					"Error: %v\n\nStack trace:\n%s\n\n",
				operation, resName, r, stackTrace,
			),
		)
	}
}

// Metadata delegates to the underlying resource
func (s *SafeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	s.underlying.Metadata(ctx, req, resp)
	if resp.TypeName != "" {
		s.nameOnce.Do(func() {
			s.resourceName.Store(resp.TypeName)
		})
	}
}

// Schema delegates to the underlying resource
func (s *SafeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	s.underlying.Schema(ctx, req, resp)
}

// Create wraps the underlying Creation with panic recovery
func (s *SafeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Create")
	s.underlying.Create(ctx, req, resp)
}

// Read wraps the underlying Read with panic recovery
func (s *SafeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Read")
	s.underlying.Read(ctx, req, resp)
}

// Update wraps the underlying Update with panic recovery
func (s *SafeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Update")
	s.underlying.Update(ctx, req, resp)
}

// Delete wraps the underlying Delete with panic recovery
func (s *SafeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Delete")
	s.underlying.Delete(ctx, req, resp)
}

// Configure delegates to the underlying resource if it implements ResourceWithConfigure
func (s *SafeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Configure")
	if rc, ok := s.underlying.(resource.ResourceWithConfigure); ok {
		rc.Configure(ctx, req, resp)
	}
}

// ImportState delegates to the underlying resource if it implements ResourceWithImportState
func (s *SafeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "ImportState")
	if ri, ok := s.underlying.(resource.ResourceWithImportState); ok {
		ri.ImportState(ctx, req, resp)
	}
	// If not implemented, the Framework handles this — do not add an error here.
}

// ValidateConfig delegates to the underlying resource if it implements ResourceWithValidateConfig
func (s *SafeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "ValidateConfig")
	if rv, ok := s.underlying.(resource.ResourceWithValidateConfig); ok {
		rv.ValidateConfig(ctx, req, resp)
	}
}

// ModifyPlan delegates to the underlying resource if it implements ResourceWithModifyPlan
func (s *SafeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "ModifyPlan")
	if rm, ok := s.underlying.(resource.ResourceWithModifyPlan); ok {
		rm.ModifyPlan(ctx, req, resp)
	}
}

// UpgradeState delegates to the underlying resource if it implements ResourceWithUpgradeState.
// NOTE: panic recovery is not possible here because the method returns a value (not a *Response),
// so there is no Diagnostics field to write to. A panic here will propagate to the Framework.
func (s *SafeResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if ru, ok := s.underlying.(resource.ResourceWithUpgradeState); ok {
		return ru.UpgradeState(ctx)
	}
	return nil
}
