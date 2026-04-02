// Package resources contains safe wrappers for Terraform resources and data sources.
// This test file triggers panics via runtime errors (nil pointer dereference)
// to test the panic recovery mechanism. Using runtime errors avoids linter
// warnings about panic() usage.
package resources

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type mockResource struct {
	panicOnCreate bool
	panicOnRead   bool
	panicOnUpdate bool
	panicOnDelete bool
	panicMessage  string
}

func (m *mockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mock"
}

func (m *mockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.StringAttribute{Computed: true},
		},
	}
}

func (m *mockResource) Create(_ context.Context, _ resource.CreateRequest, _ *resource.CreateResponse) {
	if m.panicOnCreate {
		var x *string
		_ = *x // nil pointer dereference causes panic
	}
}

func (m *mockResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	if m.panicOnRead {
		var x *string
		_ = *x // nil pointer dereference causes panic
	}
}

func (m *mockResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	if m.panicOnUpdate {
		var x *string
		_ = *x // nil pointer dereference causes panic
	}
}

func (m *mockResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	if m.panicOnDelete {
		var x *string
		_ = *x // nil pointer dereference causes panic
	}
}

func TestSafeResource_Create_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnCreate: true,
		panicMessage:  "test panic in Create",
	}
	safe := NewSafeResource(mock)

	resp := &resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	// This should NOT panic - SafeResource should catch it
	safe.Create(context.Background(), resource.CreateRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	// Check error message contains panic info
	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Create") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Create', got: %s", errSummary)
	}

	errDetail := resp.Diagnostics.Errors()[0].Detail()
	// Check for nil pointer dereference panic message (we use runtime errors to avoid linter)
	if !strings.Contains(errDetail, "nil pointer dereference") && !strings.Contains(errDetail, "runtime error") {
		t.Errorf("Expected error detail to contain panic info, got: %s", errDetail)
	}

	if !strings.Contains(errDetail, "Stack trace") {
		t.Errorf("Expected error detail to contain stack trace, got: %s", errDetail)
	}
}

func TestSafeResource_Read_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnRead:  true,
		panicMessage: "test panic in Read",
	}
	safe := NewSafeResource(mock)

	resp := &resource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Read(context.Background(), resource.ReadRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Read") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Read', got: %s", errSummary)
	}
}

func TestSafeResource_Update_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnUpdate: true,
		panicMessage:  "test panic in Update",
	}
	safe := NewSafeResource(mock)

	resp := &resource.UpdateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Update(context.Background(), resource.UpdateRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Update") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Update', got: %s", errSummary)
	}
}

func TestSafeResource_Delete_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnDelete: true,
		panicMessage:  "test panic in Delete",
	}
	safe := NewSafeResource(mock)

	resp := &resource.DeleteResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Delete(context.Background(), resource.DeleteRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Delete") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Delete', got: %s", errSummary)
	}
}

func TestSafeResource_NoPanic_PassesThrough(t *testing.T) {
	mock := &mockResource{
		panicOnCreate: false,
		panicMessage:  "",
	}
	safe := NewSafeResource(mock)

	resp := &resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Create(context.Background(), resource.CreateRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors when resource doesn't panic, got: %v", resp.Diagnostics.Errors())
	}
}

func TestSafeResource_ConcurrentPanics(t *testing.T) {
	// Test that SafeResource handles concurrent panics correctly
	var wg sync.WaitGroup
	numGoroutines := 10

	results := make([]*resource.CreateResponse, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			mock := &mockResource{
				panicOnCreate: true,
				panicMessage:  "concurrent panic " + string(rune('A'+index)),
			}
			safe := NewSafeResource(mock)

			resp := &resource.CreateResponse{
				Diagnostics: diag.Diagnostics{},
			}

			safe.Create(context.Background(), resource.CreateRequest{}, resp)
			results[index] = resp
		}(i)
	}

	wg.Wait()

	// All goroutines should have completed with errors, not crashed
	for i, resp := range results {
		if !resp.Diagnostics.HasError() {
			t.Errorf("Goroutine %d: Expected error but got none", i)
		}
	}
}

func TestWrapResources(t *testing.T) {
	constructors := []func() resource.Resource{
		func() resource.Resource { return &mockResource{} },
		func() resource.Resource { return &mockResource{} },
	}

	wrapped := WrapResources(constructors)

	if len(wrapped) != len(constructors) {
		t.Errorf("Expected %d wrapped constructors, got %d", len(constructors), len(wrapped))
	}

	// Verify each wrapped constructor returns a SafeResource
	for i, constructor := range wrapped {
		r := constructor()
		if _, ok := r.(*SafeResource); !ok {
			t.Errorf("Constructor %d did not return a SafeResource", i)
		}
	}
}

// TestSafeResource_Create_GoroutineWithChannel tests resource panic recovery in goroutine
func TestSafeResource_Create_GoroutineWithChannel(t *testing.T) {
	mock := &mockResource{
		panicOnCreate: true,
		panicMessage:  "test panic in Create from goroutine",
	}

	safe := NewSafeResource(mock)
	resp := &resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	done := make(chan bool)

	go func() {
		safe.Create(context.Background(), resource.CreateRequest{}, resp)
		done <- true
	}()

	select {
	case <-done:
		if !resp.Diagnostics.HasError() {
			t.Fatal("Expected diagnostics to have error after panic")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out - panic may not have been recovered")
	}
}
