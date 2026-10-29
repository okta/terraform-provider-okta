package expression

// Expression represents a complete expression that can be validated
type Expression struct {
	Left  *Term  `parser:"@@"`          // Left term
	Op    string `parser:"( @Operator"` // Optional operator
	Right *Term  `parser:"  @@ )?"`     // Optional right term (required if operator present)
}

// Term represents a basic unit in an expression
type Term struct {
	StringLiteral *string     `parser:"  @String"`    // String literal
	NumberLiteral *string     `parser:"| @Number"`    // Numeric literal
	BoolLiteral   *string     `parser:"| @Boolean"`   // Boolean literal
	Reference     *Reference  `parser:"| @@"`         // Dotted reference
	SubExpr       *Expression `parser:"| '(' @@ ')'"` // Nested expression
	FunctionCall  *FuncCall   `parser:"| @@"`         // Function call syntax
}

// Reference represents a dotted path (e.g., user.firstName)
// Each part must be a valid identifier according to Okta's rules
type Reference struct {
	Parts []string `parser:"@Identifier ('.' @Identifier)*"`
}

// FuncCall represents any function call syntax, without semantic validation
type FuncCall struct {
	Namespace string  `parser:"(@Identifier '.')?"`      // Optional namespace
	Name      string  `parser:"@Identifier"`             // Function name
	Arguments []*Term `parser:"'(' (@@ (',' @@)*)? ')'"` // Arguments
}

// Validate performs structural validation of an expression tree
func (e *Expression) Validate() error {
	if e.Left == nil {
		return &SyntaxError{Message: "empty expression"}
	}

	// Validate left term
	if err := e.Left.Validate(); err != nil {
		return err
	}

	// If we have an operator, we must have a right term
	if e.Op != "" {
		if e.Right == nil {
			return &SyntaxError{
				Message: "operator requires right term",
				Token:   e.Op,
			}
		}
		if e.Op == "NOT" {
			return &SyntaxError{
				Message: "NOT operator cannot be used as binary operator",
				Token:   e.Op,
			}
		}
		if err := e.Right.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate performs structural validation of a term
func (t *Term) Validate() error {
	// Only one of the fields should be non-nil
	count := 0
	if t.StringLiteral != nil {
		count++
	}
	if t.NumberLiteral != nil {
		count++
	}
	if t.BoolLiteral != nil {
		count++
	}
	if t.Reference != nil {
		count++
	}
	if t.SubExpr != nil {
		count++
	}
	if t.FunctionCall != nil {
		count++
	}

	if count != 1 {
		return &SyntaxError{Message: "invalid term structure"}
	}

	// Validate the non-nil field
	switch {
	case t.Reference != nil:
		return t.Reference.Validate()
	case t.SubExpr != nil:
		return t.SubExpr.Validate()
	case t.FunctionCall != nil:
		return t.FunctionCall.Validate()
	}

	return nil
}

// Validate checks reference structure
func (r *Reference) Validate() error {
	if len(r.Parts) == 0 {
		return &SyntaxError{Message: "empty reference"}
	}

	// MVP only validates that parts are non-empty
	for _, part := range r.Parts {
		if part == "" {
			return &SyntaxError{
				Message: "empty reference part",
				Token:   part,
			}
		}
	}

	return nil
}

// Validate checks function call structure
func (f *FuncCall) Validate() error {
	if f.Name == "" {
		return &SyntaxError{Message: "empty function name"}
	}

	// Validate each argument
	for _, arg := range f.Arguments {
		if err := arg.Validate(); err != nil {
			return err
		}
	}

	return nil
}
