<!-- generated-from:06a7eb50ee0171a2db8c865ec4dc51334a5b84d4afc551199e1a8048e78fa69c DO NOT REMOVE, DO UPDATE -->
# ACH Web Viewer
**[Purpose](README.md)** | **Configuration** | **[Running](RUNNING.md)** | **[Client](../pkg/client/README.md)**

---

## Configuration
Custom configuration for this application may be specified via an environment variable `APP_CONFIG` to a configuration file that will be merged with the default configuration file.

- [Default Configuration](../configs/config.default.yml)
- [Config Source Code](../pkg/service/model_config.go)
- Full Configuration
```yaml
  ACH Web Viewer:

    # Service configurations
    Servers:

      # Public service configuration
      Public:
        Bind:
          # Address and port to listen on.
          Address: ":8585"

      # Health/Admin service configuration.
      Admin:
        Bind:
          # Address and port to listen on.
          Address: ":9595"
```

---
**[Next - Running](RUNNING.md)**
