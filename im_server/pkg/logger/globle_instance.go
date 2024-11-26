package logger

import "sync"

var gl Logger
var lMutex sync.RWMutex

func SetGlobalLogger(logger Logger) {
	lMutex.Lock()
	defer lMutex.Unlock()
	gl = logger
}

func L() Logger {
	lMutex.RLock()
	g := gl
	lMutex.RUnlock()
	return g
}
