---
page_title: "okta_iam_assignees_user Data Source - terraform-provider-okta"
description: |-
  Lists all Okta users with IAM role assignments.
---

# okta_iam_assignees_user (Data Source)

Use this data source to list all Okta users that have IAM role assignments.

## Example Usage

```terraform
data "okta_iam_assignees_user" "example" {}
```

## Schema

### Read-Only

- `id` (String) Static identifier `okta_iam_assignees_user_list`.
- `items` (List of Object) List of users with IAM role assignments.
  - `id` (String) Unique identifier of the user.
  - `orn` (String) ORN representing the assignee.
