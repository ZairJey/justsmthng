package main

import (
	"awesomeProject/internal/app"
	"awesomeProject/internal/config"
	"log"
)

func main() {
	conf, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	app.Run(conf)
}
