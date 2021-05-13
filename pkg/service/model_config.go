// generated-from:9f67de2487f9ee84a361869beb03771ef499cc8b5e5ca3cedb827982388792b7 DO NOT REMOVE, DO UPDATE

package service

type GlobalConfig struct {
	ACHWebViewer Config
}

// Config defines all the configuration for the app
type Config struct {
	Servers ServerConfig
	Clients *ClientConfig
}

// ServerConfig - Groups all the http configs for the servers and ports that get opened.
type ServerConfig struct {
	Public HTTPConfig
	Admin  HTTPConfig
}

// HTTPConfig configuration for running an http server
type HTTPConfig struct {
	Bind BindAddress
}

// BindAddress specifies where the http server should bind to.
type BindAddress struct {
	Address string
}
