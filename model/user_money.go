package model

import "time"

type UserMoney struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	Uid       uint64    `gorm:"column:uid" json:"uid"`
	Amount    float64   `gorm:"column:amount" json:"amount"`
	Type      int8      `gorm:"column:type;comment:type 1:deposit 2:withdrawl" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (um *UserMoney) TableName() string {
	return "user_money"
}
