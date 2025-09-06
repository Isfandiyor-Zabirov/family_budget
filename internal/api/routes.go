package api

import (
	"family_budget/internal/api/handlers"
	"family_budget/internal/internal_config"
	"family_budget/internal/logger"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/julienschmidt/httprouter"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {
	router := gin.Default()

	customLogger := logger.GetLogger()

	gin.DefaultWriter = io.MultiWriter(customLogger, os.Stdout)

	logger.FormatLogger(router)

	router.Use(CORSMiddleware())

	router.Use(gin.Recovery())

	accounts := gin.Accounts{
		"sakhi":  "family_budget",
		"eraj":   "family_budget",
		"ismoil": "family_budget",
	}

	v1 := router.Group("api/v1")

	// TODO: add middleware to v1

	//////////////////////////////// Категории финансовых событий \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	financialEventCategories := v1.Group("financial_event_categories")
	financialEventCategories.POST("", handlers.CreateFinancialEventCategory)

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
