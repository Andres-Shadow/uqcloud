package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
Clase encargada de contener los elementos relacioados con la conexion a la base de datos
y posteriormente la creacion de las tablas necesarias para el funcionamiento del sistema
*/

var DB *sql.DB

// Funciòn que se encarga de realizar la conexiòn a la base de datos
// mediante el driver manual de MySql
func ManageSqlConecction() {
	fmt.Println("Conectando a la base de datos...")
	var err error

	DB, err = sql.Open("mysql", "root:root@tcp(uqcloud)/uqcloud")
	if err != nil {
		log.Fatal(err)
	}

}

var DATABASE *gorm.DB

// Funciòn que se encarga de realizar la conexiòn a la base de datos
// mediante el driver de GORM
func DBConnection() {
	var host string
	host = os.Getenv("DATABASE")
	if host == "" {
		host = "localhost"
	}


	var dsn = "root:andres_1@tcp("+host+":3306)/uqcloud?charset=utf8mb4&parseTime=True&loc=Local"

	for {
		var err error
		DATABASE, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("Failed to connect to database. Retrying in 5 seconds...")
			time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
		} else {
			log.Println("DB Connected")
			break // Exit the loop once the connection is successful
		}
	}
}