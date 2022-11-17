package entities

import "time"

type Report struct {
	ServiceID int     `json:"service_id" csv:"service_id"`
	Revenue   float64 `json:"revenue" csv:"revenue"`
}

type UserReport struct {
	Transfers    []UserTransferReport    `json:"transfers,omitempty"`
	Transactions []UserTransactionReport `json:"transactions,omitempty"`
	Reserves     []UserReserving         `json:"reserves,omitempty"`
}

type UserTransferReport struct {
	SrcID  int       `json:"src_id"`
	DestID int       `json:"dest_id"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
}

type UserTransactionReport struct {
	ServiceID int       `json:"service_id"`
	OrderID   int       `json:"order_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
}

type UserReserving struct {
	ServiceID int     `json:"service_id"`
	OrderID   int     `json:"order_id"`
	Amount    float64 `json:"amount"`
}

type ReportRequest struct {
	Date string `json:"date"`
}

type UserReportRequest struct {
	Id int `json:"id"`
}
