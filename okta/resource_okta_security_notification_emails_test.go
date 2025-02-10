package okta

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceOktaSecurityNotificationEmails_crud(t *testing.T) {
	mgr := newFixtureManager("resources", securityNotificationEmails, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", securityNotificationEmails)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaSecurityNotificationEmailsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "send_email_for_factor_enrollment_enabled", "true"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "send_email_for_factor_enrollment_enabled", "false"),
				),
			},
		},
	})
}

func checkOktaSecurityNotificationEmailsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != securityNotificationEmails {
			continue
		}
		supplimentClient := sdkSupplementClientForTest()
		oktaClient := sdkV2ClientForTest()
		oktaConfig := oktaClient.GetConfig()
		token := oktaConfig.Okta.Client.Token
		orgUrl, err := url.Parse(oktaConfig.Okta.Client.OrgUrl)
		if err != nil {
			return err
		}
		hostParts := strings.Split(orgUrl.Hostname(), ".")
		orgName := hostParts[0]
		domain := fmt.Sprintf("%s.%s", hostParts[1], hostParts[2])
		emails, err := supplimentClient.GetSecurityNotificationEmails(context.Background(), orgName, domain, token, oktaConfig.HttpClient)
		if err != nil {
			return fmt.Errorf("failed to get security notification emails: %v", err)
		}
		if !emails.SendEmailForNewDeviceEnabled ||
			!emails.SendEmailForFactorEnrollmentEnabled ||
			!emails.SendEmailForFactorResetEnabled ||
			!emails.SendEmailForPasswordChangedEnabled ||
			!emails.ReportSuspiciousActivityEnabled {
			return errors.New("all the flags should be set to true")
		}
	}
	return nil
}
