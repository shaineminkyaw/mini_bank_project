package repository

import (
	"miniproject/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	//
	FindAll(filter interface{}) ([]*model.User, error)
	FindById(id uint64) (*model.User, error)
	GetVerifyUser(email string) (*model.VerifyCode, error)
	FindByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userRepository struct {
	//
	DB *gorm.DB
}

func NewUserRepository(ds *gorm.DB) UserRepository {
	return &userRepository{
		DB: ds,
	}
}

func (ur *userRepository) FindAll(filter interface{}) ([]*model.User, error) {
	users := make([]*model.User, 0)
	err := ur.DB.Model(&model.User{}).Where(filter).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (ur *userRepository) CreateUser(user *model.User) error {
	err := ur.DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FindById(id uint64) (*model.User, error) {
	user := &model.User{}
	err := ur.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) GetVerifyUser(email string) (*model.VerifyCode, error) {
	//
	vUser := &model.VerifyCode{}
	err := ur.DB.Model(&model.VerifyCode{}).Where("email = ?", email).First(&vUser).Error
	if err != nil {
		return nil, err
	}

	return vUser, nil
}

func (ur *userRepository) FindByEmail(email string) (*model.User, error) {
	//
	user := &model.User{}
	err := ur.DB.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
