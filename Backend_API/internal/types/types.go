package types

type Product struct {
	Id       int     `json:"id"`
	Name     string  `validate:"required"`
	Price    float64 `validate:"required"`
	Quantity int     `validate:"required"`
}

type ErrorResponse struct {
	Status string
	Error  string
}
