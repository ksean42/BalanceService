package entities

type UserBalance struct {
	Id      int     `json:"ID"`
	Balance float64 `json:"balance"`
}

type Request struct {
	Id        int     `json:"id"`
	ServiceID int     `json:"service_ID"`
	OrderID   int     `json:"order_ID"`
	Amount    float64 `json:"amount"`
}
