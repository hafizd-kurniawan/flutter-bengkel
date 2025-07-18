package main

import (
	"log"
	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/database"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	log.Println("Testing PostgreSQL connection...")
	log.Printf("Database config: Host=%s, Port=%s, User=%s, Database=%s", 
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name)
	
	// Test database connection
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL database:", err)
	}
	defer db.Close()
	
	// Test basic query
	var version string
	err = db.GetDB().Get(&version, "SELECT version()")
	if err != nil {
		log.Fatal("Failed to query database:", err)
	}
	
	log.Println("✅ PostgreSQL connection successful!")
	log.Printf("Database version: %s", version)
	
	// Test UUID extension
	var uuidResult string
	err = db.GetDB().Get(&uuidResult, "SELECT uuid_generate_v4()::text")
	if err != nil {
		log.Fatal("Failed to generate UUID:", err)
	}
	
	log.Printf("✅ UUID generation test successful: %s", uuidResult)
}