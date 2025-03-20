package main

import (
	"database/sql"
	swaggerFiles "github.com/swaggo/files"
	_ "log"
	"os"

	"case/config"
	"case/handlers"
	"case/repositories"
	"case/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	_ "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.LoadConfig()

	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database connection:", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Info("Successfully connected to the database")

	repo := repositories.NewSongRepository(db, log)

	service := services.NewSongService(repo, cfg.ApiUrl)

	handler := handlers.NewSongHandler(service, log)

	r := gin.Default()

	// SWAGGER
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Наши маршруты
	r.GET("/songs", handler.GetSongs)
	r.GET("/songs/:id/lyrics", handler.GetSongLyrics)
	r.DELETE("/songs/:id", handler.DeleteSong)
	r.PUT("/songs/:id", handler.UpdateSong)
	r.POST("/songs", handler.AddSong)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
