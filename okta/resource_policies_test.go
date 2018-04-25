package okta

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaPolicies_defaultErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn_defaultErrors(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("You cannot edit a default Policy"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccOktaPolicies_nameErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_nameErrors(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the name field or type field of an existing Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaPolicies_typeErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_typeErrors(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the name field or type field of an existing Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaPolicySignOn(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_updated(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "OKTA_SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "OKTA_SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy Updated"),
				),
			},
		},
	})
}

func TestAccOktaPolicySignOn_passErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_passErrors(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("password settings options not supported in the Okta SignOn Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaPolicySignOn_authpErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicySignOn(ri)
	updatedConfig := testOktaPolicySignOn_authpErrors(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("authprovider condition options not supported in the Okta SignOn Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaPolicyPassword(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPolicyPassword(ri)
	updatedConfig := testOktaPolicyPassword_updated(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test Password Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test Password Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minlength", "12"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minlowercase", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minuppercase", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minnumber", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minsymbol", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.excludeusername", "false"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.excludeattributes.0", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.excludeattributes.1", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.dictionarylookup", "true"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.maxagedays", "60"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.expirewarndays", "15"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minageminutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.historycount", "5"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.maxlockoutattempts", "3"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.autounlockminutes", "2"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.showlockoutfailures", "true"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.questionminlength", "10"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.recoveryemailtoken", "20160"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.smsrecovery", "ACTIVE"),
				),
			},
		},
	})
}

func testOktaPolicyExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] Resource Not found: %s", name)
		}

		policyID, hasID := rs.Primary.Attributes["id"]
		if !hasID {
			return fmt.Errorf("[ERROR] No id found in state for Policy")
		}
		policyName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Policy")
		}

		err := testPolicyExists(true, policyID, policyName)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func testOktaPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_policies" {
			continue
		}

		policyID, hasID := rs.Primary.Attributes["id"]
		if !hasID {
			return fmt.Errorf("[ERROR] No id found in state for Policy")
		}
		policyName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Policy")
		}

		err := testPolicyExists(false, policyID, policyName)
		if err != nil {
			return err
		}
	}
	return nil
}

func testPolicyExists(expected bool, policyID string, policyName string) error {
	client := testAccProvider.Meta().(*Config).oktaClient

	exists := false
	_, _, err := client.Policies.GetPolicy(policyID)
	if err != nil {
		if client.OktaErrorCode != "E0000007" {
			return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
		}
	} else {
		exists = true
	}

	if expected == true && exists == false {
		return fmt.Errorf("[ERROR] Policy %v not found in Okta", policyName)
	} else if expected == false && exists == true {
		return fmt.Errorf("[ERROR] Policy %v still exists in Okta", policyName)
	}
	return nil
}

func testOktaPolicySignOn(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, rInt, rInt)
}

func testOktaPolicySignOn_updated(rInt int) string {
	return fmt.Sprintf(`
data "okta_everyone_group" "everyone-%d" {}

resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "INACTIVE"
  priority    = 999
  description = "Terraform Acceptance Test SignOn Policy Updated"
  conditions {
    groups = [ "${data.okta_everyone_group.everyone-%d.id}" ]
  }
}
`, rInt, rInt, rInt, rInt)
}

func testOktaPolicySignOn_defaultErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "Default Policy"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, rInt)
}

func testOktaPolicySignOn_nameErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAccChanged-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
}
`, rInt, rInt)
}
func testOktaPolicySignOn_typeErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "PASSWORD"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
}
`, rInt, rInt)
}

func testOktaPolicySignOn_passErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  settings {
    password {
      minlength = 12
    }
  }
}
`, rInt, rInt)
}

func testOktaPolicySignOn_authpErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  conditions {
    authprovider {
      provider = "ACTIVE_DIRECTORY"
    }
  }
}
`, rInt, rInt)
}

func testOktaPolicyPassword(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "PASSWORD"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test Password Policy"
}
`, rInt, rInt)
}

// cannot change skipunlock to "true" if the authprovider is OKTA
// unless PASSWORD_POLICY_SOFT_LOCK is enabled
// (not supported in this TF provider at this time)
func testOktaPolicyPassword_updated(rInt int) string {
	return fmt.Sprintf(`
data "okta_everyone_group" "everyone-%d" {}

resource "okta_policies" "test-%d" {
  type        = "PASSWORD"
  name        = "testAcc-%d"
  status      = "INACTIVE"
  priority    = 999
  description = "Terraform Acceptance Test Password Policy Updated"
  conditions {
    groups = [ "${data.okta_everyone_group.everyone-%d.id}" ]
  }
  settings {
    password {
      minlength = 12
      minlowercase = 0
      minuppercase = 0
      minnumber = 0
      minsymbol = 0
      excludeusername = false
      excludeattributes = [ "firstName", "lastName" ]
      dictionarylookup = true
      maxagedays = 60
      expirewarndays = 15
      minageminutes = 60
      historycount = 5
      maxlockoutattempts = 3
      autounlockminutes = 2
      showlockoutfailures = true
      questionminlength = 10
      recoveryemailtoken = 20160
      smsrecovery = "ACTIVE"
    }
  }
}
`, rInt, rInt, rInt, rInt)
}
