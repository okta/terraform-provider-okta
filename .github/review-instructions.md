# Copilot Code Review Instructions

When reviewing code in this repository, please keep the following in mind:

## Terraform Provider Specific
- Ensure all schema attributes have proper `Description` fields
- Check that `d.Set()` calls handle errors appropriately
- Verify that API responses are properly closed with `defer resp.Body.Close()`
- Ensure `TypeSet` vs `TypeList` is used correctly based on ordering requirements
- Check for proper nil checks before type assertions

## Go Best Practices
- Ensure error messages are descriptive and include context
- Check for proper resource cleanup in deferred functions
- Verify that context is properly propagated through function calls

## Testing
- Ensure acceptance tests use proper resource naming
- Check that test fixtures match the test expectations
- Verify VCR cassettes are updated when API calls change

## Security
- No sensitive values should be logged
- Ensure sensitive schema fields are marked with `Sensitive: true`