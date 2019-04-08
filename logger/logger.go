//Package logger  日志系统，提供一个默认的控制台日志
package logger

import (
	"fmt"
	"io"
	"log"
)

// Level 日志等级类型,用SetLevel设置日志等级
// 当大于此日志等级的日志将不会输出
// DEBUG>INFO>WARN>ERROR>NONE
type Level int

// NONE 无日志
const (
	NONE  Level = iota
	ERROR       // ERROR 错误日志
	WARN        // WARN 警告日志
	INFO        // INFO 信息日志
	DEBUG       // DEBUG 调试日志
)

// LevelStr 默认日志等级字串，表示此行日志等级
var LevelStr = []string{"", "ERR", "WAR", "INF", "DEB"}

// Logger 日志器
type Logger struct {
	pattern Pattern
	level   Level
	name    string
	log     *log.Logger
}

// Pattern 日志输出模板接口
type Pattern interface {
	// Format 格式化每行日志,返回格式化好的日志
	// level 输出日志等级
	// format 范式 参阅fmt
	Format(level Level, a ...interface{}) string
	// Roll 滚动日志
	// 如需滚动日志返回一个非空io.Writer
	// 否则继续使用上一次的io.Writer
	// 所以第一次Roll被调用，理论上应该始终返回非空io.Writer
	Roll() io.Writer
}

type samplePattern struct{}

func (out *samplePattern) Format(level Level, a ...interface{}) string {
	return fmt.Sprint(LevelStr[level], ": ", fmt.Sprint(a...))
}

func (out *samplePattern) Roll() io.Writer {
	return nil
}

// SetName 设置默认日志名字，此字段将出现在日志最前面
func SetName(name string) {
	logger.name = name
	if logger.name == "" {
		log.SetPrefix("")
	} else {
		log.SetPrefix(logger.name + " ")
	}
}

// SetLevel 设置默认日志等级
func SetLevel(level Level) {
	logger.level = level
}

// Debugf 输出调试日志
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

// Debug 输出调试日志
func Debug(a ...interface{}) {
	output(DEBUG, a...)
}

// Infof 输出信息日志
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

// Info 输出信息日志
func Info(a ...interface{}) {
	output(INFO, a...)
}

// Warnf 输出警告日志
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

// Warn 输出警告日志
func Warn(a ...interface{}) {
	output(WARN, a...)
}

// Errorf 输出错误日志
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

// Error 输出错误日志
func Error(a ...interface{}) {
	output(ERROR, a...)
}

//默认logger，只输出控制台
var logger = &Logger{&samplePattern{}, DEBUG, "", nil}

func init() {
	SetName(logger.name)
	log.SetFlags(log.Ltime)
}

func output(level Level, a ...interface{}) {
	if logger.level >= level {
		log.Println(logger.pattern.Format(level, a...))
	}
}

// SetName 设置日志名字，此字段将出现在日志最前面
func (l *Logger) SetName(name string) {
	l.name = name
	if l.name == "" {
		log.SetPrefix("")
	} else {
		log.SetPrefix(logger.name + " ")
	}
}

// SetLevel 设置日志等级
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Debugf 输出调试日志
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Output(DEBUG, fmt.Sprintf(format, a...))
}

// Infof 输出信息日志
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Output(INFO, fmt.Sprintf(format, a...))
}

// Warnf 输出警告日志
func (l *Logger) Warnf(format string, a ...interface{}) {
	l.Output(WARN, fmt.Sprintf(format, a...))
}

// Errorf 输出错误日志
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Output(ERROR, fmt.Sprintf(format, a...))
}

// Debug 输出调试日志
func (l *Logger) Debug(a ...interface{}) {
	l.Output(DEBUG, a...)
}

// Info 输出信息日志
func (l *Logger) Info(a ...interface{}) {
	l.Output(INFO, a...)
}

// Warn 输出警告日志
func (l *Logger) Warn(a ...interface{}) {
	l.Output(WARN, a...)
}

// Error 输出错误日志
func (l *Logger) Error(a ...interface{}) {
	l.Output(ERROR, a...)
}

// New 新日志等级
// name 新日志名字
func New(pattern Pattern, level Level, name string) *Logger {
	return &Logger{pattern, level, name, log.New(pattern.Roll(), name, 0)}
}

// Output 输出指定等级日志
// level 指定日志等级
func (l *Logger) Output(level Level, a ...interface{}) {
	if l.level >= level {
		if writer := l.pattern.Roll(); writer != nil {
			l.log.SetOutput(writer)
		}
		l.log.Println(l.pattern.Format(level, a...))
	}
}
