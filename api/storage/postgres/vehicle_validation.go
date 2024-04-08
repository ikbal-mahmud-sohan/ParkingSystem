package postgres

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"parking/api/storage"
	"parking/utility/logging"
)

func (s *Storage) ParkVehicleValidation(ctx context.Context, req *storage.ParkVehicleRequest) error {
	log := logging.FromContext(ctx).WithField("method", "storage.ParkVehicleValidation")
	if req == nil {
		return errors.New("request can not be empty")
	}

	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.VehicleNo, required),
		validation.Field(&req.Entry, required),
		validation.Field(&req.SlotID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate park vehicle request")
		return err
	}

	return nil
}

func (s *Storage) UnParkVehicleValidation(ctx context.Context, req *storage.UnParkVehicleRequest) error {
	log := logging.FromContext(ctx).WithField("method", "storage.UnParkVehicleValidation")
	if req == nil {
		return errors.New("request can not be empty")
	}

	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate unpark vehicle request")
		return err
	}

	return nil
}


func (s *Storage) GetVehicleValidation(ctx context.Context, req *storage.VehicleRequest) error {
	log := logging.FromContext(ctx).WithField("method", "storage.GetVehicleValidation")

	if req == nil {
		return errors.New("request can not be empty")
	}

	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required),
		validation.Field(&req.VehicleNo, required),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate get vehicle request")
		return err
	}

	return nil
}
