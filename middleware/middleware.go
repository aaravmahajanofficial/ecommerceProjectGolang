package middleware

import (
	"errors"
	"net/http"

	token "github.com/aaravmahajanofficial/ecommerce-project/tokens"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ClientToken := ctx.Request.Header.Get("Token")

		if ClientToken == "" {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("No authorization header provided"))
			return
		}

		claims, err := token.VerifyToken(ClientToken)

		if err != "" {

			ctx.AbortWithError(http.StatusInternalServerError, errors.New("Token validation failed"))
			return

		}

		ctx.Set("Email", claims.Email)
		ctx.Set("UID", claims.UID)
		ctx.Next()

	}

}
