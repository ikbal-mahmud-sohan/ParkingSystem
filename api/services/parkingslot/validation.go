package parkingslot

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pspb "parking/gunk/api/parkingslot"
	"parking/utility/logging"
)

func (s *Svc) UpdateParkingSlotValidation(ctx context.Context, req *pspb.UpdateParkingSlotRequest) error {
	log := logging.FromContext(ctx).WithField("method", "Service.UpdateParkingSlotValidation")

	if req == nil {
		return status.Error(codes.InvalidArgument, "request can not be empty")
	}
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.LotID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("failed to validate CreateParkingLot")
		return err
	}
	return nil
}
