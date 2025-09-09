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
	router.Any("api/v1/ping", ping)

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
	financialEventCategories := v1.Group("/financial_event_categories")
	financialEventCategories.POST("", handlers.CreateFinancialEventCategory)
	financialEventCategories.GET("", handlers.GetFinancialEventCategoryList)
	financialEventCategories.GET("/:id", handlers.GetFinancialEventCategory)
	financialEventCategories.PUT("", handlers.UpdateFinancialEventCategory)
	financialEventCategories.DELETE("/:id", handlers.DeleteFinancialEventCategory)

	//////////////////////////////// Роли и доступы \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	roles := v1.Group("/roles")
	roles.GET("", handlers.GetRoles)
	roles.GET("/:id", handlers.GetRole)
	roles.DELETE("/:id", handlers.DeleteRole)
	roles.GET("/accesses/:role_id", handlers.GetRoleWithAccesses)
	roles.POST("", handlers.CreateRoleWithAccesses)
	roles.PUT("", handlers.UpdateRoleWithAccesses)

	//////////////////////////////// Семейство (Пользователи) \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	users := v1.Group("/users")
	users.GET("", handlers.GetUserList)
	users.POST("", handlers.CreateUser)
	users.PUT("", handlers.UpdateUser)
	users.DELETE("/:id", handlers.DeleteUser)

	//////////////////////////////// Цели \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	goals := v1.Group("/goals")
	goals.GET("", handlers.GetGoalsList)
	goals.POST("", handlers.CreateGoal)
	goals.GET("/:id", handlers.GetGoal)
	goals.PUT("", handlers.UpdateGoal)
	goals.DELETE("/:id", handlers.DeleteGoal)

	//////////////////////////////// Операции \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	transactions := v1.Group("/transactions")
	transactions.POST("", handlers.CreateTransaction)
	transactions.GET("", handlers.GetTransactionList)

	//////////////////////////////// Финансовые события \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	financialEvents := v1.Group("/financial_events")
	financialEvents.GET("", handlers.GetFinancialEventList)
	financialEvents.POST("", handlers.CreateFinancialEvent)
	financialEvents.PUT("", handlers.UpdateFinancialEvent)
	financialEvents.DELETE("/:id", handlers.DeleteFinancialEvent)
	financialEvents.GET("/:id", handlers.GetFinancialEvent)

	//////////////////////////////// Отчеты \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	reports := v1.Group("/reports")
	reports.GET("/main", handlers.GetMainReport)
	reports.GET("/graph", handlers.GetGraphReport)

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

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
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

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
