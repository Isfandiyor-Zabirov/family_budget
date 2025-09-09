package handlers

import (
	"family_budget/internal/entities/role_accesses"
	"family_budget/internal/entities/roles"
	"family_budget/internal/utils/response"
	"family_budget/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetRole - получить данные роли по ID
// @Summary Получение данных роли по ID
// @ID get-role-by-id
// @Tags Роли и доступы
// @Produce json
// @Security     JWT
// @Param id path string true "id"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /role/{id} [get]
func GetRole(c *gin.Context) {
	var (
		ctxData     = getClaimsFromContext(c)
		roleID, err = strconv.Atoi(c.Param("id"))
		resp        response.ResponseModel
	)

	if err != nil {
		response.SetResponseData(&resp, roles.GetRoleResp{}, "Неверный ID роли", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !middleware.CheckAccess(middleware.Roles, middleware.READ, ctxData.UserID) {
		response.SetResponseData(&resp, roles.GetRoleResp{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}
	resp, err = roles.GetRole(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteRole - удаление роли
// @Summary Удаление роли
// @ID delete-role
// @Tags Роли и доступы
// @Produce json
// @Security     JWT
// @Param id path string true "ID"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /role/{id} [delete]
func DeleteRole(c *gin.Context) {
	var (
		ctxData     = getClaimsFromContext(c)
		roleID, err = strconv.Atoi(c.Param("id"))
		resp        response.ResponseModel
	)

	if err != nil {
		response.SetResponseData(&resp, roles.GetRoleResp{}, "Неверный ID роли", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !middleware.CheckAccess(middleware.Roles, middleware.DELETE, ctxData.UserID) {
		response.SetResponseData(&resp, roles.GetRoleResp{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}
	resp, err = roles.DeleteRole(roleID, ctxData.FamilyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetRoles - Получение данных ролей по фильтру с пагинацией и фильтрацией.
// @Summary Получение данных о ролях по фильтру с пагинацией и фильтрацией. Данные всегда только по текущей семье.
// @ID get_roles
// @Tags Роли и доступы
// @Produce json
// @Security     JWT
// @Param role_id 	 	        query integer false "ID роли"
// @Param name		 	 		query string  false "Поиск по части названия роли"
// @Param description  		 	query string  false "Поиск по части описания роли"
// @Param page 	 		 		query integer false "Страница"
// @Param page_limit 	 		query integer false "Количество рядов на странице(для пагинации)"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	var (
		filter  roles.GetRolesFilter
		resp    response.ResponseModel
		err     error
		ctxData = getClaimsFromContext(c)
	)

	if err = c.Bind(&filter); err != nil {
		log.Println("GetRoles handler cannot bind filters:", err.Error())
		response.SetResponseData(&resp, []roles.GetRoleResp{}, "Неверный фильтр", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !middleware.CheckAccess(middleware.Roles, middleware.READ, ctxData.UserID) {
		response.SetResponseData(&resp, []roles.GetRoleResp{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}
	filter.FamilyID = &ctxData.FamilyID
	resp, err = roles.GetRoles(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRoleWithAccesses - Получение роли со всеми доступами по role_id
// @Summary Получение роли и всех его доступов по role_id
// @ID get-role-with-accesses-by-id
// @Tags Роли и доступы
// @Produce json
// @Security     JWT
// @Param id path string true "id"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /accesses/{role_id} [get]
func GetRoleWithAccesses(c *gin.Context) {
	var (
		ctxData     = getClaimsFromContext(c)
		roleID, err = strconv.Atoi(c.Param("role_id"))
		resp        response.ResponseModel
	)

	if err != nil {
		log.Println("GetRoleWithAccesses handler cannot convert role ID err: ", err.Error())
		response.SetResponseData(&resp, roles.GetRolesWithAccesses{}, "Неверный ID роли", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !middleware.CheckAccess(middleware.Roles, middleware.READ, ctxData.UserID) {
		response.SetResponseData(&resp, roles.GetRolesWithAccesses{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}
	resp, err = roles.GetRoleWithAccesses(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateRoleWithAccesses - Создание роли (сразу создает список всех возможных доступов со значением false для возможности дальнейшего изменения)
// @Summary Создание Роли с доступами
// @ID create-role-with-accesses
// @Tags Роли и доступы
// @Produce json
// @Security     JWT
// @Param id body role_accesses.CreateRoleWithAccessesReq true "Даные пользователя"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /role_with_accesses [post]
func CreateRoleWithAccesses(c *gin.Context) {
	var (
		ctxData = getClaimsFromContext(c)
		request role_accesses.CreateRoleWithAccessesReq
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Roles, middleware.CREATE, ctxData.UserID) {
		response.SetResponseData(&resp, struct{}{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		log.Println("CreateRoleWithAccesses handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, struct{}{}, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	request.Role.FamilyID = ctxData.FamilyID

	resp, err = role_accesses.CreateRoleWithAccesses(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateRoleWithAccesses - Изменение роли с доступами
// @Summary Изменение Роли с доступами
// @ID update-role-with-accesses
// @Tags Роли и доступы
// @Produce json
// @Security     JWT
// @Param id body roles.UpdateRoleWithAccessesReq true "Обновленные данные роля"
// @Param Authorization header string true "Bearer + Token"
// @Success 200 {object} response.ResponseModel
// @Failure 400 {object} response.ResponseModel
// @Router /role_with_accesses [put]
func UpdateRoleWithAccesses(c *gin.Context) {
	var (
		ctx     = getClaimsFromContext(c)
		request roles.UpdateRoleWithAccessesReq
		err     error
		resp    response.ResponseModel
	)

	if !middleware.CheckAccess(middleware.Roles, middleware.UPDATE, ctx.UserID) {
		response.SetResponseData(&resp, struct{}{}, "Доступ запрещен", false, 0, 0, 0)
		c.JSON(http.StatusForbidden, resp)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println("UpdateRoleWithAccesses handler cannot bind the request:", err.Error())
		response.SetResponseData(&resp, struct{}{}, "Неверная структура запроса", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err = roles.UpdateRoleWithAccesses(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
