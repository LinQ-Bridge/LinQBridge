package network

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type Nologger struct {
}

func (n Nologger) LogMode(level logger.LogLevel) logger.Interface {
	return &Nologger{}
}

func (n Nologger) Info(ctx context.Context, s string, i ...interface{}) {

}

func (n Nologger) Warn(ctx context.Context, s string, i ...interface{}) {

}

func (n Nologger) Error(ctx context.Context, s string, i ...interface{}) {

}

func (n Nologger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

}
