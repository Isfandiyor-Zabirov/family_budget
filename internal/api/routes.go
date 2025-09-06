package api

import (
	"family_budget/internal/api/handlers"
	"family_budget/internal/internal_config"
	"family_budget/internal/logger"
	"family_budget/middleware"
	"family_budget/pkg/database"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/julienschmidt/httprouter"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {

	jwtMiddleware := &middleware.GinJWTMiddleware{
		Realm:          internal_config.InternalConfigs.Server.Realm,
		AccessKey:      []byte(internal_config.InternalConfigs.Application.AccessKey),
		RefreshKey:     []byte(internal_config.InternalConfigs.Application.RefreshKey),
		AccessTimeout:  time.Second * time.Duration(internal_config.InternalConfigs.Application.AccessTknTimeout),
		RefreshTimeout: time.Second * time.Duration(internal_config.InternalConfigs.Application.RefreshTknTimeout),
		MaxRefresh:     time.Hour * 24,
		Authenticator:  middleware.Authenticator,
		PayloadFunc:    middleware.Payload,
		DB:             database.Postgres(),
	}

	router := gin.Default()

	customLogger := logger.GetLogger()

	gin.DefaultWriter = io.MultiWriter(customLogger, os.Stdout)

	logger.FormatLogger(router)

	router.Use(CORSMiddleware())

	router.Use(gin.Recovery())

	router.POST("api/v1/register", handlers.Register)
	router.POST("api/v1/login", jwtMiddleware.LoginHandler)
	router.GET("api/v1/refresh", jwtMiddleware.RefreshToken)

	v1 := router.Group("api/v1")
	v1.Use(jwtMiddleware.MiddlewareFunc())

	v1.GET("/get_me", handlers.GetMe)

	//////////////////////////////// Категории финансовых событий \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	financialEventCategories := v1.Group("financial_event_categories")
	financialEventCategories.POST("/financial-categories", handlers.CreateFinancialEventCategory)
	financialEventCategories.GET("/financial-categories", handlers.GetFinancialEventCategoryList)
	financialEventCategories.GET("/financial-categories/{id}", handlers.GetFinancialEventCategory)
	financialEventCategories.PUT("/financial-categories", handlers.UpdateFinancialEventCategory)
	financialEventCategories.DELETE("/financial-categories/{id}", handlers.DeleteFinancialEventCategory)

	//////////////////////////////// Семья \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	/*	family := v1.Group("financial_event_categories")
		family.POST("/family", handlers.CreateFamily)
		family.GET("/family/{id}", handlers.GetFamily)
		family.PUT("/family", handlers.UpdateFamily)
		family.DELETE("/family/{id}", handlers.DeleteFamily)*/
	financialEventCategories.POST("", handlers.CreateFinancialEventCategory)
	financialEventCategories.PUT("", handlers.UpdateFinancialEventCategory)

	accounts := gin.Accounts{
		"sakhi":  "family_budget",
		"eraj":   "family_budget",
		"ismoil": "family_budget",
	}

	router.GET("/swagger/*any", gin.BasicAuth(accounts), ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "reason": "Страница не найдена"})
	})

	err := router.Run(fmt.Sprint(":", internal_config.InternalConfigs.Server.PortRun))
	if err != nil {
		fmt.Printf("Failed to start the app at the port: %d. The error is: %s", internal_config.InternalConfigs.Server.PortRun, err.Error())
		return
	}
}

// CORSMiddleware controls course middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Start-Encoding, X-CSRF-Token, Authorization, Refresh-Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		c.Writer.Header().Set("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
