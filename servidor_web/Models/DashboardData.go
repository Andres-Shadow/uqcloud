package Models

type DashboardData struct {
	Total_Maquinas_creadas    int `json:"total_Maquinas_creadas"`
	Total_maquinas_encendidas int `json:"total_maquinas_encendidas"`
	Total_usuarios            int `json:"total_usuarios"`
	Total_estudiantes         int `json:"total_estudiantes"` // TODO: Revisar si este atributo es necesario
	Total_invitados           int `json:"total_invitados"`
	Total_RAM                 int `json:"total_RAM"`
	Total_RAM_usada           int `json:"total_RAM_usada"`
	Total_CPU                 int `json:"total_CPU"`
	Total_CPU_usada           int `json:"total_CPU_usada"`
}
