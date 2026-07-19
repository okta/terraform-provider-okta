# okta_ai_agents_credentials_jwk_ec - Example Files Guide

This directory contains complete examples and documentation for the `okta_ai_agents_credentials_jwk_ec` Terraform resource.

## Files Overview

### Documentation

| File | Purpose |
|------|---------|
| **README.md** | Complete resource documentation including arguments, attributes, and usage examples |
| **TESTING.md** | Step-by-step guide for testing the resource, including key generation and troubleshooting |
| **INDEX.md** | This file - guide to all example files |

### Terraform Configuration Examples

| File | Purpose | Use Case |
|------|---------|----------|
| **resource.tf** | Basic resource definition with inline comments | Quick reference for required and optional fields |
| **real_example.tf** | Production-ready example with real EC key coordinates | Real-world implementation with multiple curve types (P-256, P-384, P-521) |
| **test.tf** | Test configuration with multiple scenarios | Comprehensive testing of different EC curves and algorithms |
| **main.tf** | Complete setup with provider configuration | Full Terraform project setup with all necessary providers and variables |

### Configuration Files

| File | Purpose |
|------|---------|
| **terraform.tfvars.example** | Example variables file (rename to terraform.tfvars for use) |

## Quick Start

### 1. Basic Usage (5 minutes)

```bash
# View basic resource definition
cat resource.tf

# Review documentation
cat README.md | head -50
```

### 2. Run a Test (15 minutes)

```bash
# Generate test keys
python3 << 'EOF'
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
import base64

pk = ec.generate_private_key(ec.SECP256R1(), default_backend())
pub = pk.public_key()
nums = pub.public_numbers()
x = base64.urlsafe_b64encode(nums.x.to_bytes(32, 'big')).decode().rstrip('=')
y = base64.urlsafe_b64encode(nums.y.to_bytes(32, 'big')).decode().rstrip('=')
print(f"x={x}\ny={y}")
EOF

# Set up environment
export TF_VAR_okta_org_name="your-org"
export TF_VAR_okta_api_token="your-token"
export TF_VAR_agent_id="your-agent-id"
export TF_VAR_key_coordinates='{"x":"YOUR_X","y":"YOUR_Y"}'

# Initialize and apply
cd . && terraform init && terraform plan
```

### 3. Production Deployment (30 minutes)

```bash
# Copy real example
cp real_example.tf production.tf

# Edit with your values
nano production.tf

# Apply carefully
terraform plan -out=prod.tfplan
terraform apply prod.tfplan
```

## File Details

### resource.tf

**Size**: ~1.9 KB  
**Content**: 
- Single resource with inline comments
- All fields documented
- Example with multiple keys
- Output examples
- Data source reference

**Best for**: Quick reference and understanding all available fields

### real_example.tf

**Size**: ~2.8 KB  
**Content**:
- Production-quality resource definitions
- Real base64url-encoded coordinates
- Multiple curve examples (P-256, P-384, P-521)
- Key generation instructions for all major languages/tools
- Import examples
- Comments explaining coordinate generation

**Best for**: Production deployments and understanding real-world usage

### test.tf

**Size**: ~2.1 KB  
**Content**:
- Three test resources for different curves
- Data source test
- Output examples
- Verification examples

**Best for**: Comprehensive testing and CI/CD pipelines

### main.tf

**Size**: ~3.2 KB  
**Content**:
- Complete Terraform project setup
- Provider configuration
- All required variables
- Comprehensive outputs
- Usage instructions embedded
- Environment variable guidance

**Best for**: Starting a new Terraform project from scratch

### README.md

**Size**: ~3.9 KB  
**Content**:
- Full resource documentation
- Argument and attribute reference
- Import instructions
- Key generation examples (OpenSSL, Python)
- Supported curves and algorithms table
- Important notes and limitations

**Best for**: Understanding resource capabilities and requirements

### TESTING.md

**Size**: ~6.3 KB  
**Content**:
- Prerequisites checklist
- Detailed key generation steps
- Environment setup options
- Step-by-step testing guide
- Multiple test scenarios
- Troubleshooting guide
- Additional resources

