var ventanaConfiguracionAbierta = false;

    function abrirVentanaEmergente() {
        ventanaConfiguracionAbierta = true;

        var ventanaEmergente = document.getElementById("ventanaEmergente");
        ventanaEmergente.style.display = "block";
    }

    function abrirVentanaEmergenteEliminacion( nombre ){
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
        if ( cpu == 1 ){
            document.getElementById("cpuSpan").textContent = cpu + " Núcleo";
        }else {
            document.getElementById("cpuSpan").textContent = cpu + " Núcleos";

        }

        // Memoria de la Máquina
        document.getElementById("memoriaSpan").textContent = ram + " MB";

        // Estado de la Máquina
        document.getElementById("estadoSpan").textContent = estado;

        document.getElementById("hostnameSpan").textContent = "uqcloud";

        document.getElementById("passwordSpan").textContent = "uqcloud1234";

        // IP de la Máquina
        if  ( ip != ""){
            document.getElementById("ipSpan").textContent = ip;
            document.getElementById("urlSpan").textContent = "http://"+ip+":4200"
        }else{
            document.getElementById("ipSpan").textContent = "No asignada";
            document.getElementById("urlSpan").textContent = "No asignada"
        }
    }

    function cerrarVentanaEmergente() {
        ventanaConfiguracionAbierta = false;

        var ventanaEmergente = document.getElementById("ventanaEmergente");
        ventanaEmergente.style.display = "none";
    }

    function cerrarVentanaEmergenteConfiguracion() {
        ventanaConfiguracionAbierta = false;

        var ventanaEmergente = document.getElementById("VentanaEmergenteConfiguracion");
        ventanaEmergente.style.display = "none";
    }

    function cerrarVentanaEmergenteInformacion(){
        ventanaConfiguracionAbierta = false;

        var ventanaEmergente = document.getElementById("VentanaEmergenteInformacion");
        ventanaEmergente.style.display = "none";
    }

    function cerrarVentanaEmergenteEliminar(){
        ventanaConfiguracionAbierta = false;

        var VentanaEmergenteEliminacion = document.getElementById("VentanaEmergenteEliminacion");
        VentanaEmergenteEliminacion.style.display = "none";
    }

    function actualizarTabla() {
        
        if (ventanaConfiguracionAbierta ) {
            return;
        }
        $.ajax({
            url: "/api/machines", // Reemplaza con la URL correcta
            method: "GET",
            dataType: "json",
            success: function(data) {
                // Limpia la tabla actual
                $("#machine-table tbody").empty();
                // Itera a través de los datos y agrega filas a la tabla
                data.forEach(function(machine) {
                    // console.log("machineestado: ", machine.vm_state);
                    // RECORDAR: el JSON que llega con los machine, viene con el modelo Models.VirtualMachine, como este modelo
                    //           usa 'TAGS' de tipo JSON para codificar y decodificar, para acceder a esos valores se debe usar 
                    //           su tag respectivo. Ej: 
                    //              * machine.Name         | NO SIRVE 
                    //              * machine.vm_name      | SIRVE  (porque en el modelo está asi: `json:"vm_name"`)
                    switch (machine.vm_state) {
                        case "Apagado":
                            // console.log("LLegó aqui-----------------1");
                            backgroundColor = "#e06666ff"; // Rojo
                            $("#machine-table tbody").append(
                            `<tr style="background-color: ${backgroundColor}">
                                <td>${machine.vm_name}</td>
                                <td>${machine.vm_ip === "" ? "No asignada" : machine.vm_ip}</td>
                                <td>${machine.vm_so_distro}</td>
                                <td>${machine.vm_state}</td>
                                
                                <td class="button-column">
                                    <button onclick="changeStateMachine('${machine.vm_name}', 'startMachine')" class="btn btn-link" style="padding: 0; margin: 0;">
                                        <img style="width: 35px;" src="/web/static/img/icons/power.png" alt="Botón 1">
                                    </button>
                                    <button type="button" class="btn btn-link" style="padding: 0; margin: 0;" onclick="abrirVentanaEmergenteInformacion('${machine.vm_name}','${machine.vm_so}','${machine.vm_so_distro}', '${machine.vm_ip}','${machine.vm_ram}','${machine.vm_cpu}', '${machine.vm_state}', '${machine.vm_hostname}')">
                                        <img style="width: 35px;" src="/web/static/img/icons/info.png" alt="Botón 4">
                                    </button>
                                    <button type="button" class="btn btn-link" style="padding: 0; margin: 0;" onclick="abrirVentanaEmergenteEliminacion('${machine.vm_name}')">
                                        <img style="width: 30px;" src="/web/static/img/icons/delete.png" alt="Botón 3">
                                    </button>
                                </td>
                                </td>
                            </tr>`
                            );
                            break;
                        case "Encendido":
                            // console.log("LLegó aqui-----------------2");
                            backgroundColor = "#93c47dff"; // Verde
                            $("#machine-table tbody").append(
                            `<tr style="background-color: ${backgroundColor}">
                                <td>${machine.vm_name}</td>
                                <td style="position: relative;">
                                    ${machine.vm_ip === "" ? "No asignada" : machine.vm_ip}
                                    
                                    <!-- Muestra el botón solo si la IP está asignada -->
                                    ${machine.vm_ip !== "" ? `

                                    <a href="http://${machine.vm_ip}:4200" target="_blank">
                                        <i class="fa-solid fa-up-right-from-square"></i>
                                    </a>                           
                                    ` : ''}
                                </td>

                                <td>${machine.vm_so_distro}</td>
                                <td>${machine.vm_state}</td>
                                
                                <td class="button-column">
                                    <button onclick="changeStateMachine('${machine.vm_name}', 'stopMachine')" class="btn btn-link" style="padding: 0; margin: 0;">
                                        <img style="width: 35px;" src="/web/static/img/icons/power.png" alt="Botón 1">
                                    </button>
                                    <button type="button" class="btn btn-link" style="padding: 0; margin: 0;" onclick="abrirVentanaEmergenteInformacion('${machine.vm_name}','${machine.vm_so}','${machine.vm_so_distro}', '${machine.vm_ip}','${machine.vm_ram}','${machine.vm_cpu}', '${machine.vm_state}', '${machine.vm_hostname}', '${machine.vm_hostname}')">
                                        <img style="width: 35px;" src="/web/static/img/icons/info.png" alt="Botón 3">
                                    </button>                                    
                                    <button type="submit" class="btn btn-link" style="padding: 0; margin: 0;" disabled>
                                        <img style="width: 30px;" src="/web/static/img/icons/delete.png" alt="Botón 4">
                                    </button>                                    
                                </td>
                            </tr>`
                            );
                            break;
                        case "Procesando":
                            // console.log("LLegó aqui-----------------3");
                            backgroundColor = "#83DEE3"; // Azul
                            $("#machine-table tbody").append(
                            `<tr style="background-color: ${backgroundColor}">
                                <td>${machine.vm_name}</td>
                                <td>${machine.vm_ip === "" ? "No asignada" : machine.vm_ip}</td>
                                <td>${machine.vm_so_distro}</td>
                                <td>${machine.vm_state}</td>
                            
                                <td class="button-column">
                                    <button onclick="changeStateMachine('${machine.vm_name}', 'startMachine')" class="btn btn-link" style="padding: 0; margin: 0;" disabled>
                                        <img style="width: 35px;" src="/web/static/img/icons/power.png" alt="Botón 1">
                                    </button>
                                    <button type="button" class="btn btn-link" style="padding: 0; margin: 0;" onclick="abrirVentanaEmergenteInformacion('${machine.vm_name}','${machine.vm_so}','${machine.vm_so_distro}', '${machine.vm_ip}','${machine.vm_ram}','${machine.vm_cpu}', '${machine.vm_state}', '${machine.vm_hostname}')">
                                        <img style="width: 35px;" src="/web/static/img/icons/info.png" alt="Botón 4">
                                    </button>                                
                                    <button class="btn btn-link" style="padding: 0; margin: 0;" disabled>
                                        <img style="width: 30px;" src="/web/static/img/icons/delete.png" alt="Botón 3">
                                    </button>
                                </td>
                            </tr>`
                            );
                            break;
                        default:
                            // console.log("LLegó aqui-----------------4");
                            backgroundColor = ""; // Puedes proporcionar un valor predeterminado si es necesario
                    }
     
                });                
            },
            error: function(error) {
                console.error("Error al obtener datos: " + error);
            }
        });

    }

    function copiarText(texto) {
        // Crea un elemento de entrada temporal
        const tempInput = document.createElement("input");
        tempInput.value = texto;
        document.body.appendChild(tempInput);
        tempInput.select();

        // Intenta copiar el texto al portapapeles
        document.execCommand("copy");

        // Elimina el elemento de entrada temporal
        document.body.removeChild(tempInput);
    }    

    function changeStateMachine(vm_name, state) {
        fetch('/api/'+state, {
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
    

    // Llama a actualizarTabla al cargar la página y periódicamente para mantener los datos actualizados
    actualizarTabla();
    setInterval(actualizarTabla, 4000);
