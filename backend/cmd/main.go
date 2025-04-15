package main

import (
	"fmt"
	"log"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/polonkoevv/ethcourse/internal/handler"
	"github.com/polonkoevv/ethcourse/internal/service"
	"github.com/polonkoevv/ethcourse/internal/storage/postgres"
)

func main() {
	// Подключение к локальному узлу IPFS
	sh := shell.NewShell("localhost:5001")

	pg, err := postgres.NewPostgres("0.0.0.0", "5432", "postgres", "postgres", "ipfs")
	if err != nil {
		log.Fatal(err)
	}

	srv := service.NewService(sh, pg)

	h := handler.NewHandler(srv)

	r := h.CreateRouter()

	// // Запуск HTTP сервера
	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8000", r))

	fmt.Println(sh.Pins())
}
