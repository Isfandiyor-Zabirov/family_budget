package handlers

import (
	"family_budget/internal/entities/users"
	"family_budget/internal/utils/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type RegistrationData struct {
	RoleID     int    `gorm:"column:role_id" json:"role_id"`
	FamilyID   int    `gorm:"column:family_id" json:"family_id"`
	Name       string `gorm:"column:name" binding:"required" json:"name"`
	Surname    string `gorm:"column:surname" binding:"required" json:"surname"`
	MiddleName string `gorm:"column:middle_name" json:"middle_name"`
	Phone      string `gorm:"column:phone" binding:"required" json:"phone"`
	Email      string `gorm:"column:email" json:"email"`
	Login      string `gorm:"column:login" binding:"required" json:"login"`
	Password   string `gorm:"column:password" binding:"required" json:"password"`
	FamilyName string `gorm:"column:family_name" binding:"required" json:"family_name"`
	HomePhone  string `gorm:"column:owner_phone" json:"owner_phone"`
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
		resp = response.SetResponseData(users.RegistrationData{}, "Неверная структура запроса", false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err = users.Register(&request)
	if err != nil {
		resp = response.SetResponseData(users.RegistrationData{}, "Что-то пошло не так", false)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp = response.SetResponseData(users.RegistrationData{}, "Регистрация прошла успешно", true)
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
