package expression

import (
	"strings"
	"testing"
)

func TestBasicExpressions(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{
			name:    "simple reference equality",
			expr:    "user.firstName == \"John\"",
			wantErr: false,
		},
		{
			name:    "simple comparison",
			expr:    "user.age > 18",
			wantErr: false,
		},
		{
			name:    "decimal number comparison",
			expr:    "user.score >= 95.5",
			wantErr: false,
		},
		{
			name:    "scientific notation",
			expr:    "user.score == 1.5e-10",
			wantErr: false,
		},
		{
			name:    "boolean comparison",
			expr:    "user.active == true",
			wantErr: false,
		},
		{
			name:    "complex boolean expression",
			expr:    "(user.age > 18) AND (user.active == true)",
			wantErr: false,
		},
		{
			name:    "custom integer attribute equals",
			expr:    "user.integer_attribute == 1",
			wantErr: false,
		},
		{
			name:    "custom integer attribute eq operator",
			expr:    "user.integer_attribute eq 1",
			wantErr: false,
		},
		{
			name:    "custom string attribute contains",
			expr:    "user.custom_string_attr co \"test\"",
			wantErr: false,
		},
		{
			name:    "custom string attribute starts with",
			expr:    "user.custom_string_attr sw \"admin\"",
			wantErr: false,
		},
		{
			name:    "attribute with numbers",
			expr:    "user.custom2_attr == \"value\"",
			wantErr: false,
		},
		{
			name:    "empty expression",
			expr:    "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			expr:    "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				t.Logf("Error: %v", err)
			}
		})
	}
}

func TestInvalidIdentifiers(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{
			name:    "starts with number",
			expr:    "user.1attribute == \"value\"",
			wantErr: "unexpected token \"1\"",
		},
		{
			name:    "contains special chars",
			expr:    "user.my-attribute == \"value\"",
			wantErr: "invalid input text",
		},
		{
			name:    "consecutive dots",
			expr:    "user..attribute == \"value\"",
			wantErr: "unexpected token \".\"",
		},
		{
			name:    "ends with dot",
			expr:    "user.attribute. == \"value\"",
			wantErr: "unexpected token",
		},
		{
			name:    "starts with dot",
			expr:    ".user.attribute == \"value\"",
			wantErr: "unexpected token \".\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseExpression(tt.expr)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got %q", tt.wantErr, err.Error())
			}
			t.Logf("Error: %v", err)
		})
	}
}

func TestParenthesisErrors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantPos int
	}{
		{
			name:    "simple unclosed",
			expr:    "(user.firstName == \"John\"",
			wantPos: 0,
		},
		{
			name:    "nested unclosed",
			expr:    "(user.firstName == (user.lastName == \"Doe\"",
			wantPos: 0,
		},
		{
			name:    "multiple with last unclosed",
			expr:    "(user.age > 18) AND (user.active == true",
			wantPos: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseExpression(tt.expr)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			syntaxErr, ok := err.(*SyntaxError)
			if !ok {
				t.Fatalf("Expected *SyntaxError but got %T", err)
			}

			if syntaxErr.OpenParenPos != tt.wantPos {
				t.Errorf("Wrong error position\nGot:  %d (%s)\nWant: %d (%s)",
					syntaxErr.OpenParenPos,
					getPositionContext(tt.expr, syntaxErr.OpenParenPos),
					tt.wantPos,
					getPositionContext(tt.expr, tt.wantPos))
			}
			t.Logf("Error: %v", err)
		})
	}
}

func TestComplexParenthesisErrors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantPos int
	}{
		{
			name:    "nested with missing middle close",
			expr:    "(a == (b == (c == \"x\" AND (d == \"y\")))",
			wantPos: 0, // First opening paren is unclosed since we're missing one )
		},
		{
			name:    "deeply nested missing inner",
			expr:    "(a == (b == (c == \"x\") AND (d == \"y\")",
			wantPos: 0, // First opening paren is unclosed
		},
		{
			name:    "multiple balanced groups with unbalanced middle",
			expr:    "(a == \"x\") AND (b == (c == \"y\") OR (d == \"z\")",
			wantPos: 15, // The opening paren after AND is unclosed
		},
		{
			name:    "multiple missing closes",
			expr:    "(a == (b == (c == \"x\"",
			wantPos: 0, // First opening paren is unclosed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseExpression(tt.expr)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			syntaxErr, ok := err.(*SyntaxError)
			if !ok {
				t.Fatalf("Expected *SyntaxError but got %T", err)
			}

			if syntaxErr.OpenParenPos != tt.wantPos {
				t.Errorf("Wrong error position\nGot:  %d (%s)\nWant: %d (%s)",
					syntaxErr.OpenParenPos,
					getPositionContext(tt.expr, syntaxErr.OpenParenPos),
					tt.wantPos,
					getPositionContext(tt.expr, tt.wantPos))
			}
			t.Logf("Error: %v", err)
		})
	}
}

func TestLexerErrors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{
			name:    "invalid operator",
			expr:    "user.age && 18",
			wantErr: "invalid input text",
		},
		{
			name:    "unclosed string",
			expr:    "user.name == \"John",
			wantErr: "invalid input text", // The lexer reports this as invalid input
		},
		{
			name:    "invalid character",
			expr:    "user.name == John$",
			wantErr: "invalid input text",
		},
		{
			name:    "invalid number format",
			expr:    "user.age > 18abc",
			wantErr: "unexpected token", // The lexer correctly identifies 18 and then fails on abc
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseExpression(tt.expr)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got %q", tt.wantErr, err.Error())
			}
			t.Logf("Error: %v", err)
		})
	}
}

// Helper function to show context around a position for test output
func getPositionContext(expr string, pos int) string {
	if pos < 0 || pos >= len(expr) {
		return "invalid position"
	}
	start := pos - 5
	if start < 0 {
		start = 0
	}
	end := pos + 5
	if end > len(expr) {
		end = len(expr)
	}
	return expr[start:end]
}
