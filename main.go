package main

import (
	"Edupay/config"
	"Edupay/database"
	"Edupay/routes"
	m "Edupay/usecase/middlewares"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load konfigurasi dari file environment
	config.LoadConfig()

	// Inisialisasi koneksi database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Inisialisasi Echo framework
	e := echo.New()

	// Setup routes
	routes.Routes(e, db)
	m.LogMiddlewares(e)

	// Jalankan server
	log.Println("Starting server on port:", config.AppConfig.AppPort)
	e.Logger.Fatal(e.Start(":" + config.AppConfig.AppPort))

}
