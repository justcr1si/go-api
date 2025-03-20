package handlers

import (
	"case/models"
	"net/http"
	"strconv"
	"time"

	"case/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongHandler struct {
	service *services.SongService
	log     *logrus.Logger
}

func NewSongHandler(service *services.SongService, log *logrus.Logger) *SongHandler {
	return &SongHandler{service: service, log: log}
}

func (h *SongHandler) GetSongs(c *gin.Context) {
	filter := make(map[string]string)
	if group := c.Query("group"); group != "" {
		filter["group"] = group
	}
	if song := c.Query("song"); song != "" {
		filter["song"] = song
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	songs, err := h.service.GetSongs(filter, page, limit)
	if err != nil {
		h.log.Errorf("Failed to get songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func (h *SongHandler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	lyrics, err := h.service.GetSongLyrics(id, page, limit)
	if err != nil {
		h.log.Errorf("Failed to get lyrics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lyrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteSong(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted"})
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	song.ID = id

	if song.ReleaseDate == "" {
		song.ReleaseDate = time.Now().Format("02.01.2006")
	}

	if err := h.service.UpdateSong(&song); err != nil {
		h.log.Errorf("Failed to update song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated", "song": song})
}

func (h *SongHandler) AddSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if song.Group == "" || song.Song == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields 'group' and 'song' are required"})
		return
	}

	if song.ReleaseDate == "" {
		song.ReleaseDate = time.Now().Format("02.01.2006")
	}
	if song.Text == "" {
		song.Text = ""
	}
	if song.Link == "" {
		song.Link = ""
	}

	if err := h.service.AddSong(&song); err != nil {
		h.log.Errorf("Failed to add song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song added", "song": song})
}
