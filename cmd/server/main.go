package main

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/config"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/logger"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	l := logger.New(cfg.LogLevel)
	l.Info("Starting book library API service")
}
