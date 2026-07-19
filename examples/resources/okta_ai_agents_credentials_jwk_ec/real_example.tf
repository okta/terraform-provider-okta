# Real-world example with actual EC P-256 key coordinates
# These are example values - replace with your actual key coordinates

resource "okta_ai_agents_credentials_jwk_ec" "prod_signing_key" {
  # Replace with your actual AI agent ID
  agent_id = "00u1a2b3c4d5e6f7g8h9i"

  # Discriminator - must be "EC" for elliptic curve keys
  kty = "EC"

  # Algorithm - ES256 (ECDSA using P-256 and SHA-256)
  alg = "ES256"

  # Curve - P-256 (also known as secp256r1 or prime256v1)
  crv = "P-256"

  # Key ID - used to identify this key in the JWKS
  kid = "prod-ec-key-2024-01"

  # Use - "sig" indicates this key is used for signing
  use = "sig"

  # X coordinate of the EC public key (base64url-encoded)
  # This represents the x component of the point on the elliptic curve
  x = "WKn-ZIGevcwGIyyrzFoZNBdaq9_TsqzGl96oc0CWuis"

  # Y coordinate of the EC public key (base64url-encoded)
  # This represents the y component of the point on the elliptic curve
  y = "y77t-RvAHRKTsSGdIYUfweuOvwrvDD-Q3Hv5J0fSKcE"

  # Optional: Set the status of the key
  # status = "ACTIVE"
}

# Example with a different algorithm (ES384)
resource "okta_ai_agents_credentials_jwk_ec" "backup_key_384" {
  agent_id = "00u1a2b3c4d5e6f7g8h9i"

  kty = "EC"

  # Algorithm - ES384 (ECDSA using P-384 and SHA-384)
  alg = "ES384"

  # Curve - P-384 (secp384r1)
  crv = "P-384"

  kid = "prod-ec-key-2024-02-backup"

  use = "sig"

  # These would be coordinates from a P-384 key (longer than P-256)
  x = "gI0GAILBdu7T53-nP0_eAfGLBLO_5K4dHo0_yF2gSMvZD8o5uZvFn5_pZQRyBFe1"

  y = "w5RiAv0tBZFz9bWGqxdZqP1i8tIDvBu1wnPAi4JNKZ4Ypxqv7Pq8mK9N3L2Q5R6S7"
}

# To generate EC keys for testing:
#
# Using OpenSSL (Linux/Mac):
#   1. Generate EC private key:
#      openssl ecparam -name prime256v1 -genkey -noout -out private-key.pem
#
#   2. Extract public key:
#      openssl ec -in private-key.pem -pubout -out public-key.pem
#
#   3. View key in JWK format:
#      openssl ec -in public-key.pem -pubin -text -noout
#
#   4. Base64url-encode the coordinates (remove padding)
#
# Using Python:
#   from cryptography.hazmat.primitives.asymmetric import ec
#   from cryptography.hazmat.backends import default_backend
#   import base64
#
#   # Generate key
#   private_key = ec.generate_private_key(ec.SECP256R1(), default_backend())
#   public_key = private_key.public_key()
#
#   # Get coordinates
#   numbers = public_key.public_numbers()
#   x = base64.urlsafe_b64encode(numbers.x.to_bytes(32, 'big')).decode().rstrip('=')
#   y = base64.urlsafe_b64encode(numbers.y.to_bytes(32, 'big')).decode().rstrip('=')
#   print(f"x: {x}")
#   print(f"y: {y}")

# Import an existing EC key
# terraform import okta_ai_agents_credentials_jwk_ec.prod_signing_key 00u1a2b3c4d5e6f7g8h9i/00c1a2b3c4d5e6f7g8h9i
