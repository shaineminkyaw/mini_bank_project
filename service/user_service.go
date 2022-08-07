package service

import (
	"miniproject/ds"
	"miniproject/model"
	"miniproject/repository"
)

type UserService interface {
	//
	FindAll(filter interface{}) ([]*model.User, error)
	FindById(id uint64) (*model.User, error)
	GetVerifyUser(email string) (*model.VerifyCode, error)
	FindByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userService struct {
	//
	UserRepository repository.UserRepository
}

type UsConfig struct {
	UserRepository repository.UserRepository
}

func NewUserService(us *UsConfig) UserService {
	return &userService{
		UserRepository: us.UserRepository,
	}
}

func (ur *userService) CreateUser(user *model.User) error {
	err := ds.DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (us *userService) FindAll(filter interface{}) ([]*model.User, error) {
	users, err := us.UserRepository.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *userService) FindById(id uint64) (*model.User, error) {
	user, err := us.UserRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (us *userService) GetVerifyUser(email string) (*model.VerifyCode, error) {
	vUser, err := us.UserRepository.GetVerifyUser(email)
	if err != nil {
		return nil, err
	}
	return vUser, nil
}
func (us *userService) FindByEmail(email string) (*model.User, error) {
	//
	user, err := us.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
