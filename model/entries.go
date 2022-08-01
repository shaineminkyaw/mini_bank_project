package model

import "time"

type MoneyEntries struct {
	ID           uint64    `gorm:"column:id" json:"id"`
	AccountID    uint64    `gorm:"column:account_id" json:"account_id"`
	AccountEmail string    `gorm:"column:account_email" json:"account_email"`
	Amount       float64   `gorm:"coulmn:amount" json:"amount"`
	User         *User     `gorm:"foreignKey:account_id"`
	Type         int8      `gorm:"column:type;comment:type 3:transfer 4:receive" json:"type"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (me *MoneyEntries) TableName() string {
	return "user_money_transfer_entries"
}
