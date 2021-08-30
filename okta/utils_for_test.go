package okta

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type checkUpstream func(string) (bool, error)

func ensureResourceExists(name string, checkUpstream checkUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}
		ID := rs.Primary.ID
		exist, err := checkUpstream(ID)
		if err != nil {
			return err
		} else if !exist {
			return missingErr
		}
		return nil
	}
}

func createCheckResourceDestroy(typeName string, checkUpstream checkUpstream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != typeName {
				continue
			}
			ID := rs.Primary.ID
			exists, err := checkUpstream(ID)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("resource still exists, ID: %s", ID)
			}
		}
		return nil
	}
}

func ensureResourceNotExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return nil
		}
		return fmt.Errorf("Resource found: %s", name)
	}
}

func condenseError(errorList []error) error {
	if len(errorList) < 1 {
		return nil
	}
	msgList := make([]string, len(errorList))
	for i, err := range errorList {
		if err != nil {
			msgList[i] = err.Error()
		}
	}
	return fmt.Errorf("series of errors occurred: %s", strings.Join(msgList, ", "))
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestHttpClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func newTestOktaClientWithResponse(response roundTripFunc) (context.Context, *okta.Client, error) {
	ctx := context.Background()

	h := newTestHttpClient(response)

	oktaCtx, c, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl("https://foo.okta.com"),
		okta.WithToken("f0oT0k3n"),
		okta.WithHttpClientPtr(h),
	)

	if err != nil {
		return nil, nil, err
	}

	return oktaCtx, c, nil
}
