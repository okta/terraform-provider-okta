
# Data source without corresponding resource
# Reads the org singleton (GET /api/v1/org); no inputs or parent resources required.
data "okta_org" "test" {
}
