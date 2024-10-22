CREATE DATABASE IF NOT EXISTS uqcloud;  -- Crea la base de datos si no existe
CREATE USER 'grafana'@'%' IDENTIFIED BY 'grafana';
GRANT SELECT ON uqcloud.* TO 'grafana'@'%';  -- Otorga permisos sobre la base de datos creada
FLUSH PRIVILEGES;