package service

import (
	"crypto/rsa"
	"log"
	"miniproject/config"
	"miniproject/ds"
	"miniproject/model"
	"miniproject/repository"
	"miniproject/util"
	"time"
)

type UserTokenService interface {
	//
	DeleteRefreshToken(userId uint64) error
	NewTokenPair(user *model.User, prvToken string) (*TokenPair, error)
}

type userTokenService struct {
	Repository repository.UserTokenRepo
	Private    *rsa.PrivateKey
	Public     *rsa.PublicKey
}
type UsTokenConfig struct {
	Repository repository.UserTokenRepo
	Private    *rsa.PrivateKey
	Public     *rsa.PublicKey
}

type AccessToken struct {
	SS string
}
type RefreshToken struct {
	Uid         uint64
	TokenID     string
	TokenString string
	Expire      time.Duration
}

type TokenPair struct {
	AccessToken
	RefreshToken
}

func NewUserToken(ut *UsTokenConfig) UserTokenService {
	//
	return &userTokenService{
		Repository: ut.Repository,
		Private:    config.PrivateKey,
		Public:     config.PublicKey,
	}
}

func (us *userTokenService) DeleteRefreshToken(userId uint64) error {
	//
	err := us.Repository.DeleteRefreshToken(userId)
	if err != nil {
		return err
	}
	return nil
}

func (us *userTokenService) NewTokenPair(user *model.User, prvToken string) (*TokenPair, error) {
	//
	if len(prvToken) > 0 {
		err := us.DeleteRefreshToken(user.Id)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	accessToken, err := util.GetAccessToken(user.Id, us.Private)
	if err != nil {
		log.Printf(err.Error())
	}
	refreshToken, err := util.GetRefreshToken(user.Id, us.Private)
	if err != nil {
		log.Printf(err.Error())
	}
	rToken, err := util.ValidRefreshToken(refreshToken.TokenString, us.Public)
	if err != nil {
		log.Printf(err.Error())
	}
	token := &model.UserToken{
		Uid:        user.Id,
		TokenID:    rToken.TokenID.String(),
		Token:      rToken.TokenString,
		ExpireTime: time.Now().Add(rToken.Expire),
	}

	err = ds.DB.Model(&model.UserToken{}).Create(&token).Error
	if err != nil {
		log.Printf(err.Error())
	}

	return &TokenPair{
		AccessToken: AccessToken{
			accessToken,
		},
		RefreshToken: RefreshToken{
			Uid:         token.Uid,
			TokenID:     token.TokenID,
			TokenString: token.Token,
			Expire:      time.Now().Sub(token.ExpireTime),
		},
	}, nil
}
