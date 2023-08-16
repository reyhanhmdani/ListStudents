package middleware

import (
	"errors"
	"ginDatabaseMhs/cfg"
	"ginDatabaseMhs/model/respErr"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// secret key untuk signing token
// middleware konsep nya adalah sesuatu yang ibaratnya intercept , request -> server,
func Authmiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// mengambil token dari header Authorization
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
				Message: "Unauthorized",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// split token dari header
		tokenString := authHeader[len("Bearer "):]

		// parsing token dengan secret key
		claims := &cfg.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return cfg.JwtKey, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
					Message: "Unauthorized",
					Status:  http.StatusUnauthorized,
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &respErr.ErrorResponse{
				Message: "invalid or expired token",
				Status:  http.StatusBadRequest,
			})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
				Message: "Unauthorized (non Valid)",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// Set data pengguna dari token ke dalam konteks
		ctx.Set("username", claims.Username)
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role) // Menambahkan data peran ke konteks

		// Pengecekan peran admin
		if claims.Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, respErr.ErrorResponse{
				Message: "Unauthorized: Only admin can access this endpoint",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// Melanjutkan ke handler jika semua pengecekan berhasil
		ctx.Next()
	}
}

func IsValidToken(tokenString string) bool {
	// Periksa apakah token sudah digunakan sebelumnya
	if cfg.UsedTokens[tokenString] {
		return false
	}

	// Parsing token dengan claims
	claims := &cfg.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return cfg.JwtKey, nil
	})

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	// Periksa apakah token kadaluwarsa
	return claims.ExpiresAt > time.Now().Unix()
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Error("Panic occurred:", r)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, respErr.ErrorResponse{
					Message: "Internal Server Error",
					Status:  http.StatusInternalServerError,
				})
			}
		}()

		ctx.Next()
	}
	//// Cek apakah token sudah digunakan sebelumnya
	//if !IsValidToken(tokenString) {
	//	ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
	//		Message: "Token has already been used or expired",
	//		Status:  http.StatusUnauthorized,
	//	})
	//	return
	//}
	//// Tandai token sebagai digunakan
	//cfg.UsedTokens[tokenString] = true
}
