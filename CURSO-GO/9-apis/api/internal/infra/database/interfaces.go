package database

import (
	entity2 "github.com/dlcdev1/pos1/9-apis/api/internal/entity"
)

type UserInterface interface {
	Create(user *entity2.User) error
	FindByEmail(email string) (*entity2.User, error)
}

type ProductInterface interface {
	Create(product *entity2.Product) error
	FindById(id string) (*entity2.Product, error)
	Update(product *entity2.Product) error
	Delete(id string) error
	FindAll(page, limit int, sort string) ([]entity2.Product, error)
}
