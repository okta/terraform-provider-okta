---
page_title: "Resource: okta_email_customization"
description: |-
  Create an email customization of an email template belonging to a brand in an Okta organization.
  Use this resource to create an email
  customization https://developer.okta.com/docs/reference/api/brands/#create-email-customization
  of an email template belonging to a brand in an Okta organization.
  ~> Okta's public API is strict regarding the behavior of the 'isdefault'
  property in an email
  customization https://developer.okta.com/docs/reference/api/brands/#email-customization.
  Make use of 'dependson' meta argument to ensure the provider navigates email customization
  language versions seamlessly. Have all secondary customizations depend on the primary
  customization that is marked default. See Example Usage.
  ~> Caveats for creating an email
  customization https://developer.okta.com/docs/reference/api/brands/#response-body-19.
  If this is the first customization being created for the email template, and
  'isdefault' is not set for the customization in its resource configuration, the
  API will respond with the created customization marked as default. The API will
  400 if the language parameter is not one of the supported languages or the body
  parameter does not contain a required variable reference. The API will error 409
  if 'isdefault' is true and a default customization exists. The API will 404 for
  an invalid 'brandid' or 'templatename'.
  ~> Caveats for updating an email
  customization https://developer.okta.com/docs/reference/api/brands/#response-body-22.
  If the 'isdefault' parameter is true, the previous default email customization
  has its 'isdefault' set to false (see previous note about mitigating this with
  'dependson' meta argument). The API will 409 if there’s already another email
  customization for the specified language or the 'isdefault' parameter is false
  and the email customization being updated is the default. The API will 400 if
  the language parameter is not one of the supported locales or the body parameter
  does not contain a required variable reference.  The API will 404 for an invalid
  'brandid' or 'templatename'.
---

# Resource: okta_email_customization

Create an email customization of an email template belonging to a brand in an Okta organization.
Use this resource to create an [email
customization](https://developer.okta.com/docs/reference/api/brands/#create-email-customization)
of an email template belonging to a brand in an Okta organization.

~> Okta's public API is strict regarding the behavior of the 'is_default'
property in [an email
customization](https://developer.okta.com/docs/reference/api/brands/#email-customization).
Make use of 'depends_on' meta argument to ensure the provider navigates email customization
language versions seamlessly. Have all secondary customizations depend on the primary
customization that is marked default. See [Example Usage](#example-usage).

~> Caveats for [creating an email
customization](https://developer.okta.com/docs/reference/api/brands/#response-body-19).
If this is the first customization being created for the email template, and
'is_default' is not set for the customization in its resource configuration, the
API will respond with the created customization marked as default. The API will
400 if the language parameter is not one of the supported languages or the body
parameter does not contain a required variable reference. The API will error 409
if 'is_default' is true and a default customization exists. The API will 404 for
an invalid 'brand_id' or 'template_name'.

~> Caveats for [updating an email
customization](https://developer.okta.com/docs/reference/api/brands/#response-body-22).
If the 'is_default' parameter is true, the previous default email customization
has its 'is_default' set to false (see previous note about mitigating this with
'depends_on' meta argument). The API will 409 if there’s already another email
customization for the specified language or the 'is_default' parameter is false
and the email customization being updated is the default. The API will 400 if
the language parameter is not one of the supported locales or the body parameter
does not contain a required variable reference.  The API will 404 for an invalid
'brand_id' or 'template_name'.

## Example Usage

```terraform
data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
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

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `brand_id` (String) Brand ID
- `template_name` (String) Template Name - Example values: `AccountLockout`,`ADForgotPassword`,`ADForgotPasswordDenied`,`ADSelfServiceUnlock`,`ADUserActivation`,`AuthenticatorEnrolled`,`AuthenticatorReset`,`ChangeEmailConfirmation`,`EmailChallenge`,`EmailChangeConfirmation`,`EmailFactorVerification`,`ForgotPassword`,`ForgotPasswordDenied`,`IGAReviewerEndNotification`,`IGAReviewerNotification`,`IGAReviewerPendingNotification`,`IGAReviewerReassigned`,`LDAPForgotPassword`,`LDAPForgotPasswordDenied`,`LDAPSelfServiceUnlock`,`LDAPUserActivation`,`MyAccountChangeConfirmation`,`NewSignOnNotification`,`OktaVerifyActivation`,`PasswordChanged`,`PasswordResetByAdmin`,`PendingEmailChange`,`RegistrationActivation`,`RegistrationEmailVerification`,`SelfServiceUnlock`,`SelfServiceUnlockOnUnlockedAccount`,`UserActivation`

### Optional

- `body` (String) The body of the customization
- `force_is_default` (String, Deprecated) Force is_default on the create and delete by deleting all email customizations. Comma separated string with values of 'create' or 'destroy' or both `create,destroy'.
- `is_default` (Boolean) Whether the customization is the default
- `language` (String) The language supported by the customization - Example values from [supported languages](https://developer.okta.com/docs/reference/api/brands/#supported-languages)
- `subject` (String) The subject of the customization

### Read-Only

- `id` (String) The ID of the customization
- `links` (String) Link relations for this object - JSON HAL - Discoverable resources related to the email template

## Import

Import is supported using the following syntax:

```shell
terraform import okta_email_customization.example &#60;customization_id&#62;/&#60;brand_id&#62;/&#60;template_name&#62;
```
