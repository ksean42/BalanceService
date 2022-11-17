package entities

type UserBalanceRequest struct {
	Id int `json:"id"`
}

type AddRequest struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
}

type Request struct {
	Id        int     `json:"id"`
	ServiceID int     `json:"service_id"`
	OrderID   int     `json:"order_id"`
	Amount    float64 `json:"amount"`
}

type TransferRequest struct {
	SrcID  int     `json:"src_id"`
	DestID int     `json:"dest_id"`
	Amount float64 `json:"amount"`
}

type ReserveReject struct {
	Id      int `json:"id"`
	OrderId int `json:"order_id"`
}

type ErrorResponse struct {
	Error string
}

type ResultResponse struct {
	Result interface{}
}
