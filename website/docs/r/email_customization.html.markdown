---
layout: 'okta'
page_title: 'Okta: okta_email_customization'
sidebar_current: 'docs-okta-resource-email-customization'
description: |-
Create an email customization of an email template belonging to a brand in an Okta organization.
---

# okta_email_customization

Use this resource to create an [email
customization](https://developer.okta.com/docs/reference/api/brands/#create-email-customization)
of an email template belonging to a brand in an Okta organization.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}

resource "okta_email_customization" "forgot_password_en_alt" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
  language      = "cs"
  is_default    = false
  subject       = "Forgot Password"
  body          = "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"
}
```

## Arguments Reference

- `brand_id` - (Required) Brand ID
- `template_name` - (Required) Template Name
  - Example values: `"AccountLockout"`,
`"ADForgotPassword"`,
`"ADForgotPasswordDenied"`,
`"ADSelfServiceUnlock"`,
`"ADUserActivation"`,
`"AuthenticatorEnrolled"`,
`"AuthenticatorReset"`,
`"ChangeEmailConfirmation"`,
`"EmailChallenge"`,
`"EmailChangeConfirmation"`,
`"EmailFactorVerification"`,
`"ForgotPassword"`,
`"ForgotPasswordDenied"`,
`"IGAReviewerEndNotification"`,
`"IGAReviewerNotification"`,
`"IGAReviewerPendingNotification"`,
`"IGAReviewerReassigned"`,
`"LDAPForgotPassword"`,
`"LDAPForgotPasswordDenied"`,
`"LDAPSelfServiceUnlock"`,
`"LDAPUserActivation"`,
`"MyAccountChangeConfirmation"`,
`"NewSignOnNotification"`,
`"OktaVerifyActivation"`,
`"PasswordChanged"`,
`"PasswordResetByAdmin"`,
`"PendingEmailChange"`,
`"RegistrationActivation"`,
`"RegistrationEmailVerification"`,
`"SelfServiceUnlock"`,
`"SelfServiceUnlockOnUnlockedAccount"`,
`"UserActivation"`
- `language` - The language supported by the customization
- `is_default` - Whether the customization is the default
  - Setting `is_default` to true when there is already a default customization will cause an error when this resource is created.
- `subject` - The subject of the customization
- `body` - The body of the customization

## Attributes Reference

- `id` - Customization ID
- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the email template
