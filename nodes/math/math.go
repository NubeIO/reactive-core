package main

func CalculateMathOperation(operation string, inputs []float64) float64 {
	if len(inputs) == 0 {
		// No input values provided, return 0 (or handle it as needed)
		return 0
	}

	// Define a map that associates each operation type with a function
	operations := map[string]func([]float64) float64{
		"add":      func(arr []float64) float64 { return sum(arr) },
		"subtract": func(arr []float64) float64 { return subtract(arr) },
		"multiply": func(arr []float64) float64 { return multiply(arr) },
		"divide":   func(arr []float64) float64 { return divide(arr) },
		"min":      func(arr []float64) float64 { return minimum(arr) },
		"max":      func(arr []float64) float64 { return maximum(arr) },
		"avg":      func(arr []float64) float64 { return avg(arr) },
	}

	operationFunc, exists := operations[operation]
	if !exists {
		// Invalid math operation, return 0 (or handle it as needed)
		return 0
	}

	return operationFunc(inputs)
}

func sum(arr []float64) float64 {
	result := 0.0
	for _, v := range arr {
		result += v
	}
	return result
}

func subtract(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result -= arr[i]
	}
	return result
}

func multiply(arr []float64) float64 {
	result := 1.0
	for _, v := range arr {
		result *= v
	}
	return result
}

func divide(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] != 0 {
			result /= arr[i]
		} else {
			// Handle division by zero gracefully (you can customize this behavior)
			return 0
		}
	}
	return result
}

func minimum(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] < result {
			result = arr[i]
		}
	}
	return result
}

func maximum(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > result {
			result = arr[i]
		}
	}
	return result
}

func avg(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	total := sum(arr)
	return total / float64(len(arr))
}
