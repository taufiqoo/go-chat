package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(redisClient *redis.Client, maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.redisClient == nil {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		ctx := context.Background()

		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.window)
		}

		if count > int64(rl.maxRequests) {
			ttl, _ := rl.redisClient.TTL(ctx, key).Result()

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":               "Rate limit exceeded",
				"message":             fmt.Sprintf("Anda sudah melebihi batas hit endpoint. IP: %s telah melakukan %d request dalam 1 jam terakhir.", clientIP, count),
				"ip":                  clientIP,
				"limit":               rl.maxRequests,
				"window":              rl.window.String(),
				"retry_after_seconds": int(ttl.Seconds()),
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.maxRequests-int(count)))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rl.window).Unix()))

		c.Next()
	}
}

func (rl *RateLimiter) LimitByIP(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rl.redisClient, maxRequests, window)
	return limiter.Middleware()
}
