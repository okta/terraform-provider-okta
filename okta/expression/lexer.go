package expression

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// Token types for structural validation
const (
	TokenWhitespace = "whitespace"
	TokenString     = "String"
	TokenIdentifier = "Identifier"
	TokenOperator   = "Operator"
	TokenPunct      = "Punct"
	TokenBoolean    = "Boolean"
	TokenNumber     = "Number"
)

// DefaultLexer is the lexer for Okta expressions
var DefaultLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: TokenWhitespace, Pattern: `\s+`},
	{Name: TokenString, Pattern: `"(?:\\.|[^"])*"`},
	// Okta operators include both symbolic and textual forms
	{Name: TokenOperator, Pattern: `==|!=|>=|<=|>|<|AND|OR|NOT|eq|ne|gt|ge|lt|le|sw|co|pr|and|or|not`},
	{Name: TokenBoolean, Pattern: `true|false`},
	{Name: TokenNumber, Pattern: `[0-9]+(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?`}, // Scientific notation supported
	{Name: TokenPunct, Pattern: `[.,(){}[\]]`},
	// Identifiers must start with a letter and can only contain alphanumeric chars and underscores
	{Name: TokenIdentifier, Pattern: `[a-zA-Z][a-zA-Z0-9_]*`},
})
