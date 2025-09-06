package handlers

import (
	"family_budget/internal/entities/family"
	"family_budget/internal/entities/financial_event_categories"
	"family_budget/internal/utils/response"
	"family_budget/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreateFinancialEventCategory - Создание категории финансовых событии
// @Summary Создание категории финансовых событии
// @ID create-financial-event-category
// @Tags Категории финансовых событий
// @Produce json
// @Security     JWT
// @Param name  		 		body string true "Название категории"
// @Param description  			body string false "Описание категории"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /financial_event_categories [post]
func CreateFinancialEventCategory(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		request financial_event_categories.FinancialEventCategories
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEventCategories, middleware.CREATE, ctxData.UserID) {
		response.SetResponseData(resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("CreateFinancialEventCategory handler cannot bind the request")
		response.SetResponseData(resp, request, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID

	resp, err = financial_event_categories.Create(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateFinancialEventCategory - Изменение категории финансовых событии
// @Summary Изменение категории финансовых событии
// @ID update-financial-event-category
// @Tags Категории финансовых событий
// @Produce json
// @Security     JWT
// @Param name  		 		body string true "Название категории"
// @Param description  			body string false "Описание категории"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /financial_event_categories [put]
func UpdateFinancialEventCategory(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		request financial_event_categories.FinancialEventCategories
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEventCategories, middleware.UPDATE, ctxData.UserID) {
		response.SetResponseData(resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("CreateFinancialEventCategory handler cannot bind the request")
		response.SetResponseData(&resp, request, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID

	resp, err = financial_event_categories.Create(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// DeleteFinancialEventCategory - Удаление категории финансовых событий
// @Summary      Удаление категории финансовых событий
// @Description  Удаляет категорию по её ID
// @ID           delete-financial-event-category
// @Tags         Категории финансовых событий
// @Produce      json
// @Security     JWT
// @Param        id   path      int  true  "ID категории для удаления"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      404  {object}  response.ResponseModel "Категория не найдена"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /financial_event_categories/{id} [delete]
func DeleteFinancialEventCategory(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEventCategories, middleware.DELETE, ctxData.UserID) {
		response.SetResponseData(resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("DeleteFinancialEventCategory handler invalid ID format: %s", idStr)
		response.SetResponseData(&resp, nil, "Неверный формат ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err = financial_event_categories.Get(id)
	if err != nil {
		log.Printf("DeleteFinancialEventCategory handler fec id not found: %d", id)
		response.SetResponseData(&resp, nil, "ID финансовой категории не найдена", false, 0, 0, 0)
		c.JSON(http.StatusNotFound, resp)
		return
	}

	familyID := ctxData.FamilyID

	_, err = family.Get(familyID)
	if err != nil {
		log.Printf("DeleteFinancialEventCategory handler family id not found: %d", familyID)
		response.SetResponseData(&resp, nil, "ID семьи не найдена", false, 0, 0, 0)
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp, err = financial_event_categories.Delete(id, familyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetFinancialEventCategory - Получении категории финансовых событий
// @Summary      Получения категории финансовых событий
// @Description  Получения категории по её ID
// @ID           get-financial-event-category
// @Tags         Категории финансовых событий
// @Produce      json
// @Security     JWT
// @Param        id   path      int  true  "ID категории для получения"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      404  {object}  response.ResponseModel "Категория не найдена"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /financial_event_categories/{id} [get]
func GetFinancialEventCategory(c *gin.Context) {
	var err error
	ctxData := getClaimsFromContext(c)
	var resp response.ResponseModel

	if !middleware.CheckAccess(middleware.FinancialEventCategories, middleware.UPDATE, ctxData.UserID) {
		response.SetResponseData(resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("DeleteFinancialEventCategory handler invalid ID format: %s", idStr)
		response.SetResponseData(&resp, nil, "Неверный формат ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	fec, err := financial_event_categories.Get(id)
	if err != nil {
		response.SetResponseData(&resp, nil, "Финансовая категория не найдена", false, 0, 0, 0)
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp.Data = fec

	c.JSON(http.StatusOK, resp)
}

// GetFinancialEventCategoryList - Получение списка категорий финансовых событий
// @Summary      Получение списка категорий финансовых событий
// @Description  Возвращает постраничный список категорий, принадлежащих семье пользователя
// @ID           get-financial-event-category-list
// @Tags         Категории финансовых событий
// @Produce      json
// @Security     JWT
// @Param        page   query     int    false  "Номер страницы" default(1)
// @Param        limit  query     int    false  "Количество элементов на странице" default(10)
// @Param        search query     string false  "Текст для поиска по названию и описанию"
// @Success      200    {object}  response.PaginatedResponse
// @Failure      400    {object}  response.ResponseModel "Неверные параметры запроса"
// @Failure      500    {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /financial_event_categories [get]
func GetFinancialEventCategoryList(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		filters financial_event_categories.Filters
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.FinancialEventCategories, middleware.UPDATE, ctxData.UserID) {
		response.SetResponseData(&resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.SetResponseData(resp, nil, "Неверный формат номера страницы", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		response.SetResponseData(resp, nil, "Неверный формат лимита на страницу", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	filters.CurrentPage = page
	filters.PageLimit = limit
	filters.FamilyID = ctxData.FamilyID

	searchStr := c.Query("search")
	if searchStr != "" {
		filters.Search = &searchStr
	}

	resp, err = financial_event_categories.GetList(filters)
	if err != nil {
		log.Printf("GetFinancialEventCategoryList handler error: %v", err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
