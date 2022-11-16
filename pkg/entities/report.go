package entities

type Report struct {
	ServiceID int     `json:"service_id" csv:"service_id"`
	Revenue   float64 `json:"revenue" csv:"revenue"`
}
