<div align="center">
<h2> A HashiCorp Vault Cluster Management tool </h2>
</div>
<table>
<tr>
<td> Usage </td> <td> vops config </td>
</tr>
<tr>
<td>

```bash
# configure
VOPS_CONFIG="./vops.yaml"
VOPS_CLUSTER="vault-dev" 
# print example config
vops config example
# list cluster and render configs
vops config validate
# initialize
vops init --cluster <cluster> # -A for all Cluster
# unseal
vops unseal -c <cluster> # optional: -n <node>
# seal
vops seal -c <cluster>
# generate root token
vops generate-root -c <cluster>
# rekey unseal/recovery keys
vops rekey -c <cluster>
# save/restory snapshots
vops snapshot save -c <cluster>
vops snapshot restore -c <cluster>
# custom commands
vops custom -c <cluster> -x <custom command>
# adhoc commands
vops adhoc -x "vault status" -c <cluster>
# open UI
vops ui -c <cluster>
# copy token
vops token -c <cluster>
# vault login 
vops login -c <cluster>
```

</td>
<td> 

```yaml
Cluster:
  - Name: vault-dev
    Addr: "http://127.0.0.1:8200"
    TokenExecCmd: "jq -r '.root_token' {{ .Keys.Path }}"
    Keys:
      Path: "{{ .Name }}.json" # https://github.com/FalcoSuessgott/vops#about-the-keypath-file
      Autounseal: false  # true if autounseal configured
      Shares: 1 # optional, default: 5
      Threshold: 1 # optional: default: 3
    SnapshotDirectory: "snapshots/"
    Nodes:
      - "{{ .Addr }}"
    ExtraEnv:
      VAULT_SKIP_VERIFY: true

  - Name: vault-prod
    Addr: "https://{{ .Name }}.example.com:8200"
    TokenExecCmd: "vault login $(cat ${{ .ENV.HOME }}/.vault_token)"
    Keys:
      Path: "{{ .Name }}.json"
      Shares: 5
      Threshold: 3
    SnapshotDirectory: "{{ .ENV.HOME }}/snapshots/"
    Nodes:
      - "{{ .Name }}-01.example.com:8200"
      - "{{ .Name }}-02.example.com:8200"
      - "{{ .Name }}-03.example.com:8200"
      - "{{ .Name }}-04.example.com:8200"
      - "{{ .Name }}-05.example.com:8200"

CustomCmds:
  list-peers: 'vault operator raft list-peers'
  status: 'vault status'
```
</td>
</tr>
</table>
<div align="center">

<img src="https://github.com/FalcoSuessgott/vops/actions/workflows/test.yml/badge.svg" alt="drawing"/>
<img src="https://github.com/FalcoSuessgott/vops/actions/workflows/lint.yml/badge.svg" alt="drawing"/>
<img src="https://codecov.io/gh/FalcoSuessgott/vops/branch/main/graph/badge.svg" alt="drawing"/>
<img src="https://img.shields.io/github/downloads/FalcoSuessgott/vops/total.svg" alt="drawing"/>
<img src="https://img.shields.io/github/v/release/FalcoSuessgott/vops" alt="drawing"/>
<img src="https://img.shields.io/docker/pulls/falcosuessgott/vops" alt="drawing"/>
</div>


***`vops` is in early stage and is likely to change***


# Background
I automate, develop and maintain a lot of Vault cluster for different clients. When automating Vault using tools such as `terraform` and `ansible` I was missing a small utility that allows me to quickly perform certain operations like generate a new root token or create a snapshot. Thus I came up with `vops`, which stands for **v**ault-**op**eration**s**

