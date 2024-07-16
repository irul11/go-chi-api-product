package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Created_at  string  `json:"created_at"`
}

type Message struct {
	Query         string
	Method        string
	Body          Product
	URLParam      int
	OutputChannel chan MessageOutput
}

type MessageOutput struct {
	Data  []Product
	Error error
}
