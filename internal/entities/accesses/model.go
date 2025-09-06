package accesses

import "time"

// Access - access to Access_Groups CRUDs for role
type Access struct {
	ID          int        `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name        string     `json:"name" gorm:"column:name"` //create
	Code        string     `json:"-" gorm:"column:code"`
	Description string     `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Access) TableName() string {
	return "accesses"
}
