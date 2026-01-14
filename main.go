package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/LouisLau-art/go-translator/api"
	"github.com/LouisLau-art/go-translator/cache"
	"github.com/LouisLau-art/go-translator/config"
	"github.com/LouisLau-art/go-translator/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set Gin mode
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize components
	doubaoClient := api.NewDoubaoClient(cfg.APIKey, cfg.APIURL)
	translatorCache := cache.NewTranslatorCache(cfg.CacheTTL, cfg.CacheMaxSize)
	limiter := rate.NewLimiter(rate.Every(2*time.Second), cfg.RateLimitBurst)
	translationHandler := handlers.NewTranslationHandler(doubaoClient, translatorCache, limiter, cfg.MaxTextLength)

	// Create router
	r := gin.Default()
	r.Use(cors.Default())

	// Serve static files
	r.Static("/static", "./static")
	r.Static("/libs", "./static/libs")
	r.StaticFile("/", "./static/index.html")

	// API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/translate", translationHandler.HandleTranslate)
		apiGroup.GET("/languages", getLanguages)
		apiGroup.GET("/health", healthCheck)
	}

	// Start server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Language map for API response
var languageMap = map[string]string{
	"zh":      "中文（简体）",
	"zh-Hant": "中文（繁体）",
	"en":      "英语",
	"ja":      "日语",
	"ko":      "韩语",
	"de":      "德语",
	"fr":      "法语",
	"es":      "西班牙语",
	"it":      "意大利语",
	"pt":      "葡萄牙语",
	"ru":      "俄语",
	"th":      "泰语",
	"vi":      "越南语",
	"ar":      "阿拉伯语",
}

// getLanguages returns supported languages
func getLanguages(c *gin.Context) {
	c.JSON(200, gin.H{
		"success":   true,
		"languages": languageMap,
	})
}

// healthCheck returns server health status
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}