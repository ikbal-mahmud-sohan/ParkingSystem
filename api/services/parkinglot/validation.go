package parkinglot

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	plpb "parking/gunk/api/parkinglot"
	"parking/utility/logging"
)

func (s *Svc) InsertParkingLotValidation(ctx context.Context, req *plpb.InsertParkingLotRequest) error {
	log := logging.FromContext(ctx).WithField("method", "Service.InsertParkingLotValidation")

	if req == nil {
		return status.Error(codes.InvalidArgument, "request can not be empty")
	}
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Name, required, validation.Length(2, 50)),
		validation.Field(&req.TotalSpace, required),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate CreateParkingLot")
		return err
	}
	return nil
}

func (s *Svc) GetParkingLotValidation(ctx context.Context, req *plpb.GetParkingLotRequest) error {
	log := logging.FromContext(ctx).WithField("method", "Service.GetParkingLotValidation")

	if req == nil {
		return status.Error(codes.InvalidArgument, "request can not be empty")
	}
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate get parking details request")
		return err
	}
	return nil
}
