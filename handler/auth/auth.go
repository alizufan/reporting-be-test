package auth

import (
	"net/http"
	"reporting/libs/util"
	"reporting/service/auth"
)

type AuthHandler struct {
	AuthService auth.AuthService
}

func (a *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		req = auth.LoginRequest{}
	)

	if ok := util.RequestBodyValidation(rw, r.Body, &req); !ok {
		return
	}

	res, err := a.AuthService.Login(ctx, req)
	if err != nil {
		util.ErrorHTTPResponse(ctx, rw, err)
		return
	}

	util.HTTPResponse(rw, http.StatusOK, "success login", res, nil)
}
