package parkingstatement

import (
	"context"
	"parking/utility/logging"

	pspb "parking/gunk/api/parkingstatement"
	uErr "parking/utility/error/error"
)

func (s *Svc) GetParkingStatement(ctx context.Context, req *pspb.GetParkingStatementRequest) (*pspb.GetParkingStatementResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.GetParkingStatement")
	log.Trace("request received")
	res, err := s.core.GetParkingStatement(ctx, req)
	if err != nil {
		errM := "failed to park vehicle in core"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}
