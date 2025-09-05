package users

import (
	"family_budget/internal/entities/family"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	RoleID     int            `gorm:"column:role_id"  json:"role_id"`
	FamilyID   int            `gorm:"column:family_id" json:"family_id"`
	Name       string         `gorm:"column:name"  json:"name"`
	Surname    string         `gorm:"column:surname"  json:"surname"`
	MiddleName string         `gorm:"column:middle_name" json:"middle_name"`
	Phone      string         `gorm:"column:phone"  json:"phone"`
	Email      string         `gorm:"column:email" json:"email"`
	Login      string         `gorm:"column:login;unique"  json:"login"`
	Password   string         `gorm:"column:password"  json:"password"`
	Limit      float64        `gorm:"column:limit" json:"limit"`
	CreatedAt  *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK   family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*User) TableName() string {
	return "users"
}
