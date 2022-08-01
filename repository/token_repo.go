package repository

import (
	"miniproject/model"

	"gorm.io/gorm"
)

type UserTokenRepo interface {
	//
	DeleteRefreshToken(userId uint64) error
}

type usertokenRepo struct {
	DB *gorm.DB
}

func NewTokenRepo(dd *gorm.DB) UserTokenRepo {
	return &usertokenRepo{
		DB: dd,
	}
}

func (ut *usertokenRepo) DeleteRefreshToken(userId uint64) error {
	//
	token := &model.UserToken{}
	db := ut.DB.Model(&model.UserToken{})
	err := db.Where("id = ?", userId).Delete(&token).Error
	if err != nil {
		return err
	}

	return nil
}
