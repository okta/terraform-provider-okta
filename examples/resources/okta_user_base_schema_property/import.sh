### User schema property of default user type can be imported via the property index.
terraform import okta_user_base_schema_property.example &#60;property name&#62;

### User schema property of custom user type can be imported via user type id and property index
terraform import okta_user_base_schema_property.example &#60;user type id&#62;.&#60;property name&#62;