package db

type Transaction struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Transaction float64 `json:"transaction"`
}
