package users

import (
	"family_budget/internal/entities/family"
	"family_budget/internal/entities/roles"
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
	RoleFK     roles.Roles    `gorm:"foreignKey:RoleID" json:"-"`
}

func (*User) TableName() string {
	return "users"
}

type UserResp struct {
	User
	Role          string `json:"role"`
	CreatedAtText string `json:"created_at_text"`
}

type Filters struct {
	FamilyID    int     // for internal only
	Search      *string `form:"search"`
	RoleID      *int    `form:"search"`
	CurrentPage int     `form:"current_page"`
	PageLimit   int     `form:"page_limit"`
}

type GetUserResponseModel struct {
	User
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

type Me struct {
	UserData   GetUserResponseModel `json:"user_data"`
	AccessList []MeAccessGroup      `json:"access_list"`
}

type MeAccessGroup struct {
	AccessGroupID   int          `json:"access_group_id"`
	AccessGroupCode string       `json:"access_group_code"`
	AccessGroupName string       `json:"access_group_name"`
	Accesses        []MeAccesses `json:"accesses"`
}

type MeAccesses struct {
	AccessID   int    `json:"access_id"`
	AccessCode string `json:"access_code"`
	AccessName string `json:"access_name"`
	Active     bool   `json:"active"`
}

type RegistrationData struct {
	RoleID     int    `gorm:"column:role_id" json:"role_id"`
	FamilyID   int    `gorm:"column:family_id" json:"family_id"`
	Name       string `gorm:"column:name" binding:"required" json:"name"`
	Surname    string `gorm:"column:surname" binding:"required" json:"surname"`
	MiddleName string `gorm:"column:middle_name" json:"middle_name"`
	Phone      string `gorm:"column:phone" binding:"required" json:"phone"`
	Email      string `gorm:"column:email" json:"email"`
	Login      string `gorm:"column:login" binding:"required" json:"login"`
	Password   string `gorm:"column:password" binding:"required" json:"password"`
	FamilyName string `gorm:"column:family_name" binding:"required" json:"family_name"`
	HomePhone  string `gorm:"column:owner_phone" json:"home_phone"`
}
