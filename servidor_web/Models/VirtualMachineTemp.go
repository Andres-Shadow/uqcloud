package Models

/*
Esta estructura temporal ayuda a manjerar las configuraciones de las maquinas virtuales de pues de ser creada, los campos de puede
Utilizar de manera que mejor convenga, es decir solo utilizar los campo que se crean necesarios deacuerdo a como lo requiera
La función
*/
type VirtualMachineTemp struct {
	VMName string `json:"vmname"`
	OS     string `json:"os"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Email  string `json:"email"`
	Ram    string `json:"ram"`
}
