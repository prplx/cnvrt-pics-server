package types

type Config struct {
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
}

type AnyMap map[string]any
