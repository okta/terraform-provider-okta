package resources

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestSafeDataSource_Read_PanicRecovery(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead:  true,
		panicMessage: "test panic in DataSource Read",
	}

	safe := NewSafeDataSource(mock)
	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}
	safe.Read(context.Background(), datasource.ReadRequest{}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	summary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(summary, "Provider Crash in Read") {
		t.Fatalf("Expected error summary to be '%s', got '%s'", "Provider Crash in Read", resp.Diagnostics.Errors()[0].Summary())
	}

	errDetail := resp.Diagnostics.Errors()[0].Detail()
	if !strings.Contains(errDetail, "test panic in DataSource Read") {
		t.Fatalf("Expected error detail to contain '%s', got '%s'", mock.panicMessage, resp.Diagnostics.Errors()[0].Detail())
	}
}

func TestSafeDataSource_Read_PassesThrough(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead: false,
	}

	safe := NewSafeDataSource(mock)
	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}
	safe.Read(context.Background(), datasource.ReadRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Expected no diagnostics error, got: %v", resp.Diagnostics)
	}
}

func TestSafeDataSource_Read_ConcurrentFails(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead:  true,
		panicMessage: "test panic in DataSource Read",
	}

	safe := NewSafeDataSource(mock)
	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}

	done := make(chan bool)
	go func() {
		safe.Read(context.Background(), datasource.ReadRequest{}, resp)
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
