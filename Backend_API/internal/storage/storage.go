package storage

import "github.com/slangeres/Vypaar/backend_API/internal/types"

type Storage interface {
	CreateProduct(name string, price float32, quantity int) (int64, error)
	GetAllProduct() ([]types.Product, error)
	GetUserById(id int64) (types.Product, error)
	DeleteUser(id int64) error
	UpdateProduct(id int64, name string, price float32, quantity int) (int64, error)
}

type UserStorage interface{
	Login()
	Signup()
}