package model

import "time"

type UserBankCard struct {
	ID            uint64    `gorm:"column:id" json:"id"`
	Uid           uint64    `gorm:"column:uid" json:"uid"`
	BanCardNumber string    `gorm:"column:bank_card" json:"bank_card"`
	User          *User     `gorm:"foreignKey:uid"`
	CreatedAt     time.Time `gorm:"coulmn:created_at" json:"created_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ub *UserBankCard) TableName() string {
	return "user_BankCard"
}

type UserCityTotalBankCard struct {
	ID            uint64 `gorm:"column:id" json:"id"`
	YangonCard    int64  `gorm:"column:yangon_total_card" json:"yangon_total_card"`
	MandalayCard  int64  `gorm:"column:mandalay_total_card" json:"mandalay_total_card"`
	NaypyitawCard int64  `gorm:"column:naypyitaw_total_card" json:"naypyitaw_total_card"`
	TaunggyiCard  int64  `gorm:"column:taunggyi_total_card" json:"taunggyi_total_card"`
}

func (uc *UserCityTotalBankCard) TableName() string {
	return "user_division_bankInfo"
}
