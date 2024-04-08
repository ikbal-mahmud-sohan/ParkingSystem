package vehicle

import (
	"context"

	pvpb "parking/gunk/api/vehicle"
	uErr "parking/utility/error/error"
	"parking/utility/logging"
)

func (s *Svc) ParkVehicle(ctx context.Context, req *pvpb.ParkVehicleRequest) (*pvpb.ParkVehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.ParkVehicle")
	log.Trace("request received")

	if err := s.ParkVehicleValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, uErr.HandleServiceErr(err)
	}

	res, err := s.core.ParkVehicle(ctx, req)
	if err != nil {
		errM := "failed to park vehicle in core"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) UnParkVehicle(ctx context.Context, req *pvpb.UnParkVehicleRequest) (*pvpb.UnParkVehicleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.UnParkVehicle")
	log.Trace("request received")

	if err := s.UnParkVehicleValidation(ctx, req); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, uErr.HandleServiceErr(err)
	}

	res, err := s.core.UnParkVehicle(ctx, req)
	if err != nil {
		errM := "failed to park vehicle in core"
		log.WithError(err).Error(errM)
		return nil, uErr.HandleServiceErr(err)
	}

	return res, nil
}
