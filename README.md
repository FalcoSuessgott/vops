<div align="center">

<h2> A HashiCorp Vault Cluster Management tool </h2>

<img src="assets/demo.gif" alt="drawing"/>

<img src="https://github.com/FalcoSuessgott/vops/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vops/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vops/branch/main/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vops/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vops" alt="drawing"/>
<img src="https://img.shields.io/docker/pulls/falcosuessgott/vops" alt="drawing"/>

</div>


***`vops` is in very early stage and is likely to change***


# Background
I automate, develop and maintain a lot of Vault cluster for different clients. When automating Vault using tools such as `terraform` and `ansible` I was missing a small utility that allows me to quickly perform certain operations like generate a new root token or create a snapshot. Thus I came up with `vops`, which stands for **v**ault-**op**eration**s**

# Features 
* define as many vault cluster as you need
* template your `vops.yaml` and be able to use clever naming convetions
* Iterate over all defined cluster for every supported option
* Initialize a Vault 
* Seal & Unseal a Vault 
* Rekey a Vault 
* Generate a new root token
* save and restore a Vault (raft storage required)
* define custom commands then can be run for any cluster

# Installation
```bash
# curl
version=$(curl -S "https://api.github.com/repos/FalcoSuessgott/vops/releases/latest" | jq -r '.tag_name[1:]')
curl -OL "https://github.com/FalcoSuessgott/vops/releases/latest/download/vops_${version}_$(uname)_$(uname -m).tar.gz"
tar xzf "vops_${version}_$(uname)_$(uname -m).tar.gz"
./vops version

# Go 
go install github.com/FalcoSuessgott/vops@latest
vops version

# Docker/Podman
docker run ghcr.io/falcosuessgott/vops version

# Ubuntu/Debian
version=$(curl -S "https://api.github.com/repos/FalcoSuessgott/vops/releases/latest" | jq -r '.tag_name[1:]')
curl -OL "https://github.com/FalcoSuessgott/vops/releases/latest/download/vops_${version}_linux_amd64.deb"
sudo dpkg -i "./vops_${version}_linux_amd64.deb"
vops version

# Fedora/CentOS/RHEL
version=$(curl -S "https://api.github.com/repos/FalcoSuessgott/vops/releases/latest" | jq -r '.tag_name[1:]')
curl -OL "https://github.com/FalcoSuessgott/vops/releases/latest/download/vops_${version}_linux_amd64.rpm"
sudo dnf localinstall "./vops_${version}_linux_amd64.rpm"
vops version

# Sources
git clone https://github.com/FalcoSuessgott/vops && cd vops
go build 
```

# Usage
`vops` looks for a `vops.yaml` configuration file in your `$PWD`, you change the location by setting `VOPS_CONFIG`.

`vops` allows you to use templates and environment variables in your configuration file: 

<table>
<tr>
<td> Default </td> <td> Templated </td>
</tr>
<tr>
<td>

```yaml
Cluster:
  - Name: dev-cluster
    Addr: "http://192.168.0.100:8200"
    TokenExecCmd: "jq -r '.root_token' dev-cluster-token-file.json"
    Keys:
      Path: "dev-cluster-token-file.json"
    SnapshotDirectory: "/home/user/snapshots/"
    Nodes:
      - "http://192.168.0.100:8200"
    ExtraEnv:
     VAULT_SKIP_VERIFY: true

  - Name: prod-cluster
    ...
```

</td>
<td> 

```yaml
Cluster:
  - Name: dev-cluster
    Addr: "http://192.168.0.100:8200"
    TokenExecCmd: "jq -r '.root_token' {{ .Keys.Path }}"
    Keys:
      Path: "{{ .Name }}-token-file.json"
    SnapshotDirectory: "{{ .Env.HOME }}/snapshots/"
    Nodes:
      - "{{ .Addr }}"
    ExtraEnv:
     VAULT_SKIP_VERIFY: true

  - Name: prod-cluster
    ....
```
</td>
</tr>
</table>

# Quickstart
## Prerequisites 
Start a Vault with Integrated Storage locally:

```bash
mkdir raft
cat <<EOF > vault-cfg.hcl
cluster_addr = "http://127.0.0.1:8201"
api_addr = "http://127.0.0.1:8200"

storage "raft" {
  path = "./raft"
  node_id = "node"
}

ui = true

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}
EOF
vault server -config=vault-cfg.hcl
```

## `vops.yaml`
In another Terminal create a `vops.yaml` example config:

```bash
vops config example > vops.yaml
cat vops.yaml
```

## Initialize
> initialize vault cluster 
```bash
$> vops init --cluster cluster-1
[ Intialization ]
using vops.yaml

[ cluster-1 ]
attempting intialization of cluster "cluster-1" with 1 shares and a threshold of 1
applying VAULT_TLS_SKIP_VERIFY
successfully initialized cluster-1 and wrote keys to cluster-1.json.
```

**Tip:** You can also specify the cluster by providing a environment variable named `VOPS_CLUSTER` or run the command for all cluster using `-A` or `--all-cluster`.

## Unseal
> unseal a vault cluster using the specified keyfile
```bash
> vops unseal --cluster cluster-1
[ Unseal ]
using vops.yaml

[ cluster-1 ]
applying VAULT_TLS_SKIP_VERIFY
using keyfile "cluster-1.json"
unsealing node "http://127.0.0.1:8200"
cluster "cluster-1" unsealed
```

## Seal
> seal a cluster
```bash
> vops seal --cluster cluster-1
[ Seal ]
using vops.yaml

[ cluster-1 ]
applying VAULT_TLS_SKIP_VERIFY
executed token exec command
cluster "cluster-1" sealed
```

## Rekey
tbd. 


## Generate Root
> generates a new root token
```bash
> vops generate-root --cluster cluster-1
[ Generate Root Token ]
using vops.yaml

[ cluster-1 ]
applying VAULT_TLS_SKIP_VERIFY
generated on OTP for root token creation
started root token generation process
root token generation completed
new root token: "hvs.dmhO9aVPT0aBB1G7nrj3UdDh" (make sure to update your token exec commands in your vops configfile if necessary.)
```

## Snapshots
### Snapshot save
```bash
> vops snapshot save --cluster cluster-1
[ Snapshot Save ]
using vops.yaml

[ cluster-1 ]
applying VAULT_TLS_SKIP_VERIFY
executed token exec command
created snapshot file "cluster-1/20230216155514" for cluster "cluster-1"
```

### Snapshot Restore
tbd.
## Custom Commands
You can run any defined custom commands:

```bash
> vops custom --list
[ Custom ]
using vops.yaml

[ Available Commands ]
"list-peers": "vault operator raft list-peers"
"status": "vault status"

run any available command with "vops custom -x <command name> -c <cluster-name>".

> vops custom -x status --cluster --cluster-1
[ Custom ]
using vops.yaml

[ cluster-1 ]
applying VAULT_TLS_SKIP_VERIFY
applying VAULT_ADDR
applying VAULT_TOKEN
token exec command successful

$> vault status
Key                     Value
---                     -----
Seal Type               shamir
Initialized             true
Sealed                  false
Total Shares            1
Threshold               1
Version                 1.12.1
Build Date              2022-10-27T12:32:05Z
Storage Type            raft
Cluster Name            vault-cluster-982e7d76
Cluster ID              10939fbf-fdfd-04b3-e037-dd044ce38fa3
HA Enabled              true
HA Cluster              https://127.0.0.1:8201
HA Mode                 active
Active Since            2023-02-16T14:54:35.541380933Z
Raft Committed Index    36
Raft Applied Index      36
```