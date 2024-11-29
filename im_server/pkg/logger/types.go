package logger

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
	err   error
}

func LoggerExample() {
	var l Logger
	phone := "187****1041"
	l.Info("用户未注册", Field{
		Key:   "phone",
		Value: phone,
	})
}
