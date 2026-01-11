package blockchain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrConnectionFailed        = errors.New("failed to connect to blockchain")
	ErrContractNotFound        = errors.New("contract not found at address")
	ErrInvalidAddress          = errors.New("invalid address format")
	ErrInsufficientGas         = errors.New("insufficient gas for transaction")
	ErrTransactionFailed       = errors.New("transaction failed")
	ErrDataNotFound            = errors.New("data not found in blockchain")
	ErrInvalidPrivateKey       = errors.New("invalid private key")
	ErrNetworkUnavailable      = errors.New("blockchain network unavailable")
	ErrContractCallFailed      = errors.New("contract call failed")
	ErrInvalidChainID          = errors.New("invalid chain ID")
	ErrInvalidNonce            = errors.New("invalid nonce")
	ErrContractExecutionFailed = errors.New("failed to execute contract function")
	ErrGasEstimationFailed     = errors.New("gas estimation failed")
	ErrTransactionReverted     = errors.New("transaction reverted")
	ErrNonceTooLow             = errors.New("nonce too low")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrNotConnected            = errors.New("not connected to blockchain")

	ErrInvalidDataLength           = errors.New("invalid data length")
	ErrCannotStoreExistingData     = errors.New("cannot store existing data")
	ErrCannotChangeNonExistentData = errors.New("cannot change non-existent data")
	ErrCannotRemoveNonExistentData = errors.New("cannot remove non-existent data")
)

type BlockchainError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *BlockchainError) Error() string {
	return e.Message
}

func NewBlockchainError(errorType, message string, code int) *BlockchainError {
	return &BlockchainError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}

func ParseContractError(err error) *BlockchainError {
	if err == nil {
		return nil
	}

	errorStr := err.Error()

	if strings.Contains(errorStr, "execution reverted") {
		if strings.Contains(errorStr, "InvalidDataLength()") {
			return NewBlockchainError("INVALID_DATA_LENGTH", "Data length must be greater than 0", 2001)
		}
		if strings.Contains(errorStr, "CannotStoreExistingData") {
			return NewBlockchainError("CANNOT_STORE_EXISTING_DATA", "Cannot store data at existing ID", 2002)
		}
		if strings.Contains(errorStr, "CannotChangeNonExistentData") {
			return NewBlockchainError("CANNOT_CHANGE_NON_EXISTENT_DATA", "Cannot change non-existent data", 2003)
		}
		if strings.Contains(errorStr, "CannotRemoveNonExistentData") {
			return NewBlockchainError("CANNOT_REMOVE_NON_EXISTENT_DATA", "Cannot remove non-existent data", 2004)
		}

		return NewBlockchainError("CONTRACT_REVERTED", "Contract execution reverted", 1001)
	}

	if strings.Contains(errorStr, "gas required exceeds allowance") {
		return NewBlockchainError("GAS_ESTIMATION_FAILED", "Gas required exceeds allowance", 1002)
	}

	if strings.Contains(errorStr, "nonce too low") {
		return NewBlockchainError("NONCE_TOO_LOW", "Transaction nonce too low", 1003)
	}

	if strings.Contains(errorStr, "insufficient funds") {
		return NewBlockchainError("INSUFFICIENT_FUNDS", "Insufficient funds for transaction", 1004)
	}

	if strings.Contains(errorStr, "connection refused") {
		return NewBlockchainError("CONNECTION_REFUSED", "Connection to blockchain refused", 1005)
	}

	if strings.Contains(errorStr, "timeout") {
		return NewBlockchainError("TIMEOUT", "Request timeout", 1006)
	}

	return NewBlockchainError("UNKNOWN_ERROR", errorStr, 9999)
}

func IsContractError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()
	return strings.Contains(errorStr, "execution reverted") ||
		strings.Contains(errorStr, "gas required exceeds allowance") ||
		strings.Contains(errorStr, "nonce too low") ||
		strings.Contains(errorStr, "insufficient funds")
}

func GetErrorCode(err error) int {
	if blockchainErr, ok := err.(*BlockchainError); ok {
		return blockchainErr.Code
	}
	return 0
}

func GetErrorType(err error) string {
	if blockchainErr, ok := err.(*BlockchainError); ok {
		return blockchainErr.Type
	}
	return "UNKNOWN"
}

func WrapError(err error, errorType, message string, code int) *BlockchainError {
	return &BlockchainError{
		Type:    errorType,
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
		Code:    code,
	}
}
