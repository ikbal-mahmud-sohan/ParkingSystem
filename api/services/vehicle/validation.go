package vehicle

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pvpb "parking/gunk/api/vehicle"
	"parking/utility/logging"
)

func (s *Svc) ParkVehicleValidation(ctx context.Context, req *pvpb.ParkVehicleRequest) error {
	log := logging.FromContext(ctx).WithField("method", "Service.UpdateParkingSlotValidation")

	if req == nil {
		return status.Error(codes.InvalidArgument, "request can not be empty")
	}
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.VehicleNo, required),
		validation.Field(&req.LotID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate park vehicle")
		return err
	}
	return nil
}

func (s *Svc) UnParkVehicleValidation(ctx context.Context, req *pvpb.UnParkVehicleRequest) error {
	log := logging.FromContext(ctx).WithField("method", "Service.UpdateParkingSlotValidation")

	if req == nil {
		return status.Error(codes.InvalidArgument, "request can not be empty")
	}
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required),
		validation.Field(&req.VehicleNo, required),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate unpark vehicle")
		return err
	}
	return nil
}
