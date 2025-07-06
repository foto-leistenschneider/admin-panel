storage "file" {
  path = "/vault/data"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}

api_addr = "http://vault.admin.ergatikos.com"
cluster_addr = "https://vault.admin.ergatikos.com:8201"
ui = true

log_level = "INFO"
log_file = "/vault/logs/vault.log"
log_rotate_duration = "24h"
log_rotate_max_files = 5
