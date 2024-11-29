package logger

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Error(key string, err error) Field {
	return Field{
		Key:   key,
		Value: err,
	}
}
