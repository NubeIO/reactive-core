package errormap

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestErrorMap(t *testing.T) {
	// Create a new Logger instance for testing (you can configure it as needed)
	logger := logrus.New()

	// Create a new ErrorMap instance for testing with a specified file path and logger
	filePath := "test_errors.json"
	errorMap := NewErrorMap("test", 10, "info", logger, filePath, 10)

	// Add some test errors
	if err := errorMap.AddError("1", "TypeA", "Error 1"); err != nil {
		t.Errorf("Error adding error 1: %v", err)
	}
	if err := errorMap.AddError("2", "TypeB", "Error 2"); err != nil {
		t.Errorf("Error adding error 2: %v", err)
	}

	// Retrieve errors from disk
	errorsFromDisc, _ := errorMap.GetErrorsFromDisc()

	for i, e := range errorsFromDisc {
		fmt.Println(i)
		fmt.Println(e)
	}

}
