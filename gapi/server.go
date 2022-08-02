package gapi

import (
	"fmt"
	"log"
	"miniproject/config"
	"miniproject/controller"
	"miniproject/ds"
	"miniproject/pb"
	"miniproject/repository"
	"miniproject/service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//

type Server struct {
	pb.UnimplementedSimpleBankServer
	H *controller.Handler
}

func NewgRPC(router *gin.Engine) (*Server, error) {

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
	hConfig := &controller.HConfig{
		Router:       router,
		UserService:  userService,
		TokenService: tokenService,
	}
	h := NewHandler1(hConfig)
	return &Server{
		H: h,
	}, nil

}

func NewHandler1(r *controller.HConfig) *controller.Handler {
	h := &controller.Handler{
		Router:       r.Router,
		UserService:  r.UserService,
		TokenService: r.TokenService,
	}

	//user controller
	userController := controller.NewUserController(h)
	userController.Register()

	//user money
	transferMoney := controller.NewTransferMoney(h)
	transferMoney.Register()

	return h
}

func RunGRPCServer() {
	//
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.10.199"})
	host := "localhost"
	gRPCServer := grpc.NewServer()
	server, err := NewgRPC(router)
	if err != nil {
		log.Fatalf(err.Error())
	}
	pb.RegisterSimpleBankServer(gRPCServer, server)
	reflection.Register(gRPCServer)

	starter := fmt.Sprintf("%v:%v", host, config.GRPC.Port)
	router.Run(starter)
}
