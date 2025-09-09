package handlers

import (
	"family_budget/internal/entities/users"
	"family_budget/internal/utils/response"
	"family_budget/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreateUser - Добавление членов семьи (пользователей)
// @Summary Добавление членов семьи (пользователей)
// @ID create-user
// @Tags Семейство (Пользователи)
// @Produce json
// @Security     JWT
// @Param name  		 	body string true "Имя"
// @Param surname  			body string true "Фамилия"
// @Param middle_name  		body string false "Отчетсво"
// @Param role_id	  		body integer true "ID роли"
// @Param phone  			body string false "Номер телефона"
// @Param email  			body string false "Электронная почта"
// @Param login  			body string true "Логин"
// @Param password  		body string true "Пароль"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var (
		request users.User
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Users, middleware.CREATE, ctxData.UserID) {
		response.SetResponseData(&resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("CreateUser handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, []users.UserResp{}, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
	}

	request.FamilyID = ctxData.FamilyID

	resp, err = users.CreateUser(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateUser - Изменение членов семьи (пользователей)
// @Summary Изменение членов семьи (пользователей)
// @ID update-user
// @Tags Семейство (Пользователи)
// @Produce json
// @Security     JWT
// @Param id	  		 	body integer true "ID пользователя"
// @Param name  		 	body string true "Имя"
// @Param surname  			body string true "Фамилия"
// @Param middle_name  		body string false "Отчетсво"
// @Param role_id	  		body integer true "ID роли"
// @Param phone  			body string false "Номер телефона"
// @Param email  			body string false "Электронная почта"
// @Param login  			body string true "Логин"
// @Param password  		body string true "Пароль"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /users [put]
func UpdateUser(c *gin.Context) {
	var (
		request users.User
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Users, middleware.UPDATE, ctxData.UserID) {
		response.SetResponseData(&resp, request, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("UpdateUser handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, []users.UserResp{}, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
	}

	resp, err = users.UpdateUser(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser - Удаление членов семьи (пользователей)
// @Summary Удаление членов семьи (пользователей)
// @ID delete-user
// @Tags Семейство (Пользователи)
// @Produce json
// @Security     JWT
// @Param id path string true "ID пользователя"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var (
		ctxData     = getClaimsFromContext(c)
		resp        response.ResponseModel
		userID, err = strconv.Atoi(c.Param("id"))
	)

	if err != nil {
		log.Println("DeleteUser handler cannot convert the the ID:", err.Error())
		response.SetResponseData(&resp, nil, "Неверный ID", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !middleware.CheckAccess(middleware.Users, middleware.DELETE, ctxData.UserID) {
		response.SetResponseData(&resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	resp, err = users.DeleteUser(userID, ctxData.FamilyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserList - Получение списка всех членов семьи (пользователей)
// @Summary Получение списка всех членов семьи (пользователей)
// @ID get-user-list
// @Tags Семейство (Пользователи)
// @Produce json
// @Security     JWT
// @Param search 	    query string 	false "Поиск по ФИО, номеру телефона, эл. почты и логина"
// @Param role_id 	 	query integer  	false "ID роли"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /users [get]
func GetUserList(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
		filters users.Filters
		err     error
	)

	if !middleware.CheckAccess(middleware.Users, middleware.READ, ctxData.UserID) {
		response.SetResponseData(&resp, nil, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	err = c.Bind(&filters)
	if err != nil {
		log.Println("GetUserList handler cannot bind filters:", err.Error())
		response.SetResponseData(&resp, nil, "Неверные фильтры", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
	}

	filters.FamilyID = ctxData.FamilyID

	resp, err = users.GetUserList(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Register - Регистрация
// @Summary Регистарция новых семей
// @ID register
// @Tags Регистрация
// @Produce json
// @Security     JWT
// @Param name  		 	body string true "Имя"
// @Param surname  			body string true "Фамилия"
// @Param middle_name  		body string false "Отчетсво"
// @Param phone  			body string false "Номер телефона"
// @Param email  			body string false "Электронная почта"
// @Param login  			body string true "Логин"
// @Param password  		body string true "Пароль"
// @Param family_name  		body string true "Название семейство"
// @Param home_phone  		body string false "Номер домашнего телефона"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /register [post]
func Register(c *gin.Context) {
	var (
		request users.RegistrationData
		err     error
		resp    response.ResponseModel
	)

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("Register handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, users.RegistrationData{}, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err = users.Register(&request)
	if err != nil {
		response.SetResponseData(&resp, users.RegistrationData{}, "Что-то пошло не так", false, 0, 0, 0)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	response.SetResponseData(&resp, request, "Регистрация прошла успешно", true, 0, 0, 0)
	c.JSON(http.StatusOK, resp)
}

// GetMe - получить данные текущего пользователя
// @Summary Получение данные текущего пользователя
// @ID get-me
// @Tags Пользователи
// @Produce json
// @Security     JWT
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /get_me [get]
func GetMe(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
	)

	resp, err := users.GetMe(ctxData.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
