package log

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var logger Logger

func SetLogger(log Logger) {
	logger = log
}

func Debugf(format string, args ...interface{}) {
	if logger != nil {
		logger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger != nil {
		logger.Infof(format, args...)
	}
}

func Printf(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warnf(format, args...)
		errsNotif.i++
	}
}

func Errorf(format string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(format, args...)
		errsNotif.i++
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger != nil {
		logger.Fatalf(format, args...)
	}
}

func Panicf(format string, args ...interface{}) {
	if logger != nil {
		logger.Panicf(format, args...)
	}
}

func Debug(args ...interface{}) {
	if logger != nil {
		logger.Debug(args...)
	}
}

func Print(args ...interface{}) {
	if logger != nil {
		logger.Print(args...)
	}
}
func Info(args ...interface{}) {
	if logger != nil {
		logger.Info(args...)
	}
}
func Warn(args ...interface{}) {
	if logger != nil {
		logger.Warn(args...)
		errsNotif.i++
	}
}
func Error(args ...interface{}) {
	if logger != nil {
		logger.Error(args...)
		errsNotif.i++
	}
}
func Fatal(args ...interface{}) {
	if logger != nil {
		logger.Fatal(args...)
	}
}
func Panic(args ...interface{}) {
	if logger != nil {
		logger.Panic(args...)
	}
}
