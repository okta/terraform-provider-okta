---
layout: 'okta' 
page_title: 'Okta: okta_user_security_questions' 
sidebar_current: 'docs-okta-datasource-user-security-questions' 
description: |- 
  Get a list of user's security questions.
---

# okta_user_security_questions

Use this data source to retrieve a list of user's security questions.

## Example Usage

```hcl
resource "okta_user" "example" {
  first_name = "John"
  last_name = "Smith"
  login = "john.smith@example.com"
  email = "john.smith@example.com"
}

data "okta_user_security_questions" example {
  user_id = okta_user.example.id
}
```

## Arguments Reference

- `user_id` - (Required) User ID.

## Attributes Reference

- `id` - User ID.

- `questions` - collection of user's security question retrieved from Okta with the following properties:
    - `key` - Security question unique key.
    - `text` - Display text for security question.
