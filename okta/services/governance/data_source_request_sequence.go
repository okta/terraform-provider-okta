package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// requestSequenceDataSource is an alias for approvalSequenceDataSource
// registered under the original "okta_request_sequence" name.

var _ datasource.DataSource = &requestSequenceDataSource{}

func newRequestSequencesDataSource() datasource.DataSource {
	return &requestSequenceDataSource{}
}

type requestSequenceDataSource struct {
	approvalSequenceDataSource
}

func (d *requestSequenceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_sequence"
}

func (d *requestSequenceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.approvalSequenceDataSource.Schema(ctx, req, resp)
}
