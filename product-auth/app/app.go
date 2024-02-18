package app

import (
	"fmt"
	"product_auth/domain"
	"product_auth/handlers"
	"product_auth/repositories"
	"product_auth/services"
	"product_auth/utils/config"
	"time"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func GetDbClient() *sqlx.DB {
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",viper.Get("DB_USER"), viper.Get("DB_PASSWORD"), viper.Get("DB_NAME"));
    client, err := sqlx.Open("postgres",dataSource)
    if err != nil {
		panic(err);
	}

	client.SetConnMaxLifetime((time.Minute * 3))
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Service, X-Api-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Start() { 
	config.LoadConfig(".")
	router := gin.Default()
	router.Use(CORS())
	dbClient := GetDbClient()

	authRepositoryDb := repositories.NewAuthRepository(dbClient)
	
	authService := services.NewAuthService(authRepositoryDb, domain.GetRolePermissions())

	ah := handlers.AuthHandler{
		AuthService: authService,
	}
	
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"message": "Not Found",
		})
	})

	auth := router.Group("v1/")
	auth.GET("/auth", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome To Our Service",
        })
	})
	auth.POST("/auth/login",ah.Login)
	auth.GET("/auth/verify", ah.Verify)
	auth.POST("/auth/refresh", ah.Refresh)
	

	router.Run(":5001")
}