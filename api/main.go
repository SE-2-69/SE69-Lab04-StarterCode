package main

import (
	"log"

	"se69-lab04-starter/internal/config"
	"se69-lab04-starter/internal/models"
	"se69-lab04-starter/internal/routes"
)

func main() {
	// load settings
	cfg := config.Load()
	log.Printf("config: port=%s db=%s token_ttl=%s", cfg.Port, cfg.DBPath, cfg.TokenTTL)

	// open database, create tables, add sample data
	db, err := models.ConnectDB(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// register routes and start server
	r := routes.Setup(db, cfg)
	log.Printf("server: Lab Equipment Loan API listening on http://localhost:%s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
