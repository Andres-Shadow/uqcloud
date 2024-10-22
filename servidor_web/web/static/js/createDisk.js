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
                          <td><input type="checkbox"></td>
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

// Función para obtener los hosts asociados a cada disco
function getHostsForDisk(diskName) {
    $.ajax({
        url: `/getHostOfDisk/${diskName}`,  // El endpoint para obtener los hosts del disco
        method: 'GET',
        success: function(hosts) {
            let hostsText = 'Ninguno';  // Valor por defecto si no hay hosts
            if (hosts.length > 0) {
                hostsText = hosts.map(host => host.hst_name).join(', ');  // Unir los nombres de los hosts
            }
            $(`#hosts-${diskName}`).text(hostsText);  // Actualizar la columna de hosts
        },
        error: function() {
            $(`#hosts-${diskName}`).text('Error al cargar los hosts');
        }
    });
}

// Llamar a la función para cargar los discos al cargar la página
$(document).ready(function() {
    getDisks();
});



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


