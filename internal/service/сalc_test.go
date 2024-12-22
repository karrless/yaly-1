package service

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"yaly-1/pkg/errs"
)

func TestGetPriority(t *testing.T) {
	tests := []struct {
		name     string
		input    *token
		expected []int
	}{
		{
			name:     "Division operator",
			input:    &token{value: "/"},
			expected: []int{1, 2},
		},
		{
			name:     "Multiplication operator",
			input:    &token{value: "*"},
			expected: []int{1, 2},
		},
		{
			name:     "Addition operator",
			input:    &token{value: "+"},
			expected: []int{1, 1},
		},
		{
			name:     "Subtraction operator",
			input:    &token{value: "-"},
			expected: []int{1, 1},
		},
		{
			name:     "Unknown operator",
			input:    &token{value: "^"},
			expected: []int{0, 0},
		},
		{
			name:     "Empty value",
			input:    &token{value: ""},
			expected: []int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPriority(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getPriority(%v) = %v, expected %v", tt.input.value, result, tt.expected)
			}
		})
	}
}

func TestExecute(t *testing.T) {

	tests := []struct {
		name        string
		op          *token
		lnum, rnum  float64
		expected    float64
		expectError error
	}{
		{
			name:        "Addition",
			op:          &token{value: "+"},
			lnum:        5,
			rnum:        10,
			expected:    15,
			expectError: nil,
		},
		{
			name:        "Subtraction",
			op:          &token{value: "-"},
			lnum:        5,
			rnum:        10,
			expected:    5,
			expectError: nil,
		},
		{
			name:        "Multiplication",
			op:          &token{value: "*"},
			lnum:        5,
			rnum:        10,
			expected:    50,
			expectError: nil,
		},
		{
			name:        "Division",
			op:          &token{value: "/"},
			lnum:        5,
			rnum:        10,
			expected:    2,
			expectError: nil,
		},
		{
			name:        "Division by zero",
			op:          &token{value: "/"},
			lnum:        0,
			rnum:        10,
			expected:    0,
			expectError: errs.ErrDivisionByZero,
		},
		{
			name:        "Unknown operator",
			op:          &token{value: "^"},
			lnum:        5,
			rnum:        10,
			expected:    0,
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := execute(tt.op, tt.lnum, tt.rnum)

			if result != tt.expected {
				t.Errorf("execute(%v, %v, %v) = %v, expected %v", tt.op.value, tt.lnum, tt.rnum, result, tt.expected)
			}

			if !errors.Is(err, tt.expectError) {
				t.Errorf("execute(%v, %v, %v) error = %v, expected %v", tt.op.value, tt.lnum, tt.rnum, err, tt.expectError)
			}
		})
	}
}

func TestGetNextToken(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedToken  string
		expectedRemain string
	}{
		{
			name:           "Single character string",
			input:          "a",
			expectedToken:  "a",
			expectedRemain: "",
		},
		{
			name:           "Multiple character string",
			input:          "abc",
			expectedToken:  "a",
			expectedRemain: "bc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := tt.input
			token := getNextToken(&ex)

			if token != tt.expectedToken {
				t.Errorf("getNextToken() = %v, expected %v", token, tt.expectedToken)
			}

			if ex != tt.expectedRemain {
				t.Errorf("remaining string = %v, expected %v", ex, tt.expectedRemain)
			}
		})
	}
}

func TestGetTokens(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    tokens
		expectError error
	}{
		{
			name:  "Valid expression with operators and numbers",
			input: "3 + 5 * ( 10 - 20 )",
			expected: tokens{
				&token{"float", "3"},
				&token{"op", "+"},
				&token{"float", "5"},
				&token{"op", "*"},
				&token{"op", "("},
				&token{"float", "10"},
				&token{"op", "-"},
				&token{"float", "20"},
				&token{"op", ")"},
			},
			expectError: nil,
		},
		{
			name:        "Empty input string",
			input:       "",
			expected:    nil,
			expectError: errs.ErrExpressionNotValid,
		},
		{
			name:        "Invalid character in expression",
			input:       "3 + 5 * a",
			expected:    nil,
			expectError: errs.ErrExpressionNotValid,
		},
		{
			name:  "Valid float numbers",
			input: "3.14 + 2.71",
			expected: tokens{
				&token{"float", "3.14"},
				&token{"op", "+"},
				&token{"float", "2.71"},
			},
			expectError: nil,
		},
		{
			name:        "Single valid token",
			input:       "42",
			expected:    nil,
			expectError: errs.ErrExpressionNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getTokens(tt.input)
			fmt.Printf("%+v\n", result)

			if !errors.Is(err, tt.expectError) {
				t.Errorf("getTokens(%q) error = %v, expected %v", tt.input, err, tt.expectError)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getTokens(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalc(t *testing.T) {
	service := NewCalcService()

	tests := []struct {
		name        string
		expression  string
		expected    float64
		expectError error
	}{
		{
			name:        "Simple addition",
			expression:  "2 + 2",
			expected:    4,
			expectError: nil,
		},
		{
			name:        "Simple subtraction",
			expression:  "5.5 - 3.5",
			expected:    2,
			expectError: nil,
		},
		{
			name:        "Multiplication",
			expression:  "4 * 2",
			expected:    8,
			expectError: nil,
		},
		{
			name:        "Division",
			expression:  "8 / 4",
			expected:    2,
			expectError: nil,
		},
		{
			name:        "Division by zero",
			expression:  "10 / 0",
			expected:    0,
			expectError: errs.ErrExpressionNotValid,
		},
		{
			name:        "Expression with parentheses",
			expression:  "(2 + 3) * 4",
			expected:    20,
			expectError: nil,
		},
		{
			name:        "Complex expression",
			expression:  "10 + 2 * (5 - 3)",
			expected:    14,
			expectError: nil,
		},
		{
			name:        "Complex expression 2",
			expression:  "10 + 2 * (5 - (2 + 1))",
			expected:    14,
			expectError: nil,
		},
		{
			name:        "Invalid expression",
			expression:  "2 + + 3",
			expected:    0,
			expectError: errs.ErrExpressionNotValid,
		},
		{
			name:        "Invalid expression 2",
			expression:  "((2 + 3)",
			expected:    0,
			expectError: errs.ErrExpressionNotValid,
		},
		{
			name:        "Empty expression",
			expression:  "",
			expected:    0,
			expectError: errs.ErrExpressionNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Calc(tt.expression)

			if err != tt.expectError {
				t.Errorf("Calc(%q) error = %v, expected error = %v", tt.expression, err, tt.expectError)
			}

			if result != tt.expected {
				t.Errorf("Calc(%q) = %v, expected = %v", tt.expression, result, tt.expected)
			}
		})
	}
}
