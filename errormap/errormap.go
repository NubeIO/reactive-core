package errormap

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	current      = "current"
	acknowledged = "acknowledged"
)

// Error represents an individual error.
type Error struct {
	ID           string
	Type         string
	State        string // "current" or "acknowledged"
	ErrorMessage string
	LoggerType   string // "info", "debug", "error", "warning", "disabled"
	Timestamp    time.Time
}

// ErrorMap represents a map of errors.
type ErrorMap struct {
	Key            string
	Errors         map[string]Error
	Size           int
	MaxSize        int
	removedCount   int
	MaxMessageSize int // Maximum error message size to save to disk
	mux            sync.Mutex
	loggerType     string         // "info", "debug", "error", "warning", "disabled"
	logger         *logrus.Logger // Logger for logging
	savePath       string         // Path to save errors JSON file
	filePath       string         // File path for storing errors

}

// NewErrorMap creates a new ErrorMap with the specified key, maximum size, logger type, save path, and maximum message size.
func NewErrorMap(key string, maxSize int, loggerType string, logger *logrus.Logger, savePath string, maxMessageSize int) *ErrorMap {
	return &ErrorMap{
		Key:            key,
		Errors:         make(map[string]Error),
		Size:           0,
		MaxSize:        maxSize,
		loggerType:     loggerType,
		logger:         logger,
		savePath:       savePath,
		MaxMessageSize: maxMessageSize,
	}
}

// SetLoggerType updates the logger type for the ErrorMap.
func (em *ErrorMap) SetLoggerType(loggerType string) {
	em.mux.Lock()
	defer em.mux.Unlock()

	em.loggerType = loggerType
}

// AddError adds a new error to the ErrorMap and returns an error if any.
func (em *ErrorMap) AddError(errorID, errorType, errorMessage string) error {
	em.mux.Lock()
	defer em.mux.Unlock()

	if errorID == "" {
		return errors.New("error-id cannot be empty")
	}

	if len(errorMessage) > em.MaxMessageSize {
		return fmt.Errorf("error message exceeds the maximum size of %d", em.MaxMessageSize)
	}

	// Check if the ErrorMap has reached its maximum size
	if em.Size >= em.MaxSize {
		// Remove the oldest error if the limit is exceeded
		if errMsg := em.removeOldestError(); errMsg != nil {
			return errMsg
		}
	}

	newError := Error{
		ID:           errorID,
		Type:         errorType,
		State:        current,
		ErrorMessage: errorMessage,
		LoggerType:   em.loggerType, // Assign the current logger type
	}
	em.Errors[errorID] = newError
	em.Size++

	// Log the addition of the error with the appropriate logger type
	switch em.loggerType {
	case "info":
		em.logger.Infof("Added new error: ErrorID=%s, ErrorType=%s, ErrorMessage=%s", errorID, errorType, errorMessage)
	case "debug":
		em.logger.Debugf("Added new error: ErrorID=%s, ErrorType=%s, ErrorMessage=%s", errorID, errorType, errorMessage)
	case "error":
		em.logger.Errorf("Added new error: ErrorID=%s, ErrorType=%s, ErrorMessage=%s", errorID, errorType, errorMessage)
	case "warning":
		em.logger.Warningf("Added new error: ErrorID=%s, ErrorType=%s, ErrorMessage=%s", errorID, errorType, errorMessage)
	}

	// Save errors to the file after adding a new error
	if err := em.saveErrorsToFile(); err != nil {
		em.logger.Errorf("Error saving errors to file: %v", err)
	}

	return nil
}

// removeOldestError removes the oldest error from the ErrorMap.
func (em *ErrorMap) removeOldestError() error {
	if em.Size <= 0 {
		return errors.New("no errors to remove")
	}

	// Find the oldest error and remove it
	for key := range em.Errors {
		delete(em.Errors, key)
		em.Size--
		em.removedCount++
		if em.Size <= 0 {
			return nil
		}
	}
	return nil
}

// DeleteError removes an error from the ErrorMap by its ID.
func (em *ErrorMap) DeleteError(errorID string) error {
	em.mux.Lock()
	defer em.mux.Unlock()

	_, exists := em.Errors[errorID]
	if !exists {
		return errors.New("error not found")
	}

	delete(em.Errors, errorID)

	// Save errors to the file after deleting an error
	if err := em.saveErrorsToFile(); err != nil {
		em.logger.Errorf("Error saving errors to file: %v", err)
	}

	return nil
}

