package storage

import "github.com/slangeres/Vypaar/backend_API/internal/types"

type Storage interface {
	CreateProduct(name string, price float32, quantity int, shopID string) (int64, error)
	GetAllProduct(shopID string, offset int64, limit int64, sortOrder string, sortField string) ([]types.Product, error)
	GetUserById(id int64, shopID string) (types.Product, error)
	DeleteUser(id int64, shopID string) error
	UpdateProduct(id int64, name string, price float32, quantity int, shopID string) (int64, error)
}

type UserStorage interface {
	Login(email string, password string) (string, error)
	Signup(name string, email string, password string, shopID string) (int64, error)
	VerifyEmail(email string) (int64, error)
}
