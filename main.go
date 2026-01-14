package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// --- 结构体定义 ---

type TranslateRequest struct {
	Text   string `json:"text" binding:"required"`
	Source string `json:"source"`
	Target string `json:"target" binding:"required"`
}

type DoubaoRequest struct {
	Model string               `json:"model"`
	Input []DoubaoInputMessage `json:"input"`
}

type DoubaoInputMessage struct {
	Role    string          `json:"role"`
	Content []DoubaoContent `json:"content"`
}

type DoubaoContent struct {
	Type               string           `json:"type"`
	Text               string           `json:"text"`
	TranslationOptions *TranslationOpts `json:"translation_options,omitempty"`
}

type TranslationOpts struct {
	SourceLanguage string `json:"source_language,omitempty"`
	TargetLanguage string `json:"target_language"`
}

type DoubaoResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type CacheItem struct {
	Value     string
	Timestamp time.Time
}

// --- 全局变量 ---

var (
	apiKey      string
	apiURL      string
	cache       sync.Map
	cacheTTL    = 3600 * time.Second
	limiter     *rate.Limiter
	maxTextLen  = 5000
	httpClient  = &http.Client{Timeout: 30 * time.Second}

	languageMap = map[string]string{
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
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	apiKey = os.Getenv("ARK_API_KEY")
	if apiKey == "" {
		log.Fatal("ARK_API_KEY not set")
	}

	apiURL = os.Getenv("ARK_API_URL")
	if apiURL == "" {
		apiURL = "https://ark.cn-beijing.volces.com/api/v3/responses"
	}

	limiter = rate.NewLimiter(rate.Every(2*time.Second), 30)
	go cleanCache()
}

func main() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/static", "./static")
	r.Static("/libs", "./static/libs")
	r.StaticFile("/", "./static/index.html")

	api := r.Group("/api")
	{
		api.POST("/translate", handleTranslate)
		api.GET("/languages", getLanguages)
		api.GET("/health", healthCheck)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// 以下函数保持不变（已修复）
// handleTranslate、callDoubaoAPI、smartSplit、getLanguages、healthCheck、
// getCacheKey、getFromCache、saveToCache、cleanCache

func handleTranslate(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(429, gin.H{
			"success": false,
			"error":   "请求过于频繁，请稍后再试",
		})
		return
	}

	var req TranslateRequest
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

	if len(req.Text) > maxTextLen {
		c.JSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("文本长度超过限制（最大%d字符）", maxTextLen),
		})
		return
	}

	cacheKey := getCacheKey(req.Text, req.Source, req.Target)
	if cached, ok := getFromCache(cacheKey); ok {
		log.Printf("Cache hit for key: %s", cacheKey)
		c.JSON(200, gin.H{
			"success": true,
			"text":    cached,
			"cached":  true,
		})
		return
	}

	chunks := smartSplit(req.Text, 800)
	log.Printf("Split text into %d chunks", len(chunks))
	
	results := make([]string, len(chunks))
	
	for i, chunk := range chunks {
		result, err := callDoubaoAPI(chunk, req.Source, req.Target)
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

	finalText := strings.Join(results, "\n")
	saveToCache(cacheKey, finalText)

	c.JSON(200, gin.H{
		"success": true,
		"text":    finalText,
		"cached":  false,
	})
}

// 只保留一个 callDoubaoAPI 函数
// 新增响应结构体
type DoubaoNewResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Status string `json:"status"`
}

// 修复后的 callDoubaoAPI 函数
func callDoubaoAPI(text, source, target string) (string, error) {
	opts := &TranslationOpts{
		TargetLanguage: target,
	}
	if source != "" {
		opts.SourceLanguage = source
	}

	reqBody := DoubaoRequest{
		Model: "doubao-seed-translation-250915",
		Input: []DoubaoInputMessage{
			{
				Role: "user",
				Content: []DoubaoContent{
					{
						Type:               "input_text",
						Text:               text,
						TranslationOptions: opts,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("API Response Status: %d", resp.StatusCode)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	// 尝试新格式解析
	var newResult DoubaoNewResponse
	if err := json.Unmarshal(body, &newResult); err == nil && newResult.Status == "completed" {
		for _, output := range newResult.Output {
			if output.Type == "message" && output.Role == "assistant" {
				for _, content := range output.Content {
					if content.Type == "output_text" {
						return content.Text, nil
					}
				}
			}
		}
		log.Printf("New format parsed successfully")
		return "", fmt.Errorf("no output_text in new format")
	}

	// 回退到旧格式
	var oldResult DoubaoResponse
	if err := json.Unmarshal(body, &oldResult); err == nil {
		if len(oldResult.Choices) > 0 {
			return oldResult.Choices[0].Message.Content, nil
		}
	}
	
	log.Printf("Response body: %s", string(body))
	return "", fmt.Errorf("无法解析API响应")
}

// 修复 smartSplit 函数
func smartSplit(text string, maxChars int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}

	var chunks []string
	for i := 0; i < len(text); i += maxChars {
		end := i + maxChars
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	return chunks
}

func getLanguages(c *gin.Context) {
	c.JSON(200, gin.H{
		"success":   true,
		"languages": languageMap,
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}

func getCacheKey(text, source, target string) string {
	data := fmt.Sprintf("%s:%s:%s", source, target, text)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func getFromCache(key string) (string, bool) {
	if val, ok := cache.Load(key); ok {
		item := val.(CacheItem)
		if time.Since(item.Timestamp) < cacheTTL {
			return item.Value, true
		}
		cache.Delete(key)
	}
	return "", false
}

func saveToCache(key, value string) {
	cache.Store(key, CacheItem{
		Value:     value,
		Timestamp: time.Now(),
	})
}

func cleanCache() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cache.Range(func(key, value interface{}) bool {
			item := value.(CacheItem)
			if time.Since(item.Timestamp) > cacheTTL {
				cache.Delete(key)
			}
			return true
		})
	}
}