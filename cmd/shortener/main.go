package main

import (
	"github.com/keyjin88/shortener/internal/app/api"
	"log"
)

func main() {
	server := api.New()
	//api server start
	log.Fatal(server.Start())
}
