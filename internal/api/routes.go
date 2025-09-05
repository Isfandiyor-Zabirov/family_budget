package api

import (
	"family_budget/internal/internal_config"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/julienschmidt/httprouter"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func Init() {
	router := gin.Default()

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
		fmt.Println("Failed to run port: ", err.Error())
		return
	}
}
