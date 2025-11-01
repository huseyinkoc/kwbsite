package middlewares

import (
	"admin-panel/configs"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Yeni JWT Claims yapısı
type Claims struct {
	UserID            string    `json:"userID"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PreferredLanguage string    `json:"preferred_language"`
	Roles             []string  `json:"roles"`
	LastLogin         time.Time `json:"last_login,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT doğrulama middleware'i
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization başlığını kontrol et
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// 'Bearer ' kısmını temizle
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Token'ı parse et ve doğrula
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Signing method kontrolü (önemli)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			// jwt secret'i merkezi configs üzerinden al
			return []byte(configs.GetJWTSecret()), nil
		}, jwt.WithLeeway(5*time.Second))

		if err != nil || !token.Valid {
			log.Printf("JWT verify error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Claims validation ekle
		if claims.UserID == "" || len(claims.Roles) == 0 {
			log.Printf("JWT invalid claims: missing required fields")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Kullanıcı bilgilerini context'e ekle
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("preferred_language", claims.PreferredLanguage)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// Basit in-memory per-IP rate limiter (prod için Redis/cluster-safe önerilir)
var (
	visitors   = make(map[string]*visitor)
	visitorsMu sync.Mutex
)

type visitor struct {
	limiter  *tokenBucket
	lastSeen time.Time
}

// Basit token-bucket implementasyonu (dışa bağımlılık yok)
type tokenBucket struct {
	mu     sync.Mutex
	tokens float64
	rate   float64 // tokens per second
	burst  float64
	last   time.Time
}

func newTokenBucket(ratePerSec float64, burst int) *tokenBucket {
	return &tokenBucket{
		tokens: float64(burst),
		rate:   ratePerSec,
		burst:  float64(burst),
		last:   time.Now(),
	}
}

func (tb *tokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()
	tb.last = now

	// refill
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.burst {
		tb.tokens = tb.burst
	}

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	return false
}

func addVisitor(ip string, rps float64, burst int) *tokenBucket {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		tb := newTokenBucket(rps, burst)
		visitors[ip] = &visitor{limiter: tb, lastSeen: time.Now()}
		return tb
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func getVisitorLimiter(ip string) *tokenBucket {
	// örnek: 1 request/sec, burst 5 (uyarlayın)
	return addVisitor(ip, 1.0, 5)
}

func cleanupVisitors(interval time.Duration, olderThan time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			visitorsMu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > olderThan {
					delete(visitors, ip)
				}
			}
			visitorsMu.Unlock()
		}
	}()
}

// RateLimitMiddleware returns a gin middleware that limits requests per IP.
// Note: production için Redis-based limiter kullanın (distributed).
func RateLimitMiddleware() gin.HandlerFunc {
	// başlat: map temizleyici
	cleanupVisitors(1*time.Minute, 5*time.Minute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := getVisitorLimiter(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
