<!-- generated-from:5eaae48b1b4eacced51fe9c8a724bf608d4c39edfb4a2fa3df458b8a7ea5e763 DO NOT REMOVE, DO UPDATE -->
# ACH Web Viewer
**[Purpose](README.md)** | **[Configuration](CONFIGURATION.md)** | **Running**

---

## Running

### Getting Started

More tutorials to come on how to use this as other pieces required to handle authorization are in place!

- [Using docker-compose](#local-development)
- [Using our Docker image](#docker-image)

No configuration is required to serve on `:8200` and metrics at `:8201/metrics` in Prometheus format.

### Docker image

You can download [our docker image `moov/ach-web-viewer`](https://hub.docker.com/r/moov/ach-web-viewer/) from Docker Hub or use this repository.

### Local Development

```
go run ./cmd/ach-web-viewer
```
