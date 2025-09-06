package access_groups

import "time"

// AccessGroup - access to entity for role
type AccessGroup struct {
	ID          int        `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
	Code        string     `gorm:"column:code" json:"code"`
	Name        string     `gorm:"column:name" json:"name"`               //users
	Description string     `gorm:"column:description" json:"description"` //detail for users
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *AccessGroup) TableName() string {
	return "access_groups"
}
