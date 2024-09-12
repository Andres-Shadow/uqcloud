package database

import "log"

// Función que crea los triggers en la base de datos
func CreateTriggers() {
	createUsersDeleteTrigger()
}

// Función encargada de inicializar el trigger de eliminación programada de usuarios con antiguedad mayor a 6 horas
func createUsersDeleteTrigger() {
	triggerSQL := `
	CREATE TRIGGER automatic_delete_users
	AFTER INSERT ON Persona
	FOR EACH ROW
	BEGIN
		IF NEW.rol = 0 THEN
			DELETE FROM Persona
			WHERE rol = 0
			AND id != NEW.id
			AND created_at < DATE_SUB(NOW(), INTERVAL 6 HOUR);
		END IF;
	END;
	`

	err := DATABASE.Exec(triggerSQL).Error
	if err != nil {
		log.Println("Error al crear el trigger de eliminación de usuarios:", err)
	} else {
		log.Println("Trigger de eliminación de usuarios creado correctamente.")
	}
}
