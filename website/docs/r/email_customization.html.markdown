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
Make use of `depends_on` meta argument to ensure the provider navigates email customization
language versions seamlessly. Have all secondary customizations depend on the primary
customization that is marked default. See [Example Usage](#example-usage).

~> Caveats for [creating an email
customization](https://developer.okta.com/docs/reference/api/brands/#response-body-19).
If this is the first customization being created for the email template, and
`is_default` is not set for the customization in its resource configuration, the
API will respond with the created customization marked as default. The API will
400 if the language parameter is not one of the supported languages or the body
parameter does not contain a required variable reference. The API will error 409
if `is_default` is true and a default customization exists. The API will 404 for
an invalid `brand_id` or `template_name`.

~> Caveats for [updating an email
customization](https://developer.okta.com/docs/reference/api/brands/#response-body-22).
If the `is_default` parameter is true, the previous default email customization
has its `is_default` set to false (see previous note about mitigating this with
`depends_on` meta argument). The API will 409 if there’s already another email
customization for the specified language or the `is_default` parameter is false
and the email customization being updated is the default. The API will 400 if
the language parameter is not one of the supported locales or the body parameter
does not contain a required variable reference.  The API will 404 for an invalid
`brand_id` or `template_name`.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}

resource "okta_email_customization" "forgot_password_en" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
  language      = "en"
  is_default    = true
  subject       = "Account password reset"
  body          = "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"
}

resource "okta_email_customization" "forgot_password_es" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
  language      = "es"
  subject       = "Restablecimiento de contraseña de cuenta"
  body          = "Hola $$user.firstName,<br/><br/>Haga clic en este enlace para restablecer tu contraseña: $$resetPasswordLink"

  depends_on = [
    okta_email_customization.forgot_password_en
  ]
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
- `subject` - The subject of the customization
- `body` - The body of the customization
- `force_is_default` (Deprecated) `force_is_default` is deprecated and now is a no-op in behavior. Rely upon the `depends_on` meta argument to force dependency of secondary templates to the default template",

## Attributes Reference

- `id` - Customization ID
- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the email template

## Import

An email customization can be imported using the customization ID, brand ID and template name.

```
$ terraform import okta_email_customization.example &#60;customization_id&#62;/&#60;brand_id&#62;/&#60;template_name&#62;
```
