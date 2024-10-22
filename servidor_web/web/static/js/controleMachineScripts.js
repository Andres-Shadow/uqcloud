var ventanaConfiguracionAbierta = false;

function abrirVentanaEmergenteEliminacion(nombre) {
    ventanaConfiguracionAbierta = true;

    var VentanaEmergenteEliminacion = document.getElementById("VentanaEmergenteEliminacion");

    document.getElementById("vmnameDelete").value = nombre;

    VentanaEmergenteEliminacion.style.display = "block";
}


function abrirVentanaEmergenteInformacion(nombre, sistemaOperativo, distribucion, ip, ram, cpu, estado, hostname) {
    ventanaConfiguracionAbierta = true;

    var ventanaEmergente = document.getElementById("VentanaEmergenteInformacion");
    ventanaEmergente.style.display = "block";

    // Asignar valores a los elementos <span> con sus respectivos ids

    // Nombre de la Máquina
    document.getElementById("nombreSpan").textContent = nombre;

    // Sistema Operativo de la Máquina
    document.getElementById("sistemaOperativoSpan").textContent = "GNU/" + sistemaOperativo;

    // Distribución de la Máquina
    document.getElementById("distribucionSpan").textContent = distribucion;

    // Modelo de CPU
    if (cpu == 1) {
        document.getElementById("cpuSpan").textContent = cpu + " Núcleo";
    } else {
        document.getElementById("cpuSpan").textContent = cpu + " Núcleos";

    }

    // Memoria de la Máquina
    document.getElementById("memoriaSpan").textContent = ram + " MB";

    // Estado de la Máquina
    document.getElementById("estadoSpan").textContent = estado;

    document.getElementById("hostnameSpan").textContent = "uqcloud";

    document.getElementById("passwordSpan").textContent = "uqcloud1234";

    // IP de la Máquina
    if (ip != "") {
        document.getElementById("ipSpan").textContent = ip;
        document.getElementById("urlSpan").textContent = "http://" + ip + ":4200"
    } else {
        document.getElementById("ipSpan").textContent = "No asignada";
        document.getElementById("urlSpan").textContent = "No asignada"
    }
}

function cerrarVentanaEmergenteConfiguracion() {
    ventanaConfiguracionAbierta = false;

    var ventanaEmergente = document.getElementById("VentanaEmergenteConfiguracion");
    ventanaEmergente.style.display = "none";
}

function cerrarVentanaEmergenteInformacion() {
    ventanaConfiguracionAbierta = false;

    var ventanaEmergente = document.getElementById("VentanaEmergenteInformacion");
    ventanaEmergente.style.display = "none";
}

function cerrarVentanaEmergenteEliminar() {
    ventanaConfiguracionAbierta = false;

    var VentanaEmergenteEliminacion = document.getElementById("VentanaEmergenteEliminacion");
    VentanaEmergenteEliminacion.style.display = "none";
}

function actualizarTabla() {

    if (ventanaConfiguracionAbierta) {
        return;
    }
    $.ajax({
        url: "/api/machines",
        method: "GET",
        dataType: "json",
        success: function (data) {
            $("#machine-table tbody").empty();

            data.forEach(function (machine) {
                let backgroundColor;
                switch (machine.vm_state) {
                    case "Apagado":
                        backgroundColor = "#e06666ff";
                        break;
                    case "Encendido":
                        backgroundColor = "#93c47dff";
                        break;
                    case "Procesando":
                        backgroundColor = "#83DEE3";
                        break;
                    default:
                        backgroundColor = "";
                }

                const ipDisplay = machine.vm_ip ? machine.vm_state === 'Encendido' ? `${machine.vm_ip} <a href="http://${machine.vm_ip}:4200" target="_blank"><i class="fa-solid fa-up-right-from-square ms-2"></i></a>` : `${machine.vm_ip}` : "No asignada";

                const actionButtons = `
                    <div class="">
                        <button class="btn btn-link p-0 me-1" onclick="changeStateMachine('${machine.vm_name}', '${machine.vm_state === 'Apagado' ? 'startMachine' : 'stopMachine'}')">
                            <i class="fa ${machine.vm_state === 'Apagado' ? 'fa-power-off' : 'fa-stop-circle'} icon-large" title="${machine.vm_state === 'Apagado' ? 'Start' : 'Stop'}"></i>
                        </button>
                        <button class="btn btn-link p-0 me-1" onclick="abrirVentanaEmergenteInformacion('${machine.vm_name}', '${machine.vm_so}', '${machine.vm_so_distro}', '${machine.vm_ip}', '${machine.vm_ram}', '${machine.vm_cpu}', '${machine.vm_state}', '${machine.vm_hostname}')">
                            <i class="fa fa-info-circle icon-large" title="Info"></i>
                        </button>
                        <button class="btn btn-link p-0 me-1" onclick="abrirVentanaEmergenteEliminacion('${machine.vm_name}')" ${machine.vm_state !== 'Apagado' ? 'disabled' : ''}>
                            <i class="fa fa-trash icon-large" title="Delete"></i>
                        </button>
                    </div>`;

                const row = `
                    <tr style="background-color: ${backgroundColor}">
                        <td>${machine.vm_name}</td>
                        <td>${ipDisplay}</td>
                        <td>${machine.vm_so_distro}</td>
                        <td>${machine.vm_state}</td>
                        <td class="button-column">${actionButtons}</td>
                    </tr>`;

                $("#machine-table tbody").append(row);
            });
        },
        error: function (error) {
            $("#machine-table tbody").empty();
            console.error("Error al obtener datos: " + error);
        }
    });
}

function changeStateMachine(vm_name, state) {
    actualizarTabla();
    fetch('/api/' + state, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', },
        body: JSON.stringify({ vm_name: vm_name, }),
    })
        .then(response => {
            if (response.status === 200) {
                return response.json();
            } else if (response.status === 400 || response.status === 500) {
                return response.json();
            } else {
                throw new Error('Error en el servidor web | response status raro');
            }
        })
        .then(data => {
            if (data.SuccessMessage) {
                const successMessage = data.SuccessMessage;
                showAlert(successMessage, "success");
            } else if (data.ErrorMessage) {
                const errorMessage = data.ErrorMessage;
                showAlert(errorMessage, "danger");
            }
        })
        .catch(error => {
            showAlert("Error al realizar la solicitud al servidor: " + error, "danger");
            console.error('Error: ' + error);
        })
}

setInterval(actualizarTabla, 4000);