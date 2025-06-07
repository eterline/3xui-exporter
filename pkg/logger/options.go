package logger

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type (
	EnvValue         int
	LoggerOptionFunc func(c *LoggerOptions)
)

func (v EnvValue) String() string {
	l := [3]string{"LOCAL", "DEVELOP", "PRODUCTION"}
	return l[v+1]
}

type LoggerOptions struct {
	level logrus.Level
	env   EnvValue

	path     string
	filename string

	isTimestamp bool
	isFormat    bool
	json        bool

	fmtLogs string
	fmtTime string
}

func WithJSONFormat() LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.json = true
	}
}

func WithJSONFormatValue(v bool) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.json = v
	}
}

func WithFmtLogs(f string) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.fmtLogs = fmt.Sprintf("%s\n", f)
	}
}

func WithFmtTime(f string) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.fmtTime = f
	}
}

func WithEnv(env EnvValue) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		switch env {
		case LOCAL:
			o.level = logrus.InfoLevel
		case DEVELOP:
			o.level = logrus.TraceLevel
		case PRODUCTION:
			o.level = logrus.WarnLevel
		}
	}
}

func WithDevEnvBool(value bool) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		if value {
			o.level = logrus.TraceLevel
			return
		}
		o.level = logrus.InfoLevel
	}
}

func WithLevel(lvl string) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		if l, err := logrus.ParseLevel(strings.ToLower(lvl)); err == nil {
			o.level = l
		}
	}
}

func WithPretty() LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.isFormat = true
	}
}

func WithPrettyValue(v bool) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.isFormat = v
	}
}

func WithTimestamp() LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.isTimestamp = true
	}
}

// WithTimestampValue set timestamp bool.
func WithTimestampValue(v bool) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.isTimestamp = v
	}
}

// WithPath set log file directory
func WithPath(path string) LoggerOptionFunc {
	return func(o *LoggerOptions) {
		o.path = filepath.Clean(path)
	}
}

func mustOptions(options ...LoggerOptionFunc) *LoggerOptions {

	o := &LoggerOptions{
		level: logrus.TraceLevel,

		path:     filepath.Clean("./"),
		filename: "trace",

		isTimestamp: false,
		isFormat:    false,

		fmtLogs: "[%lvl%]: %time% - %msg%\n",
		fmtTime: "2006-01-02 15:04:05",
	}

	for _, option := range options {
		option(o)
	}

	return o
}

func returnName(opts *LoggerOptions) (name string) {

	if opts.isTimestamp {
		name = fmt.Sprintf(
			"%s.%s.log", opts.filename, time.Now().Format(time.RFC3339),
		)
	} else {
		name = fmt.Sprintf("%s.log", opts.filename)
	}

	if opts.path != "" {
		name = filepath.Join(opts.path, name)
	}

	return
}

func returnFormatter(opts *LoggerOptions) (fmt logrus.Formatter) {

	if opts.isFormat && !opts.json {
		fmt = &easy.Formatter{
			TimestampFormat: opts.fmtTime,
			LogFormat:       opts.fmtLogs,
		}
	}

	if opts.json {
		fmt = &logrus.JSONFormatter{}
	}

	return
}
