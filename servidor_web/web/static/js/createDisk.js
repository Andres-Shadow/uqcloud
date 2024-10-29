// Función para obtener los discos y llenar la primera columna
function getDisks() {
    $.ajax({
        url: '/getDisks',  // El endpoint para obtener los discos
        method: 'GET',
        success: function(disks) {
            $('#diskTableBody').empty();  // Limpiar cualquier fila anterior
            disks.forEach(disk => {
                const diskName = disk.dsk_so_distro;
                const newRow = `<tr>
                          <td>${diskName}</td>
                          <td id="hosts-${diskName}">Cargando...</td>
                        </tr>`;
                $('#diskTableBody').append(newRow);
                getHostsForDisk(diskName);  // Llamar para obtener los hosts asociados
            });
        },
        error: function() {
            alert('Error al obtener los discos.');
        }
    });
}

// Función para obtener los hosts de un disco
function getHostsForDisk(diskName) {
    $.ajax({
        url: `/hostOfDisk/${diskName}`,
        method: 'GET',
        success: function(hosts) {
            let hostsText = 'Ninguno';
            if (hosts.length > 0) {
                hostsText = hosts.map(host => {
                    return `<span>${host.hst_name}
                                <button class="btn btn-danger btn-sm delete-host-btn" data-disk="${diskName}" data-hostid="${host.hst_id}" data-hostname="${host.hst_name}">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </span>`;
                }).join(', '); // Separa cada host con una coma
            }
            document.querySelector(`#hosts-${diskName}`).innerHTML = hostsText;

            // Asignar eventos de eliminación a los botones
            assignDeleteHostEvents();
        },
        error: function() {
            console.error(`Error al cargar los hosts para ${diskName}`);
        }
    });
}


// Función para agregar el disco a la tabla
function addDiskToTable(diskName) {
    // Crear una nueva fila para la tabla
    const newRow = `<tr>
                        <td>${diskName}</td>
                        <td id="hosts-${diskName}">Cargando...</td>
                    </tr>`;

    // Agregar la nueva fila al cuerpo de la tabla
    document.getElementById('diskTableBody').insertAdjacentHTML('beforeend', newRow);

    // Llamar a la función que obtiene los hosts asociados al disco recién creado
    getHostsForDisk(diskName);
}

// Llamar a la función para cargar los discos al cargar la página
$(document).ready(function() {
    getDisks();
});

// Función para asignar eventos de eliminación a cada botón de host
function assignDeleteHostEvents() {
    const deleteButtons = document.querySelectorAll('.delete-host-btn');
    deleteButtons.forEach(button => {
        button.addEventListener('click', function() {
            const diskName = this.getAttribute('data-disk');
            const hostId = this.getAttribute('data-hostid');
            const hostName = this.getAttribute('data-hostname');

            // Mostrar confirmación antes de eliminar
            const confirmation = confirm(`¿Está seguro de que desea eliminar el host ${hostName} del disco ${diskName}?`);
            if (confirmation) {
                deleteHostFromDisk(diskName, hostId, this);
            }
        });
    });
}

// Función para eliminar el host del disco
function deleteHostFromDisk(diskName, hostId, buttonElement) {
    // Realizar la petición de eliminación
    fetch(`/hostOfDisk/${diskName}?host_id=${hostId}`, {
        method: 'DELETE'
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => {
                    throw new Error(err.error);
                });
            }
            return response.json();
        })
        .then(result => {
            alert(`Host eliminado exitosamente: ${result.message}`);

            // Eliminar el host de la tabla visualmente
            const hostElement = buttonElement.parentNode;
            hostElement.remove();
        })
        .catch(error => {
            alert(`Error al eliminar el host: ${error.message}`);
        });
}

document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("buttonCreateDisc").addEventListener("click", function (event) {
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
                    const modalElement = document.getElementById('diskModal');
                    const modal = bootstrap.Modal.getInstance(modalElement);
                    modal.hide();

                    getDisks(data,document.getElementById("diskTableBody").childElementCount)

                    const successMessage = document.getElementById("successMessage")
                    successMessage.innerText = result.message;
                    successMessage.style.display = "block";
                    setTimeout(()=>{
                        successMessage.style.display = "none"
                    }, 5000)

                }
            })
            .catch(error => {
                // Mostrar mensaje de error
                const errorMessage = document.getElementById("errorMessage");
                errorMessage.innerText = error.message;
                errorMessage.style.display = "block";
                setTimeout(()=>{
                    errorMessage.style.display = "none";
                }, 5000)
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
                option.text = `${host.id} - ${host.hst_name}`
                selectHost.appendChild(option);
            });
        })
        .catch(error => console.error('Error fetching hosts:', error));
});


