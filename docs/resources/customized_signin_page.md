---
page_title: "Resource: okta_customized_signin_page"
description: |-
  Manage the customized signin page of a brand
---

# Resource: okta_customized_signin_page

Manage the customized signin page of a brand

## Example Usage

```terraform
resource "okta_brand" "test" {
  name   = "testBrand"
  locale = "en"
}

resource "okta_customized_signin_page" "test" {
  brand_id       = resource.okta_brand.test.id
  page_content   = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">\n<html>\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n    <meta name=\"robots\" content=\"noindex,nofollow\" />\n    <!-- Styles generated from theme -->\n    <link href=\"{{themedStylesUrl}}\" rel=\"stylesheet\" type=\"text/css\">\n    <!-- Favicon from theme -->\n    <link rel=\"shortcut icon\" href=\"{{faviconUrl}}\" type=\"image/x-icon\"/>\n\n    <title>{{pageTitle}}</title>\n    {{{SignInWidgetResources}}}\n\n    <style nonce=\"{{nonceValue}}\">\n        #login-bg-image-id {\n            background-image: {{bgImageUrl}}\n        }\n    </style>\n</head>\n<body>\n    <div id=\"login-bg-image-id\" class=\"login-bg-image tb--background\"></div>\n    <div id=\"okta-login-container\"></div>\n\n    <!--\n        \"OktaUtil\" defines a global OktaUtil object\n        that contains methods used to complete the Okta login flow.\n     -->\n    {{{OktaUtil}}}\n\n    <script type=\"text/javascript\" nonce=\"{{nonceValue}}\">\n        // \"config\" object contains default widget configuration\n        // with any custom overrides defined in your admin settings.\n        var config = OktaUtil.getSignInWidgetConfig();\n\n        // Render the Okta Sign-In Widget\n        var oktaSignIn = new OktaSignIn(config);\n        oktaSignIn.renderEl({ el: '#okta-login-container' },\n            OktaUtil.completeLogin,\n            function(error) {\n                // Logs errors that occur when configuring the widget.\n                // Remove or replace this with your own custom error handler.\n                console.log(error.message, error);\n            }\n        );\n    </script>\n</body>\n</html>\n"
  widget_version = "^6"
  widget_customizations {
    widget_generation = "G3"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `brand_id` (String) brand id of the preview signin page
- `page_content` (String) page content of the preview signin page
- `widget_version` (String) widget version specified as a Semver. The following are currently supported
			*, ^1, ^2, ^3, ^4, ^5, ^6, ^7, 1.6, 1.7, 1.8, 1.9, 1.10, 1.11, 1.12, 1.13, 2.1, 2.2, 2.3, 2.4,
			2.5, 2.6, 2.7, 2.8, 2.9, 2.10, 2.11, 2.12, 2.13, 2.14, 2.15, 2.16, 2.17, 2.18, 2.19, 2.20, 2.21,
			3.0, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 4.0, 4.1, 4.2, 4.3, 4.4, 4.5, 5.0, 5.1, 5.2, 5.3,
			5.4, 5.5, 5.6, 5.7, 5.8, 5.9, 5.10, 5.11, 5.12, 5.13, 5.14, 5.15, 5.16, 6.0, 6.1, 6.2, 6.3, 6.4, 6.5,
			6.6, 6.7, 6.8, 6.9, 7.0, 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 7.9, 7.10, 7.11, 7.12, 7.13.

### Optional

- `content_security_policy_setting` (Block, Optional) (see [below for nested schema](#nestedblock--content_security_policy_setting))
- `widget_customizations` (Block, Optional) (see [below for nested schema](#nestedblock--widget_customizations))

### Read-Only

- `id` (String) placeholder id

<a id="nestedblock--content_security_policy_setting"></a>
### Nested Schema for `content_security_policy_setting`

Optional:

- `mode` (String) enforced or report_only
- `report_uri` (String)
- `src_list` (List of String)


<a id="nestedblock--widget_customizations"></a>
### Nested Schema for `widget_customizations`

Required:

- `widget_generation` (String)

Optional:

- `authenticator_page_custom_link_label` (String)
- `authenticator_page_custom_link_url` (String)
- `classic_recovery_flow_email_or_username_label` (String)
- `custom_link_1_label` (String)
- `custom_link_1_url` (String)
- `custom_link_2_label` (String)
- `custom_link_2_url` (String)
- `forgot_password_label` (String)
- `forgot_password_url` (String)
- `help_label` (String)
- `help_url` (String)
- `password_info_tip` (String)
- `password_label` (String)
- `show_password_visibility_toggle` (Boolean)
- `show_user_identifier` (Boolean)
- `sign_in_label` (String)
- `unlock_account_label` (String)
- `unlock_account_url` (String)
- `username_info_tip` (String)
- `username_label` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import okta_customized_signin_page.example <brand_id>
```