# Features 
* template your `vops.yaml` using [clever naming conventions](https://github.com/FalcoSuessgott/vops#usage)
* Iterate over all defined cluster (`--all-cluster/-A`) for every supported option
* [Initialize](https://github.com/FalcoSuessgott/vops#initialize) a Vault 
* [Seal](https://github.com/FalcoSuessgott/vops#seal) & [Unseal](https://github.com/FalcoSuessgott/vops#unseal) a Vault 
* [Generate a new root token](https://github.com/FalcoSuessgott/vops#generate-root)
* [save](https://github.com/FalcoSuessgott/vops#snapshot-save) a Vault (raft storage required) Snapshot
* [open the UI](https://github.com/FalcoSuessgott/vops#ui) in your default browser
* [perform a vault login](https://github.com/FalcoSuessgott/vops#login) to a specified cluster in order to continue working with the vault CLI
* [copy the token](https://github.com/FalcoSuessgott/vops#token) from a the token exec command to your clipboard buffer for Vault UI login
* define [custom commands](https://github.com/FalcoSuessgott/vops#custom-commands) then can be run for any cluster
* run [adhoc commands](https://github.com/FalcoSuessgott/vops#adhoc-commands)

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

created the following `vops.yaml`:

```yaml
Cluster:
  - Name: cluster-1
    Addr: http://127.0.0.1:8200
    TokenExecCmd: jq -r '.root_token' {{ .Keys.Path }}
    Keys:
      Path: '{{ .Name }}.json'
      Shares: 1
      Threshold: 1
    SnapshotDirectory: '{{ .Name }}/'
    Nodes:
      - '{{ .Addr }}'
    ExtraEnv:
      VAULT_TLS_SKIP_VERIFY: true
CustomCmds:
  list-peers: vault operator raft list-peers
  status: vault status
```

## List Cluster
> lists all defined clusters and validates their settings, performs a vault login
```bash
$> vops config validate      
[ Validate ]
using ./assets/vops.yaml

Name:                   cluster-1
Address:                http://127.0.0.1:8200
TokenExecCmd:           jq -r '.root_token' cluster-1.json
TokenExecCmd Policies:  [root]
Nodes:                  [http://127.0.0.1:8200]
Key Config:             {Path: cluster-1.json, Shares: 5, Threshold: 5}
Snapshot Directory:     snapshots/
```


## Initialize
> initialize vault cluster 
```bash
$> vops init -c <cluster>
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
> vops unseal -c <cluster>
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
> vops seal -c <cluster>
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
> vops generate-root -c <cluster>
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
> vops snapshot save -c <cluster>
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

> vops custom -x <custom command> -c <cluster>
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

## Adhoc Commands
> similar to Ansible, run a adhoc command for a cluster
```bash
vops adhoc -x "vault status" -c <cluster>
[ Adhoc ]
reading ./assets/vops.yaml

[ cluster-1 ]
applying VAULT_SKIP_VERIFY
applying VAULT_ADDR
applying VAULT_TOKEN
token exec command successful

$> vault status
Key                     Value
---                     -----
Seal Type               shamir
Initialized             true
Sealed                  false
Total Shares            5
Threshold               5
Version                 1.12.1
Build Date              2022-10-27T12:32:05Z
Storage Type            raft
Cluster Name            vault-cluster-39a4f8c3
Cluster ID              c9f326fa-9e59-b527-04a0-c4270995cc7f
HA Enabled              true
HA Cluster              https://127.0.0.1:8201
HA Mode                 active
Active Since            2023-03-02T17:51:12.542083099Z
Raft Committed Index    38
Raft Applied Index      38
```

## UI
> opens the Vault Address in your default browser

```bash
vops ui -c <cluster>
[ UI ]
using ./assets/vops.yaml

[ cluster-1 ]
opening http://127.0.0.1:8200
```

## Login 
> performs a vault token login command in order to work with the vault CLI

```bash
vops login -c <cluster>
[ Login ]
using ./assets/vops.yaml

[ cluster-1 ]
performing a vault login to http://127.0.0.1:8200
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
applying VAULT_ADDR
applying VAULT_TOKEN
executed token exec command

$> vault login $(jq -r '.root_token' cluster-1.json)
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                  Value
---                  -----
token                hvs.5aYsklTzIYR26iWRttqFoCm1
token_accessor       R7l5fidrohVB5Tc4pXBcUtoO
token_duration       ∞
token_renewable      false
token_policies       ["root"]
identity_policies    []
policies             ["root"]
```

## Token
> copy the token from the token exec command to your clipboard buffer

```bash
vops token -c <cluster>

[ Token ]
using ./assets/vops.yaml

[ cluster-1 ]
copying token for cluster cluster-1
applying VAULT_SKIP_VERIFY
applying VAULT_TLS_CA
applying VAULT_ADDR
applying VAULT_TOKEN
token for cluster cluster-1 copied to clipboard buffer.
```

---

## About the `Key.Path`-file
for now, `vops` expect the JSON format output from a `vault operator init` command:

```json
{
  "unseal_keys_b64": [
    "YrnZCLIdwKDNn9RYkUx3A7J9/I4ogORIXYcTtJ/AWtg="
  ],
  "unseal_keys_hex": [
    "62b9d908b21dc0a0cd9fd458914c7703b27dfc8e2880e4485d8713b49fc05ad8"
  ],
  "unseal_shares": 1,
  "unseal_threshold": 1,
  "recovery_keys_b64": [],
  "recovery_keys_hex": [],
  "recovery_keys_shares": 0,
  "recovery_keys_threshold": 0,
  "root_token": "hvs.EhCMSSb1uCW1y0aHI1IZ3feO"
}
```

Later you can also just list the unseal/recovery keys in the `vops.yml` aswell, or specifiy pgp encrypted key files.