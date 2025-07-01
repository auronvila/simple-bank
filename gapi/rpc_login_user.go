package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/pb"
	"github.com/auronvila/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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

	md, ok := metadata.FromIncomingContext(context)
	if !ok {
		return nil, status.Errorf(codes.Internal, "cannot get metadata from context")
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

	userAgents := md.Get("user-agent")
	clientIPs := md.Get("client-ip")

	var userAgent string
	if len(userAgents) > 0 {
		userAgent = userAgents[0]
	}

	var clientIP string
	if len(clientIPs) > 0 {
		clientIP = clientIPs[0]
	} else {
		// optionally infer from peer info
		if p, ok := peer.FromContext(context); ok {
			clientIP = p.Addr.String()
		}
	}

	session, err := server.store.CreateSession(context, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
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
