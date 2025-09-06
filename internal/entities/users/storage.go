package users

import (
	"errors"
	"family_budget/internal/entities/family"
	"family_budget/internal/entities/financial_event_categories"
	"family_budget/internal/entities/role_accesses"
	"family_budget/internal/entities/roles"
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
	"log"
	"regexp"
)

func createUser(user *User) (User, error) {
	repo := crud.NewRepository[User]()
	db := database.Postgres()

	return repo.Create(db, user)
}

func updateUser(user *User) (User, error) {
	repo := crud.NewRepository[User]()
	db := database.Postgres()
	return repo.Update(db, user, "Password", "Login")
}

func deleteUser(user *User) error {
	repo := crud.NewRepository[User]()
	db := database.Postgres()
	return repo.Delete(db, user)
}

func getUser(id int) (User, error) {
	repo := crud.NewRepository[User]()
	db := database.Postgres()
	return repo.Get(db, id)
}

func updatePassword(newHashedPassword string, userID int) error {
	sqlQuery := `update users set password = ? where id = ?`

	if err := database.Postgres().Exec(sqlQuery, newHashedPassword, userID).Error; err != nil {
		log.Println("updatePassword func query error:", err.Error())
		return err
	}
	return nil
}

func getMe(userID int) (me Me, err error) {
	userSqlQuery := `select u.id, u.name, u.surname, u.middle_name, u.surname || ' ' || u.name || ' ' || u.middle_name as full_name, u.login, u.phone,
						r.name as role, u.role_id, u.family_id, u.email
						from users u, roles r
						where u.role_id = r.id and u.id = ?`
	if err = database.Postgres().Raw(userSqlQuery, userID).Scan(&me.UserData).Error; err != nil {
		log.Println("getMe func user data query error:", err.Error())
		return Me{}, errors.New("Ошибка сервера ")
	}

	accessGroupSql := `select ag.id as access_group_id, ag.code as access_group_code, ag.name as access_group_name from access_groups ag`

	var accessGroupInfo []struct {
		AccessGroupID   int    `json:"access_group_id"`
		AccessGroupCode string `json:"access_group_code"`
		AccessGroupName string `json:"access_group_name"`
	}

	if err = database.Postgres().Raw(accessGroupSql).Scan(&accessGroupInfo).Error; err != nil {
		log.Println("getMe func accessGroupSql query error:", err.Error())
		return Me{}, errors.New("Ошибка сервера ")
	}

	accessSqlQuery := `select 
    					a.id as access_id, 
    					a.code as access_code, 
    					ra.active as active, 
    					a.name as access_name
						from accesses a, access_groups ag, role_accesses ra, roles r, users u
						where u.role_id = r.id and ra.access_id = a.id and ra.access_group_id = ag.id and ra.role_id = r.id and ag.id = ? and u.id = ?`

	for _, row := range accessGroupInfo {
		var accesses MeAccessGroup
		accesses.AccessGroupID = row.AccessGroupID
		accesses.AccessGroupCode = row.AccessGroupCode
		accesses.AccessGroupName = row.AccessGroupName
		if err = database.Postgres().Raw(accessSqlQuery, row.AccessGroupID, userID).Scan(&accesses.Accesses).Error; err != nil {
			log.Println("getMe func access list query error:", err.Error())
			return Me{}, errors.New("Ошибка сервера ")
		}
		me.AccessList = append(me.AccessList, accesses)
	}

	return
}

func register(d *RegistrationData) (err error) {
	var (
		u User
		f family.Family
	)

	u.Phone, err = validatePhoneTj(d.Phone)
	if err != nil {
		return err
	}
	u.Login, err = checkLogin(d.Login)
	log.Println(d.Login)
	if err != nil {
		return err
	}

	f.Name = d.FamilyName
	f.Phone = d.HomePhone

	tx := database.Postgres().Begin()
	err = tx.Create(&f).Error
	if err != nil {
		log.Println("create family while registration err:", err.Error())
		tx.Rollback()
		return err
	}

	u.FamilyID = f.ID
	u.RoleID = 1 // Глава семьи
	u.Name = d.Name
	u.Surname = d.Surname
	u.MiddleName = d.MiddleName
	u.Email = d.Email

	u.Password = hashPassword(d.Password)
	err = tx.Create(&u).Error
	if err != nil {
		log.Println("create user while registration err:", err.Error())
		tx.Rollback()
		return err
	}

	ras := role_accesses.AssignAccessesToRole(roles.Roles{ID: 1})
	err = tx.Create(&ras).Error
	if err != nil {
		log.Println("create role accesses while registration err:", err.Error())
		tx.Rollback()
		return err
	}

	fecs := []financial_event_categories.FinancialEventCategories{
		{
			FamilyID:    f.ID,
			Name:        "Еда и напитки",
			Description: "Расходы на еду и напитки",
		},
		{
			FamilyID:    f.ID,
			Name:        "Транспорт",
			Description: "Расходы на транспорт",
		},
		{
			FamilyID:    f.ID,
			Name:        "Домашные расходы",
			Description: "Домашные расходы",
		},
		{
			FamilyID:    f.ID,
			Name:        "ЖКХ",
			Description: "Расходы на свет, воду, отопление и т.д.",
		},
		{
			FamilyID:    f.ID,
			Name:        "Здоровье",
			Description: "Расходы на здоровье",
		},
		{
			FamilyID:    f.ID,
			Name:        "Шоппинг",
			Description: "Расходы на одежду и другие",
		},
		{
			FamilyID:    f.ID,
			Name:        "Зарплата",
			Description: "Доход от зарплаты",
		},
		{
			FamilyID:    f.ID,
			Name:        "Инвестиции",
			Description: "Инвестиции",
		},
		{
			FamilyID:    f.ID,
			Name:        "Развелечение",
			Description: "Расходы на кино, концерты и другие",
		},
		{
			FamilyID:    f.ID,
			Name:        "Другие",
			Description: "Другие доходы или расходы",
		},
	}

	for _, v := range fecs {
		fec := v

		err = tx.Create(&fec).Error
		if err != nil {
			log.Println("create financial entity categories while registration err:", err.Error())
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func loginExists(login string) bool {
	count := 0
	_ = database.Postgres().Raw("SELECT count(*) from users where login = ?", login).Scan(&count).Error
	if count > 0 {
		return true
	}
	return false
}

var phoneLengthTj = 13
var phoneRegex = "^[0-9]{9}$"

// validatePhoneTj - чек номера, добавляем 992, код для TJ, пока так, далее при выходе в интернейшнал по ip вычисляем страну
func validatePhoneTj(phone string) (string, error) {
	if ok, err := regexp.MatchString(phoneRegex, phone); !ok {
		if err != nil {
			return "", err
		}
		log.Println("phone doesn't match regex :", phone)
		return "", errors.New("Неверная структура номера (подсказка: 900-00-00-00) ")
	}
	if len(phone) == phoneLengthTj-4 {
		phone = "992" + phone
	}
	return phone, nil
}
