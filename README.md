# vops

# CLI

## keys
* vops keys unseal [ rotate, status, ...]
* vops keys recovery [rotate, status, ...]
* vops root [ regen, revoke, ...]
* vops unseal 
* vops seal

## replication
* vops pr [ enable, disable, status ] ---primary=... --secondary=...
* vops dr [ enable, disable, status ] ---primary=... --secondary=...


# Config file
```yaml
vaults:
    - name: cluster-1
      vault_addr: ...
      vault_token: ....
      vault_token_exec: vault login {{ vault_addr }} ... 
      envs:
        - VAULT_TLS_SKIP_VERIY: true
        - ...
    - name: cluster-2
      vault_addr: ...
      vault_token: ....
      vault_token_exec: vault login {{ vault_addr }} ... 
      envs:
        - VAULT_TLS_SKIP_VERIY: true
        - ... 
```