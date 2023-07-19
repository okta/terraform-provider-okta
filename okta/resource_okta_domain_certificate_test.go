package okta

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaDomainCertificate(t *testing.T) {
	t.Skip("This test is bespoke and has to be run by hand. We need to spend some time automating this test.")

	pwd, err := os.Getwd()
	if err != nil {
		t.Skip("can't get working directory from OS")
	}

	// NOTE: Setting up reading cert, pk, and chain files via TF file() so that
	// in the future we can set up CI to get a certificate from something like
	// letsencrypt.
	// TF file() needs an absolute path so set up the file names and pass into a
	// sprintf to interpolate a working config
	domainFile := filepath.Join(pwd, "../test/fixtures/okta_domain_certificate/domain.txt")
	certFile := filepath.Join(pwd, "../test/fixtures/okta_domain_certificate/cert.pem")
	pkFile := filepath.Join(pwd, "../test/fixtures/okta_domain_certificate/privkey.pem")
	chainFile := filepath.Join(pwd, "../test/fixtures/okta_domain_certificate/chain.pem")

	for _, path := range []string{domainFile, certFile, pkFile, chainFile} {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			t.Skip(path, "appears to not exist, skipping")
		}
	}

	config := fmt.Sprintf(`
data "okta_domain" "test" {
  domain_id_or_name = file("%s")
}
  
resource "okta_domain_certificate" "test" {
  domain_id   = data.okta_domain.test.id
  type        = "PEM"
  
  certificate = file("%s")
  private_key = file("%s")
  certificate_chain = file("%s")
}`,
		domainFile, certFile, pkFile, chainFile)
	resourceName := fmt.Sprintf("%s.test", domainCertificate)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkIsCertificate(resourceName, "certificate"),
					checkIsCertificate(resourceName, "certificate_chain"),
					checkIsPrivateKey(resourceName, "private_key"),
				),
			},
		},
	})
}

func checkIsCertificate(resourceName, attribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if isVCRPlayMode() {
			return nil
		}
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		attr, ok := rs.Primary.Attributes[attribute]
		if !ok {
			return fmt.Errorf("resource attribute not found: %s.%s", resourceName, attribute)
		}
		ok, _ = regexp.MatchString("^-----BEGIN CERTIFICATE-----\n", attr)
		if !ok {
			return fmt.Errorf("resource %s.%s does not appear to be certificate", resourceName, attribute)
		}

		return nil
	}
}

func checkIsPrivateKey(resourceName, attribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if isVCRPlayMode() {
			return nil
		}
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		attr, ok := rs.Primary.Attributes[attribute]
		if !ok {
			return fmt.Errorf("resource attribute not found: %s.%s", resourceName, attribute)
		}
		ok, _ = regexp.MatchString("^-----BEGIN PRIVATE KEY-----\n", attr)
		if !ok {
			return fmt.Errorf("resource %s.%s does not appear to be a private key", resourceName, attribute)
		}
		return nil
	}
}
