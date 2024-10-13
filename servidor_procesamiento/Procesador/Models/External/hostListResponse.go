package external

type HostListResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Count   int                      `json:"count"`
	Data    []map[string]interface{} `json:"data"`
}
