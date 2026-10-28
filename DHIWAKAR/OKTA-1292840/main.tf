terraform {
  required_providers {
    okta = {
      source = "okta/okta"
    }
  }
}

# Credentials come from the environment:
# OKTA_ORG_NAME, OKTA_BASE_URL, OKTA_API_TOKEN (or OKTA_API_CLIENT_ID/OKTA_API_PRIVATE_KEY)
#
# OKTA-1292840 — okta_request_condition.priority non-zero mismatch.
#
# Background: the governance API treats `priority` as a position within the
# resource's ordered list of conditions. It clamps an out-of-range request to
# the length of that list and does not echo back the value that was sent. The
# guard on master only restored the planned priority when the API returned
# exactly 0, so a non-zero reassignment (planned 5 -> API 3) fell through and
# the apply failed with "Provider produced inconsistent result after apply".
#
# Replace resource_id and approval_sequence_id with values from your org.
# `resource_id` forces replacement when changed.

variable "resource_id" {
  type        = string
  description = "ID of the resource (e.g. an app) the request conditions hang off."
  default     = "0oa14j7fa2peilw0n1d8"
}

variable "approval_sequence_id" {
  type        = string
  description = "ID of an existing approval sequence on that resource."
  default     = "68d224058c0cff364ca377e8"
}

# ---------------------------------------------------------------------------
# Case 1 — sparse priority, the regression this ticket fixes.
#
# Only one condition exists, so a requested priority of 5 is clamped by the
# API to a lower non-zero value. Before the fix this apply failed outright.
# Expect: apply succeeds, state holds 5, and a "Priority reassigned by API"
# warning is printed naming the value the server chose.
# ---------------------------------------------------------------------------
resource "okta_request_condition" "sparse" {
  resource_id          = var.resource_id
  approval_sequence_id = var.approval_sequence_id
  name                 = "okta-1292840-sparse"
  priority             = 0

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

# ---------------------------------------------------------------------------
# Case 2 — contiguous priorities, the scheme that already worked.
#
# Included as a control: these two should apply cleanly with no warning,
# confirming the fix did not disturb the common path.
# ---------------------------------------------------------------------------
resource "okta_request_condition" "dense_first" {
  resource_id          = var.resource_id
  approval_sequence_id = var.approval_sequence_id
  name                 = "okta-1292840-dense-first"
  priority             = 2

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }
}

resource "okta_request_condition" "dense_second" {
  resource_id          = var.resource_id
  approval_sequence_id = var.approval_sequence_id
  name                 = "okta-1292840-dense-second"
  priority             = 1

  access_scope_settings {
    type = "RESOURCE_DEFAULT"
  }
  requester_settings {
    type = "EVERYONE"
  }

}

output "sparse_priority_in_state" {
  description = "Should be 5 — the configured value, not the server's clamped one."
  value       = okta_request_condition.sparse.priority
}
