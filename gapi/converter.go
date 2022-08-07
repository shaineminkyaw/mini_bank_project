package gapi

import (
	"miniproject/model"
	"miniproject/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func converter(user *model.User) *pb.User {
	//
	return &pb.User{
		Email:     user.Email,
		BankCard:  user.BankCardNumber,
		Balance:   user.Balance,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
