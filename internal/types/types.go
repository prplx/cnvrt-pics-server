package types

type Config struct {
	Env    string
	Server struct {
		AllowOrigins    string `yaml:"allowOrigins"`
		AllowHeaders    string `yaml:"allowHeaders"`
		AllowMethods    string `yaml:"allowMethods"`
		BodyLimit       int    `yaml:"bodyLimit"`
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		ShutdownTimeout int    `yaml:"shutdownTimeout"`
	}
	Process struct {
		UploadDir string `yaml:"uploadDir"`
	}
	DB struct {
		DSN string
	}
	App struct {
		Name               string `yaml:"name"`
		JobFlushTimeout    int    `yaml:"jobFlushTimeout"`
		MetricsUser        string
		MetricsPassword    string
		SupportedFileTypes string `yaml:"supportedFileTypes"`
		MaxFileCount       int    `yaml:"maxFileCount"`
	}
	Firebase struct {
		AppCheckHeader string `yaml:"appCheckHeader"`
		ProjectID      string
	}
}

type ImageProcessInput struct {
	JobID    int64
	FileID   int64
	FileName string
	Format   string
	Quality  int
	Width    int
	Height   int
	Buffer   []byte
}

type SuccessResult struct {
	SourceFileName string
	SourceFileID   int64
	TargetFileName string
	SourceFileSize int64
	TargetFileSize int64
	Width          int
	Height         int
	OriginalWidth  int
	OriginalHeight int
	Format         string
	Quality        int
}

type AnyMap map[string]any

type WebsocketConnection interface {
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
	Params(key string, defaultValue ...string) string
	WriteJSON(v interface{}) error
}
