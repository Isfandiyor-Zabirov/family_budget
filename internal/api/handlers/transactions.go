package handlers

import (
	"family_budget/internal/entities/transactions"
	"family_budget/internal/utils/response"
	"family_budget/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// CreateTransaction - Добавление операции
// @Summary Добавление членов семьи (пользователей)
// @ID create-transaction
// @Tags Операции
// @Produce json
// @Security     JWT
// @Param financial_event_id  	body integer false "ID дохода или расхода"
// @Param goal_id  				body integer false "ID цели"
// @Param amount  				body number true "Сумма"
// @Param description	  		body string false "Описание"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /transactions [post]
func CreateTransaction(c *gin.Context) {
	var (
		ctx     = getClaimsFromContext(c)
		request transactions.Transactions
		resp    response.ResponseModel
		err     error
	)

	if !middleware.CheckAccess(middleware.Transactions, middleware.CREATE, ctx.UserID) {
		response.SetResponseData(&resp, struct{}{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Println("CreateTransaction handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, struct{}{}, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.UserID = ctx.UserID
	request.FamilyID = ctx.FamilyID

	resp, statusCode, err := transactions.CreateTransaction(&request)
	if err != nil {
		c.JSON(statusCode, resp)
		return
	}

	c.JSON(statusCode, resp)
}

// GetTransactionList - Получение списка операции
// @Summary Получение списка всех членов семьи (пользователей)
// @ID get-transaction-list
// @Tags Операции
// @Produce json
// @Security     JWT
// @Param search 	    			query string 	false "Поиск"
// @Param user_id 	 				query integer  	false "ID пользователя"
// @Param financial_event_id 	 	query integer  	false "ID финансовой событии"
// @Param goal_id 	 				query integer  	false "ID цели"
// @Param date_from 	 			query string  	false "Дата с (YYYY-MM-DD)"
// @Param date_to 	 				query string  	false "Дата до (YYYY-MM-DD)"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /transactions [get]
func GetTransactionList(c *gin.Context) {
	var (
		ctx     = getClaimsFromContext(c)
		resp    response.ResponseModel
		filters transactions.Filters
		err     error
	)

	err = c.Bind(&filters)
	if err != nil {
		log.Println("GetTransactionList handler cannot bind filters:", err.Error())
		response.SetResponseData(&resp, []transactions.Response{}, "Неверный фильтр", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !middleware.CheckAccess(middleware.Transactions, middleware.READ, ctx.UserID) {
		if !middleware.CheckAccess(middleware.Transactions, middleware.CREATE, ctx.UserID) {
			response.SetResponseData(&resp, []transactions.Response{}, "Доступ запрещен", false, 0, 0, 0)
			c.JSON(http.StatusForbidden, resp)
			return
		} else {
			filters.UserID = &ctx.UserID
		}
	}

	filters.FamilyID = ctx.FamilyID
	resp, err = transactions.TransactionList(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
