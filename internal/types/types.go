package types

type Config struct {
	Env    string
	Server struct {
		AllowOrigins string `yaml:"allowOrigins"`
		AllowHeaders string `yaml:"allowHeaders"`
		AllowMethods string `yaml:"allowMethods"`
		BodyLimit    int    `yaml:"bodyLimit"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
	}
	Process struct {
		UploadDir string `yaml:"uploadDir"`
	}
	DB struct {
		DSN string
	}
	App struct {
		Name            string `yaml:"name"`
		JobFlushTimeout int    `yaml:"jobFlushTimeout"`
		MetricsUser     string
		MetricsPassword string
	}
	Firebase struct {
		AppCheckHeader string `yaml:"appCheckHeader"`
		ProjectID      string
	}
}

type ImageProcessInput struct {
	JobID    int
	FileID   int
	FileName string
	Format   string
	Quality  int
	Width    int
	Height   int
	Buffer   []byte
}

type SuccessResult struct {
	SourceFileName string
	SourceFileID   int
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
