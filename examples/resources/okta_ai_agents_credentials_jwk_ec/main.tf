# Example: Complete setup with provider configuration

terraform {
  required_version = ">= 1.0"

  required_providers {
    okta = {
      source  = "okta/okta"
      version = ">= 4.0"
    }
  }
}

provider "okta" {
  org_name  = var.okta_org_name
  base_url  = var.okta_base_url
  api_token = var.okta_api_token
}

# Variables
variable "okta_org_name" {
  description = "Okta organization name (e.g., 'dev', 'staging', 'prod')"
  type        = string
}

variable "okta_base_url" {
  description = "Okta base URL (e.g., 'okta.com', 'oktapreview.com')"
  type        = string
  default     = "okta.com"
}

variable "okta_api_token" {
  description = "Okta API token"
  type        = string
  sensitive   = true
}

variable "agent_id" {
  description = "ID of the AI agent in Okta"
  type        = string
}

variable "key_coordinates" {
  description = "EC public key coordinates"
  type = object({
    x = string
    y = string
  })
}

# Create EC JSON Web Key for the AI Agent
resource "okta_ai_agents_credentials_jwk_ec" "agent_signing_key" {
  agent_id = var.agent_id

  # Key type - must be "EC" for elliptic curve
  kty = "EC"

  # Signing algorithm
  alg = "ES256"

  # Elliptic curve
  crv = "P-256"

  # Key identifier
  kid = "terraform-managed-key-${formatdate("YYYY-MM-DD", timestamp())}"

  # Key use
  use = "sig"

  # Public key coordinates
  x = var.key_coordinates.x
  y = var.key_coordinates.y

  # Key status
  status = "ACTIVE"

  # Tag for identification
  # Note: tags might not be supported - check your provider version
  # tags = ["terraform-managed", "signing-key"]
}

# Outputs
output "jwk_ec_id" {
  description = "The ID of the created EC JSON Web Key"
  value       = okta_ai_agents_credentials_jwk_ec.agent_signing_key.id
}

output "jwk_ec_kid" {
  description = "The key ID (kid) of the EC JSON Web Key"
  value       = okta_ai_agents_credentials_jwk_ec.agent_signing_key.kid
}

output "jwk_ec_created" {
  description = "Timestamp of when the key was created"
  value       = okta_ai_agents_credentials_jwk_ec.agent_signing_key.created
}

output "jwk_ec_status" {
  description = "Current status of the key"
  value       = okta_ai_agents_credentials_jwk_ec.agent_signing_key.status
}

# Usage instructions:
#
# 1. Set environment variables:
#    export TF_VAR_okta_org_name="your-org"
#    export TF_VAR_okta_api_token="your-api-token"
#    export TF_VAR_agent_id="00u1a2b3c4d5e6f7g8h9i"
#    export TF_VAR_key_coordinates='{"x":"WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis","y":"y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKcE"}'
#
# 2. Or create a terraform.tfvars file:
#    okta_org_name = "your-org"
#    okta_api_token = "your-api-token"
#    agent_id = "00u1a2b3c4d5e6f7g8h9i"
#    key_coordinates = {
#      x = "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"
#      y = "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKcE"
#    }
#
# 3. Initialize, plan, and apply:
#    terraform init
#    terraform plan
#    terraform apply
#
# 4. To import an existing key:
#    terraform import okta_ai_agents_credentials_jwk_ec.agent_signing_key 00u1a2b3c4d5e6f7g8h9i/00c1a2b3c4d5e6f7g8h9i
#
# 5. To destroy:
#    terraform destroy