// AcknowledgeError changes the state of an error from "current" to "acknowledged".
func (em *ErrorMap) AcknowledgeError(errorID string) error {
	em.mux.Lock()
	defer em.mux.Unlock()

	errMsg, exists := em.Errors[errorID]
	if !exists {
		return errors.New("error not found")
	}
	errMsg.State = acknowledged
	em.Errors[errorID] = errMsg

	// Save errors to the file after acknowledging an error
	if err := em.saveErrorsToFile(); err != nil {
		em.logger.Errorf("Error saving errors to file: %v", err)
	}

	return nil
}

func (em *ErrorMap) JSONDump() ([]byte, error) {
	data, err := json.MarshalIndent(em.Errors, "", "  ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAllErrors returns all errors in the ErrorMap.
func (em *ErrorMap) GetAllErrors() map[string]Error {
	em.mux.Lock()
	defer em.mux.Unlock()

	return em.Errors
}

// GetAllErrorsByType returns all errors in the ErrorMap of a specific type.
func (em *ErrorMap) GetAllErrorsByType(errorType string) map[string]Error {
	em.mux.Lock()
	defer em.mux.Unlock()

	result := make(map[string]Error)

	for _, errMsg := range em.Errors {
		if errMsg.Type == errorType {
			result[errMsg.ID] = errMsg
		}
	}

	return result
}

// GetAllErrorsCurrent returns all current errors in the ErrorMap.
func (em *ErrorMap) GetAllErrorsCurrent() map[string]Error {
	em.mux.Lock()
	defer em.mux.Unlock()

	result := make(map[string]Error)

	for _, errMsg := range em.Errors {
		if errMsg.State == "current" {
			result[errMsg.ID] = errMsg
		}
	}

	return result
}

// DeleteAllErrors removes all errors from the ErrorMap.
func (em *ErrorMap) DeleteAllErrors() {
	em.mux.Lock()
	defer em.mux.Unlock()

	em.Errors = make(map[string]Error)
	em.Size = 0

	// Save errors to the file after deleting all errors
	if err := em.saveErrorsToFile(); err != nil {
		em.logger.Errorf("Error saving errors to file: %v", err)
	}
}

// StartCronJob starts a cron job to periodically save errors to disk and clear in-memory errors.
func (em *ErrorMap) StartCronJob(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			em.SaveToDisk(10)
			em.ClearInMemoryErrors() // Clear in-memory errors after saving to disk
		}
	}()
}

// ClearInMemoryErrors clears all in-memory errors.
func (em *ErrorMap) ClearInMemoryErrors() {
	em.mux.Lock()
	defer em.mux.Unlock()

	em.Errors = make(map[string]Error)
	em.Size = 0
}

// SaveToDisk saves errors to a file, limiting the number of errors stored to 'limit'.
func (em *ErrorMap) SaveToDisk(limit int) error {
	em.mux.Lock()
	defer em.mux.Unlock()

	if limit <= 0 {
		return errors.New("limit must be greater than zero")
	}

	// Sort errors by timestamp (if applicable) or in any other way to prioritize which ones to keep
	var errorsToSave []Error
	for _, errMsg := range em.Errors {
		errorsToSave = append(errorsToSave, errMsg)
	}

	// Sort errors by timestamp in descending order (most recent first)
	sort.Slice(errorsToSave, func(i, j int) bool {
		return errorsToSave[i].Timestamp.After(errorsToSave[j].Timestamp)
	})

	// If the limit is reached, truncate the errors slice
	if len(errorsToSave) > limit {
		errorsToSave = errorsToSave[:limit]
	}

	// Serialize errors to JSON
	data, err := json.Marshal(errorsToSave)
	if err != nil {
		return fmt.Errorf("error serializing errors: %v", err)
	}

	// Write errors data to the file
	err = os.WriteFile(em.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing errors to file: %v", err)
	}

	return nil
}

// saveErrorsToFile saves errors to a JSON file.
func (em *ErrorMap) saveErrorsToFile() error {
	data, err := em.JSONDump()
	if err != nil {
		return err
	}

	err = os.WriteFile(em.savePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetErrorsFromDisc retrieves errors from a file on disk.
func (em *ErrorMap) GetErrorsFromDisc() ([]Error, error) {
	data, err := os.ReadFile(em.filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading errors from file: %v", err)
	}

	var errorsFromDisc []Error
	err = json.Unmarshal(data, &errorsFromDisc)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling errors from file: %v", err)
	}

	return errorsFromDisc, nil
}

// GetErrorsFromDiscByType retrieves errors of a specific type from a file on disk.
func (em *ErrorMap) GetErrorsFromDiscByType(errorType string) ([]Error, error) {
	allErrors, err := em.GetErrorsFromDisc()
	if err != nil {
		return nil, err
	}

	var errorsByType []Error
	for _, err := range allErrors {
		if err.Type == errorType {
			errorsByType = append(errorsByType, err)
		}
	}

	return errorsByType, nil
}
