# Testing the okta_ai_agents_credentials_jwk_ec Resource

This guide explains how to test the EC JSON Web Key resource for Okta AI Agents.

## Prerequisites

1. **Okta Account**: With AI Agents enabled
2. **Okta API Token**: From Admin Console > API > Tokens
3. **Terraform**: Version 1.0 or higher
4. **Okta Terraform Provider**: Version 4.0 or higher
5. **EC Key Pair**: For testing (see "Generating Test Keys" below)

## Step 1: Generate Test EC Keys

### Option A: Using OpenSSL

```bash
# Generate P-256 private key
openssl ecparam -name prime256v1 -genkey -noout -out test-private-key.pem

# Extract public key
openssl ec -in test-private-key.pem -pubout -out test-public-key.pem

# View the key details
openssl ec -in test-public-key.pem -pubin -text -noout

# Extract the x and y coordinates and base64url-encode them
# Look for the "pub:" section in the output
```

### Option B: Using Python

```python
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
import base64

# Generate P-256 key
private_key = ec.generate_private_key(ec.SECP256R1(), default_backend())
public_key = private_key.public_key()

# Get the numbers
numbers = public_key.public_numbers()

# Base64url-encode without padding
x = base64.urlsafe_b64encode(numbers.x.to_bytes(32, 'big')).decode().rstrip('=')
y = base64.urlsafe_b64encode(numbers.y.to_bytes(32, 'big')).decode().rstrip('=')

print(f'x = "{x}"')
print(f'y = "{y}"')
```

### Option C: Using jwk-cli (Node.js)

```bash
npm install -g jwk-cli
jwk-cli generate -t EC -c P-256 -f PEM
```

## Step 2: Set Up Your Test Environment

### Option A: Environment Variables

```bash
# Set Okta credentials
export OKTA_ORG_NAME="your-org"
export OKTA_API_TOKEN="your-api-token"

# Set Terraform variables
export TF_VAR_okta_org_name="your-org"
export TF_VAR_okta_api_token="your-api-token"
export TF_VAR_agent_id="00u1a2b3c4d5e6f7g8h9i"  # Replace with your agent ID
export TF_VAR_key_coordinates='{"x":"YOUR_X_COORDINATE","y":"YOUR_Y_COORDINATE"}'
```

### Option B: terraform.tfvars

```bash
# Copy the example file
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars with your values
nano terraform.tfvars

# Note: Don't commit terraform.tfvars to version control!
```

## Step 3: Initialize Terraform

```bash
# Initialize Terraform in the examples directory
cd examples/resources/okta_ai_agents_credentials_jwk_ec
terraform init

# Verify initialization
terraform version
```

## Step 4: Run Terraform Plan

```bash
# Dry-run to see what will be created
terraform plan

# Save plan to file for review
terraform plan -out=tfplan
```

## Step 5: Apply Configuration

```bash
# Create the resources
terraform apply

# Or apply the saved plan
terraform apply tfplan

# Wait for completion and note the outputs
```

## Step 6: Verify in Okta

1. Go to Admin Console
2. Navigate to Applications > AI Agents
3. Click on your agent
4. Go to the "Keys" or "Credentials" section
5. Verify your new key appears in the JWKS

## Step 7: Test Import

```bash
# First, destroy the local state
terraform destroy

# Then import the key back
# Format: {agent_id}/{key_id}
terraform import okta_ai_agents_credentials_jwk_ec.example 00u1a2b3c4d5e6f7g8h9i/00c1a2b3c4d5e6f7g8h9i

# Verify the import
terraform state show okta_ai_agents_credentials_jwk_ec.example
```

## Step 8: Clean Up

```bash
# Destroy all created resources
terraform destroy

# Verify resources are deleted
# Check the Okta console to confirm
```

## Test Scenarios

### Scenario 1: Create P-256 Key

**File**: `test.tf` (resource: `test_basic`)

**Test Steps**:
```bash
terraform apply -target=okta_ai_agents_credentials_jwk_ec.test_basic
terraform state show okta_ai_agents_credentials_jwk_ec.test_basic
terraform destroy -target=okta_ai_agents_credentials_jwk_ec.test_basic
```

### Scenario 2: Create P-384 Key

**File**: `test.tf` (resource: `test_p384`)

**Test Steps**:
```bash
terraform apply -target=okta_ai_agents_credentials_jwk_ec.test_p384
terraform state show okta_ai_agents_credentials_jwk_ec.test_p384
terraform destroy -target=okta_ai_agents_credentials_jwk_ec.test_p384
```

### Scenario 3: Create P-521 Key

**File**: `test.tf` (resource: `test_p521`)

**Test Steps**:
```bash
terraform apply -target=okta_ai_agents_credentials_jwk_ec.test_p521
terraform state show okta_ai_agents_credentials_jwk_ec.test_p521
terraform destroy -target=okta_ai_agents_credentials_jwk_ec.test_p521
```

### Scenario 4: Data Source

**File**: `test.tf` (data source: `test_read`)

**Test Steps**:
```bash
# Create the key first
terraform apply -target=okta_ai_agents_credentials_jwk_ec.test_basic

# Then read it back
terraform apply -target=data.okta_ai_agents_credentials_jwk_ec.test_read

# Check the data source
terraform state show data.okta_ai_agents_credentials_jwk_ec.test_read
```

## Troubleshooting

### Error: "API returned a response for a different variant type"

**Cause**: The resource/data source is trying to parse an RSA key as EC, or vice versa.

**Solution**: Ensure you're using the correct resource for your key type:
- EC keys → use `okta_ai_agents_credentials_jwk_ec`
- RSA keys → use `okta_ai_agents_credentials_jwk_rsa`

### Error: "Invalid base64url encoding"

**Cause**: The x/y coordinates are not properly base64url-encoded.

**Solution**: Ensure coordinates are:
- Base64url-encoded (not standard Base64)
- Without padding (`=` characters removed)
- Correct length (32 bytes for P-256, 48 for P-384, 66 for P-521)

### Error: "agent_id not found"

**Cause**: The specified agent ID doesn't exist or you don't have access to it.

**Solution**:
1. Verify the agent ID from Okta console
2. Verify your API token has AI Agent permissions
3. Check your org name is correct

### Error: "Unauthorized"

**Cause**: Invalid API token or insufficient permissions.

**Solution**:
1. Regenerate your API token
2. Ensure token has "Admin" or appropriate scope
3. Check token hasn't expired

## Additional Resources

- [Okta AI Agents Documentation](https://developer.okta.com/)
- [JSON Web Key (JWK) Specification](https://tools.ietf.org/html/rfc7517)
- [ECDSA Specification](https://tools.ietf.org/html/rfc6090)
- [Terraform Okta Provider](https://registry.terraform.io/providers/okta/okta/latest/docs)
