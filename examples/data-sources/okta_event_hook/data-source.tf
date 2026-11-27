data "okta_event_hook" "example" {
  id = "<event_hook_id>"
}

output "event_hook_uri" {
  value = data.okta_event_hook.example.channel.config.uri
}
