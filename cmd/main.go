package main

import (
	"avito_test_task/pkg"
	"avito_test_task/pkg/handler"
	"avito_test_task/pkg/repository"
	"avito_test_task/pkg/server"
	"avito_test_task/pkg/services"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
