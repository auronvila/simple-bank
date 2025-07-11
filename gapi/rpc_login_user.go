package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/pb"
	"github.com/auronvila/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (server *Server) LoginUser(context context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(context, req.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user could not be found")
		}
		return nil, status.Errorf(codes.Internal, "internal server err!!")
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	duration, _ := time.ParseDuration(server.config.AccessTokenDuration)
	accessToken, accessTokenPayload, err := server.tokenMaker.GenerateToken(user.Username, duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server err!! %s", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.GenerateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server err!! %s", err)
	}

	mtdt := server.extractMetadata(context)

	session, err := server.store.CreateSession(context, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
		IsBlocked:    false,
		ExpiredAt:    refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server err!! %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	// you can print the response that is being returned as bytes so the dates are printed correctly
	//jsonBytes, err := protojson.Marshal(rsp)
	//if err != nil {
	//	log.Fatal("Failed to marshal proto response:", err)
	//}
	//fmt.Println(string(jsonBytes))

	return rsp, nil
}
