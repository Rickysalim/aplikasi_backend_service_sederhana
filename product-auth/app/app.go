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

func Start() { 
	config.LoadConfig(".")
	router := gin.Default()
	dbClient := GetDbClient()

	authRepositoryDb := repositories.NewAuthRepository(dbClient)
	
	authService := services.NewAuthService(authRepositoryDb, domain.GetRolePermissions())

	ah := handlers.AuthHandler{
		AuthService: authService,
	}

	auth := router.Group("v1/")

	auth.POST("/auth/login",ah.Login)
	auth.GET("/auth/verify", ah.Verify)
	auth.POST("/auth/refresh", ah.Refresh)
    
	router.Run(":5001")
}