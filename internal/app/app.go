package app

import (
	grpcapp "auth-application/internal/app/grpc"
	"auth-application/internal/service/auth"
	"auth-application/internal/service/token"
	"auth-application/internal/storage/postgresql"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(grpcPort int, connectionString string, secret string, tokenTTL time.Duration) *App {
	tokenService := token.New(secret, tokenTTL)
	storage, err := postgresql.New(connectionString)
	if err != nil {
		panic("dont connected to database")
	}
	authService := auth.New(storage, storage, storage, tokenService)
	if err != nil {
		return nil
	}
	server := grpcapp.New(grpcPort, authService)

	return &App{
		GRPCServer: server,
	}
}
