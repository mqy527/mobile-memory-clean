package log

import (
	"os"
	"path"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

func GetLogger(name string) Logger {
	lg := logrus.New()
	lg.AddHook(&loggerNameHook{
		name: name,
	})
	initLogger(lg)
	return lg
}

type loggerNameHook struct {
	name string
}

func (h *loggerNameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *loggerNameHook) Fire(entry *logrus.Entry) error {
	entry.Data["NAME"] = h.name
	return nil
}

func initLogger(logger *logrus.Logger) {
	logLevel := os.Getenv("logLevel")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	logPath := os.Getenv("logPath")
	logFileName := "mobile-memory-clean.log"
	maxAge := 30 * 24 * time.Hour
	rotationTime := 24 * time.Hour
	configLocalFilesystemLogger(logger, logPath, logFileName, maxAge, rotationTime)
}

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(logger *logrus.Logger, logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	fm := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldsOrder:     []string{"NAME"},
		NoColors:        true,
	}
	logger.SetFormatter(fm)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, fm)
	logger.AddHook(lfHook)
}
