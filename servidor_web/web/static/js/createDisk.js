function showSuccessMessage() {
    var successMessage = document.getElementById('successMessage');
    successMessage.style.display = 'block';
    setTimeout(function () {
        successMessage.style.display = 'none';
    }, 3000); // Desaparecer después de 5 segundos
}

// Mostrar el mensaje de éxito al cargar la página si es necesario
window.onload = function () {
    var message = "{{ .message }}";
    if (message.trim() !== "") {
        showSuccessMessage();
    }
};



document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("createForm").addEventListener("submit", function (event) {
        event.preventDefault(); // Prevenir el comportamiento por defecto del formulario

        // Crear el objeto con los datos del formulario
        const data = {
            name: document.getElementById("inputSelectName").value,
            ruta_ubicacion: document.getElementById("inputRouteDisk").value,
            sistema_operativo: document.getElementById("inputSelectSystem").value,
            distrubucion_so: document.getElementById("inputSelectDistribution").value,
            arquitectura: parseInt(document.getElementById("inputSelectArchitecture").value),
            host_id: parseInt(document.getElementById("inputSelectHost").value)
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


