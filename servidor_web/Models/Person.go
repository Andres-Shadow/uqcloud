package Models

// Clase persona que ayuda a la gestion de usuarios
type Person struct {
	Nombre      string `json:"usr_name"`
	Apellido    string `json:"usr_surname"`
	Email       string `json:"usr_email" gorm:"unique"`
	Contrasenia string `json:"usr_password"`
	Rol         byte   `json:"usr_role"`
}
