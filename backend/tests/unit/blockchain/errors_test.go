package blockchain_test

import (
	"errors"
	"testing"

	"encryptkeep-backend/internal/blockchain"
)

func TestParseContractError(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		expectedType   string
		expectedCode   int
		expectedPrefix string
	}{
		{
			name:           "InvalidDataLength error",
			inputError:     errors.New("execution reverted: InvalidDataLength()"),
			expectedType:   "INVALID_DATA_LENGTH",
			expectedCode:   2001,
			expectedPrefix: "Data length must be greater than 0",
		},
		{
			name:           "CannotStoreExistingData error",
			inputError:     errors.New("execution reverted: CannotStoreExistingData(0x123, 1)"),
			expectedType:   "CANNOT_STORE_EXISTING_DATA",
			expectedCode:   2002,
			expectedPrefix: "Cannot store data at existing ID",
		},
		{
			name:           "CannotChangeNonExistentData error",
			inputError:     errors.New("execution reverted: CannotChangeNonExistentData(0x123, 1)"),
			expectedType:   "CANNOT_CHANGE_NON_EXISTENT_DATA",
			expectedCode:   2003,
			expectedPrefix: "Cannot change non-existent data",
		},
		{
			name:           "CannotRemoveNonExistentData error",
			inputError:     errors.New("execution reverted: CannotRemoveNonExistentData(0x123, 1)"),
			expectedType:   "CANNOT_REMOVE_NON_EXISTENT_DATA",
			expectedCode:   2004,
			expectedPrefix: "Cannot remove non-existent data",
		},
		{
			name:           "Gas estimation failed",
			inputError:     errors.New("gas required exceeds allowance"),
			expectedType:   "GAS_ESTIMATION_FAILED",
			expectedCode:   1002,
			expectedPrefix: "Gas required exceeds allowance",
		},
		{
			name:           "Nonce too low",
			inputError:     errors.New("nonce too low"),
			expectedType:   "NONCE_TOO_LOW",
			expectedCode:   1003,
			expectedPrefix: "Transaction nonce too low",
		},
		{
			name:           "Insufficient funds",
			inputError:     errors.New("insufficient funds"),
			expectedType:   "INSUFFICIENT_FUNDS",
			expectedCode:   1004,
			expectedPrefix: "Insufficient funds for transaction",
		},
		{
			name:           "Connection refused",
			inputError:     errors.New("connection refused"),
			expectedType:   "CONNECTION_REFUSED",
			expectedCode:   1005,
			expectedPrefix: "Connection to blockchain refused",
		},
		{
			name:           "Timeout error",
			inputError:     errors.New("timeout"),
			expectedType:   "TIMEOUT",
			expectedCode:   1006,
			expectedPrefix: "Request timeout",
		},
		{
			name:           "Unknown error",
			inputError:     errors.New("some unknown error"),
			expectedType:   "UNKNOWN_ERROR",
			expectedCode:   9999,
			expectedPrefix: "some unknown error",
		},
		{
			name:           "Nil error",
			inputError:     nil,
			expectedType:   "",
			expectedCode:   0,
			expectedPrefix: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := blockchain.ParseContractError(tt.inputError)

			if tt.inputError == nil {
				if result != nil {
					t.Errorf("Expected nil for nil input, got %v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result for non-nil input")
			}

			if result.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, result.Type)
			}

			if result.Code != tt.expectedCode {
				t.Errorf("Expected code %d, got %d", tt.expectedCode, result.Code)
			}

			if tt.expectedPrefix != "" && result.Message != tt.expectedPrefix {
				t.Errorf("Expected message to contain '%s', got '%s'", tt.expectedPrefix, result.Message)
			}
		})
	}
}

func TestIsContractError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Execution reverted",
			err:      errors.New("execution reverted"),
			expected: true,
		},
		{
			name:     "Gas required exceeds allowance",
			err:      errors.New("gas required exceeds allowance"),
			expected: true,
		},
		{
			name:     "Nonce too low",
			err:      errors.New("nonce too low"),
			expected: true,
		},
		{
			name:     "Insufficient funds",
			err:      errors.New("insufficient funds"),
			expected: true,
		},
		{
			name:     "Regular error",
			err:      errors.New("regular error"),
			expected: false,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := blockchain.IsContractError(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetErrorCode(t *testing.T) {
	blockchainErr := &blockchain.BlockchainError{
		Type:    "TEST_ERROR",
		Message: "Test error message",
		Code:    1234,
	}

	code := blockchain.GetErrorCode(blockchainErr)
	if code != 1234 {
		t.Errorf("Expected code 1234, got %d", code)
	}

	// Test with regular error
	regularErr := errors.New("regular error")
	code = blockchain.GetErrorCode(regularErr)
	if code != 0 {
		t.Errorf("Expected code 0 for regular error, got %d", code)
	}
}

func TestGetErrorType(t *testing.T) {
	blockchainErr := &blockchain.BlockchainError{
		Type:    "TEST_ERROR",
		Message: "Test error message",
		Code:    1234,
	}

	errorType := blockchain.GetErrorType(blockchainErr)
	if errorType != "TEST_ERROR" {
		t.Errorf("Expected type 'TEST_ERROR', got '%s'", errorType)
	}

	// Test with regular error
	regularErr := errors.New("regular error")
	errorType = blockchain.GetErrorType(regularErr)
	if errorType != "UNKNOWN" {
		t.Errorf("Expected type 'UNKNOWN' for regular error, got '%s'", errorType)
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := blockchain.WrapError(originalErr, "WRAP_TEST", "Wrapped error", 5000)

	if wrappedErr.Type != "WRAP_TEST" {
		t.Errorf("Expected type 'WRAP_TEST', got '%s'", wrappedErr.Type)
	}

	if wrappedErr.Code != 5000 {
		t.Errorf("Expected code 5000, got %d", wrappedErr.Code)
	}

	expectedMessage := "Wrapped error: original error"
	if wrappedErr.Message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, wrappedErr.Message)
	}
}

func TestNewBlockchainError(t *testing.T) {
	err := blockchain.NewBlockchainError("TEST_TYPE", "Test message", 1234)

	if err.Type != "TEST_TYPE" {
		t.Errorf("Expected type 'TEST_TYPE', got '%s'", err.Type)
	}

	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", err.Message)
	}

	if err.Code != 1234 {
		t.Errorf("Expected code 1234, got %d", err.Code)
	}
}

func TestBlockchainError_Error(t *testing.T) {
	err := &blockchain.BlockchainError{
		Type:    "TEST_TYPE",
		Message: "Test error message",
		Code:    1234,
	}

	errorString := err.Error()
	if errorString != "Test error message" {
		t.Errorf("Expected 'Test error message', got '%s'", errorString)
	}
}
