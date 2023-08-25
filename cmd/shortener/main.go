package main

import (
	"fmt"
	"github.com/keyjin88/shortener/internal/app"
	"github.com/keyjin88/shortener/internal/app/logger"
)

// Переменные для вывода флагов
// go build -ldflags "-X main.buildVersion=0.20.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=Iteration20" main.go
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()
	server := app.New()
	//api server start
	logger.Log.Panic(server.Start())
}

func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}
