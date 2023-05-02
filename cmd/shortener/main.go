package main

import (
	"github.com/keyjin88/shortener/internal/app"
	"log"
)

func main() {
	server := app.New()
	//api server start
	log.Panicln(server.Start())
}
