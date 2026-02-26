# okta_event_hook

This resource represents an Okta Event Hook. For more information see
the [API docs](https://developer.okta.com/docs/api/resources/event-hooks)

- Example of a simple user create/delete hook [can be found here](./basic.tf)
- Example of a simple inactive user CRUD hook [can be found here](./basic_updated.tf)
- Example of an event hook with filter [can be found here](./basic_with_filter.tf)

## Event Hook Filters

Event hook filters allow you to reduce the number of event hook calls by filtering events based on Okta Expression Language conditions. This is a self-service Early Access (EA) feature.

The `event_filter` block supports:
- `type` - (Optional) The type of filter. Currently only supports 'EXPRESSION_LANGUAGE'. Defaults to 'EXPRESSION_LANGUAGE'.
- `event_filter_map` - (Optional) Array of objects that map events to filter conditions.
  - `event` - (Required) The event type to filter.
  - `condition` - (Required) The condition object containing the filter expression.
    - `expression` - (Required) The Okta Expression Language statement that filters the event type.

### Example Filter Expression

To filter group membership additions to only the "Sales" group:
```
event.target.?[type eq 'UserGroup'].size()>0 && event.target.?[displayName eq 'Sales'].size()>0
```
