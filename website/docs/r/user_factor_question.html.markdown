---
layout: 'okta'
page_title: 'Okta: okta_user_factor_question'
sidebar_current: 'docs-okta-resource-user-factor-question'
description: |-
    Creates security question factor for a user.
---

# okta_user_factor_question

Creates security question factor for a user.

This resource allows you to create and configure security question factor for a user.

## Example Usage

```hcl
data "okta_user_security_questions" example {
  user_id = okta_user.example.id
}

resource "okta_user" "example" {
  first_name = "John"
  last_name = "Smith"
  login = "john.smith@example.com"
  email = "john.smith@example.com"
}

resource "okta_factor" "example" {
  provider_id = "okta_question"
  active = true
}

resource "okta_user_factor_question" "example" {
  user_id = okta_user.example.id
  key = data.okta_user_security_questions.example.questions[0].key
  answer = "meatball"
  depends_on = [
    okta_factor.example]
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required) ID of the user. Resource will be recreated when `user_id` changes.

- `key` - (Required) Security question unique key. 

- `answer` - (Required) Security question answer. Note here that answer won't be set during the resource import.

## Attributes Reference

- `id` - ID of the security question factor.

- `status` - The status of the security question factor.

- `text` - Display text for security question.

## Import

Security question factor for a user can be imported via the `user_id` and the `factor_id`.

```
$ terraform import okta_user_factor_question.example &#60;user id&#62;/&#60;question factor id&#62;
```
