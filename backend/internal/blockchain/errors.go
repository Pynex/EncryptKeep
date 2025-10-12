package blockchain

import "errors"

// Предопределённые ошибки блокчейна
var (
	ErrConnectionFailed   = errors.New("failed to connect to blockchain")
	ErrContractNotFound   = errors.New("contract not found at address")
	ErrInvalidAddress     = errors.New("invalid address format")
	ErrInsufficientGas    = errors.New("insufficient gas for transaction")
	ErrTransactionFailed  = errors.New("transaction failed")
	ErrDataNotFound       = errors.New("data not found in blockchain")
	ErrInvalidPrivateKey  = errors.New("invalid private key")
	ErrNetworkUnavailable = errors.New("blockchain network unavailable")
	ErrContractCallFailed = errors.New("contract call failed")
	ErrInvalidChainID     = errors.New("invalid chain ID")
)

// BlockchainError представляет ошибку блокчейна с дополнительным контекстом
type BlockchainError struct {
	Type    string `json:"type"`    // Тип ошибки
	Message string `json:"message"` // Сообщение об ошибке
	Code    int    `json:"code"`    // Код ошибки
}

func (e *BlockchainError) Error() string {
	return e.Message
}

// NewBlockchainError создаёт новую ошибку блокчейна
func NewBlockchainError(errorType, message string, code int) *BlockchainError {
	return &BlockchainError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}

// Задачи:
// Определить все предопределённые ошибки
// Создать структуру BlockchainError с дополнительным контекстом
// Реализовать метод Error() для BlockchainError
// Создать функцию NewBlockchainError() для создания ошибок
