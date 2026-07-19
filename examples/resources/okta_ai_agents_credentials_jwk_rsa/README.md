# okta_ai_agents_credentials_jwk_rsa - Sample Terraform Configs

Sample Terraform configurations for the `okta_ai_agents_credentials_jwk_rsa` resource.

## Files Included

- **sample.tf** - Standalone configuration with 3 examples (RS256, PS256, RSA-OAEP)
- **main.tf** - Complete project setup with provider and variables
- **README.md** - This file

## Quick Start

### 1. Copy sample.tf

```bash
cp sample.tf my-rsa-keys.tf
```

### 2. Update Values

Replace placeholders:
- `your-org` → Your Okta organization
- `00u1a2b3c4d5e6f7g8h9i` → Your AI Agent ID
- `n` value → Your RSA public key modulus (base64url)

### 3. Set API Token

```bash
export TF_VAR_okta_api_token="your-api-token"
```

### 4. Deploy

```bash
terraform init
terraform plan
terraform apply
```

## RSA Key Algorithms

| Algorithm | Padding | Use Case | Security |
|-----------|---------|----------|----------|
| RS256 | PKCS#1 v1.5 | Signing | Good |
| RS384 | PKCS#1 v1.5 | Signing | Better |
| RS512 | PKCS#1 v1.5 | Signing | Best |
| PS256 | PSS | Signing | Better |
| PS384 | PSS | Signing | Better |
| PS512 | PSS | Signing | Best |
| RSA-OAEP | OAEP | Encryption | Recommended |

**Recommendation**: Use PSS algorithms (PS256/PS384/PS512) for new deployments.

## Generating RSA Keys

### Using OpenSSL

```bash
# Generate RSA-2048 private key
openssl genrsa -out private-key.pem 2048

# Extract public key
openssl rsa -in private-key.pem -pubout -out public-key.pem

# View key details
openssl rsa -in public-key.pem -pubin -text -noout

# Extract and encode n and e values (manual)
# Use Python script below for automation
```

### Using Python (Recommended)

```python
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
import base64

# Load public key from PEM file
with open('public-key.pem', 'rb') as f:
    public_key = serialization.load_pem_public_key(
        f.read(), 
        backend=default_backend()
    )

# Get RSA numbers
public_numbers = public_key.public_numbers()

# Encode n (modulus)
n_bytes = public_numbers.n.to_bytes(
    (public_numbers.n.bit_length() + 7) // 8, 
    'big'
)
n = base64.urlsafe_b64encode(n_bytes).decode().rstrip('=')

# Encode e (exponent) - usually 65537 (0x10001) = "AQAB"
e_bytes = public_numbers.e.to_bytes(
    (public_numbers.e.bit_length() + 7) // 8, 
    'big'
)
e = base64.urlsafe_b64encode(e_bytes).decode().rstrip('=')

print(f'n = "{n}"')
print(f'e = "{e}"')
```

### Using Node.js

```bash
npm install -g jwk-cli
jwk-cli generate -t RSA -s 2048 -f PEM
```

## Example Configurations

### Simple RSA-2048 with RS256

```hcl
resource "okta_ai_agents_credentials_jwk_rsa" "my_key" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"
  
  kty = "RSA"
  alg = "RS256"
  kid = "my-signing-key"
  use = "sig"
  
  n = "your-modulus-here"
  e = "AQAB"
  
  status = "ACTIVE"
}
```

### Advanced: Multiple Algorithms

```hcl
resource "okta_ai_agents_credentials_jwk_rsa" "primary" {
  agent_id = var.agent_id
  kty      = "RSA"
  alg      = "PS256"    # Recommended
  kid      = "primary-key"
  use      = "sig"
  n        = var.rsa_n
  e        = "AQAB"
}

resource "okta_ai_agents_credentials_jwk_rsa" "backup" {
  agent_id = var.agent_id
  kty      = "RSA"
  alg      = "RS256"    # Backup
  kid      = "backup-key"
  use      = "sig"
  n        = var.rsa_n_backup
  e        = "AQAB"
}
```

## Resource Arguments

### Required

- `agent_id` - AI Agent ID (immutable)
- `kty` - Key type, must be "RSA" (immutable)

### Optional

- `alg` - Algorithm (RS256, RS384, RS512, PS256, PS384, PS512, RSA-OAEP)
- `kid` - Key ID for JWKS
- `use` - Key usage ("sig" for signing, "enc" for encryption)
- `n` - RSA modulus (base64url-encoded)
- `e` - RSA exponent (base64url-encoded, default: "AQAB")
- `status` - Key status (ACTIVE, INACTIVE)

## Resource Attributes

- `id` - Unique key identifier
- `created` - Creation timestamp
- `last_updated` - Last update timestamp

## Import

```bash
terraform import okta_ai_agents_credentials_jwk_rsa.my_key {agent_id}/{key_id}

# Example:
terraform import okta_ai_agents_credentials_jwk_rsa.my_key 00u1a2b3c4d5e6f7g8h9i/00c1a2b3c4d5e6f7g8h9i
```

## Differences: RSA vs EC

| Feature | RSA | EC |
|---------|-----|-----|
| Key Size | 2048, 3072, 4096 bits | P-256, P-384, P-521 |
| Algorithm | RS256/RS384/RS512, PS256/PS384/PS512, RSA-OAEP | ES256/ES384/ES512 |
| Use Case | General purpose, widely supported | Modern, smaller keys |
| Performance | Slower | Faster |
| Key Size | Larger | Smaller (more secure/bit) |
| Compatibility | Maximum | Newer systems only |

## Testing

```bash
# Plan only
terraform plan

# Apply to staging
terraform apply -var-file=staging.tfvars

# Import existing key
terraform import okta_ai_agents_credentials_jwk_rsa.existing {agent_id}/{key_id}

# Destroy keys
terraform destroy
```

## Security Best Practices

1. **Use PSS algorithms** (PS256/PS384/PS512) for new deployments
2. **Rotate keys regularly** - Create new keys and archive old ones
3. **Never commit secrets** - Use environment variables for `n` and API tokens
4. **Use separate keys** for signing and encryption
5. **Store private keys** securely (not in Terraform)
6. **Monitor key usage** in Okta audit logs

## Troubleshooting

### "Invalid base64url encoding"
- Ensure `n` and `e` are base64url-encoded (no `=` padding)
- Check key generation script output

### "agent_id not found"
- Verify agent ID exists in your Okta org
- Check API token permissions

### "API returned a different variant type"
- Using RSA config for EC key or vice versa
- Verify key type matches resource type

## Additional Resources

- [RSA Specification (RFC 3447)](https://tools.ietf.org/html/rfc3447)
- [JSON Web Key (JWK) Spec (RFC 7517)](https://tools.ietf.org/html/rfc7517)
- [Okta AI Agents Docs](https://developer.okta.com/)
