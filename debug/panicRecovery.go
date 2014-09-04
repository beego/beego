package debug

import (
	"fmt"
)

type interfaceRecoveryLogger interface {
	Emergency(format string, v ...interface{})
}

var recoveryLogger interfaceRecoveryLogger

func GoRoutineRecovered(anonFunc func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if recoveryLogger != nil {
					recoveryLogger.Emergency("%+v", r)
				} else {
					fmt.Println("EMERGENCY:", r)
				}
			}
		}()
		anonFunc()
	}()
}

func SetRecoveryLogger(newRecoveryLogger interfaceRecoveryLogger) {
	recoveryLogger = newRecoveryLogger
}
