package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*
Clase encargada de contener los elementos relacioados con la conexion a la base de datos
y posteriormente la creacion de las tablas necesarias para el funcionamiento del sistema
*/

var DATABASE *gorm.DB

// Funciòn que se encarga de realizar la conexiòn a la base de datos
// mediante el driver de GORM
func DBConnection() {
	var host string
	host = os.Getenv("DATABASE")
	if host == "" {
		host = "localhost"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	var dsn = dbUser + ":" + dbPassword + "@tcp(" + host + ":3306)/uqcloud?charset=utf8mb4&parseTime=True&loc=Local"

	for {
		var err error
		DATABASE, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			log.Println("Failed to connect to database. Retrying in 10 seconds...")
			time.Sleep(10 * time.Second) // Wait for 5 seconds before retrying
		} else {
			log.Println("DB Connected")
			break // Exit the loop once the connection is successful
		}
	}

}
