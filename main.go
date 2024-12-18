package main

import (
	"admin-panel/middlewares"
	"admin-panel/routes"
	"admin-panel/services"
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// MongoDB bağlantısını başlat
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Servisleri başlat
	services.InitUserService(client)
	services.InitPostService(client)
	services.InitPageService(client) // Sayfa servisini başlat
	services.InitCategoryService(client)
	services.InitTagService(client)
	services.InitMediaService(client) // Medya servisini başlat

	// Gin başlat
	r := gin.Default()

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

	r.Static("/uploads", "./uploads")

	// Sunucuyu başlat
	r.Run(":8080")
}
