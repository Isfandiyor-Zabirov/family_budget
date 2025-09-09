package handlers

import (
	"family_budget/internal/entities/goals"
	"family_budget/internal/utils/response"
	"family_budget/middleware"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateGoal Создание новой цели
// @Summary      Создание новой цели
// @Description  Создает новую финансовую цель для семьи
// @ID           create-goal
// @Tags         Цели
// @Accept       json
// @Produce      json
// @Security     JWT
// @Param        {object}  body      goals.Goals  true  "Данные для создания цели"
// @Success      201   {object}  response.ResponseModel
// @Failure      400   {object}  response.ResponseModel "Неверные входные данные"
// @Failure      500   {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /goals [post]
func CreateGoal(c *gin.Context) {
	var (
		request goals.Goals
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Goals, middleware.CREATE, ctxData.UserID) {
		response.SetResponseData(&resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("CreateGoal handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, request, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID
	request.RemainingBudget = request.TotalBudget
	request.Status = goals.StatusPlanned

	resp, err = goals.CreateGoal(&request)
	if err != nil {
		log.Println("CreateGoal handler cannot create goal:", err.Error())
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetGoal 		 Получение информации о цели
// @Summary      Получение информации о цели
// @Description  Получает детальную информацию о цели по ее ID
// @ID           get-goal
// @Tags         Цели
// @Produce      json
// @Security     JWT
// @Param        id   path      int  true  "ID Цели"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      403  {object}  response.ResponseModel "Доступ запрещен"
// @Failure      404  {object}  response.ResponseModel "Цель не найдена"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /goals/{id} [get]
func GetGoal(c *gin.Context) {
	var (
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
		id      int
	)

	if !middleware.CheckAccess(middleware.Goals, middleware.READ, ctxData.UserID) {
		response.SetResponseData(&resp, goals.Goals{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("GetGoal handler invalid ID format:", err.Error())
		response.SetResponseData(&resp, goals.Goals{}, "Неверный формат ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err = goals.GetGoal(id, ctxData.FamilyID)
	if err != nil {
		log.Println("GetGoal handler cannot get goal:", err.Error())
		response.SetResponseData(&resp, goals.Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetGoalsList 		 Получение информации о цели
// @Summary      Получение списка целей
// @Description  Получает список целей по фильтрам
// @ID           get-goals-list
// @Tags         Цели
// @Produce      json
// @Security     JWT
// @Param        search        query  string  false  "Поиск по названию и описанию"
// @Param        status        query  string  false  "Статус цели"
// @Param        due_date_from query  string  false  "Дата выполнения 'от' (YYYY-MM-DD)"
// @Param        due_date_to   query  string  false  "Дата выполнения 'до' (YYYY-MM-DD)"
// @Param        current_page  query  int     false  "Номер страницы"
// @Param        page_limit    query  int     false  "Количество на странице"
// @Success      200           {object}  response.ResponseModel
// @Failure      403           {object}  response.ResponseModel "Доступ запрещен"
// @Failure      500           {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /goals [get]
func GetGoalsList(c *gin.Context) {
	var (
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
		filter  goals.Filters
	)

	if !middleware.CheckAccess(middleware.Goals, middleware.READ, ctxData.UserID) {
		response.SetResponseData(&resp, goals.Goals{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	err = c.Bind(&filter)
	if err != nil {
		log.Println("GetGoals handler cannot get params:", err.Error())
		response.SetResponseData(&resp, []goals.Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	filter.FamilyID = ctxData.FamilyID

	resp, err = goals.GetGoalsList(filter)
	if err != nil {
		log.Println("GetGoals handler cannot get goals list:", err.Error())
		response.SetResponseData(&resp, []goals.Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateGoal godoc
// @Summary      Обновление цели
// @Description  Обновляет данные существующей цели по ее ID из тела запроса
// @ID           update-goal
// @Tags         Цели
// @Accept       json
// @Produce      json
// @Security     JWT
// @Param        {object}   body      goals.Goals  true  "Данные для обновления цели (включая ID)"
// @Success      200   {object}  response.ResponseModel
// @Failure      400   {object}  response.ResponseModel "Неверные данные"
// @Failure      403   {object}  response.ResponseModel "Доступ запрещен"
// @Failure      404   {object}  response.ResponseModel "Цель не найдена"
// @Failure      500   {object}  response.ResponseModel "Внутренняя ошибка"
// @Router       /goals [put]
func UpdateGoal(c *gin.Context) {
	var (
		request goals.Goals
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Goals, middleware.UPDATE, ctxData.UserID) {
		response.SetResponseData(&resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("UpdateGoal handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, nil, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.FamilyID = ctxData.FamilyID
	resp, err = goals.UpdateGoal(&request)
	if err != nil {
		log.Println("UpdateGoal handler internal server error:", err.Error())
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteGoal 		 Удаление цели
// @Summary      Удаление цели
// @Description  Удаляет цель по ее ID
// @ID           delete-goal
// @Tags         Цели
// @Produce      json
// @Security     JWT
// @Param        id   path      int  true  "ID Цели"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      403  {object}  response.ResponseModel "Доступ запрещен"
// @Failure      404  {object}  response.ResponseModel "Цель не найдена"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /goals/{id} [delete]
func DeleteGoal(c *gin.Context) {
	var (
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
		id      int
	)

	if !middleware.CheckAccess(middleware.Goals, middleware.DELETE, ctxData.UserID) {
		response.SetResponseData(&resp, goals.Goals{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("GetGoal handler invalid ID format:", err.Error())
		response.SetResponseData(&resp, goals.Goals{}, "Неверный формат ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err = goals.DeleteGoal(id, ctxData.FamilyID)
	if err != nil {
		log.Println("DeleteGoal handler cannot delete goal:", err.Error())
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
