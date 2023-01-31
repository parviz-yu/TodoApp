package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pyuldashev912/todoapp/internal/app/apiserver"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	config := apiserver.NewConfig()
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
