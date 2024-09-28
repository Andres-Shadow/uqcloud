FROM alpine:3.14

# Instalar curl y jq
RUN apk add --no-cache curl jq

# Copiar el script de actualización
COPY update_prometheus_config.sh /etc/prometheus/update_prometheus_config.sh
RUN chmod +x /etc/prometheus/update_prometheus_config.sh

# Ejecutar el script al iniciar el contenedor
ENTRYPOINT ["/bin/sh", "-c", "/etc/prometheus/update_prometheus_config.sh"]
