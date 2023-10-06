package types

type SuccessResult struct {
	SourceFileName string
	SourceFileID   int
	TargetFileName string
	SourceFileSize int64
	TargetFileSize int64
}

type AnyMap map[string]any
