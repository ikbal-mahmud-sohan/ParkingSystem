package postgres

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"parking/api/storage"
	"parking/utility/logging"
)

func (s *Storage) UpdateParkingSlotValidation(ctx context.Context, req *storage.UpdateParkingSlotRequest) error {
	log := logging.FromContext(ctx).WithField("method", "storage.UpdateParkingSlotValidation")

	if req == nil {
		return errors.New("request can not be empty")
	}

	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.LotID, required, is.UUID),
		validation.Field(&req.SlotNumber, required),
		validation.Field(&req.Status, required),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate Upsert ParkinSlot")
		return err
	}

	return nil
}
