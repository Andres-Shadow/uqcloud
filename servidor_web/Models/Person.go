package Models

// Clase persona que ayuda a la gestion de usuarios
type Person struct {
	Name     string `json:name`
	LastName string `json:lastName`
	Email    string `json:email`
	Password string `json:password`
	Role     string `json:role`
}
