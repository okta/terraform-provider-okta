# TODO

Use the TODO file to capture greater implementation and architecture concepts
that should be considered in the provider.

## `403 Forbidden` / `401 Unauthorized` errors (November 10, 2022)

The Okta API can return `403 Forbidden` and `401 Unauthorized` errors on some
endpoints. This typically isn't noticeable if the API token was created by a
Super Admin on the org. If the token is made with a lesser role, for instance
Org Admin, or the org doesn't have certain feature enabled, then it interferes
with consistent behavior in the provider.  Presently, Okta API doesn't document
which endpoints have permission levels and there isn't any kind of endpoint
that presents that information in a discoverable way.

Possible solutions.

### New config variable `OTKA_API_TOKEN_ROLE=[super-admin|org-admin|etc]`

Allow the operator to manually set a provider configuration variable
corresponding to their token's role. The variable could be named somethign like
Okta API Token Role. Resource code paths can then gate on that variable's value
and act accordingly for known endpoints having permission levels.

### golang http roundtripper intercepting 401/403 errors

Add in a golang roundtripper https://pkg.go.dev/net/http#RoundTripper that
intercepts 401/403 responses and returns a soft error.

### Simple guards in code

Proactively guard known endpoint errors in code like is done with
suppressErrorOn404(resp, err) 

### Additional arguments on resources

Define arguments on resources like `skip_roles` in data source okta_user
https://registry.terraform.io/providers/okta/okta/latest/docs/data-sources/user
that allows the operator to explicitly guard behaviors known to have permission
issues.

