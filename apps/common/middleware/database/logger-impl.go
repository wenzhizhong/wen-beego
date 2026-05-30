package database

import (
	"WenBeego/apps/common/global"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// "WenBeego/apps/common/global"

// 实现接口
type loggerImpl struct {
	Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

var _ logger.Interface = (*loggerImpl)(nil)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogLevel log level
type LogLevel int

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Config logger config
type Config struct {
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  logger.LogLevel
}

var (
	// Discard logger will print any log to io.Discard
	Discard = New(log.New(io.Discard, "", log.LstdFlags), Config{})
	// Default Default logger
	Default = New(log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	// Recorder logger records running SQL into a recorder instance
	Recorder = traceRecorder{Interface: Default, BeginAt: time.Now()}

	// RecorderParamsFilter defaults to no-op, allows to be run-over by a different implementation
	RecorderParamsFilter = func(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
		return sql, params
	}
)

// New initialize logger
func New(writer Writer, config Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &loggerImpl{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *loggerImpl) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l *loggerImpl) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *loggerImpl) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *loggerImpl) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		global.Log.Error(filterColor(l.errStr+msg), append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l *loggerImpl) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)

			if !errors.Is(err, ErrRecordNotFound) {
				fmt.Println(fmt.Sprintf(filterColor(l.traceErrStr), utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			}
		} else {
			l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				global.Log.Error(fmt.Sprintf(filterColor(l.traceErrStr), utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		}

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// ParamsFilter filter params
func (l *loggerImpl) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

type traceRecorder struct {
	logger.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

// New trace recorder
func (l *traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

// Trace implement logger interface
func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}

func (l *traceRecorder) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if RecorderParamsFilter == nil {
		return sql, params
	}
	return RecorderParamsFilter(ctx, sql, params...)
}

func filterColor(str string) string {
	colorArr := []string{
		logger.Reset,       // = "\033[0m"
		logger.Red,         // = "\033[31m"
		logger.Green,       // = "\033[32m"
		logger.Yellow,      // = "\033[33m"
		logger.Blue,        // = "\033[34m"
		logger.Magenta,     // = "\033[35m"
		logger.Cyan,        // = "\033[36m"
		logger.White,       // = "\033[37m"
		logger.BlueBold,    // = "\033[34;1m"
		logger.MagentaBold, // = "\033[35;1m"
		logger.RedBold,     // = "\033[31;1m"
		logger.YellowBold,  // = "\033[33;1m"
	}
	for _, color := range colorArr {
		str = strings.ReplaceAll(str, color, "")
	}
	return str
}
