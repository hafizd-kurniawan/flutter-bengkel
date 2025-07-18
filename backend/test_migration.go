package main

import (
	"log"
	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/database"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	log.Println("Testing PostgreSQL migrations...")
	
	// Test database connection
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL database:", err)
	}
	defer db.Close()
	
	log.Println("✅ PostgreSQL connection successful!")
	
	// Enable UUID extension first
	log.Println("Enabling UUID extension...")
	_, err = db.GetDB().Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		log.Fatal("Failed to enable UUID extension:", err)
	}
	
	// Test UUID generation
	var uuidResult string
	err = db.GetDB().Get(&uuidResult, "SELECT uuid_generate_v4()::text")
	if err != nil {
		log.Fatal("Failed to generate UUID:", err)
	}
	
	log.Printf("✅ UUID generation test successful: %s", uuidResult)
	
	// Test basic migration structure
	log.Println("Testing basic tables creation...")
	
	// Try to run just the foundation migration
	migrationSQL := `
		-- Create database schema for Workshop Management System
		-- Foundation & Security Tables with PostgreSQL + Soft Delete
		
		-- Enable UUID extension
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		
		-- Create trigger function for updating updated_at
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
		    NEW.updated_at = NOW();
		    RETURN NEW;
		END;
		$$ language 'plpgsql';
		
		-- Outlets table for multi-branch support
		CREATE TABLE IF NOT EXISTS outlets (
		    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		    name VARCHAR(255) NOT NULL,
		    address TEXT,
		    phone VARCHAR(20),
		    email VARCHAR(255),
		    is_active BOOLEAN DEFAULT TRUE,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		    deleted_at TIMESTAMP WITH TIME ZONE NULL,
		    deleted_by UUID NULL
		);
		
		-- Create trigger for updated_at
		DROP TRIGGER IF EXISTS update_outlets_updated_at ON outlets;
		CREATE TRIGGER update_outlets_updated_at BEFORE UPDATE ON outlets
		    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
		
		-- Test insert
		INSERT INTO outlets (name, address, phone, email) 
		VALUES ('Main Workshop', 'Jl. Raya No 123', '021-12345', 'main@bengkel.com')
		ON CONFLICT DO NOTHING;
	`
	
	_, err = db.GetDB().Exec(migrationSQL)
	if err != nil {
		log.Fatal("Failed to run migration:", err)
	}
	
	log.Println("✅ Basic migration successful!")
	
	// Test query
	type Outlet struct {
		ID      string `db:"id"`
		Name    string `db:"name"`
		Address string `db:"address"`
	}
	
	var outlets []Outlet
	err = db.GetDB().Select(&outlets, "SELECT id, name, address FROM outlets WHERE deleted_at IS NULL")
	if err != nil {
		log.Fatal("Failed to query outlets:", err)
	}
	
	log.Printf("✅ Found %d outlets:", len(outlets))
	for _, outlet := range outlets {
		log.Printf("   - %s: %s (%s)", outlet.ID, outlet.Name, outlet.Address)
	}
	
	log.Println("🎉 PostgreSQL with UUID and soft delete is working!")
}