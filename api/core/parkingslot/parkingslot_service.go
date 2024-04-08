package parkingslot

import (
	"context"

	"github.com/sirupsen/logrus"

	"parking/api/storage"
)

type ParkingslotStore interface {
	UpdateParkingSlot(ctx context.Context, req *storage.UpdateParkingSlotRequest) (*storage.ParkingSlotResponse, error)
	GetParkinglotByID(ctx context.Context, lotID string) (*storage.ParkingLotResponse, error)
}

type Svc struct {
	store  ParkingslotStore
	logger *logrus.Entry
}

func New(store ParkingslotStore, logger *logrus.Entry) *Svc {
	return &Svc{
		store:  store,
		logger: logger,
	}
}
