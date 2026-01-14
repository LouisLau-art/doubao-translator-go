package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/LouisLau-art/go-translator/api"
	"github.com/LouisLau-art/go-translator/cache"
)

// TranslationHandler handles translation requests
type TranslationHandler struct {
	apiClient *api.DoubaoClient
	cache     *cache.TranslatorCache
	limiter   *rate.Limiter
	maxLength int
}

// NewTranslationHandler creates a new translation handler
func NewTranslationHandler(apiClient *api.DoubaoClient, cache *cache.TranslatorCache, limiter *rate.Limiter, maxLength int) *TranslationHandler {
	return &TranslationHandler{
		apiClient: apiClient,
		cache:     cache,
		limiter:   limiter,
		maxLength: maxLength,
	}
}

// HandleTranslate processes translation requests
func (h *TranslationHandler) HandleTranslate(c *gin.Context) {
	// Check rate limit
	if !h.limiter.Allow() {
		c.JSON(429, gin.H{
			"success": false,
			"error":   "请求过于频繁，请稍后再试",
		})
		return
	}

	// Parse request
	var req api.TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Request bind error: %v", err)
		c.JSON(400, gin.H{
			"success": false,
			"error":   "请求格式错误: " + err.Error(),
		})
		return
	}

	log.Printf("Translation request: text length=%d, source=%s, target=%s",
		len(req.Text), req.Source, req.Target)

	// Validate text length
	if len(req.Text) > h.maxLength {
		c.JSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("文本长度超过限制（最大%d字符）", h.maxLength),
		})
		return
	}

	// Check cache
	cacheKey := cache.GetCacheKey(req.Text, req.Source, req.Target)
	if cached, ok := h.cache.Get(cacheKey); ok {
		log.Printf("Cache hit for key: %s", cacheKey)
		c.JSON(200, gin.H{
			"success": true,
			"text":    cached,
			"cached":  true,
		})
		return
	}

	// Split text into chunks for long documents
	chunks := smartSplit(req.Text, 800)
	log.Printf("Split text into %d chunks", len(chunks))

	results := make([]string, len(chunks))

	// Process each chunk
	for i, chunk := range chunks {
		result, err := h.apiClient.Translate(chunk, req.Source, req.Target)
		if err != nil {
			log.Printf("API call error for chunk %d: %v", i, err)
			c.JSON(500, gin.H{
				"success": false,
				"error":   "翻译失败: " + err.Error(),
			})
			return
		}
		results[i] = result
	}

	// Combine results
	finalText := strings.Join(results, "\n")

	// Save to cache
	if err := h.cache.Set(cacheKey, finalText); err != nil {
		log.Printf("Cache set error: %v", err)
	}

	c.JSON(200, gin.H{
		"success": true,
		"text":    finalText,
		"cached":  false,
	})
}

// smartSplit splits text into chunks, trying to preserve paragraph boundaries
func smartSplit(text string, maxChars int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}

	// Try to split at paragraph boundaries first
	paragraphs := strings.Split(text, "\n\n")
	var chunks []string
	currentChunk := ""

	for _, paragraph := range paragraphs {
		// If paragraph itself is too long, split it further
		if len(paragraph) > maxChars {
			if currentChunk != "" {
				chunks = append(chunks, currentChunk)
				currentChunk = ""
			}

			// Split long paragraph by sentences or fixed size
			for i := 0; i < len(paragraph); i += maxChars {
				end := i + maxChars
				if end > len(paragraph) {
					end = len(paragraph)
				}
				chunks = append(chunks, paragraph[i:end])
			}
		} else {
			// Check if adding this paragraph would exceed maxChars
			if len(currentChunk)+len(paragraph)+2 > maxChars && currentChunk != "" {
				chunks = append(chunks, currentChunk)
				currentChunk = paragraph
			} else {
				if currentChunk != "" {
					currentChunk += "\n\n"
				}
				currentChunk += paragraph
			}
		}
	}

	// Add the last chunk if not empty
	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}