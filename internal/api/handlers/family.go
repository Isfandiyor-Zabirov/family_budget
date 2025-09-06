package handlers

import (
	"errors"
	"family_budget/internal/entities/family"
	"family_budget/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateFamily godoc
// @Summary      Создание новой семьи
// @Description  Создает новую запись о семье
// @ID           create-family
// @Tags         Семьи
// @Accept       json
// @Produce      json
// @Param        family  body      CreateFamilyRequest  true  "Данные для создания семьи"
// @Success      201     {object}  response.ResponseModel
// @Failure      400     {object}  response.ResponseModel "Неверные входные данные"
// @Failure      500     {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /api/v1/families [post]
func CreateFamily(c *gin.Context) {
	var request family.Family
	if err := c.ShouldBindJSON(&request); err != nil {
		resp := response.SetResponseData(nil, err.Error(), false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	fam := &family.Family{
		Name:  request.Name,
		Phone: request.Phone,
	}

	resp, err := family.Create(fam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetFamily godoc
// @Summary      Получение информации о семье
// @Description  Получает информацию о семье по ее ID
// @ID           get-family
// @Tags         Семьи
// @Produce      json
// @Param        id   path      int  true  "ID Семьи"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      404  {object}  response.ResponseModel "Семья не найдена"
// @Router       /api/v1/families/{id} [get]
func GetFamily(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp := response.SetResponseData(nil, "Неверный формат ID", false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	family, err := family.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp := response.SetResponseData(nil, "Семья не найдена", false)
			c.JSON(http.StatusNotFound, resp)
			return
		}
		resp := response.SetResponseData(nil, "Внутренняя ошибка сервера", false)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := response.SetResponseData(family, "Успех", true)
	c.JSON(http.StatusOK, resp)
}

// UpdateFamily godoc
// @Summary      Обновление информации о семье
// @Description  Обновляет информацию о семье по ее ID
// @ID           update-family
// @Tags         Семьи
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "ID Семьи"
// @Param        family  body      UpdateFamilyRequest  true  "Данные для обновления"
// @Success      200     {object}  response.ResponseModel
// @Failure      400     {object}  response.ResponseModel "Неверные входные данные"
// @Failure      404     {object}  response.ResponseModel "Семья не найдена"
// @Failure      500     {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /api/v1/families/{id} [put]
func UpdateFamily(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp := response.SetResponseData(nil, "Неверный формат ID", false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err = family.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp := response.SetResponseData(nil, "Семья для обновления не найдена", false)
			c.JSON(http.StatusNotFound, resp)
			return
		}
		resp := response.SetResponseData(nil, "Ошибка при поиске семьи", false)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	var request family.Family
	if err := c.ShouldBindJSON(&request); err != nil {
		resp := response.SetResponseData(nil, err.Error(), false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	familyToUpdate := &family.Family{
		ID:    id,
		Name:  request.Name,
		Phone: request.Phone,
	}

	resp, err := family.Update(familyToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteFamily godoc
// @Summary      Удаление семьи
// @Description  Удаляет семью по ее ID (мягкое удаление)
// @ID           delete-family
// @Tags         Семьи
// @Produce      json
// @Param        id   path      int  true  "ID Семьи"
// @Success      200  {object}  response.ResponseModel
// @Failure      400  {object}  response.ResponseModel "Неверный формат ID"
// @Failure      500  {object}  response.ResponseModel "Внутренняя ошибка сервера"
// @Router       /api/v1/families/{id} [delete]
func DeleteFamily(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp := response.SetResponseData(nil, "Неверный формат ID", false)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	familyToDelete := &family.Family{ID: id}

	resp, err := family.Delete(familyToDelete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
