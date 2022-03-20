package logx

type Fields map[string]interface{}

// Level type
type Level int

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}
