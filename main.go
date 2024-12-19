package main

import (
	"admin-panel/configs"
	"admin-panel/middlewares"
	"admin-panel/routes"
	"admin-panel/services"
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Println(".env dosyası yüklenemedi, ortam değişkenleri kullanılacak")
	}
}

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Veritabanı bağlantısını başlat
	if err := configs.Init(); err != nil {
		log.Fatalf("Veritabanı başlatılamadı: %v", err)
	}

	// Uygulama sonlandığında bağlantıyı kapat
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
		if err := configs.DB.Disconnect(ctx); err != nil {
			log.Fatalf("MongoDB bağlantısı kapatılamadı: %v", err)
		}
	}()

	// Servisleri başlat
	services.InitUserService(configs.DB)
	services.InitPostService(configs.DB)
	services.InitPageService(configs.DB)
	services.InitCategoryService(configs.DB)
	services.InitTagService(configs.DB)
	services.InitMediaService(configs.DB)
	services.InitCommentService(configs.DB)
	services.InitNotificationService(configs.DB)
	services.InitRolesService(configs.DB)

	log.Println("Tüm servisler başarıyla başlatıldı.")

	// Gin başlat
	r := gin.Default()

	// Dil ve SEO dostu rotalar
	r.GET("/:lang/:slug", func(c *gin.Context) {
		lang := c.Param("lang")
		slug := c.Param("slug")

		// Burada dil ve slug'a göre içeriği getirme işlemi yapılacak
		// Örneğin:
		// content := getContentByLangAndSlug(lang, slug)

		c.JSON(http.StatusOK, gin.H{
			"language": lang,
			"slug":     slug,
			// "content":  content,
		})
	})

	// Logger middleware, tüm rotalar için etkin
	r.Use(middlewares.LoggerMiddleware())

	// Rotaları yükle
	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	routes.PostRoutes(r)
	routes.PageRoutes(r) // Sayfa rotalarını yükle
	routes.CategoryRoutes(r)
	routes.TagRoutes(r)
	routes.MediaRoutes(r) // Medya rotalarını ekle
	routes.RegisterCommentRoutes(r)
	routes.RegisterNotificationRoutes(r)
	routes.RoleRoutes(r)

	r.Static("/uploads", "./uploads")

	// Sunucuyu başlat
	r.Run(":8080")
}
