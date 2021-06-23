package sdk

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func (m *ApiSupplement) UploadAppLogo(ctx context.Context, appID, filename string) (*okta.Response, error) {
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
	req, err := m.RequestExecutor.WithContentType(writer.FormDataContentType()).NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

func GetAppLogoHash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer func() {
		_ = file.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