**Best for**: Testing and troubleshooting

### terraform.tfvars.example

**Size**: ~1.6 KB  
**Content**:
- Variable definitions
- Example values
- Usage instructions
- Security notes
- Comments for each variable

**Best for**: Creating your terraform.tfvars file

## Common Workflows

### Scenario 1: I want to understand what this resource does

```bash
1. Read: README.md (Argument Reference section)
2. Review: resource.tf (see field examples)
3. Test: Follow TESTING.md (Step 1-3)
```

### Scenario 2: I need to deploy this to production

```bash
1. Read: README.md (all sections)
2. Copy: real_example.tf to production.tf
3. Edit: production.tf with your values
4. Setup: terraform.tfvars using terraform.tfvars.example
5. Deploy: terraform plan && terraform apply
6. Verify: Check Okta console for the key
```

### Scenario 3: I need to test multiple scenarios

```bash
1. Follow: TESTING.md (Steps 1-7)
2. Use: test.tf for comprehensive testing
3. Review: All three scenarios in TESTING.md
```

### Scenario 4: I'm having issues

```bash
1. Troubleshoot: TESTING.md (Troubleshooting section)
2. Verify: All prerequisites are met
3. Check: Your key coordinates are valid
4. Test: With TESTING.md key generation steps
```

## Key Generation Quick Reference

### Python (Recommended)

```bash
pip install cryptography

python3 << 'EOF'
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.backends import default_backend
import base64

# P-256
pk = ec.generate_private_key(ec.SECP256R1(), default_backend())
pub = pk.public_key()
nums = pub.public_numbers()
x = base64.urlsafe_b64encode(nums.x.to_bytes(32, 'big')).decode().rstrip('=')
y = base64.urlsafe_b64encode(nums.y.to_bytes(32, 'big')).decode().rstrip('=')
print(f"P-256:\nx={x}\ny={y}")

# P-384
pk384 = ec.generate_private_key(ec.SECP384R1(), default_backend())
pub384 = pk384.public_key()
nums384 = pub384.public_numbers()
x384 = base64.urlsafe_b64encode(nums384.x.to_bytes(48, 'big')).decode().rstrip('=')
y384 = base64.urlsafe_b64encode(nums384.y.to_bytes(48, 'big')).decode().rstrip('=')
print(f"\nP-384:\nx={x384}\ny={y384}")

# P-521
pk521 = ec.generate_private_key(ec.SECP521R1(), default_backend())
pub521 = pk521.public_key()
nums521 = pub521.public_numbers()
x521 = base64.urlsafe_b64encode(nums521.x.to_bytes(66, 'big')).decode().rstrip('=')
y521 = base64.urlsafe_b64encode(nums521.y.to_bytes(66, 'big')).decode().rstrip('=')
print(f"\nP-521:\nx={x521}\ny={y521}")
EOF
```

### OpenSSL

```bash
# Generate and view
openssl ecparam -name prime256v1 -genkey -out key.pem
openssl ec -in key.pem -pubout -text -noout

# Extract coordinates (manual base64url encoding)
# See TESTING.md for detailed instructions
```

## Next Steps

1. **Read the documentation**: Start with README.md
2. **Generate your keys**: Use Python script above
3. **Choose your deployment method**: 
   - Quick test? Use `test.tf`
   - Production? Use `real_example.tf` or `main.tf`
4. **Follow the testing guide**: TESTING.md has complete instructions
5. **Deploy**: `terraform apply`
6. **Verify**: Check Okta console for the new key

## Support

For issues or questions:

1. Check **TESTING.md** Troubleshooting section
2. Review **README.md** for limitations and notes
3. Verify your key generation using the quick reference above
4. Check the [Okta AI Agents documentation](https://developer.okta.com/)
5. Review [Terraform Okta Provider documentation](https://registry.terraform.io/providers/okta/okta/latest/docs)

---

**Last Updated**: 2024-07-13  
**Resource**: okta_ai_agents_credentials_jwk_ec  
**Terraform Provider**: Okta v4.0+
