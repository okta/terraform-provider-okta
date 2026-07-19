# okta_ai_agents_credentials_jwk_ec Resource

Manage EC (Elliptic Curve) JSON Web Keys for Okta AI Agents.

This resource allows you to create and manage public elliptic curve keys used for signing operations with Okta AI Agents.

## Example Usage

### Basic EC Key (P-256)

```hcl
resource "okta_ai_agents_credentials_jwk_ec" "example" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"
  kty      = "EC"
  alg      = "ES256"
  crv      = "P-256"
  kid      = "my-signing-key"
  use      = "sig"
  x        = "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"
  y        = "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKcE"
  status   = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

### Required

- `agent_id` - (Required) The ID of the AI agent in Okta. This field is immutable and forces replacement on change.
- `kty` - (Required) The key type. Must be set to `"EC"` for elliptic curve keys. This field is immutable and forces replacement on change.

### Optional

- `alg` - (Optional) The algorithm used for signing. Examples:
  - `"ES256"` - ECDSA using P-256 and SHA-256
  - `"ES384"` - ECDSA using P-384 and SHA-384
  - `"ES512"` - ECDSA using P-521 and SHA-512

- `crv` - (Optional) The cryptographic curve used for the key pair. Examples:
  - `"P-256"` - secp256r1 (most common)
  - `"P-384"` - secp384r1
  - `"P-521"` - secp521r1

- `kid` - (Optional) The key ID (kid). A unique identifier for this key in the JWKS.

- `use` - (Optional) The intended use of the public key. Examples:
  - `"sig"` - The key is used for signing
  - `"enc"` - The key is used for encryption

- `x` - (Optional) The public X coordinate of the elliptic curve point (base64url-encoded).

- `y` - (Optional) The public Y coordinate of the elliptic curve point (base64url-encoded).

- `status` - (Optional) The status of the key. Examples:
  - `"ACTIVE"`
  - `"INACTIVE"`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the JSON Web Key
- `created` - Timestamp of when the key was created
- `last_updated` - Timestamp of when the key was last updated

## Import

You can import an existing EC key using the format `{agent_id}/{key_id}`:

```shell
terraform import okta_ai_agents_credentials_jwk_ec.example 00u1a2b3c4d5e6f7g8h9i/00c1a2b3c4d5e6f7g8h9i
```

## Generating EC Keys

### Using OpenSSL

```bash
# Generate a P-256 key pair
openssl ecparam -name prime256v1 -genkey -noout -out private-key.pem

# Extract public key
openssl ec -in private-key.pem -pubout -out public-key.pem

# View key in JWK format
openssl ec -in public-key.pem -pubin -text -noout
```

### Using Python

```python
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
import base64

# Generate P-256 key
private_key = ec.generate_private_key(ec.SECP256R1(), default_backend())
public_key = private_key.public_key()

# Extract and encode coordinates
numbers = public_key.public_numbers()
x = base64.urlsafe_b64encode(numbers.x.to_bytes(32, 'big')).decode().rstrip('=')
y = base64.urlsafe_b64encode(numbers.y.to_bytes(32, 'big')).decode().rstrip('=')

print(f"x: {x}")
print(f"y: {y}")
```

## Supported Curves and Algorithms

| Curve | Algorithm | Use Case |
|-------|-----------|----------|
| P-256 (secp256r1) | ES256 | Default, most widely supported |
| P-384 (secp384r1) | ES384 | Higher security level |
| P-521 (secp521r1) | ES512 | Maximum security level |

## Notes

- Updates are not supported. Changes to any field will require resource replacement.
- The resource enforces immutability on `agent_id` and `kty` fields using `RequiresReplace`.
- Key coordinates (`x`, `y`) must be base64url-encoded as per RFC 7517 (JSON Web Key standard).
- Only EC keys are supported by this resource. For RSA keys, use `okta_ai_agents_credentials_jwk_rsa`.
