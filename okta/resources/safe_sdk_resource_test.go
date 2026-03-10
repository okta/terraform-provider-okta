package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestWrapSDKResource_CreateContext_PanicRecovery tests that CreateContext panics are recovered
func TestWrapSDKResource_CreateContext_PanicRecovery(t *testing.T) {
	// Create a resource with a CreateContext that panics
	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			// Trigger a nil pointer dereference panic
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	// Wrap the resource
	wrapped := WrapSDKResource(r)

	// Call CreateContext - should not panic, should return diagnostics
	diags := wrapped.CreateContext(context.Background(), nil, nil)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics from panic recovery, got none")
	}

	if diags[0].Severity != diag.Error {
		t.Errorf("expected Error severity, got %v", diags[0].Severity)
	}

	if !strings.Contains(diags[0].Summary, "Provider Crash in Create") {
		t.Errorf("expected summary to contain 'Provider Crash in Create', got %s", diags[0].Summary)
	}

	if !strings.Contains(diags[0].Detail, "nil pointer dereference") {
		t.Errorf("expected detail to contain 'nil pointer dereference', got %s", diags[0].Detail)
	}

	if !strings.Contains(diags[0].Detail, "Stack trace:") {
		t.Errorf("expected detail to contain 'Stack trace:', got %s", diags[0].Detail)
	}
}

// TestWrapSDKResource_ReadContext_PanicRecovery tests that ReadContext panics are recovered
func TestWrapSDKResource_ReadContext_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	diags := wrapped.ReadContext(context.Background(), nil, nil)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics from panic recovery, got none")
	}

	if !strings.Contains(diags[0].Summary, "Provider Crash in Read") {
		t.Errorf("expected summary to contain 'Provider Crash in Read', got %s", diags[0].Summary)
	}
}

// TestWrapSDKResource_UpdateContext_PanicRecovery tests that UpdateContext panics are recovered
func TestWrapSDKResource_UpdateContext_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	diags := wrapped.UpdateContext(context.Background(), nil, nil)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics from panic recovery, got none")
	}

	if !strings.Contains(diags[0].Summary, "Provider Crash in Update") {
		t.Errorf("expected summary to contain 'Provider Crash in Update', got %s", diags[0].Summary)
	}
}

// TestWrapSDKResource_DeleteContext_PanicRecovery tests that DeleteContext panics are recovered
func TestWrapSDKResource_DeleteContext_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	diags := wrapped.DeleteContext(context.Background(), nil, nil)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics from panic recovery, got none")
	}

	if !strings.Contains(diags[0].Summary, "Provider Crash in Delete") {
		t.Errorf("expected summary to contain 'Provider Crash in Delete', got %s", diags[0].Summary)
	}
}

// TestWrapSDKResource_LegacyCreate_PanicRecovery tests that legacy Create panics are recovered
func TestWrapSDKResource_LegacyCreate_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	err := wrapped.Create(nil, nil)

	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}

	if !strings.Contains(err.Error(), "provider crashed during the Create operation") {
		t.Errorf("expected error to contain 'provider crashed during the Create operation', got %s", err.Error())
	}

	if !strings.Contains(err.Error(), "nil pointer dereference") {
		t.Errorf("expected error to contain 'nil pointer dereference', got %s", err.Error())
	}
}

// TestWrapSDKResource_LegacyRead_PanicRecovery tests that legacy Read panics are recovered
func TestWrapSDKResource_LegacyRead_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	err := wrapped.Read(nil, nil)

	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}

	if !strings.Contains(err.Error(), "provider crashed during the Read operation") {
		t.Errorf("expected error to contain 'provider crashed during the Read operation', got %s", err.Error())
	}
}

// TestWrapSDKResource_LegacyUpdate_PanicRecovery tests that legacy Update panics are recovered
func TestWrapSDKResource_LegacyUpdate_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		Update: func(d *schema.ResourceData, meta interface{}) error {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	err := wrapped.Update(nil, nil)

	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}

	if !strings.Contains(err.Error(), "provider crashed during the Update operation") {
		t.Errorf("expected error to contain 'provider crashed during the Update operation', got %s", err.Error())
	}
}

