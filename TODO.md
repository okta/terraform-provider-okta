# TODO

Use the TODO file to capture greater implementation and architecture concepts
that should be considered in the provider.

## (November 10, 2022) `403 Forbidden` / `401 Unauthorized` errors

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

## (February 18, 2025) Test Package Names

https://stackoverflow.com/questions/19998250/proper-package-naming-for-testing-with-the-go-language/31443271#31443271

Test Code Package Comparison
 - Black-box Testing: Use package myfunc_test, which will ensure you're only using the exported identifiers.
 - White-box Testing: Use package myfunc so that you have access to the non-exported identifiers. Good for unit tests that require access to non-exported variables, functions, and methods.

## (February 18, 2025) Address all FIXMEs

```
find . -name \*.go -exec grep -q "FIXME" {} \; -print
```

## (February 18, 2025) Address all TODOs

```
find . -name \*.go -exec grep -q "TODO" {} \; -print
```

## (February 18, 2025) public/private functions

Double check all function declarations if they should be public/private

## (February 20, 2025) examples

Move all examples/* to IDaaS dir structure exmaples/idaas/* and update all
references to that directory (docs and code).
