package main

import (
	"fmt"

	"github.com/keyjin88/shortener/internal/app"
	"github.com/keyjin88/shortener/internal/app/logger"
)

// Переменные для вывода флагов
// go build -ldflags "-X main.buildVersion=0.20.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'
// -X main.buildCommit=Iteration20" main.go.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()
	server := app.New()
	err := server.Start()
	if nil != err {
		logger.Log.Info("Error starting api server")
		return
	}
}

func printBuildInfo() {
	const NotAssigned = "N/A"
	if buildVersion == "" {
		buildVersion = NotAssigned
	}
	if buildDate == "" {
		buildDate = NotAssigned
	}
	if buildCommit == "" {
		buildCommit = NotAssigned
	}

	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}
