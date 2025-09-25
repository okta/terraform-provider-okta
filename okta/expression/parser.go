package expression

import (
	"strings"

	"github.com/alecthomas/participle/v2"
)

var parser = participle.MustBuild[Expression](
	participle.Lexer(DefaultLexer),
	participle.Unquote("String"),
	participle.Elide("whitespace"),
	participle.UseLookahead(2),
)

// parenInfo tracks information about parentheses for error reporting
type parenInfo struct {
	openPos  int // Position of opening parenthesis
	closePos int // Position of matching closing parenthesis (-1 if not found)
}

// findProblematicParenthesis finds the most likely problematic opening parenthesis
func findProblematicParenthesis(expr string) int {
	stack := []parenInfo{}
	inString := false

	// First pass: collect information about parentheses and their nesting
	for i, c := range expr {
		if c == '"' && (i == 0 || expr[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}

		switch c {
		case '(':
			// For each opening parenthesis, track its position
			stack = append(stack, parenInfo{openPos: i, closePos: -1})
		case ')':
			if len(stack) > 0 {
				// Mark the innermost unclosed parenthesis as closed
				for j := len(stack) - 1; j >= 0; j-- {
					if stack[j].closePos == -1 {
						stack[j].closePos = i
						break
					}
				}
			}
		}
	}

	if len(stack) == 0 {
		return -1
	}

	// Find the first unclosed parenthesis (scanning from left to right)
	for _, info := range stack {
		if info.closePos == -1 {
			return info.openPos
		}
	}

	return -1
}

// ParseExpression performs structural validation of an Okta expression
func ParseExpression(expr string) error {
	if strings.TrimSpace(expr) == "" {
		return &SyntaxError{Message: ErrEmptyExpression}
	}

	// Check for unbalanced parentheses first
	if pos := findProblematicParenthesis(expr); pos >= 0 {
		return &SyntaxError{
			Message:      ErrUnbalancedParens,
			Position:     pos,
			OpenParenPos: pos,
			Context:      getErrorContext(expr, pos),
		}
	}

	// Parse and validate structure
	expression, err := parser.ParseString("", expr)
	if err != nil {
		// Convert participle errors to our error type
		if pe, ok := err.(participle.Error); ok {
			return &SyntaxError{
				Message:  pe.Error(),
				Position: pe.Position().Offset,
				Context:  getErrorContext(expr, pe.Position().Offset),
			}
		}
		return err
	}

	return expression.Validate()
}

// getErrorContext returns a snippet of the expression around the error position
func getErrorContext(expr string, pos int) string {
	if pos < 0 {
		return ""
	}

	start := pos - 10
	if start < 0 {
		start = 0
	}

	end := pos + 10
	if end > len(expr) {
		end = len(expr)
	}

	return expr[start:end]
}
