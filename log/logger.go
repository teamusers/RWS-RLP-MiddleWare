package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

// InitLogger initializes Logger, supporting daily automatic backups and retaining logs for the past three days.
func InitLogger(path string) {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true) // Enable caller information.

	// 创建日志目录
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("❌ Unable to create log directory: %v\n", err)
		return
	}

	// Log file rotation configuration (rotate logs daily, retaining the most recent 3 days).
	infoLog := &lumberjack.Logger{
		Filename:   path + "/info.log",
		MaxSize:    10,   // Single log file maximum 10MB
		MaxAge:     3,    // Only retain logs for the most recent 3 days.
		MaxBackups: 3,    // Keep at most 3 backups."
		Compress:   true, // Enable compression.
	}

	errorLog := &lumberjack.Logger{
		Filename:   path + "/error.log",
		MaxSize:    10,
		MaxAge:     3,
		MaxBackups: 3,
		Compress:   true,
	}

	// 日志输出到文件和控制台
	multiWriterInfo := io.MultiWriter(os.Stdout, infoLog)
	multiWriterError := io.MultiWriter(os.Stdout, errorLog)

	logger.SetOutput(os.Stdout)

	// Hook: Write INFO level logs to info.log.
	logger.AddHook(&FileHook{
		Writer: multiWriterInfo,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.WarnLevel,
		},
	})

	// Hook: Write ERROR level logs to error.log.
	logger.AddHook(&FileHook{
		Writer: multiWriterError,
		LogLevels: []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		},
	})

	logger.Info("✅ Logger Initialization successful, logs are automatically rotated, backed up daily, and the most recent 3 days are retained.")
}

// FileHook custom hook, supports writing to log files.
type FileHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire implements a logrus Hook that writes logs to a file.
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	entry.Data["file"], entry.Data["func"] = getCaller() // Record caller information.
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels returns the log levels applicable to the hook.
func (hook *FileHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Get the code location of the log call.
func getCaller() (string, string) {
	for i := 3; i < 15; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(pc).Name()
		if !isLogrusFile(file) {
			return fmt.Sprintf("%s:%d", trimGoPath(file), line), funcName
		}
	}
	return "unknown", "unknown"
}

// Filter logrus's own calls.
func isLogrusFile(file string) bool {
	return strings.Contains(file, "logrus") || strings.Contains(file, "logger.go")
}

// Remove GOPATH or module path to optimize file display.
func trimGoPath(file string) string {
	goPath := os.Getenv("GOPATH")
	if goPath != "" {
		goPathPrefix := goPath + "/src/"
		if strings.HasPrefix(file, goPathPrefix) {
			return strings.TrimPrefix(file, goPathPrefix)
		}
	}
	return file
}

// **Log method encapsulation.**
func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Get logger instance.
func GetLogger() *logrus.Logger {
	return logger
}
