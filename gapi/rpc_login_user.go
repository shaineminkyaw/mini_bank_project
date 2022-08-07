package gapi

import (
	"context"
	"miniproject/pb"
	"miniproject/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	//

	user, err := server.H.UserService.FindByEmail(req.GetEmail())
	if err == gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found in DB")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}
	ok, err := util.ValidateHashedPassword(user.Password, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "password not match")
	}
	token, err := server.H.TokenService.NewTokenPair(user, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error on token generate")
	}
	resp := &pb.LoginUserResponse{
		User:       converter(user),
		AccessToke: token.AccessToken.SS,
	}
	return resp, nil

	//

}
