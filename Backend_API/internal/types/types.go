package types

type Product struct {
	Id       int     `json:"id"`
	ShopID   string  `json:"shopID"`
	Name     string  `validate:"required"`
	Price    float64 `validate:"required"`
	Quantity int     `validate:"required"`
}

type ErrorResponse struct {
	Status string
	Error  string
}
type User struct {
	Id       int    `json:"id"`
	Name     string `validate:"required"`
	Email    string `validate:"required"`
	ShopID   string `json:"shopID"`
	Password string `validate:"required"`
	CreateBy string `json:"created_at"`
}
type Login struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

var ValidSortFields = map[string]string{
	"id":         "id",
	"price":      "product_price",
	"name":       "product_name",
	"created_at": "created_at",
}
