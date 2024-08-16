package main

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token
	}{
		{
			name:  "Basic arithmetic",
			input: "(add 10 (subtract 10 6))",
			expected: []token{
				{kind: "paren", value: "("},
				{kind: "name", value: "add"},
				{kind: "number", value: "10"},
				{kind: "paren", value: "("},
				{kind: "name", value: "subtract"},
				{kind: "number", value: "10"},
				{kind: "number", value: "6"},
				{kind: "paren", value: ")"},
				{kind: "paren", value: ")"},
			},
		},
		{
			name:  "String and boolean",
			input: `(print "Hello" true)`,
			expected: []token{
				{kind: "paren", value: "("},
				{kind: "name", value: "print"},
				{kind: "string", value: "Hello"},
				{kind: "boolean", value: "true"},
				{kind: "paren", value: ")"},
			},
		},
		{
			name:  "Operators",
			input: "(+ 5 (* 3 2))",
			expected: []token{
				{kind: "paren", value: "("},
				{kind: "operator", value: "+"},
				{kind: "number", value: "5"},
				{kind: "paren", value: "("},
				{kind: "operator", value: "*"},
				{kind: "number", value: "3"},
				{kind: "number", value: "2"},
				{kind: "paren", value: ")"},
				{kind: "paren", value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tokenizer(tt.input)
			if err != nil {
				t.Fatalf("Tokenizer error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected: %v, Got: %v", tt.expected, result)
			}
		})
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    []token
		expected ast
	}{
		{
			name: "Basic arithmetic",
			input: []token{
				{kind: "paren", value: "("},
				{kind: "name", value: "add"},
				{kind: "number", value: "10"},
				{kind: "paren", value: "("},
				{kind: "name", value: "subtract"},
				{kind: "number", value: "10"},
				{kind: "number", value: "6"},
				{kind: "paren", value: ")"},
				{kind: "paren", value: ")"},
			},
			expected: ast{
				kind: "Program",
				body: []node{
					{
						kind: "CallExpression",
						name: "add",
						params: []node{
							{kind: "NumberLiteral", value: "10"},
							{
								kind: "CallExpression",
								name: "subtract",
								params: []node{
									{kind: "NumberLiteral", value: "10"},
									{kind: "NumberLiteral", value: "6"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "String and boolean",
			input: []token{
				{kind: "paren", value: "("},
				{kind: "name", value: "print"},
				{kind: "string", value: "Hello"},
				{kind: "boolean", value: "true"},
				{kind: "paren", value: ")"},
			},
			expected: ast{
				kind: "Program",
				body: []node{
					{
						kind: "CallExpression",
						name: "print",
						params: []node{
							{kind: "StringLiteral", value: "Hello"},
							{kind: "BooleanLiteral", value: "true"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser(tt.input)
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected: %+v, Got: %+v", tt.expected, result)
			}
		})
	}
}

func TestCompiler(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic arithmetic",
			input:    "(add 2 (subtract 4 2))",
			expected: "add(2, subtract(4, 2));",
		},
		{
			name:     "String and boolean",
			input:    `(print "Hello" true)`,
			expected: `print("Hello", true);`,
		},
		{
			name:     "Operators",
			input:    "(+ 5 (* 3 2))",
			expected: "+(5, *(3, 2));",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compiler(tt.input)
			if err != nil {
				t.Fatalf("Compiler error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

func TestTokenizerErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "Invalid character", input: "(add 2 @)"},
		{name: "Unclosed string", input: `(print "Hello)`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tokenizer(tt.input)
			if err == nil {
				t.Errorf("Expected error for input: %s, but got none", tt.input)
			}
		})
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		name  string
		input []token
	}{
		{
			name: "Unmatched parentheses",
			input: []token{
				{kind: "paren", value: "("},
				{kind: "name", value: "add"},
				{kind: "number", value: "2"},
				{kind: "number", value: "3"},
			},
		},
		{
			name: "Invalid token sequence",
			input: []token{
				{kind: "paren", value: "("},
				{kind: "number", value: "2"},
				{kind: "name", value: "add"},
				{kind: "paren", value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser(tt.input)
			if err == nil {
				t.Errorf("Expected error for input: %v, but got none", tt.input)
			}
		})
	}
}
