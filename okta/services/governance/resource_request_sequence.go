package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// requestSequenceResource is an alias for approvalSequenceResource
// registered under the original "okta_request_sequence" name.

var (
	_ resource.Resource                = &requestSequenceResource{}
	_ resource.ResourceWithConfigure   = &requestSequenceResource{}
	_ resource.ResourceWithImportState = &requestSequenceResource{}
)

func newRequestSequenceResource() resource.Resource {
	return &requestSequenceResource{}
}

type requestSequenceResource struct {
	approvalSequenceResource
}

func (r *requestSequenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_sequence"
}

func (r *requestSequenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = approvalSequenceResourceSchema()
}
