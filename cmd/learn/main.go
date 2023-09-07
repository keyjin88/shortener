package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// GenSecretKey генерирует секретный ключ и возвращает его в виде Base64-строки.
func GenSecretKey(n int) (string, error) {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func main() {
	secretKey := func(n int) string {
		s, err := GenSecretKey(n)
		if err != nil {
			log.Fatal(err)
		}
		return s
	}
	fmt.Println(secretKey(16), secretKey(32))
}
