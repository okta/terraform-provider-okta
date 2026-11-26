package idaas

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

// getAppSAMLMetadata fetches an app's IdP SAML metadata like
// APISupplement.GetSAMLMetadata, but reports found=false with no error when
// the API responds 404. That is expected for apps where Okta does not publish
// the metadata, for example catalog apps where Okta is the service provider.
func getAppSAMLMetadata(ctx context.Context, s *sdk.APISupplement, appID, keyID string) (raw []byte, root *saml.EntityDescriptor, found bool, err error) {
	path := fmt.Sprintf("/api/v1/apps/%s/sso/saml/metadata", url.PathEscape(appID))
	if keyID != "" {
		path += "?kid=" + url.QueryEscape(keyID)
	}

	// Copy the executor so WithAccept does not change the shared one.
	re := *s.RequestExecutor
	req, err := re.WithAccept("application/xml").NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, false, err
	}
	metadataRoot := saml.EntityDescriptor{}
	resp, err := re.Do(ctx, req, &metadataRoot)
	if utils.Is404(resp) {
		return nil, nil, false, nil
	}
	if err != nil {
		return nil, nil, false, err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, false, err
	}
	// The request executor does not always decode XML, so decode it here if needed.
	if metadataRoot.EntityID == "" {
		if err := xml.Unmarshal(raw, &metadataRoot); err != nil {
			return nil, nil, false, err
		}
	}
	return raw, &metadataRoot, true, nil
}
