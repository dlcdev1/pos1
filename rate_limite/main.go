package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dlcdev1/pos1/rate_limite/limiter"
	"github.com/dlcdev1/pos1/rate_limite/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	ipLimit, err := strconv.Atoi(os.Getenv("LIMIT_IP_RPS"))
	if err != nil || ipLimit <= 0 {
		ipLimit = 5
	}

	tokenLimit, err := strconv.Atoi(os.Getenv("LIMIT_TOKEN_RPS"))
	if err != nil || tokenLimit <= 0 {
		tokenLimit = 10
	}

	blockDurationStr := os.Getenv("BLOCK_DURATION_MIN")
	blockDuration, err := time.ParseDuration(blockDurationStr)
	if err != nil || blockDuration <= 0 {
		blockDuration = 5 * time.Minute
	}

	tokenLimits := os.Getenv("TOKEN_LIMITS")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	storage := limiter.NewRedisStorage(redisAddr, redisPass, redisDB)

	l := limiter.NewLimiter(storage, ipLimit, tokenLimit, blockDuration, tokenLimits)

	r := gin.Default()

	r.Use(middleware.RateLimiterMiddleware(l))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	fmt.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
