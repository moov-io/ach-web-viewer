<!-- generated-from:06a7eb50ee0171a2db8c865ec4dc51334a5b84d4afc551199e1a8048e78fa69c DO NOT REMOVE, DO UPDATE -->
# ACH Web Viewer
**[Purpose](README.md)** | **Configuration** | **[Running](RUNNING.md)**

---

## Configuration
Custom configuration for this application may be specified via an environment variable `APP_CONFIG` to a configuration file that will be merged with the default configuration file.

- [Default Configuration](../configs/config.default.yml)
- [Config Source Code](../pkg/service/model_config.go)
- Full Configuration
```yaml
ACHWebViewer:
  Servers:
    Public:
      BasePath: "/ach"
      Bind:
        Address: ":8585"

  # Formatting for ACH files
  Display:
    Format: "human-readable"
    Masking:
      AccountNumbers: true
      Names: false
      CorrectedData: true
      PrettyAmounts: false

  Sources:
    - id: "mergable"
      filesystem:
        paths:
          - "./testdata/"

    - id: "audittrail"
      bucket:
        url: "gs://ach-audittrail/"
        paths:
          - "files"
      encryption:
        gpg:
          files:
            - keyFile: "/conf/keys/audittrail.priv"
              keyPassword: "secret"
```

---
**[Next - Running](RUNNING.md)**
