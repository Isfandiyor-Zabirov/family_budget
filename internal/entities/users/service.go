package users

import (
	"errors"
	"family_budget/internal/utils/response"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func GetMe(userID int) (resp response.ResponseModel, err error) {
	me, err := getMe(userID)
	if err != nil {
		resp = response.SetResponseData(Me{}, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(me, "Успех", true)
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
