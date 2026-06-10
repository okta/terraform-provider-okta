package expression

import (
	"fmt"
	"strings"
)

// SyntaxError represents a structural error in the expression
type SyntaxError struct {
	Message      string
	Token        string // The problematic token, if any
	Position     int    // Position in the input where the error occurred
	Context      string // Surrounding context of the error
	OpenParenPos int    // Position of opening parenthesis for unmatched paren errors
}

func (e *SyntaxError) Error() string {
	var parts []string

	// For unmatched parenthesis errors, use the opening parenthesis position
	if strings.Contains(e.Message, "unexpected token \"<EOF>\" (expected \")\"") && e.OpenParenPos >= 0 {
		parts = append(parts, fmt.Sprintf("at position %d", e.OpenParenPos))
		e.Message = "unclosed parenthesis"
	} else if e.Position >= 0 {
		parts = append(parts, fmt.Sprintf("at position %d", e.Position))
	}

	if e.Token != "" {
		parts = append(parts, fmt.Sprintf("near '%s'", e.Token))
	}

	if e.Context != "" {
		parts = append(parts, fmt.Sprintf("in '%s'", e.Context))
	}

	if len(parts) > 0 {
		return fmt.Sprintf("syntax error %s: %s", strings.Join(parts, " "), e.Message)
	}

	return fmt.Sprintf("syntax error: %s", e.Message)
}

// Common structural error messages
const (
	ErrEmptyExpression  = "empty expression"
	ErrUnbalancedParens = "unclosed parenthesis"
	ErrUnclosedString   = "unclosed string literal"
	ErrInvalidStructure = "invalid expression structure"
	ErrMissingOperand   = "missing operand"
	ErrUnexpectedToken  = "unexpected token"
)
