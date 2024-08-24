package main

import (
	"Teller/internal/config"
	"Teller/internal/server"
)

func main() {
	cfg := config.ParseFlags()

	srv := server.New(cfg)
	srv.Start()
}
