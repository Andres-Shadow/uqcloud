package config

// Quick Virtual Machine Constants
// for ram, cpu and distro
// the distro name has to be the same as the one in the database
// in the field "distribucion_sistema_operativo"
const (
	DEFAULT_QUICK_VM_CPU    = 1
	DEFAULT_QUICK_VM_RAM    = 1024
)

// Host fast register constants
// used for default values that repeats in all the hosts

const (
	QUICK_HOST_HOSTNAME  = "Admin"
	QUICK_HOST_RAM       = 32000
	QUICK_HOST_CPU       = 16
	QUICK_HOST_STORAGE   = 500
	QUICK_HOST_NETWORK   = "Realtek 8821CE Wireless LAN 802.11ac PCI-E NIC"
	QUICK_HOST_SO        = "Windows"
	QUICK_HOST_SO_DISTRO = "64"
)
