package entities

type UserBalance struct {
	Id      int     `json:"ID"`
	Balance float64 `json:"balance"`
}

type Request struct {
	Id        int     `json:"id"`
	ServiceID int     `json:"serviceID"`
	OrderID   int     `json:"orderID"`
	Amount    float64 `json:"amount"`
}
