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

~> Okta's public API is strict regarding the behavior of the `is_default`
property in [an email
customization](https://developer.okta.com/docs/reference/api/brands/#email-customization).
When a customization is
[created](https://developer.okta.com/docs/reference/api/brands/#create-email-customization)
it can not be created with an `is_default` value of `true` if there is already a
default customization. If an email customization is the last of the template
type it can not be
[deleted](https://developer.okta.com/docs/reference/api/brands/#delete-email-customization).
And the `is_default` value can't be set to false when updating the last
remaining customization. **To allow this resource to be more flexible** set the
`force_is_default` property to `create`, `destroy`, or `create,destroy`. This
will cause all the customizations to be
[reset/deleted](https://developer.okta.com/docs/reference/api/brands/#delete-all-email-customizations)
for a create when there is a `create` value in `force_is_default` and
`is_default` is `true`.  Likewise reset will be called for a delete when there
is a `delete` value in `force_is_default` and `is_default` is `true`.

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
  language      = "en"
  is_default    = true
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
  - Example values from [supported languages](https://developer.okta.com/docs/reference/api/brands/#supported-languages): 
    `"cs"`,
    `"da"`,
    `"de"`,
    `"el"`,
    `"en"`,
    `"es"`,
    `"fi"`,
    `"fr"`,
    `"hu"`,
    `"id"`,
    `"it"`,
    `"ja"`,
    `"ko"`,
    `"ms"`,
    `"nb"`,
    `"nl-NL"`,
    `"pl"`,
    `"pt-BR"`,
    `"ro"`,
    `"ru"`,
    `"sv"`,
    `"th"`,
    `"tr"`,
    `"uk"`,
    `"vi"`,
    `"zh-CN"`,
    `"zh-TW"`
- `is_default` - Whether the customization is the default
  - Setting `is_default` to true when there is already a default customization will cause an error when this resource is created.
- `subject` - The subject of the customization
- `body` - The body of the customization
- `force_is_default` Force `is_default` on the create and delete operation by
   deleting all email customizations. See Note above explaing email customization API
   behavior and [API
   documentation](https://developer.okta.com/docs/reference/api/brands/#list-email-customizations).
   Valid values `create`, `delete`, `create,delete`.

## Attributes Reference

- `id` - Customization ID
- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the email template

## Import

An email customization can be imported using the customization ID, brand ID and template name.

```
$ terraform import okta_email_customization.example &#60;customization_id&#62;/&#60;brand_id&#62;/&#60;template_name&#62;
```
