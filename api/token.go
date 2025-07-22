package api

import (
	"errors"
	"fmt"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	refreshTokenPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshTokenPayload.ID)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if session.Username != refreshTokenPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("missmatched session token")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	duration, _ := time.ParseDuration(server.config.AccessTokenDuration)
	accessToken, accessTokenPayload, err := server.tokenMaker.GenerateToken(refreshTokenPayload.Username, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	rsp := renewAccessTokenResponse{
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
		AccessToken:          accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}
