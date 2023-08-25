package main

import (
	"fmt"
	"time"
)

var (
	Version   string
	BuildTime string
)

func main() {
	fmt.Printf("version=%s, time=%s\n", Version, BuildTime)
	time.Now()
}
