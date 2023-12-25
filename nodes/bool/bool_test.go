package main

import (
	"fmt"
	"testing"
)

func TestCalculateLogicalOperation(t *testing.T) {
	// Test cases for logical AND
	boolResult := CalculateLogicalOperation("AND", []interface{}{true, true, 1.1})
	fmt.Println("AND Result (bool):", boolResult)

	boolResult = CalculateLogicalOperation("AND", []interface{}{true, true, 1})
	fmt.Println("AND Result (bool):", boolResult)

	intResult := CalculateLogicalOperation("OR", []interface{}{2, 0, 0})
	fmt.Println("OR Result (int):", intResult)

	floatResult := CalculateLogicalOperation("OR", []interface{}{1, 0.0, 3.0})
	fmt.Println("OR Result (float64):", floatResult)
}
