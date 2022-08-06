package okta

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSecurityNotificationEmails(t *testing.T) {
	mgr := newFixtureManager(securityNotificationEmails, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updated := mgr.GetFixtures("basic_updated.tf", t)
	resourceName := fmt.Sprintf("%s.test", securityNotificationEmails)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkOktaSecurityNotificationEmailsDestroy(),
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

func checkOktaSecurityNotificationEmailsDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != securityNotificationEmails {
				continue
			}
			c := testAccProvider.Meta().(*Config)
			emails, err := getSupplementFromMetadata(testAccProvider.Meta()).GetSecurityNotificationEmails(context.Background(), c.orgName, c.domain, c.apiToken, c.client)
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
}
