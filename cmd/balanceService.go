package main

import (
	"context"
	"github.com/ksean42/BalanceService/pkg"
	"github.com/ksean42/BalanceService/pkg/handler"
	"github.com/ksean42/BalanceService/pkg/repository"
	"github.com/ksean42/BalanceService/pkg/server"
	"github.com/ksean42/BalanceService/pkg/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Balance Service API
// @version 1.0
// @description API service to manage user balance, payments and get revenue reports
// @host localhost:8071
func main() {

	cfg := pkg.NewConfig()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	serv := &server.Server{}
	pg, err := repository.NewPostgresClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	balanceService := services.NewBalanceService(pg)
	h := handler.NewHandler(balanceService)
	go gracefulShutdown(ctx, cancel, serv)

	if err := serv.Start(cfg.Port, h.InitRouter()); err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, server *server.Server) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	<-exit
	cancel()
	log.Println("Server shutting down...")
	if err := server.Stop(ctx); err != nil {
		log.Println(err)
	}
}
