package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *writerHook) Fire(entr *logrus.Entry) error {
	str, err := entr.String()

	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		w.Write([]byte(str))
	}

	return err
}

func (h *writerHook) Levels() []logrus.Level {
	return h.LogLevels
}
