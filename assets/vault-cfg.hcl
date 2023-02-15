cluster_addr = "http://127.0.0.1:8201"
api_addr = "http://127.0.0.1:8200"

storage "raft" {
  path = "./assets/raft"
  node_id = "node"
}

ui = true

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}
