package Models

type VirtualMachineTemp struct {
	VMName string `json:"vmname"`
	OS     string `json:"os"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
