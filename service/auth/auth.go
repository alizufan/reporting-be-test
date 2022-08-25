package auth

import (
	"context"
	"crypto/md5"
	"fmt"
	"reporting/libs/util"
	"reporting/repository/merchant"
	"reporting/repository/user"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rotisserie/eris"
)

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

type Auth struct {
	UserRepo     user.UserRepository
	MerchantRepo merchant.MerchantRepository
}

type (
	LoginRequest struct {
		Username string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	LoginResponse struct {
		Token string `json:"token"`
	}
)

func (a *Auth) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	userRes, err := a.UserRepo.Find(ctx, user.FindFilter{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	merRes, err := a.MerchantRepo.Find(ctx, merchant.FindFilter{
		UserID: userRes.ID,
	})
	if err != nil {
		return nil, err
	}

	var (
		reqPass  = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
		userPass = userRes.Password
	)
	if reqPass != userPass {
		return nil, eris.Wrap(util.ErrInvalid, "password un-match")
	}

	claims := util.JWTPayload{
		UserID:     userRes.ID,
		MerchantID: merRes.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "ReportingAPI",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	if err != nil {
		return nil, eris.Wrap(err, "sign jwt token, an error occurred")
	}

	return &LoginResponse{
		Token: token,
	}, nil
}
