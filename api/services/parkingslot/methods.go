package parkingslot

import (
	"context"

	pspb "parking/gunk/api/parkingslot"
	uErr "parking/utility/error/error"
	"parking/utility/logging"
)

func (s *Svc) UpdateParkingSlot(ctx context.Context, req *pspb.UpdateParkingSlotRequest) (*pspb.UpdateParkingSlotResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.UpdateParkingSlot")
	log.Trace("request received")

	if err := s.UpdateParkingSlotValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, uErr.HandleServiceErr(err)
	}

	res, err := s.core.UpdateParkingSlot(ctx, req)
	if err != nil {
		errM := "failed to Update parking slot"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}
