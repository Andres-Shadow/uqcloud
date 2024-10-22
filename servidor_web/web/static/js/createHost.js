// Función para agregar un nuevo host a la tabla
function addHostToTable(host) {
    // Obtener la instancia de DataTable
    const dataTable = $('#hostsTable').DataTable();

    // Crear una nueva fila con los datos del host
    const newRow = [
        '<input type="checkbox" class="host-checkbox" data-id="${host.id}" value="${host.id}">',
        host.id || "N/A",
        host.hst_name || "N/A",
        host.hst_hostname || "N/A",
        host.ip || "N/A",
        host.estado || "N/A",
        host.cpu_total || "N/A",
        host.almacenamiento_total || "N/A",
        host.ram_total || "N/A",
        host.hst_network || "N/A",
        host.hst_so || "N/A"
    ];

    // Agregar la nueva fila a la tabla
    dataTable.row.add(newRow).draw(false); // Agregar fila y actualizar tabla
}

// Llamar a esta función al cargar la página
function loadHosts() {
    fetch("/GetHost") // Cambia esto a la ruta donde obtienes los hosts
        .then(response => response.json())
        .then(data => {
            const tbody = document.getElementById("hostTableBody")
            tbody.innerHTML = ""
            data.forEach((host) => {
                addHostToTable(host);
            });

            if ($.fn.DataTable.isDataTable('#hostsTable')) {
                $('#hostsTable').DataTable().destroy();
            }

            // Inicializar DataTables
            $('#hostsTable').DataTable({
                "paging": true, // Activar la paginación
                "searching": false, // Activar la búsqueda
                "ordering": true, // Activar el ordenamiento
                "info": true, // Mostrar información de paginación
                "lengthMenu": [5, 10, 20] // Opciones de cantidad de filas por página
            });

        })
        .catch(error => {
            console.error("Error al cargar los hosts:", error);
        });
}

document.addEventListener("DOMContentLoaded", loadHosts);

document.getElementById('selectAll').addEventListener('change', function() {
    const checkboxes = document.querySelectorAll('.host-checkbox');
    checkboxes.forEach(checkbox => {
        checkbox.checked = this.checked;
    });
});

document.getElementById('botonEliminar').addEventListener('click', function() {
    // Obtener todos los checkboxes seleccionados
    const checkboxes = document.querySelectorAll('input[type="checkbox"]:checked');
    const hostIds = [];

    // Recorrer los checkboxes y obtener los IDs de los hosts seleccionados
    checkboxes.forEach(checkbox => {
        const row = checkbox.closest('tr');
        const id = row.querySelector('td:nth-child(2)').innerText; // Asegurarse de que este selector obtiene la columna de ID
        hostIds.push(parseInt(id));
    });

    // Verificar si hay hosts seleccionados
    if (hostIds.length === 0) {
        alert('Por favor, selecciona al menos un host para eliminar.');
        return;
    }

    // Confirmar eliminación
    const confirmDelete = confirm(`¿Estás seguro de que deseas eliminar los siguientes hosts: ${hostIds.join(', ')}?`);
    if (!confirmDelete) {
        return;
    }

    // Crear el JSON con los IDs de los hosts seleccionados
    const data = {
        hostIds: hostIds
    };

    // Hacer la petición POST a la ruta /deleteHosts
    fetch('/deleteHosts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
            alert('Los hosts seleccionados han sido eliminados exitosamente.');
            // Aquí puedes agregar código para actualizar la tabla después de eliminar los hosts
        })
        .catch((error) => {
            console.error('Error:', error);
            alert('Hubo un error al intentar eliminar los hosts. Por favor, inténtalo de nuevo.');
        });
});




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
            document.getElementById('inputIpHost').value = jsonData.hst_ip;
            document.getElementById('inputUserName').value = jsonData.hst_hostname;
            document.getElementById('inputRAM').value = jsonData.hst_ram;
            document.getElementById('inputCPUHost').value = jsonData.hst_cpu;
            document.getElementById('inputAlmacenamiento').value = jsonData.hst_storage;
            document.getElementById('inputAdaptadorRed').value = jsonData.hst_network;
            document.getElementById('inputSistemaOperativo').value = jsonData.hst_so;
        } catch (err) {
            alert('Error al cargar el archivo JSON: ' + err.message);
        }
    };
    reader.readAsText(file);
});


document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("buttonCreateHost").addEventListener("click", function (event) {
        event.preventDefault(); // Prevenir el comportamiento por defecto del formulario

        // Crear el objeto con los datos del formulario
        const data = {
            hst_name:       document.getElementById('inputNameHost').value,
            hst_ip:         document.getElementById('inputIpHost').value,
            hst_hostname:   document.getElementById('inputUserName').value,
            hst_ram:        parseInt (document.getElementById('inputRAM').value),
            hst_cpu:        parseInt (document.getElementById('inputCPUHost').value),
            hst_storage:    parseInt (document.getElementById('inputAlmacenamiento').value),
            hst_network:    document.getElementById('inputAdaptadorRed').value,
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
                    const modalElement = document.getElementById('hostModal');
                    const modal = bootstrap.Modal.getInstance(modalElement);
                    modal.hide();

                    // Agregar el nuevo host a la tabla
                    addHostToTable(data, document.getElementById("hostTableBody").childElementCount);

                    const successMessage = document.getElementById("successMessage");
                    successMessage.innerText = result.message;
                    successMessage.style.display = "block"; // Mostrar el mensaje
                    setTimeout(() => {
                        successMessage.style.display = "none"; // Ocultar después de 5 segundos
                    }, 5000); // 5000 ms = 5 segundos
                }
            })
            .catch(error => {
                const errorMessage = document.getElementById("errorMessage");
                errorMessage.innerText = error.message;
                errorMessage.style.display = "block"; // Mostrar el mensaje
                setTimeout(() => {
                    errorMessage.style.display = "none"; // Ocultar después de 5 segundos
                }, 5000); // 5000 ms = 5 segundos
            });
    });
});

