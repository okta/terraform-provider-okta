package sdk

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadOrgLogo uploads app's logo
func (m *APISupplement) UploadOrgLogo(ctx context.Context, filename string) (*Response, error) {
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
	re := m.cloneRequestExecutor()
	req, err := re.WithContentType(writer.FormDataContentType()).NewRequest(http.MethodPost, "/api/v1/org/logo", body)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
