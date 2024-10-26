package apiutilities

import external "servidor_procesamiento/Procesador/Models/External"

func BuildGenericResponse[T any](data T, status, message string) external.GenericResponse[T] {
	return external.GenericResponse[T]{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
