package handlers

import (
	"family_budget/internal/entities/reports"
	"family_budget/internal/utils/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func GetMainReport(c *gin.Context) {
	var (
		ctx = getClaimsFromContext(c)
	)

	resp, err := reports.GetMainReport(ctx.FamilyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetGraphReport(c *gin.Context) {
	var (
		ctx     = getClaimsFromContext(c)
		filters reports.Filter
		resp    response.ResponseModel
		err     error
	)

	err = c.Bind(&filters)
	if err != nil {
		log.Println("GetGraphReport handler cannot bind filters", err.Error())
		response.SetResponseData(&resp, []reports.GraphReport{}, "Неверный фильтр", false, 0, 0, 0)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	filters.FamilyID = ctx.FamilyID
	if filters.DateFrom != "" {
		filters.DateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	}
	if filters.DateTo != "" {
		filters.DateTo = time.Now().Format("2006-01-02")
	}
	resp, err = reports.GetGraphReport(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
