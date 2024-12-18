package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dinamik yetki kontrolü için yapı
func AuthorizeRolesMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role") // Context'ten rol bilgisi alınır
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}
		// Yetkisiz erişim
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
			"role":  role,
		})
		c.Abort()
	}
}
