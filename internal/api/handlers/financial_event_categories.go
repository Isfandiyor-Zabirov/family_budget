package handlers

import (
	"family_budget/internal/entities/financial_event_categories"
	"family_budget/internal/utils/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// CreateFinancialEventCategory - Создание категории финансовых событии
// @Summary Создание катгории финансовых событии
// @ID create-financial-event-category
// @Tags Категории финансовых событий
// @Produce json
// @Security     JWT
// @Param name  		 		body string true "Название категории"
// @Param description  			body string false "Описание категории"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /api/v1/financial_event_categories [post]
func CreateFinancialEventCategory(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		request financial_event_categories.FinancialEventCategories
		err     error
	)

	// TODO: check access by roleID or userID

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("CreateFinancialEventCategory handler cannot bind the request")
		resp := response.SetResponseData(request, "Неверная структура запроса", false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID

	resp, err := financial_event_categories.Create(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}
