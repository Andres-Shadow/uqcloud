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
    const messageBox = document.getElementById(type === 'success' ? 'successMessage' : 'errorMessage');
    messageBox.textContent = message;
    messageBox.style.display = 'block'; // Muestra el mensaje

    // Oculta el mensaje después de 5 segundos
    setTimeout(() => {
        messageBox.style.display = 'none';
    }, 5000);
}

//Cargar archvio JSON para rellenar los campos
document.getElementById('fileInputJSON').addEventListener('change', function (e) {
    var file = e.target.files[0];
    if (!file) return;

    var reader = new FileReader();
    reader.onload = function (e) {
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
            hst_ram:        document.getElementById('inputRAM').value,
            hst_cpu:        document.getElementById('inputCPUHost').value,
            hst_storage:    document.getElementById('inputAlmacenamiento').value,
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
