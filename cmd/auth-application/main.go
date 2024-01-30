package main

import (
	"auth-application/internal/app"
	"auth-application/internal/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name)
	fmt.Printf(connectionString)
	application := app.New(cfg.Grpc.Port, connectionString, cfg.Secret, cfg.TokenTTL)
	go func() {
		application.GRPCServer.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
}
