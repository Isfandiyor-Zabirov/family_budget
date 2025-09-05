package financial_event_categories

func Create(fec *FinancialEventCategories) (FinancialEventCategories, error) {
	return createFec(fec)
}

func Update(fec *FinancialEventCategories) (FinancialEventCategories, error) {
	return updateFec(fec)
}

func Delete(id int) error {
	fec := &FinancialEventCategories{ID: id}
	return deleteFec(fec)
}

func Get(id int) (FinancialEventCategories, error) {
	return getFec(id)
}

func GetList(filters Filters) {

}
