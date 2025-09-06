package handlers

import (
	"family_budget/internal/entities/users"
	"family_budget/internal/utils/response"
	"family_budget/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var (
		request users.User
		err     error
		ctxData = getClaimsFromContext(c)
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Users, middleware.CREATE, ctxData.UserID) {
		response.SetResponseData(resp, request, "Доступ запрещен", false, 0, 0, 0)
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
// @Router /api/v1/register [post]
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

	response.SetResponseData(&resp, users.RegistrationData{}, "Регистрация прошла успешно", true, 0, 0, 0)
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
// @Router /api/v1/get_me [get]
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
