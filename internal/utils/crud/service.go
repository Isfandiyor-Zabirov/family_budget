package crud

import "gorm.io/gorm"

type Repository[T any] struct{}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{}
}

func (c *Repository[T]) Create(tx *gorm.DB, entity *T) (T, error) {
	err := tx.Create(&entity).Error
	return *entity, err
}

func (c *Repository[T]) Update(tx *gorm.DB, entity *T) (T, error) {
	err := tx.Save(&entity).Error
	return *entity, err
}

func (c *Repository[T]) Delete(tx *gorm.DB, entity *T) error {
	return tx.Delete(&entity).Error
}

func (c *Repository[T]) Get(tx *gorm.DB, id int) (T, error) {
	var entity T
	err := tx.First(&entity, id).Error
	return entity, err
}
