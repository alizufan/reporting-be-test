package middleware

import (
	"context"
	"net/http"

	"reporting/libs/util"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/xid"
)

func Tracker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), util.CTXTrackerID, xid.New().String())
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func JWTValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if t := r.Header.Get("Authorization"); t != "" {
			c, err := VerifyJWT(t)
			if err != nil {
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Write([]byte(http.StatusText(http.StatusUnauthorized)))
				return
			}

			ctx := context.WithValue(r.Context(), util.CTXJWTPayload, c)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return
		}

		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	})
}

func KeyFunc(t *jwt.Token) (any, error) {
	return []byte("secret"), nil
}

func VerifyJWT(tokenStr string) (payload util.JWTPayload, err error) {
	if len(tokenStr) == 0 {
		return
	}

	var token *jwt.Token
	token, err = jwt.ParseWithClaims(tokenStr[7:], &util.JWTPayload{}, KeyFunc)
	if p, ok := token.Claims.(*util.JWTPayload); ok && token.Valid {
		return *p, nil
	}

	return
}
