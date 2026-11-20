resource "okta_ai_agents_schema_default" "test" {
  description = "Default AI agent schema"

  definitions {
    custom {
      properties = jsonencode({
        department = {
          title       = "Department"
          type        = "string"
          description = "The department that owns this AI agent"
          mutability  = "READ_WRITE"
          required    = false
        }
      })
    }
  }
}
