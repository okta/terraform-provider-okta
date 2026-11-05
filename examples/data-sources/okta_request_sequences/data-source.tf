data "okta_request_sequences" "example" {
  resource_id = "<okta_resource_id>"
}

# Look up a sequence by name
locals {
  manager_approval_sequence_id = one([
    for seq in data.okta_request_sequences.example.request_sequences : seq.id
    if seq.name == "Requester's Manager Approval"
  ])
}
