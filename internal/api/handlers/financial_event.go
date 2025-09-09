package handlers

import (
	"family_budget/internal/entities/financial_events"
	"family_budget/internal/utils/response"
	"family_budget/middleware"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateFinancialEvent - Создание финансовой событии
// @Summary Создание финансовой событии
// @ID create-financial-event
// @Tags Финансовая события
// @Produce json
// @Security     JWT
// @Param   category  body      financial_events.FinancialEvent  true  "Данные для создания финансовой событии"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /financial_events [post]
func CreateFinancialEvent(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		request financial_events.FinancialEvent
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEvents, middleware.CREATE, ctxData.UserID) {
		response.SetResponseData(&resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("CreateFinancialEvent handler cannot bind the request")
		response.SetResponseData(&resp, request, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID

	resp, err = financial_events.Create(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateFinancialEvent - Изменение категории финансовой событии
// @Summary Изменение финансовой событии
// @ID update-financial-event
// @Tags Финансовая события
// @Produce json
// @Security     JWT
// @Param   category  body      financial_events.FinancialEvent  true  "Данные для обновлении финансовой событии"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /financial_events [put]
func UpdateFinancialEvent(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		request financial_events.FinancialEvent
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEvents, middleware.UPDATE, ctxData.UserID) {
		response.SetResponseData(&resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("UpdateFinancialEvent handler cannot bind the request")
		response.SetResponseData(&resp, request, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID

	resp, err = financial_events.Update(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteFinancialEvent - Удаление
// @Summary      Удаление
// @Description  Удаляет по её ID
// @ID           delete-financial-event
// @Tags         Финансовая события
// @Produce      json
// @Security     JWT
// @Param        id   path      int  true  "ID для удаления"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      404  {object}  response.ResponseModel "Финансовая события не найдена"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /financial_events/{id} [delete]
func DeleteFinancialEvent(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEvents, middleware.DELETE, ctxData.UserID) {
		response.SetResponseData(&resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SetResponseData(&resp, nil, "Неверный формат ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err = financial_events.Delete(id, ctxData.FamilyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetFinancialEvent - Получении
// @Summary Получение финансового события
// @Description  Получения по её ID
// @ID           get-financial-event
// @Tags 		 Финансовые события
// @Produce      json
// @Security     JWT
// @Param        id   path      int  true  "ID для получения"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      404  {object}  response.ResponseModel "не найдена"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /financial_events/{id} [get]
func GetFinancialEvent(c *gin.Context) {
	var err error
	ctxData := getClaimsFromContext(c)
	var resp response.ResponseModel

	if !middleware.CheckAccess(middleware.FinancialEvents, middleware.READ, ctxData.UserID) {
		response.SetResponseData(
			&resp,
			financial_events.FinancialEvent{},
			"Доступ запрещен",
			false,
			0,
			0,
			0,
		)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("GetFinancialEvent handler invalid ID format: %s", idStr)
		response.SetResponseData(
			&resp,
			financial_events.FinancialEvent{},
			"Неверный формат ID",
			false,
			0,
			0,
			0,
		)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err = financial_events.Get(id, ctxData.FamilyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetFinancialEventList - Получение списка финансовых событий
// @Summary      Получение списка финансовых событий
// @Description  Возвращает постраничный список, принадлежащих семье пользователя
// @ID           get-financial-event-list
// @Tags 		 Финансовые события
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int    false  "Номер страницы" default(1)
// @Param        limit  query     int    false  "Количество элементов на странице" default(10)
// @Param        search query     string false  "Текст для поиска по названию и описанию"
// @Success      200    {object}  response.ResponseModel
// @Failure      400    {object}  response.ResponseModel "Неверные параметры запроса"
// @Failure      500    {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /financial_events [get]
func GetFinancialEventList(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		filters financial_events.Filters
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEvents, middleware.READ, ctxData.UserID) {
		response.SetResponseData(
			&resp,
			[]financial_events.FinancialEvent{},
			"Доступ запрещен",
			false,
			0,
			0,
			0,
		)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	err := c.Bind(&filters)
	if err != nil {
		log.Println("GetFinancialEventList handler cannot get params:", err.Error())
		response.SetResponseData(
			&resp,
			[]financial_events.FinancialEvent{},
			"Что-то пошло не так",
			false,
			0,
			0,
			0,
		)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	filters.FamilyID = ctxData.FamilyID

	resp, err = financial_events.GetList(filters)
	if err != nil {
		log.Printf("GetFinancialEventList handler error: %v", err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
