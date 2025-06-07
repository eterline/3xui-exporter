package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type StorageLogger struct {
	logger *logrus.Logger
}

func InitStorageLogger() *StorageLogger {
	return &StorageLogger{
		logger: entry.Logger,
	}
}

func (sl *StorageLogger) LogMode(lvl logger.LogLevel) logger.Interface {
	return sl
}

func (sl *StorageLogger) Info(ctx context.Context, s string, args ...interface{}) {
	sl.logger.Infof(s, args...)
}

func (sl *StorageLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	sl.logger.Warnf(s, args...)
}

func (sl *StorageLogger) Error(ctx context.Context, s string, args ...interface{}) {
	sl.logger.Errorf(s, args...)
}

func (sl *StorageLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}
