package sdk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

// UploadAppLogo uploads app's logo
func (m *APISupplement) UploadAppLogo(ctx context.Context, appID, filename string) (*okta.Response, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()
	url := fmt.Sprintf("/api/v1/apps/%s/logo", appID)
	req, err := m.RequestExecutor.WithContentType(writer.FormDataContentType()).NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
