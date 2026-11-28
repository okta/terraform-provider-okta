// Basic structural validation of Okta Expression Language.
//
// Basic Usage:
//
//	err := expression.ParseExpression("user.firstName == \"John\"")
//	if err != nil {
//	    log.Fatalf("Invalid expression: %v", err)
//	}
//
// Expression Syntax:
//
// The parser validates the basic structure of Okta expressions, which can include:
//   - References to attributes (e.g., user.firstName)
//   - String literals (e.g., "John")
//   - Binary operators (==, !=, >, <, >=, <=, AND, OR)
//   - Function calls (e.g., String.startsWith(user.email, "admin"))
//   - Parentheses for grouping
//
// Examples of valid expressions:
//   - user.firstName == "John"
//   - String.startsWith(user.email, "admin")
//   - (user.department == "Sales" AND user.location == "SF")
//
// Validation Scope:
//
// The current implementation performs structural validation only, including:
//   - Basic syntax checking
//   - Expression structure validation
//   - Reference path validation (non-empty parts)
//   - Function call structure validation
//
// Note that semantic validation is not supported for maintainability reasons.
// This means that we don't valid things like:
//   - if a function exists
//   - if an attribute reference is referencing a valid attribute
//   - if the types of two values in a comparison are valid
//
// Error Handling:
//
// The package provides error messages for invalid expressions, including:
//   - Syntax errors
//   - Empty expressions
//   - Invalid term structure
//   - Empty reference parts
//   - Basic operator validation
//
// For more information about Okta Expression Language, see:
// https://developer.okta.com/docs/reference/okta-expression-language/
package expression
