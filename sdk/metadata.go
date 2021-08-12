package sdk

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/crewjam/saml"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func (m *APISupplement) GetSAMLMetadata(ctx context.Context, id, keyID string) ([]byte, *saml.EntityDescriptor, error) {
	var query string
	if keyID != "" {
		query = fmt.Sprintf("?kid=%s", keyID)
	}
	return m.getXML(ctx, fmt.Sprintf("/api/v1/apps/%s/sso/saml/metadata%s", id, query))
}

func (m *APISupplement) GetSAMLIdpMetadata(ctx context.Context, id string) ([]byte, *saml.EntityDescriptor, error) {
	return m.getXML(ctx, fmt.Sprintf("/api/v1/idps/%s/metadata.xml", id))
}

func (m *APISupplement) getXML(ctx context.Context, url string) ([]byte, *saml.EntityDescriptor, error) {
	re := &okta.RequestExecutor{}
	*re = *m.RequestExecutor
	re = re.WithAccept("application/xml")
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	metadataRoot := saml.EntityDescriptor{}
	resp, err := re.Do(ctx, req, &metadataRoot)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	// this means, that RequestExecutor didn't decode the data, so doing it manually
	if metadataRoot.EntityID == "" {
		copyRawBytes := make([]byte, len(raw))
		copy(copyRawBytes, raw)
		err = xml.Unmarshal(copyRawBytes, &metadataRoot)
		if err != nil {
			return nil, nil, err
		}
	}
	return raw, &metadataRoot, nil
}
