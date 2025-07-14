# VCR Domain Rewriting Fixes

## Problem Summary

The VCR (Video Cassette Recorder) implementation for acceptance tests was having issues with domain rewriting when recording and playing back cassettes. The main problems were:

1. **Inconsistent domain replacement**: Domain rewriting was happening at different times and in different ways during recording vs playback
2. **Incomplete domain patterns**: Not all domain references were being properly replaced
3. **Fixture variable handling**: Terraform variables in test fixtures weren't being consistently replaced
4. **URL matching issues**: The VCR matcher wasn't properly handling domain differences between recording and playback

## Changes Made

### 1. Improved `rewriteDomainsInCassette()` function (`acctest.go`)

- **Better pattern matching**: Now replaces domains in order of specificity (most specific first)
- **More comprehensive replacement**: Handles admin hostnames, regular hostnames, org names, and base URLs
- **Consistent ordering**: Ensures replacements happen in the correct order to avoid conflicts

### 2. Enhanced `vcrHook()` function (`acctest.go`)

- **Comprehensive domain replacement**: Now handles all parts of HTTP interactions (Host, URL, Body, Response Body, Headers)
- **Better organization**: Replaces domains in order of specificity
- **Improved Link header handling**: More robust replacement of domain references in Link headers
- **Consistent logic**: Same replacement logic for both request and response bodies

### 3. Improved `ConfigReplace()` method (`helpers_for_test.go`)

- **Enhanced variable replacement**: Better handling of Terraform variable references
- **Hardcoded domain replacement**: Also replaces any hardcoded domain references in fixtures
- **Fallback logic**: Multiple fallback sources for domain information

### 4. Enhanced `vcrMatcher()` function (`acctest.go`)

- **Domain normalization**: Normalizes URLs for comparison by replacing domain differences
- **VCR-aware matching**: Only applies domain normalization when in VCR mode
- **Comprehensive replacement**: Handles admin hostnames, regular hostnames, and base URLs

## Testing the Fixes

### Manual Testing

1. **Test live (no VCR)**:

   ```bash
   export TF_LOG=DEBUG
   export TF_LOG_PATH="../tflog.log"
   unset OKTA_VCR_CASSETTE
   unset OKTA_VCR_TF_ACC
   rm test-output.log
   rm tflog.log

   op run --account="NDM32WGBUND3XEEL5OEN22MMOE" --env-file="./env_trial-5309542.okta.com.env" -- \
     make testacc TEST=./okta TESTARGS='-run=TestAccDataSourceOktaResourceSets_read' 2>&1 | tee test-output.log
   ```

2. **Record VCR**:

   ```bash
   export TF_LOG=DEBUG
   export TF_LOG_PATH="../tflog.log"
   export OKTA_VCR_CASSETTE=oie-00
   export OKTA_VCR_TF_ACC=record
   rm -f test-output.log
   rm -f tflog.log

   op run --account="NDM32WGBUND3XEEL5OEN22MMOE" --env-file="./env_trial-5309542.okta.com.env" -- \
     make testacc TEST=./okta TESTARGS='-run=TestAccDataSourceOktaResourceSets_read' 2>&1 | tee test-output.log
   ```

3. **Play VCR**:

   ```bash
   export TF_LOG=DEBUG
   export TF_LOG_PATH="../tflog.log"
   export OKTA_VCR_CASSETTE=oie-00
   export OKTA_VCR_TF_ACC=play
   rm -f test-output.log
   rm -f tflog.log

   op run --account="NDM32WGBUND3XEEL5OEN22MMOE" --env-file="./env_trial-5309542.okta.com.env" -- \
     make testacc TEST=./okta TESTARGS='-run=TestAccDataSourceOktaResourceSets_read' 2>&1 | tee test-output.log
   ```

### Automated Testing

Use the provided test script:

```bash
./test-vcr.sh
```

This script will:

1. Run each test in live mode (no VCR)
2. Record VCR cassettes for each test
3. Play back the recorded cassettes
4. Report success/failure for each step

## Test Cases Covered

The fixes should resolve issues with these test prefixes:

- `TestAccDataSourceOktaResourceSets_`
- `TestAccDataSourceOktaResourceSet_`
- `TestAccDataSourceOktaResourceSetResources_`

## Key Improvements

1. **Consistent domain handling**: All domain references are now handled consistently across recording and playback
2. **Better error handling**: More robust error handling and logging
3. **Comprehensive replacement**: Handles edge cases and different domain formats
4. **Improved debugging**: Better logging and error messages for troubleshooting

## Expected Results

After these fixes:

- Tests should pass in live mode (no VCR)
- VCR recording should complete successfully
- VCR playback should work without domain-related errors
- Cassettes should contain properly rewritten domains
- GitHub CI should be able to run tests using the recorded cassettes

## Troubleshooting

If issues persist:

1. **Check environment variables**: Ensure `OKTA_ORG_NAME` and `OKTA_BASE_URL` are set correctly
2. **Verify cassette files**: Check that cassette files are created in `test/fixtures/vcr/idaas/`
3. **Review logs**: Look at `test-output.log` and `tflog.log` for detailed error messages
4. **Clean cassettes**: Remove existing cassettes and re-record if needed
5. **Check domain patterns**: Verify that the domain patterns in your environment match the expected format
