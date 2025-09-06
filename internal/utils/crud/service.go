package crud

import "gorm.io/gorm"

// Repository is a generic CRUDer to make the routine CRUD operations
type Repository[T any] struct{}

// NewRepository to get the new instance of the generic Repository
func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{}
}

func (c *Repository[T]) Create(tx *gorm.DB, entity *T) (T, error) {
	err := tx.Create(entity).Error
	return *entity, err
}

func (c *Repository[T]) Update(tx *gorm.DB, entity *T, ignoredFields ...string) (T, error) {
	err := tx.Model(entity).Omit(ignoredFields...).Updates(entity).Error
	return *entity, err
}

func (c *Repository[T]) Delete(tx *gorm.DB, entity *T) error {
	return tx.Delete(entity).Error
}

func (c *Repository[T]) Get(tx *gorm.DB, id int) (T, error) {
	var entity T
	err := tx.First(&entity, id).Error
	return entity, err
}
