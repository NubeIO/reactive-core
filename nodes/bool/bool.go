package main

import (
	"fmt"
)

func CalculateLogicalOperation(operation string, inputs []interface{}) interface{} {
	if len(inputs) == 0 {
		return false
	}

	// Define a map that associates each operation type with a function
	operations := map[string]func([]float64) float64{
		"and": func(arr []float64) float64 { return logicalAND(arr) },
		"or":  func(arr []float64) float64 { return logicalOR(arr) },
	}

	operationFunc, exists := operations[operation]
	if !exists {
		// Invalid operation type, return false (or handle it as needed)
		return false
	}

	// Convert and demo input values
	convertedInputs := make([]float64, len(inputs))
	for i, input := range inputs {
		converted, err := convertToFloat64(input)
		if err != nil {
			// Handle the error (unsupported type) as needed
			return false
		}
		convertedInputs[i] = converted
	}

	result := operationFunc(convertedInputs)

	if result == 1.0 {
		return true
	}
	return false
}

func convertToFloat64(input interface{}) (float64, error) {
	switch v := input.(type) {
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0.0, fmt.Errorf("unsupported type: %T", input)
	}
}

func logicalAND(inputs []float64) float64 {
	for _, input := range inputs {
		if input != 1.0 {
			return 0.0
		}
	}
	return 1.0
}

func logicalOR(inputs []float64) float64 {
	for _, input := range inputs {
		if input == 1.0 {
			return 1.0
		}
	}
	return 0.0
}
