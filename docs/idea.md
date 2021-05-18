## Example configs

```yaml
ACHWebViewer:
  servers:
    public:
      bind:
        address: ":8585"
    admin:
      bind:
        address: ":9595"

display:
  format: "human-readable"
  masking:
    accountNumbers: true
    names: false

sources:
  - id: "audittrail"
    bucket:
      url: "gs://moov-paygate-audittrail/"
      paths:
        - "files"
    encryption:
      gpg:
        keyFile: "/keys/audit.priv"
        keyPassword: "<secret>"

  - id: "mergable"
    filesystem:
      paths:
        - "/storage/mergable/"
```
