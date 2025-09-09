package reports

import (
	"family_budget/pkg/database"
	"log"
)

func getMainReport(familyID int) (resp MainReport, err error) {
	sqlQuery := `SELECT case
           when fe.inflow = false || g.id != 0 then sum(t.amount)
           else sum(t.amount)
           end as expenses
FROM transactions t
         LEFT JOIN goals g on t.goal_id = g.id
         LEFT JOIN financial_events fe on t.financial_event_id = fe.id
WHERE t.family_id = ?`
	err = database.Postgres().Raw(sqlQuery, familyID).Scan(&resp).Error
	if err != nil {
		log.Println("getMainReport func query error:", err.Error())
	}

	return
}

func getGraphReport(filter Filter) (resp []GraphReport, err error) {
	sqlQuery := `SELECT t.created_at::date as date,
       case
           when fe.inflow = false || g.id != 0 then sum(t.amount)
           else sum(t.amount)
           end            as expenses
FROM transactions t
         LEFT JOIN goals g on t.goal_id = g.id
         LEFT JOIN financial_events fe on t.financial_event_id = fe.id
WHERE t.family_id = ?
  AND t.created_at::date >= ?
  AND t.created_at::date <= ?
GROUP BY t.created_at::date
ORDER BY t.created_at::date;`

	err = database.Postgres().Raw(sqlQuery, filter.FamilyID, filter.DateFrom, filter.DateTo).Scan(&resp).Error
	if err != nil {
		log.Println("getMainReport func query error:", err.Error())
	}

	return
}
