# Configuration Examples

Here lies the examples that will aid you on your Okta Terraform journey 🚵🏄. In the current directory you will find an example stack using this provider!

## Test Fixture Examples

Anything that lies underneath a resource directory is config we use as fixtures to our tests. This means we pass our tests through a string formatter that replaces the format placeholders with the actual values but regardless seeing these might be helpful when trying to implement this provider. Just remember, it is not valid config until you replace or remove `%[1]d` :blush:.
