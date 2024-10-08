package main

import (
	"ClipsArchiver/internal/config"
	"ClipsArchiver/internal/db"
	"ClipsArchiver/internal/rest/clips"
	"ClipsArchiver/internal/rest/files"
	"ClipsArchiver/internal/rest/legends"
	"ClipsArchiver/internal/rest/maps"
	"ClipsArchiver/internal/rest/tags"
	"ClipsArchiver/internal/rest/transcodeRequests"
	"ClipsArchiver/internal/rest/trimRequests"
	"ClipsArchiver/internal/rest/users"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const logFileLocation = "clipsarchiver.log"

var logger *slog.Logger

func main() {

	options := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	file, err := os.OpenFile(logFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to get log file handle: %s", err.Error())
	}

	var handler slog.Handler = slog.NewJSONHandler(file, options)
	logger = slog.New(handler)

	err = db.SetupDb(logger)
	if err != nil {
		log.Fatalf("Failed to setup database: %s", err.Error())
	}

	router := gin.Default()

	router.GET("/clips/:clipId", clips.Get)
	router.PUT("/clips/:clipId", clips.Update)
	router.DELETE("/clips/:clipId", clips.Delete)
	router.GET("/clips/date/:date", clips.GetForDate)
	router.GET("/clips/filename/:filename", clips.GetByFilename)
	router.GET("/users", users.GetAll)
	router.GET("/tags", tags.GetAll)
	router.GET("/maps", maps.GetAll)
	router.GET("/legends", legends.GetAll)
	router.GET("/clips/queue", transcodeRequests.GetAll)
	router.GET("/clips/queue/:clipId", transcodeRequests.GetById)
	router.GET("/clips/download/:clipId", files.DownloadClipById)
	router.GET("/clips/download/thumbnail/:clipId", files.DownloadClipThumbnailById)
	router.POST("/clips/upload/:ownerId", files.UploadClip)
	router.POST("/clips/trim/:clipId", trimRequests.Create)
	router.GET("clips/trim/:clipId", trimRequests.GetByClipId)
	//router.POST("/clips/combine/:firstId/:secondId", files.CombineClips)
	router.StaticFS("/clips/archive", http.Dir(config.GetOutputPath()))
	router.StaticFS("/clips/thumbnails", http.Dir(config.GetThumbnailsPath()))
	router.StaticFS("/resources", http.Dir(config.GetResourcesPath()))

	routerErr := router.Run()
	if routerErr != nil {
		log.Fatalf("Failed to start gin router: %s", routerErr)
	}
}
