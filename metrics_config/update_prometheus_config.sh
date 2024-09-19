#!/bin/sh

PROMETHEUS_CONFIG="/etc/prometheus/prometheus.yml"
HOSTS_URL="http://procesamiento:8081/api/v1/hosts"

# Función para actualizar la configuración de Prometheus
update_prometheus_config() {
    # Obtener la lista de hosts desde el servidor de procesamiento
    response=$(curl -s "$HOSTS_URL")
    http_code=$(echo "$response" | tail -n1)
    hosts=$(echo "$response" | head -n1)

    if [ "$http_code" -ne 200 ]; then
        echo "Error: No se pudo obtener los hosts. Código HTTP: $http_code"
        return
    fi

    # Extraer las IPs de los hosts y agregar el puerto, rodear con comillas simples
    ips=$(echo "$hosts" | jq -r '.[] | .ip + ":9182"' | awk '{printf " '\''%s'\'',", $0}' | sed 's/,$//')
    # Crear el nuevo archivo de configuración
    cat > "$PROMETHEUS_CONFIG" <<EOL
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'windows'
    scrape_interval: 15s
    static_configs:
      - targets: [$ips]
EOL

    echo "Archivo prometheus.yml actualizado."
}

# Bucle infinito para actualizar cada 10 segundos
while true; do
    update_prometheus_config
    echo "Prometheus ha recargado su configuración."

    # Esperar 15 segundos antes de la próxima actualización
    sleep 15
done
