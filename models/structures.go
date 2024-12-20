package models

import (
	"github.com/google/uuid"
)

type Wallet struct {
	WalletId      uuid.UUID `json:"valletId" gorm:"type:uuid;primaryKey;uniqueIndex"`
	OperationType string    `json:"operationType" gorm:"-"`
	Amount        int       `json:"amount" gorm:"-"`
	Balance       int       `json:"-" gorm:"type:int"`
}
