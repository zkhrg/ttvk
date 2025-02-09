package entities

type Entity struct {
	ID          string `json:"id"`
	Address     string `json:"ip_address"`
	PingTime    int    `json:"ping_time"`
	LastSuccess string `json:"last_success"`
}

type EntityRequest struct {
	ID          string `json:"id"`
	Address     string `json:"ip_address"`
	PingTime    int    `json:"ping_time"`
	LastSuccess string `json:"last_success"`
}
