function showMessageIfNeeded() {
    const urlParams = new URLSearchParams(window.location.search);
    const message = urlParams.get('message');

    if (message) {
        showMessage(message, 'success');
    }
}

// Llama a la función cuando se carga la página
document.addEventListener('DOMContentLoaded', showMessageIfNeeded);

function showMessage(message, type) {
    const messageDiv = document.createElement("div");
    messageDiv.className = `alert alert-${type === "success" ? "success" : "danger"} alert-dismissible fade show alert-message`;
    messageDiv.role = "alert";
    messageDiv.textContent = message;

    // Botón de cerrar para el mensaje
    const closeButton = document.createElement("button");
    closeButton.type = "button";
    closeButton.className = "btn-close";
    closeButton.setAttribute("data-bs-dismiss", "alert");
    closeButton.setAttribute("aria-label", "Close");
    messageDiv.appendChild(closeButton);

    // Añadir el mensaje al cuerpo del documento
    document.body.appendChild(messageDiv);

    // Ocultar el mensaje después de 5 segundos
    setTimeout(() => {
        // Utilizar la clase de Bootstrap para la transición de desvanecimiento
        messageDiv.classList.remove("show");
        messageDiv.classList.add("fade-out");

        // Eliminar el mensaje del DOM después de que la transición se complete
        messageDiv.addEventListener("transitionend", () => {
            if (messageDiv.parentElement) {
                messageDiv.parentElement.removeChild(messageDiv);
            }
        });
    }, 5000); // 5000 ms = 5 segundos
}

//Cargar archvio JSON para rellenar los campos
document.getElementById('fileInputJSON').addEventListener('change', function (e) {
    var file = e.target.files[0];
    if (!file) return;

    var reader = new FileReader();
    reader.onload = function (e) {
        console.log("Contenido del archivo JSON:", e.target.result);
        try {
            var jsonData = JSON.parse(e.target.result);

            document.getElementById('inputNameHost').value = jsonData.hst_name;
            document.getElementById('inputMacHost').value = jsonData.hst_mac;
            document.getElementById('inputIpHost').value = jsonData.hst_ip;
            document.getElementById('inputUserName').value = jsonData.hst_hostname;
            document.getElementById('inputRAM').value = jsonData.hst_ram;
            document.getElementById('inputCPUHost').value = jsonData.hst_cpu;
            document.getElementById('inputAlmacenamiento').value = jsonData.hst_storage;
            document.getElementById('inputAdaptadorRed').value = jsonData.hst_network;
            document.getElementById('inputRutaSSH').value = jsonData.hst_sshroute;
            document.getElementById('inputSistemaOperativo').value = jsonData.hst_so;
        } catch (err) {
            alert('Error al cargar el archivo JSON: ' + err.message);
        }
    };
    reader.readAsText(file);
});


document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("createHost").addEventListener("submit", function (event) {
        event.preventDefault(); // Prevenir el comportamiento por defecto del formulario

        // Crear el objeto con los datos del formulario
        const data = {
            hst_name:       document.getElementById('inputNameHost').value,
            hst_mac:        document.getElementById('inputMacHost').value,
            hst_ip:         document.getElementById('inputIpHost').value,
            hst_hostname:   document.getElementById('inputUserName').value,
            hst_ram:        parseInt (document.getElementById('inputRAM').value),
            hst_cpu:        parseInt (document.getElementById('inputCPUHost').value),
            hst_storage:    parseInt (document.getElementById('inputAlmacenamiento').value),
            hst_network:    document.getElementById('inputAdaptadorRed').value,
            hst_sshroute:   document.getElementById('inputRutaSSH').value,
            hst_so:         document.getElementById('inputSistemaOperativo').value,
        };

        console.log(data)

        // Enviar los datos al servidor
        fetch("/createHost", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(err => {
                        throw new Error(err.error); // Lanzar el error recibido del servidor
                    });
                }
                return response.json();
            })
            .then(result => {
                // Mostrar mensaje de éxito
                if (result.message) {
                    document.getElementById("successMessage").innerText = result.message;
                    document.getElementById("successMessage").style.display = "block";
                }
            })
            .catch(error => {
                // Mostrar mensaje de error
                const errorMessage = document.getElementById("errorMessage");
                errorMessage.innerText = error.message;
                errorMessage.style.display = "block";
            });
    });
});
