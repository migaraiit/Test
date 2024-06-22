package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type video_lessons struct {
	ID  int    `json:"id" gorm:"primarykey"`
	URL string `json:"url"`
}

func main() {
	dsn := "host=localhost user=root password=password sslmode=disable dbname=studyflow"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("database connection failed", err)
	} else {
		log.Println("database connection success")
	}

	db.Table("video_lessons")
	db.AutoMigrate(&video_lessons{})

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/video_lessons/:id", getVideoURL) // Modified route

	r.Run()
}

func getVideoURL(c *gin.Context) {
	var video video_lessons
	id := c.Param("id") // Get the ID from the URL parameter

	if err := db.First(&video, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": video.URL})
}
