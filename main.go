package main

import (
	"database/sql"
	_ "log"
	"os"

	"case/config"
	"case/handlers"
	"case/repositories"
	"case/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, user, password, dbname)

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

	repo := repositories.NewSongRepository(db)

	service := services.NewSongService(repo, cfg.ApiUrl)

	handler := handlers.NewSongHandler(service, log)

	r := gin.Default()
	r.GET("/songs", handler.GetSongs)
	r.GET("/songs/:id/lyrics", handler.GetSongLyrics)
	r.DELETE("/songs/:id", handler.DeleteSong)
	r.PUT("/songs/:id", handler.UpdateSong)
	r.POST("/songs", handler.AddSong)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
