<div align="center">

`vops` - A HashiCorp Vault cluster management tool

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
`vops` comes as `RPM`, `DEB`, `APK`, Container and CLI-tool:

```bash
# using curl
curl ...
# go
go get github.com/FalcoSuessgott/vops
vops status

# docker/podman
docker pull ghcr.io/falcosuessgott/vops
docker run falcosuessgott/vops -v `$(PWD):/ooo` status


...
```

# Quickstart
Imagine the following `vops.yaml` configuration file:

```yaml
cluster:
  - name: dev-cluster
    addr: "http://192.168.0.100:8200"
    tokenExecCmd: "cat dev-cluster-token-file.json"
    keyfilePath: "dev-cluster-token-file.json"
    snapshotDirectory: "./snapshots/"
    nodes:
      - "http://192.168.0.100:8200"
    extraEnv:
     VAULT_SKIP_VERIFY: true
```

using `vops` ability to use Go-Template Syntax that file could be rewritten as:

```yaml
cluster:
  - name: dev-cluster
    addr: "http://192.168.0.100:8200"
    tokenExecCmd: "cat {{ .KeyFilePath }}"
    keyfilePath: "{{ .Name }}-token-file.json"
    snapshotDirectory: "./snapshots/"
    nodes:
      - "{{ .Addr }}"
    extraEnv:
     VAULT_SKIP_VERIFY: true
```

You will understand how this makes maintaining multiple cluster more convenient.

`vops` also allows to render any existing environment variable for example:

```yaml
cluster:
  - name: dev-cluster
    keyfilePath: "{{ .Env.PWD}}/{{ .Name }}-token-file.json"
    snapshotDirectory: "{{ .Env.Home }}/snapshots/"
```

`vops` looks for a `vops.yaml` in the current working directory. You can specify a file by using the `--config` CLI arg or setting the `VOPS_CONFIG` env var.

# Operations
## Initialization

```bash
> vops init -c cluster-1
[ VOPS ]
using vops.yaml
attempting intialization of cluster "cluster-1" with 5 shares and a threshold of 3
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
successfully initialized cluster-1 and wrote keys to cluster-1.json.
```

**Tip:** You can also specify the used cluster be providing a environment variable named: `VOPS_CLUSTER`
## Unsealing
```bash
> vops unseal -c cluster-1    
[ VOPS ]
using vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
using keyfile "cluster-1.json"
unsealing node "http://127.0.0.1:8200"
unsealing node "http://127.0.0.1:8200"
unsealing node "http://127.0.0.1:8200"
cluster "cluster-1" unsealed
```

## Sealing
```bash
> vops seal -c cluster-1  
[ VOPS ]
using vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
executed token exec command
cluster "cluster-1" sealed
```

## Rekey
```bash
> vops rekey -c cluster-1            
[ VOPS ]
using vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
using keyfile "cluster-1.json"
rekeying successfully completed
renamed keyfile "cluster-1.json" for cluster "cluster-1" to "cluster-1.json_20230210233133" (snapshots depend on the unseal/recovery keys from the moment the snapshot has been created. This way you always have the matching unseal/recovery keys ready.)
```


## Generate Root Token
```bash
> vops  generate-root -c cluster-1        
[ VOPS ]
using vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
generated on OTP for root token creation
started root token generation process
root token generation completed
new root token: "hvs.byNOU9DVxCbvgatIMHAwXOKS"
make sure to uvopspdate your token exec commands in your vops configfile if necessary.
```

## Snapshots
### Snapshot save
```bash
> vops snapshot save -c cluster-1
[ VOPS ]
using vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
executed token exec command
created snapshot file "snapshots/20230210232954" for cluster "cluster-1"
```

## Custom Commands
You can define custom commands:

```yml
CustomCmds:
  list-peers: 'vault operator raft list-peers'
  status: 'vault status'

Cluster:
  - Name: cluster-1
    Addr: "http://127.0.0.1:8200"
    TokenExecCmd: "jq -r '.root_token' {{ .Keys.Path }}"
    Keys:
      Path: "{{ .Name }}.json"
      Shares: 1
      Threshold: 1
    SnapshotDirectory: "snapshots/"
    Nodes:
      - "{{ .Addr }}"
    ExtraEnv:
     VAULT_SKIP_VERIFY: true
     VAULT_TLS_CA: "ok"
```

and run them for all cluster:

```bash
$> vops custom --list
[ Custom ]
using ./assets/vops.yaml

[ Available Commands ]
"list-peers": "vault operator raft list-peers"
"status": "vault status"

run any available command with "vops custom -x <command name> -c <cluster-name>".

$> vops custom -x custom -x status --all-cluster  
[ Custom ]
using ./assets/vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
applying VAULT_ADDR
applying VAULT_TOKEN
token exec command successful

$> vault status
Key                     Value
---                     -----
...

[ cluster-2 ]
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
applying VAULT_ADDR
applying VAULT_TOKEN
token exec command successful

$> vault status
Key                     Value
---                     -----
....
```

# Global Flags
* `--config` (`VOPS_CONFIG`) - `vops` configuration file (default: `vops.yaml`)
* `--all-cluster` (`-A`) - perform the chosen operation for all defined clusters.
