package gapi

import (
	"context"
	"miniproject/model"
	"miniproject/pb"
	"miniproject/util"

	"github.com/mazen160/go-random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	//hash password
	userPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {

		return nil, status.Errorf(codes.InvalidArgument, "password not match")
	}
	charset := "12345678"
	length := 8
	usr, _ := random.Random(length, charset, true)
	userName := "U" + usr
	currency := "USD"
	bankCard, err := util.GetBankCardNumber(req.City)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error on bank card generate")
	}
	user := &model.User{
		Username:       userName,
		Password:       userPassword,
		Email:          req.GetEmail(),
		NationID:       req.GetNationId(),
		BankCardNumber: bankCard,
		City:           req.GetCity(),
		Balance:        0,
		Currency:       currency,
		Type:           int8(req.GetType()),
	}
	mail := req.GetEmail()
	_, err = server.H.UserService.FindByEmail(mail)
	if err == gorm.ErrRecordNotFound {
		err := server.H.UserService.CreateUser(user)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error on user create")
		}
		err = util.SaveUserBankCard(req.GetCity(), user.Id, bankCard)
		if err != nil {
			return nil, status.Errorf(codes.Unimplemented, "error on create user_bankcard")
		}
	} else {
		return nil, status.Errorf(codes.AlreadyExists, "User already exists!")
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, "Something went wrong")
	}

	resp := &pb.CreateUserResponse{
		User: converter(user),
	}
	//
	return resp, nil
}
