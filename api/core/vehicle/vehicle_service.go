package vehicle

import (
	"context"

	"github.com/sirupsen/logrus"

	"parking/api/storage"
)

const FEE_PER_HOUR = int32(10)

type VehicleStore interface {
	ParkVehicle(ctx context.Context, req *storage.ParkVehicleRequest) (*storage.VehicleResponse, error)
	UnParkVehicle(ctx context.Context, req *storage.UnParkVehicleRequest) (*storage.VehicleResponse, error)
	GetParkingSlotsByLotID(ctx context.Context, lotID string) ([]*storage.ParkingSlotResponse, error)
	UpdateParkingSlot(ctx context.Context, req *storage.UpdateParkingSlotRequest) (*storage.ParkingSlotResponse, error)
	GetVehicle(ctx context.Context, req *storage.VehicleRequest) (*storage.VehicleResponse, error)
}

type Svc struct {
	store  VehicleStore
	logger *logrus.Entry
}

func New(store VehicleStore, logger *logrus.Entry) *Svc {
	return &Svc{
		store:  store,
		logger: logger,
	}
}
