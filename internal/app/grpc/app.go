package grpcapp

import (
	"auth-application/internal/grpc/auth"
	authService "auth-application/internal/service/auth"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(port int, aService *authService.Auth) *App {
	grpcServer := grpc.NewServer()
	auth.Register(grpcServer, aService)

	return &App{
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}
	log.Printf("grpc server is running %s", l.Addr().String())
	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
