package gapi

import (
	"log"
	"miniproject/config"
	"miniproject/controller"
	"miniproject/ds"
	"miniproject/pb"
	"miniproject/repository"
	"miniproject/service"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

//

type Server struct {
	pb.UnimplementedSimpleBankServer
	H    *controller.Handler
	Data *gorm.DB
}

func NewgRPC() (*Server, error) {

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
		UserService:  userService,
		TokenService: tokenService,
	}
	return &Server{
		H:    (*controller.Handler)(hConfig),
		Data: db,
	}, nil
	// h := NewHandler1(hConfig)
	// return &Server{
	// 	H: h,
	// }, nil

}

func RunGRPCServer() {
	//
	server, err := NewgRPC()
	if err != nil {
		log.Fatalf(err.Error())
	}
	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, server)
	reflection.Register(gRPCServer)

	lis, err := net.Listen("tcp", config.GRPC.GRPCAddress)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("started gRPC server ...%v", config.GRPC.GRPCAddress)

	err = gRPCServer.Serve(lis)
	if err != nil {
		log.Fatal(err.Error())
	}

}
