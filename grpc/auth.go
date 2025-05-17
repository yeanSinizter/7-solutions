package grpc

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	JWTSecret string
}

func NewAuthInterceptor(secret string) *AuthInterceptor {
	return &AuthInterceptor{JWTSecret: secret}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		userID, err := a.authorize(ctx)
		if err != nil {
			return nil, err
		}

		newCtx := context.WithValue(ctx, "userID", userID)
		return handler(newCtx, req)
	}
}

func (a *AuthInterceptor) authorize(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return "", errors.New("authorization token not provided")
	}

	tokenParts := strings.SplitN(authHeaders[0], " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("invalid authorization format")
	}

	tokenString := tokenParts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		log.Println("Invalid token:", err)
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token claims")
	}

	return userID, nil
}
