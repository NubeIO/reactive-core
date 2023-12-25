package main

import (
	"fmt"
	"testing"
)

func TestCalculateMathOperation(t *testing.T) {
	// Replace with the actual node instance if available

	// Define demo cases
	testCases := []struct {
		operation      string
		inputs         []float64
		expectedResult float64
	}{
		{"add", []float64{1, 2, 3, 4, 5}, 15},
		{"subtract", []float64{10, 2, 3}, 5},
		{"multiply", []float64{2, 3, 4}, 24},
		{"divide", []float64{10, 2, 5}, 1},
		{"min", []float64{5, 3, 9, 1}, 1},
		{"max", []float64{5, 3, 9, 1}, 9},
		{"avg", []float64{1, 2, 3, 4, 5}, 3},
		{"invalid", []float64{1, 2, 3}, 0}, // Invalid operation
		{"add", []float64{}, 0},            // Empty inputs
	}

	for _, testCase := range testCases {
		result := CalculateMathOperation(testCase.operation, testCase.inputs)
		fmt.Println(result, testCase.operation)
		if result != testCase.expectedResult {
			t.Errorf("Operation: %s, Inputs: %v, Expected: %f, Got: %f", testCase.operation, testCase.inputs, testCase.expectedResult, result)
		}
	}
}