// TestWrapSDKResource_LegacyDelete_PanicRecovery tests that legacy Delete panics are recovered
func TestWrapSDKResource_LegacyDelete_PanicRecovery(t *testing.T) {
	r := &schema.Resource{
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKResource(r)
	err := wrapped.Delete(nil, nil)

	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}

	if !strings.Contains(err.Error(), "provider crashed during the Delete operation") {
		t.Errorf("expected error to contain 'provider crashed during the Delete operation', got %s", err.Error())
	}
}

// TestWrapSDKResource_NoPanic_PassesThrough tests that normal operations work correctly
func TestWrapSDKResource_NoPanic_PassesThrough(t *testing.T) {
	createCalled := false
	readCalled := false
	updateCalled := false
	deleteCalled := false

	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			createCalled = true
			return nil
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			readCalled = true
			return nil
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			updateCalled = true
			return nil
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			deleteCalled = true
			return nil
		},
	}

	wrapped := WrapSDKResource(r)

	// Call each operation
	if diags := wrapped.CreateContext(context.Background(), nil, nil); len(diags) > 0 {
		t.Errorf("unexpected diagnostics from CreateContext: %v", diags)
	}
	if diags := wrapped.ReadContext(context.Background(), nil, nil); len(diags) > 0 {
		t.Errorf("unexpected diagnostics from ReadContext: %v", diags)
	}
	if diags := wrapped.UpdateContext(context.Background(), nil, nil); len(diags) > 0 {
		t.Errorf("unexpected diagnostics from UpdateContext: %v", diags)
	}
	if diags := wrapped.DeleteContext(context.Background(), nil, nil); len(diags) > 0 {
		t.Errorf("unexpected diagnostics from DeleteContext: %v", diags)
	}

	// Verify all operations were called
	if !createCalled {
		t.Error("CreateContext was not called")
	}
	if !readCalled {
		t.Error("ReadContext was not called")
	}
	if !updateCalled {
		t.Error("UpdateContext was not called")
	}
	if !deleteCalled {
		t.Error("DeleteContext was not called")
	}
}

// TestWrapSDKResource_NilResource tests that nil resources are handled
func TestWrapSDKResource_NilResource(t *testing.T) {
	wrapped := WrapSDKResource(nil)
	if wrapped != nil {
		t.Error("expected nil result for nil input")
	}
}

// TestWrapSDKResources_Map tests that a map of resources is wrapped correctly
func TestWrapSDKResources_Map(t *testing.T) {
	panicCalled := 0

	resources := map[string]*schema.Resource{
		"okta_resource_a": {
			CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
				panicCalled++
				var nilPtr *string
				_ = *nilPtr
				return nil
			},
		},
		"okta_resource_b": {
			CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
				panicCalled++
				var nilPtr *string
				_ = *nilPtr
				return nil
			},
		},
	}

	wrapped := WrapSDKResources(resources)

	// Verify both resources are wrapped
	if len(wrapped) != 2 {
		t.Errorf("expected 2 wrapped resources, got %d", len(wrapped))
	}

	// Call both and verify panic recovery
	for name, r := range wrapped {
		diags := r.CreateContext(context.Background(), nil, nil)
		if len(diags) == 0 {
			t.Errorf("expected diagnostics from panic recovery for %s", name)
		}
	}

	if panicCalled != 2 {
		t.Errorf("expected panicCalled to be 2, got %d", panicCalled)
	}
}

// TestWrapSDKResource_PreservesDiagnostics tests that diagnostics from the original function are preserved
func TestWrapSDKResource_PreservesDiagnostics(t *testing.T) {
	expectedDiag := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Test warning",
		Detail:   "This is a test warning",
	}

	r := &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.Diagnostics{expectedDiag}
		},
	}

	wrapped := WrapSDKResource(r)
	diags := wrapped.CreateContext(context.Background(), nil, nil)

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}

	if diags[0].Summary != expectedDiag.Summary {
		t.Errorf("expected summary %q, got %q", expectedDiag.Summary, diags[0].Summary)
	}
}
