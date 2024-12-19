package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dinamik yetki kontrolü için middleware
func AuthorizeRolesMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles") // Context'ten roller dizisi alınır
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid roles format"})
			c.Abort()
			return
		}

		// Kullanıcının rolleri ile izin verilen roller arasında eşleşme kontrolü
		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					c.Next()
					return
				}
			}
		}

		// Yetkisiz erişim
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
			"roles": userRoles,
		})
		c.Abort()
	}
}
