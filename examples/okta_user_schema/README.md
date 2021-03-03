# okta_user_schemas

Represents an Okta User Profile Attribute Schema. [See Okta documentation for more details](https://developer.okta.com/docs/api/resources/users).

- An example of a user with multiple custom attributes, [can be found here](../okta_user/custom_attributes.tf). Note the `depends_on` see https://github.com/okta/terraform-provider-okta/issues/144 for more info.
