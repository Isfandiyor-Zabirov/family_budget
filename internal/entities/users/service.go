package users

import (
	"errors"
	"family_budget/internal/utils/response"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func CreateUser(user *User) (resp response.ResponseModel, err error) {
	created, err := createUser(user)
	if err != nil {
		response.SetResponseData(&resp, User{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, created, "Пользователь успешно добавлен", true, 0, 0, 0)
	return
}

func UpdateUser(user *User) (resp response.ResponseModel, err error) {
	updated, err := updateUser(user)
	if err != nil {
		response.SetResponseData(&resp, User{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, updated, "Пользователь успешно обновлен", true, 0, 0, 0)
	return
}

func DeleteUser(id, familyID int) (resp response.ResponseModel, err error) {
	userDB, err := getUser(id)
	if err != nil {
		response.SetResponseData(&resp, User{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if userDB.FamilyID != familyID {
		log.Println("User DeleteUser func trying to delete user with wrong familyID")
		err = errors.New("user DeleteUser func trying to delete user with wrong familyID")
		response.SetResponseData(&resp, User{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		return
	}

	err = deleteUser(&userDB)
	if err != nil {
		response.SetResponseData(&resp, User{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, userDB, "Пользователь успешно удален", true, 0, 0, 0)
	return
}

func GetUserList(filters Filters) (resp response.ResponseModel, pagination response.Pagination, err error) {
	list, total, err := getList(filters)
	if err != nil {
		response.SetResponseData(&resp, []UserResp{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, list, "Успех", true, filters.PageLimit, total, filters.CurrentPage)
	return
}

func GetMe(userID int) (resp response.ResponseModel, err error) {
	me, err := getMe(userID)
	if err != nil {
		response.SetResponseData(&resp, Me{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, me, "Успех", true, 0, 0, 0)
	return
}

func Register(data *RegistrationData) error {
	return register(data)
}

func UpdatePassword(oldPassword, newPassword string, userID int) error {
	existingUser, err := getUser(userID)
	if err != nil {
		return err
	}

	if existingUser.ID == 0 {
		return errors.New("пользователь не найден ")
	}

	if !checkPassword(existingUser.Password, oldPassword) {
		return errors.New("неверный пароль ")
	}

	return updatePassword(hashPassword(newPassword), userID)
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return ""
	}
	password = string(bytes)
	return password
}

// checkPassword - создает хэш
func checkPassword(provided, existing string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(provided), []byte(existing))
	return err == nil
}

func validateString(str string, min, max int) error {
	symbols := []rune(str)

	if len(symbols) < min || len(symbols) > max {
		return fmt.Errorf("длина должна быть между %v и %v", min, max)
	}

	for _, val := range symbols {
		if (val >= 'а' && val <= 'я') || (val >= 'А' && val <= 'Я') {
			return fmt.Errorf("логин содержит кирилицу")
		}
	}

	return nil
}

func checkLogin(login string) (string, error) {
	log.Println("login: ", login)
	login = strings.TrimSpace(login)
	login = strings.ToLower(login)
	err := validateString(login, 4, 30)
	if err != nil {
		return "", err
	}
	if loginExists(login) {
		return "", errors.New("логин занят, введите другой логин")
	}

	return login, nil
}
