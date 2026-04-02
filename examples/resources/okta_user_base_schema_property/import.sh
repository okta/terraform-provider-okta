### User schema property of default user type can be imported via the property index.
terraform import okta_user_base_schema_property.example <property_name>

### User schema property of custom user type can be imported via user type id and property index
terraform import okta_user_base_schema_property.example <user_type_id>.<property name>