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


document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("createForm").addEventListener("submit", function (event) {
        event.preventDefault(); // Prevenir el comportamiento por defecto del formulario

        // Crear el objeto con los datos del formulario
        const data = {
            dsk_name: document.getElementById("inputSelectName").value,
            dsk_route: document.getElementById("inputRouteDisk").value,
            dsk_so: document.getElementById("inputSelectSystem").value,
            dsk_so_distro: document.getElementById("inputSelectDistribution").value,
            dsk_arch: parseInt(document.getElementById("inputSelectArchitecture").value),
            dsk_host_id: parseInt(document.getElementById("inputSelectHost").value)
        };

        console.log(data)

        // Enviar los datos al servidor
        fetch("/createDisk", {
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


document.addEventListener("DOMContentLoaded", function() {
    fetch('/GetHost')
        .then(response => response.json())
        .then(data => {
            const selectHost = document.getElementById('inputSelectHost');

            // Limpiar opciones existentes (excepto la primera)
            while (selectHost.options.length > 1) {
                selectHost.remove(1);
            }

            // Añadir las nuevas opciones
            data.forEach(host => {
                const option = document.createElement('option');
                option.value = host.id;  // Usar ID como valor
                option.text = `${host.id} - ${host.nombre}`
                selectHost.appendChild(option);
            });
        })
        .catch(error => console.error('Error fetching hosts:', error));
});


