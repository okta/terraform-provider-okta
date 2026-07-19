# Complete RSA JWK Configuration with Provider Setup

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
  description = "Okta organization name"
  type        = string
}

variable "okta_base_url" {
  description = "Okta base URL"
  type        = string
  default     = "okta.com"
}

variable "okta_api_token" {
  description = "Okta API token"
  type        = string
  sensitive   = true
}

variable "agent_id" {
  description = "AI Agent ID"
  type        = string
}

variable "rsa_modulus" {
  description = "RSA public key modulus (n) - base64url encoded"
  type        = string
  sensitive   = true
}

variable "rsa_exponent" {
  description = "RSA public key exponent (e) - typically AQAB"
  type        = string
  default     = "AQAB"
}

# Create RSA signing key
resource "okta_ai_agents_credentials_jwk_rsa" "signing_key" {
  agent_id = var.agent_id

  kty = "RSA"
  alg = "RS256"
  kid = "terraform-rsa-key-${formatdate("YYYY-MM-DD", timestamp())}"
  use = "sig"

  n = var.rsa_modulus
  e = var.rsa_exponent

  status = "ACTIVE"
}

# Create RSA encryption key
resource "okta_ai_agents_credentials_jwk_rsa" "encryption_key" {
  agent_id = var.agent_id

  kty = "RSA"
  alg = "RSA-OAEP"
  kid = "terraform-rsa-enc-key-${formatdate("YYYY-MM-DD", timestamp())}"
  use = "enc"

  n = var.rsa_modulus
  e = var.rsa_exponent

  status = "ACTIVE"
}

# Outputs
output "signing_key_id" {
  description = "ID of the RSA signing key"
  value       = okta_ai_agents_credentials_jwk_rsa.signing_key.id
}

output "signing_key_kid" {
  description = "Key ID of the RSA signing key"
  value       = okta_ai_agents_credentials_jwk_rsa.signing_key.kid
}

output "encryption_key_id" {
  description = "ID of the RSA encryption key"
  value       = okta_ai_agents_credentials_jwk_rsa.encryption_key.id
}

output "encryption_key_kid" {
  description = "Key ID of the RSA encryption key"
  value       = okta_ai_agents_credentials_jwk_rsa.encryption_key.kid
}
