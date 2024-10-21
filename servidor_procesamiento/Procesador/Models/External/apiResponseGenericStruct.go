package external

type GenericResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
