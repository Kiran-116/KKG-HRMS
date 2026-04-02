package middleware

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configures CORS middleware
func SetupCORS() gin.HandlerFunc {
	// Strict allowlist without wildcards; compare normalized origins exactly
	allowed := map[string]struct{}{
		"http://localhost:3000":       {},
		"http://localhost:5173":       {},
		"https://kkg-hrms.vercel.app": {},
	}

	normalize := func(origin string) string {
		u, err := url.Parse(strings.TrimSpace(origin))
		if err != nil || u.Scheme == "" || u.Host == "" {
			return ""
		}
		// Keep scheme://host[:port] only
		host := u.Host
		return u.Scheme + "://" + host
	}

	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	cfg.AllowOriginFunc = func(origin string) bool {
		n := normalize(origin)
		if n == "" {
			return false
		}
		_, ok := allowed[n]
		return ok
	}

	return cors.New(cfg)
}
