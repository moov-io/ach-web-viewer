// generated-from:9f67de2487f9ee84a361869beb03771ef499cc8b5e5ca3cedb827982388792b7 DO NOT REMOVE, DO UPDATE

package service

import (
	"time"
)

type GlobalConfig struct {
	ACHWebViewer Config
}

// Config defines all the configuration for the app
type Config struct {
	Servers ServerConfig
	Clients *ClientConfig

	Display DisplayConfig
	Sources Sources
}

// ServerConfig - Groups all the http configs for the servers and ports that get opened.
type ServerConfig struct {
	Public HTTPConfig
	Admin  HTTPConfig
}

// HTTPConfig configuration for running an http server
type HTTPConfig struct {
	Bind     BindAddress
	BasePath string
}

// BindAddress specifies where the http server should bind to.
type BindAddress struct {
	Address string
}

type DisplayConfig struct {
	Format       string // e.g. "human-readable"
	Masking      MaskingConfig
	HelpfulLinks HelpfulLinks
}

type MaskingConfig struct {
	AccountNumbers bool
	CorrectedData  bool
	Names          bool

	PrettyAmounts bool
}

type HelpfulLinks struct {
	Corrections string
	Returns     string
}

type Sources []Source

type Source struct {
	ID string

	ACHGateway *ACHGatewayConfig
	Bucket     *BucketConfig
	Filesystem *FilesystemConfig

	Encryption *EncryptionConfig
}

type ACHGatewayConfig struct {
	Endpoint string
	Timeout  time.Duration
	Shards   []string
}

type BucketConfig struct {
	URL   string
	Paths []string
}

type FilesystemConfig struct {
	Paths []string
}

type EncryptionConfig struct {
	GPG *GPG
}

type GPG struct {
	Files []GPGFile

	// KeyFile and KeyPassword are deprecated, use .Files
	KeyFile     string
	KeyPassword string
}

type GPGFile struct {
	KeyFile     string
	KeyPassword string
}
