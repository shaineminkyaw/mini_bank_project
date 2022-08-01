package controller

import (
	"miniproject/config"
	"miniproject/ds"
	"miniproject/repository"
	"miniproject/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router       *gin.Engine
	UserService  service.UserService
	TokenService service.UserTokenService
}

type HConfig struct {
	Router       *gin.Engine
	UserService  service.UserService
	TokenService service.UserTokenService
}

func Inject(router *gin.Engine) *Handler {
	db := ds.DB

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(&service.UsConfig{
		UserRepository: userRepo,
	})

	tokenRepo := repository.NewTokenRepo(db)
	tokenService := service.NewUserToken(&service.UsTokenConfig{
		Repository: tokenRepo,
		Private:    config.PrivateKey,
		Public:     config.PublicKey,
	})
	hConfig := &HConfig{
		Router:       router,
		UserService:  userService,
		TokenService: tokenService,
	}
	h := NewHandler(hConfig)
	return h
}

func NewHandler(r *HConfig) *Handler {
	h := &Handler{
		Router:       r.Router,
		UserService:  r.UserService,
		TokenService: r.TokenService,
	}

	//user controller
	userController := newUserController(h)
	userController.Register()

	//user money
	transferMoney := NewTransferMoney(h)
	transferMoney.Register()

	return h
}
