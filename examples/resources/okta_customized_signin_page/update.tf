resource "okta_brand" "test" {
  name   = "testBrand"
  locale = "en"
}

resource "okta_customized_signin_page" "test" {
  brand_id       = resource.okta_brand.test.id
  page_content   = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">\n<html>\n<head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n    <meta name=\"robots\" content=\"noindex,nofollow\" />\n    <!-- Styles generated from theme -->\n    <link href=\"{{themedStylesUrl}}\" rel=\"stylesheet\" type=\"text/css\">\n    <!-- Favicon from theme -->\n    <link rel=\"shortcut icon\" href=\"{{faviconUrl}}\" type=\"image/x-icon\"/>\n\n    <title>{{pageTitle}}</title>\n    {{{SignInWidgetResources}}}\n\n    <style nonce=\"{{nonceValue}}\">\n        #login-bg-image-id {\n            background-image: {{bgImageUrl}}\n        }\n    </style>\n</head>\n<body>\n    <div id=\"login-bg-image-id\" class=\"login-bg-image tb--background\"></div>\n    <div id=\"okta-login-container\"></div>\n\n    <!--\n        \"OktaUtil\" defines a global OktaUtil object\n        that contains methods used to complete the Okta login flow.\n     -->\n    {{{OktaUtil}}}\n\n    <script type=\"text/javascript\" nonce=\"{{nonceValue}}\">\n        // \"config\" object contains default widget configuration\n        // with any custom overrides defined in your admin settings.\n        var config = OktaUtil.getSignInWidgetConfig();\n\n        // Render the Okta Sign-In Widget\n        var oktaSignIn = new OktaSignIn(config);\n        oktaSignIn.renderEl({ el: '#okta-login-container' },\n            OktaUtil.completeLogin,\n            function(error) {\n                // Logs errors that occur when configuring the widget.\n                // Remove or replace this with your own custom error handler.\n                console.log(error.message, error);\n            }\n        );\n    </script>\n</body>\n</html>\n"
  widget_version = "^6"
  widget_customizations {
    widget_generation = "G2"
  }
  content_security_policy_setting {
    mode       = "report_only"
    report_uri = ""
    src_list   = ["https://idp.example.com/authorize", "https://idp.example.com/authoriz"]
  }
}


resource "okta_customized_signin_page" "test-2" {
  brand_id       = "bndnkw1sc3flIOTF51d7"
  widget_customizations {
    widget_generation = "G3"
    help_url          = "https://helpurltestupdated.com"
    help_label        = "Help URL Test Updated"
    custom_link_1_url = "https://customlink1url.com"
    custom_link_1_label = "Custom Link 1 URL"
  }
}