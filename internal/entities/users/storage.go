package users

import (
	"gorm.io/gorm"
)

type UserStorage interface {
	Create(user *User) error
	GetByLogin(login string) (*User, error)
	GetByID(id int) (*User, error) 
}

type storage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return &storage{db: db}
}

func (s *storage) Create(user *User) error {
	result := s.db.Create(user)
	return result.Error
}

func (s *storage) GetByLogin(login string) (*User, error) {
	var user User
	if err := s.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *storage) GetByID(id int) (*User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
