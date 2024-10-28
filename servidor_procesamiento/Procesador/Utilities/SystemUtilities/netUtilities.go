package systemutilities

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	database "servidor_procesamiento/Procesador/Database"
	models "servidor_procesamiento/Procesador/Models/Entities"
	"time"

	"golang.org/x/crypto/ssh"
)

/*
Clase encarga de contener todos los metodos que representan una funcionalidad
interna de la aplicacion que son llamadas dentro del programa y se encuentran relacionaodas con la
gestion de elementos relacionados con la red
*/

var logger = log.New(os.Stdout, "Logger: ", log.Ldate|log.Ltime|log.Lshortfile)

func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil // La IP es válida si no es nil
}

func Pacemaker(rutallaveprivata string, usuario string, ip string) bool {
	if !ValidateIP(ip) {
		logger.Println("IP no válida:", ip)
		return false
	}

	salida := false
	config := &ssh.ClientConfig{
		User: usuario,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(rutallaveprivata),
		},
		Timeout:         2 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ip = ip + ":22"
	conn, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		logger.Println("Error al establecer la conexión SSH:", ip)
		return salida
	}

	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, _, err := conn.SendRequest("heartbeat", true, nil)
		if err != nil {
			logger.Println("La conexión SSH está inactiva:", err)
			// Implementar lógica de reconexión
		} else {
			logger.Println("La conexión SSH está activa.")
			salida = true
			return salida
		}
	}

	session, err := conn.NewSession()
	if err != nil {
		logger.Println("Error al abrir sesión SSH:", err)
		return false
	}

	defer session.Close()

	command := "ls " + "C:/Uqcloud"
	if err := session.Run(command); err != nil {
		logger.Println("El archivo no existe en la ruta especificada:", err)
		return false
	}

	return salida
}

/*
Funciòn que se encarga de enviar los comandos a travès de la conexiòn SSH con el host

@host Paràmetro que contien la direcciòn IP del host al cual le va a enviar los comandos
@comando Paràmetro que contiene la instrucciòn que se desea ejecutar en el host
@config Paràmetro que contiene la configuraciòn SSH
@return Retorna la respuesta del host si la hay
*/
func SendSSHCommand(host string, comando string, config *ssh.ClientConfig) (salida string, err error) {

	//Establece la conexiòn SSH
	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		log.Println("Error al establecer la conexiòn SSH: ", err)
		return "", err
	}
	defer conn.Close()

	//Crea una nueva sesiòn SSH
	session, err := conn.NewSession()
	if err != nil {
		log.Println("Error al crear la sesiòn SSH: ", err)
		return "", err
	}
	defer session.Close()
	//Ejecuta el comando remoto
	output, err := session.CombinedOutput(comando)
	if err != nil {
		log.Println("Error al ejecutar el comando remoto: " + string(output))
		return "", err
	}
	return string(output), nil
}

/*
Funciòn que dada una direcciòn IP permite conocer si pertenece o no a un host registrado en la base de datos.
@ip Paràmetro que contiene la direcciòn IP a consultar
@Return Retorna el host en caso de que la IP estè en la base de datos
*/
func IsAHostIp(ip string) (models.Host, error) {

	host, err := database.GetHostByIp(ip)
	if err != nil {
		if err == sql.ErrNoRows {
			return host, err
		}
		return host, err // Otro error de la base de datos
	}
	return host, nil // IP encontrada en la base de datos, devuelve el objeto Host correspondiente
}

/*
Funciòn que se encarga de realizar la configuraciòn SSH con el host por medio de la contrasenia
@user Paràmetro que contiene el nombre del usuario al cual se va a conectar
@privateKeyPath Paràmetro que contiene la ruta de la llave privada SSH
*/

func ConfigureSSHPassword(user string) (*ssh.ClientConfig, error) {

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password("uqcloud"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config, nil
}

func ProcessSshPublicKeyConfiguration(fileRoute, ip string) error {
	bitSize := 4096

	//crear la carpeta en la ruta
	err := os.MkdirAll(fileRoute, 0700)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	savePrivateFileTo := fileRoute + "/id_rsa_test"
	savePublicFileTo := fileRoute + "/id_rsa_test.pub"

	privateKey, err := GeneratePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	publicKeyBytes, err := GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	privateKeyBytes := EncodePrivateKeyToPEM(privateKey)

	err = WriteKeyToFile(privateKeyBytes, savePrivateFileTo)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	err = WriteKeyToFile([]byte(publicKeyBytes), savePublicFileTo)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	// Enviar la clave pública por SSH
	err = SendPublicKeyViaSSH("uqcloud", ip, "./id_rsa", savePublicFileTo)
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writeKeyToFile writes keys to a file
func WriteKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := os.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}

// sendPublicKeyViaSSH sends the public key to a remote server using SSH
func SendPublicKeyViaSSH(user, host, privateKeyPath, publicKeyPath string) error {
	cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", "-i", privateKeyPath, publicKeyPath, fmt.Sprintf("%s@%s:.", user, host))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error sending public key: %v, output: %s", err, string(output))
	}

	cmd2 := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i", privateKeyPath, fmt.Sprintf("%s@%s", user, host), "bash", "-s", "<", "./init.sh")
	output2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		return fmt.Errorf("error executing init.sh: %v, output: %s", err2, string(output2))
	}

	log.Printf("Public key sent to %s@%s", user, host)
	return nil
}
