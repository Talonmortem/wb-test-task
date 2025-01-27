package main

import (
	"fmt"
	"log"

	"github.com/Talonmortem/wb-test-task/config"
	"github.com/Talonmortem/wb-test-task/internal/handler"
	"github.com/Talonmortem/wb-test-task/internal/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg.User, cfg.Password, cfg.Host)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Pool.Close()

	h := handler.NewHandler(db, cfg.Schema)

	r := gin.Default()
	r.Any("/"+cfg.Endpoint+"/*path", h.HandleRequest)

	log.Printf("Server started on port %s", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
