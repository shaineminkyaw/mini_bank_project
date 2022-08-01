package model

import "time"

type MoneyTransfer struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	FromID    uint64    `gorm:"column:from_id" json:"from_id"`
	FromEmail string    `gorm:"column:from_email" json:"from_email"`
	ToID      uint64    `gorm:"column:to_id" json:"to_id"`
	ToEmail   string    `gorm:"to_email" json:"to_email"`
	Amount    float64   `gorm:"amount" json:"amount"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	DeletedAt time.Time `gorm:"deleted_at" json:"deleted_at"`
}

func (mt *MoneyTransfer) TableName() string {
	return "user_money_transfer"
}
