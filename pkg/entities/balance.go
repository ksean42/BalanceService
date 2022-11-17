package entities

type UserBalanceRequest struct {
	Id int `json:"ID"`
}

type AddRequest struct {
	Id     int     `json:"ID"`
	Amount float64 `json:"amount"`
}

type Request struct {
	Id        int     `json:"id"`
	ServiceID int     `json:"service_ID"`
	OrderID   int     `json:"order_ID"`
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

type ReportRequest struct {
	Date string `json:"date"`
}

type ErrorResponse struct {
	Error string
}

type ResultResponse struct {
	Result interface{}
}
