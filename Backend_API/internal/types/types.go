package types

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"product_name"`
	Price    float64 `json:"product price"`
	Quantity int     `json:"product_quantity"`
}

type ErrorResponse struct {
	Status string
	Error  string
}
