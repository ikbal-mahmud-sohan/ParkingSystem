package parkingstatement

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"parking/api/storage"
)

type StatementStore interface {
	GetParkingStatement(ctx context.Context, startDate, endDate time.Time) (*storage.Statement, error)
}

type Svc struct {
	store  StatementStore
	logger *logrus.Entry
}

func New(store StatementStore, logger *logrus.Entry) *Svc {
	return &Svc{
		store:  store,
		logger: logger,
	}
}
